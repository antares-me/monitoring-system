package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/heroku/x/hmetrics/onload"

	"antares-me/monitoring-system/internal/config"
	"antares-me/monitoring-system/internal/delivery/http"
	"antares-me/monitoring-system/internal/repository"
	"antares-me/monitoring-system/internal/server"
	"antares-me/monitoring-system/internal/service"

	"antares-me/monitoring-system/pkg/cache"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		log.Println(err)
		return
	}
	cache := cache.NewCache(cfg.Cache.Expiration, cfg.Cache.CleanupInterval)
	repos := repository.NewRepositories(cfg, cache)
	monService := service.NewMonitoringService(repos)
	handlers := http.NewHandler(*monService)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init())
	go func() {
		if err := srv.Run(); err != nil {
			log.Println(err)
		}
	}()

	log.Println("Server started...")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}
