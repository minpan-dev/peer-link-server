package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"peer-link-server/config"
	"peer-link-server/internal/database"
	"peer-link-server/internal/handler"
	"peer-link-server/internal/repository"
	"peer-link-server/internal/router"
	"peer-link-server/internal/service"
	"peer-link-server/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 支持通过 -config 指定配置文件路径
	cfgFile := flag.String("config", "", "config file path")
	flag.Parse()

	// 1. 加载配置
	cfg, err := config.Load(*cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	log := logger.MustNew(&cfg.Log)
	defer log.Sync() //nolint:errcheck

	// 3. 初始化数据库
	db, err := database.New(&cfg.Database, cfg.IsProd())
	if err != nil {
		log.Fatal("init database", zap.Error(err))
	}

	// 4. 依赖注入（手动 wire）
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, log)

	handlers := &router.Handlers{
		Health: handler.NewHealthHandler(),
		User:   handler.NewUserHandler(userSvc),
		LiveKit: handler.NewLiveKitHandler(
			service.NewLiveKitService(&cfg.LiveKit, repository.NewRoomRepository(db)),
		),
	}

	// 5. 初始化路由
	r := router.New(handlers, log, cfg.JWT.Secret)

	// 6. 启动 HTTP Server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("server starting", zap.String("addr", addr), zap.String("env", cfg.App.Env))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	// 7. 优雅关闭：等待系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("shutting down", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", zap.Error(err))
	}
	log.Info("server exited")
}
