package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucheng0127/courier/internal/client"
	"github.com/lucheng0127/courier/internal/handler"
	"github.com/lucheng0127/courier/internal/middleware"
	"github.com/lucheng0127/courier/internal/repository"
	"github.com/lucheng0127/courier/internal/router"
	"github.com/lucheng0127/courier/internal/service"
	"github.com/lucheng0127/courier/pkg/config"
	"github.com/lucheng0127/courier/pkg/logger"
	"go.uber.org/zap"
)

var configPath = flag.String("config", "config.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志器
	zapLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("初始化日志器失败: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("启动服务", zap.String("config", *configPath))

	// 初始化数据库
	db, err := repository.InitDB(cfg.DB.DataSourceName)
	if err != nil {
		zapLogger.Fatal("初始化数据库失败", zap.Error(err))
	}

	// 初始化仓储
	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	requestLogRepo := repository.NewRequestLogRepository(db)

	// 初始化服务
	userService := service.NewUserService(userRepo)
	apiKeyService := service.NewAPIKeyService(userRepo, apiKeyRepo)
	modelService := service.NewModelService(cfg.Models)

	// 初始化模型客户端（30秒超时）
	modelClient := client.NewModelClient(30 * time.Second)

	// 初始化聊天服务
	chatService := service.NewChatService(modelService, modelClient, requestLogRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeyService)
	modelHandler := handler.NewModelHandler(modelService)
	chatHandler := handler.NewChatHandler(chatService)

	// 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(apiKeyRepo, userRepo)

	// 设置路由
	r := router.SetupRouter(userHandler, apiKeyHandler, modelHandler, chatHandler, authMiddleware)

	// 启动 HTTP 服务器
	addr := ":" + cfg.Server.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		zapLogger.Info("服务器启动", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("服务器关闭失败", zap.Error(err))
	}

	zapLogger.Info("服务器已关闭")
	fmt.Println("服务已停止")
}
