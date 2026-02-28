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

	// 初始化 Service
	providerSvc := service.NewProviderService(providerRepo)

	// 初始化 Providers（导入 Adapter 包触发 init 注册）
	ctx := context.Background()
	if err := bootstrap.InitProviders(ctx, providerSvc); err != nil {
		log.Printf("Warning: Provider initialization had issues: %v", err)
	}

	// 启动 HTTP 服务器
	router := gin.Default()
	api := router.Group("/api/v1")

	// Provider 管理 API
	providerCtrl := controller.NewProviderController(providerSvc)
	providerCtrl.RegisterRoutes(api)

	// Provider 重载 API（管理员）
	reloadCtrl := controller.NewProviderReloadController(providerSvc)
	adminGroup := api.Group("")
	adminGroup.Use(middleware.AdminAuth())
	reloadCtrl.RegisterRoutes(adminGroup)

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

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
