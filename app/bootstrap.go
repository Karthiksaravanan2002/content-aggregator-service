package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/handler"
	"dev.azure.com/daimler-mic/content-aggregator/server"
	"dev.azure.com/daimler-mic/content-aggregator/service"
	"dev.azure.com/daimler-mic/content-aggregator/service/cache"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"

	"go.uber.org/zap"
)

func Start() {

	// 1. Load props
	cfg, err := props.LoadConfig()
	if err != nil {
		panic(err)
	}

	// 2. Logger
	logger, err := server.NewZapLogger(&cfg.Logging, "content-aggregator")
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// 3. Cache
	redisCache := cache.NewRedisCache(cfg.Cache)

	// 4. Provider Factory
	providerFactory := service.NewProviderFactory(
		cfg.Providers,
		cfg.Cache,
		redisCache,
		logger,
	)

	// 5. Content Service
	timeout := time.Duration(cfg.Server.Timeout) * time.Second
	contentSvc := service.NewContentService(providerFactory, logger, timeout)

	// 6. Content Handler
	contentHandler := handler.NewContentHandler(contentSvc, logger, cfg)

	// 7. Server Setup
	srv := server.NewServer(cfg, logger)
	server.ConfigureRoutes(srv, cfg, contentHandler, logger)

	// 8. Start server
	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("server stopped", zap.Error(err))
		}
	}()

	logger.Info("server started", zap.String("address", cfg.Server.Address))

	// 9. Graceful shutdown
	waitForShutdown(srv, logger)
}

func waitForShutdown(srv *server.Server, logger *zap.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown signal received")

	srv.Shutdown()
	logger.Info("server stopped gracefully")
}
