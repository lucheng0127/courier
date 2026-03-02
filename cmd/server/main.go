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

	// TODO: 从配置文件加载数据库连接
	// 数据库连接字符串格式: host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable"
	}

	// 连接数据库
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
	providerSvc := service.NewProviderService(providerRepo)
	authSvc := service.NewAuthService(userRepo)
	usageSvc := service.NewUsageService(usageRepo, userRepo)
	routerSvc := service.NewRouterService()

	// 初始化 Providers（导入 Adapter 包触发 init 注册）
	ctx := context.Background()
	if err := bootstrap.InitProviders(ctx, providerSvc); err != nil {
		log.Printf("Warning: Provider initialization had issues: %v", err)
	}

	// 启动 HTTP 服务器
	router := gin.Default()
	api := router.Group("/api/v1")
	v1 := router.Group("/v1")

	// Provider 管理 API
	providerCtrl := controller.NewProviderController(providerSvc)
	providerCtrl.RegisterRoutes(api)

	// Provider 重载 API（管理员）
	reloadCtrl := controller.NewProviderReloadController(providerSvc)
	adminGroup := api.Group("")
	adminGroup.Use(middleware.AdminAuth())
	reloadCtrl.RegisterRoutes(adminGroup)

	// 用户和 API Key 管理 API（管理员）
	userCtrl := controller.NewUserController(authSvc)
	userGroup := v1.Group("/users")
	userGroup.Use(middleware.AdminAuth())
	userGroup.POST("", userCtrl.CreateUser)
	userGroup.GET("/:id", userCtrl.GetUser)
	userGroup.POST("/:id/api-keys", userCtrl.CreateAPIKey)
	userGroup.GET("/:id/api-keys", userCtrl.ListAPIKeys)
	userGroup.DELETE("/:id/api-keys/:key_id", userCtrl.RevokeAPIKey)

	// 使用统计 API（管理员）
	usageCtrl := controller.NewUsageController(usageSvc)
	usageGroup := v1.Group("/usage")
	usageGroup.Use(middleware.AdminAuth())
	usageGroup.GET("", usageCtrl.GetUsageStats)

	// Chat Completions API（需要 API Key 鉴权）
	chatCtrl := controller.NewChatController(routerSvc, usageSvc)
	chatGroup := v1.Group("")
	chatGroup.Use(middleware.APIKeyAuth(authSvc), middleware.TraceID())
	chatGroup.POST("/chat/completions", chatCtrl.ChatCompletions)

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
