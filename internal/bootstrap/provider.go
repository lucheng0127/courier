package bootstrap

import (
	"context"
	"log"

	"github.com/lucheng0127/courier/internal/service"
)

// InitProviders 初始化所有 Provider
func InitProviders(ctx context.Context, svc *service.ProviderService) error {
	log.Println("Initializing providers...")

	if err := svc.InitProviders(ctx); err != nil {
		log.Printf("Warning: Some providers failed to initialize: %v", err)
		// 不中断启动，继续运行
	}

	log.Println("Provider initialization completed")
	return nil
}
