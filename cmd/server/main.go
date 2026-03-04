package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lucheng0127/courier/internal/adapter/openai"
	_ "github.com/lucheng0127/courier/internal/adapter/vllm"
	"github.com/lucheng0127/courier/internal/bootstrap"
	"github.com/lucheng0127/courier/internal/controller"
	"github.com/lucheng0127/courier/internal/logger"
	"github.com/lucheng0127/courier/internal/migrate"
	"github.com/lucheng0127/courier/internal/middleware"
	"github.com/lucheng0127/courier/internal/repository"
	"github.com/lucheng0127/courier/internal/service"
)

const schemaVersion = "v1.0.0"

func main() {
	// 1. 初始化 Logger
	logger.InitFromEnv()
	defer logger.Sync()
	logger.L.Info("Starting Courier LLM Gateway...")

	// 2. 数据库连接
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		logger.L.Fatal("Failed to connect to database",
			zap.Error(err))
	}
	defer db.Close()

	logger.L.Info("Connected to database")

	// 3. 数据库自动迁移
	autoMigrate := os.Getenv("AUTO_MIGRATE")
	if autoMigrate != "false" {
		// 创建 GORM DB 用于迁移
		gormDB, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err != nil {
			logger.L.Fatal("Failed to create GORM DB for migration",
				zap.Error(err))
		}

		migrator := migrate.NewMigrator(gormDB, schemaVersion)
		if err := migrator.Run(); err != nil {
			logger.L.Fatal("Database migration failed",
				zap.Error(err))
		}
	} else {
		logger.L.Warn("Auto migration disabled by AUTO_MIGRATE=false")
	}

	// 4. 初始化 Repository
	providerRepo := repository.NewProviderRepository(db)
	userRepo := repository.NewUserRepository(db)
	usageRepo := repository.NewUsageRepository(db)

	// 5. 初始化 Service
	jwtSvc, err := service.NewJWTService()
	if err != nil {
		logger.L.Fatal("Failed to initialize JWT service",
			zap.Error(err))
	}

	providerSvc := service.NewProviderService(providerRepo)
	authSvc := service.NewAuthService(userRepo, jwtSvc)
	usageSvc := service.NewUsageService(usageRepo, userRepo)
	routerSvc := service.NewRouterService()

	// 6. 确保存在初始管理员用户
	if err := authSvc.EnsureInitialAdmin(context.Background()); err != nil {
		logger.L.Warn("Failed to ensure initial admin",
			zap.Error(err))
	}

	// 7. 初始化 Providers
	ctx := context.Background()
	if err := bootstrap.InitProviders(ctx, providerSvc); err != nil {
		logger.L.Warn("Provider initialization had issues",
			zap.Error(err))
	}

	// 8. 创建路由
	router := gin.Default()

	// 设置路由
	setupRoutes(router, providerSvc, authSvc, usageSvc, routerSvc, jwtSvc)

	// 9. 启动服务器
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.L.Info("HTTP server listening",
			zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatal("Failed to start server",
				zap.Error(err))
		}
	}()

	logger.L.Info("Courier LLM Gateway started")

	// 10. 等待信号退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.L.Info("Shutting down...")

	// 关闭 Usage Service
	if err := usageSvc.Close(); err != nil {
		logger.L.Error("Failed to close usage service",
			zap.Error(err))
	}

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.L.Error("Server forced to shutdown",
			zap.Error(err))
	}

	logger.L.Info("Server exited")
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

	// Provider 管理操作（仅管理员）
	providerCtrl := controller.NewProviderController(providerSvc)
	adminOnly.POST("/providers", providerCtrl.CreateProvider)
	adminOnly.PUT("/providers/:name", providerCtrl.UpdateProvider)
	adminOnly.DELETE("/providers/:name", providerCtrl.DeleteProvider)
	adminOnly.GET("/providers/:name", providerCtrl.GetProvider)

	// Provider 运维（仅管理员）
	reloadCtrl := controller.NewProviderReloadController(providerSvc)
	reloadCtrl.RegisterRoutes(adminOnly)

	// ========== Provider 查询操作（所有认证用户）==========
	jwtAuth.GET("/providers", providerCtrl.ListProviders)
	jwtAuth.GET("/providers/:name/models", providerCtrl.ListProviderModels)

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
