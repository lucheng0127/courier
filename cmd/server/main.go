package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	_ "github.com/lucheng0127/courier/internal/adapter/openai"
	_ "github.com/lucheng0127/courier/internal/adapter/vllm"
	"github.com/lucheng0127/courier/internal/bootstrap"
	"github.com/lucheng0127/courier/internal/controller"
	"github.com/lucheng0127/courier/internal/middleware"
	"github.com/lucheng0127/courier/internal/repository"
	"github.com/lucheng0127/courier/internal/service"
)

func main() {
	log.Println("Starting Courier LLM Gateway...")

	// 数据库连接
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	// 初始化 Repository
	providerRepo := repository.NewProviderRepository(db)
	userRepo := repository.NewUserRepository(db)
	usageRepo := repository.NewUsageRepository(db)

	// 初始化 Service
	jwtSvc, err := service.NewJWTService()
	if err != nil {
		log.Fatalf("Failed to initialize JWT service: %v", err)
	}

	providerSvc := service.NewProviderService(providerRepo)
	authSvc := service.NewAuthService(userRepo, jwtSvc)
	usageSvc := service.NewUsageService(usageRepo, userRepo)
	routerSvc := service.NewRouterService()

	// 确保存在初始管理员用户
	if err := authSvc.EnsureInitialAdmin(context.Background()); err != nil {
		log.Printf("Warning: Failed to ensure initial admin: %v", err)
	}

	// 初始化 Providers
	ctx := context.Background()
	if err := bootstrap.InitProviders(ctx, providerSvc); err != nil {
		log.Printf("Warning: Provider initialization had issues: %v", err)
	}

	// 创建路由
	router := gin.Default()

	// 设置路由
	setupRoutes(router, providerSvc, authSvc, usageSvc, routerSvc, jwtSvc)

	// 启动服务器
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("HTTP server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Courier LLM Gateway started")

	// 等待信号退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	// 关闭 Usage Service
	if err := usageSvc.Close(); err != nil {
		log.Printf("Failed to close usage service: %v", err)
	}

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// setupRoutes 设置所有路由
func setupRoutes(router *gin.Engine, providerSvc *service.ProviderService, authSvc *service.AuthService, usageSvc *service.UsageService, routerSvc *service.RouterService, jwtSvc service.JWTService) {
	// API v1 组（管理接口）
	api := router.Group("/api/v1")

	// ========== 认证接口（无需鉴权） ==========
	authCtrl := controller.NewAuthController(authSvc)
	authCtrl.RegisterRoutes(api)

	// ========== 需要 JWT 鉴权的组 ==========
	jwtAuth := api.Group("")
	jwtAuth.Use(middleware.JWTAuth(jwtSvc))

	// ========== 需要 Admin 角色的组 ==========
	adminOnly := jwtAuth.Group("")
	adminOnly.Use(middleware.RequireAdmin())

	// Provider 管理（仅管理员）
	providerCtrl := controller.NewProviderController(providerSvc)
	providerCtrl.RegisterRoutes(adminOnly)

	// Provider 运维（仅管理员）
	reloadCtrl := controller.NewProviderReloadController(providerSvc)
	reloadCtrl.RegisterRoutes(adminOnly)

	// ========== 用户管理接口 ==========
	// 管理员可管理所有用户，普通用户可查看自己、管理自己的 API Key
	userCtrl := controller.NewUserController(authSvc)
	userCtrl.RegisterRoutes(jwtAuth)

	// ========== 使用统计接口 ==========
	// 管理员可查看所有用户，普通用户只能查看自己的
	usageCtrl := controller.NewUsageController(usageSvc)
	usageCtrl.RegisterRoutes(jwtAuth)

	// ========== Chat API（使用 API Key 鉴权） ==========
	v1 := router.Group("/v1")
	chatCtrl := controller.NewChatController(routerSvc, usageSvc)
	chatGroup := v1.Group("")
	chatGroup.Use(middleware.APIKeyAuth(authSvc), middleware.TraceID())
	chatCtrl.RegisterRoutes(chatGroup)
}
