package main

import (
	"gateWay/internal/config"
	connection "gateWay/internal/grpc"
	"gateWay/internal/routers"
	authv1 "github.com/DgekoTT/protos/gen/go/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	userServiceConn, err := connection.ConnectToService(cfg.Auth.Host, cfg.Auth.Port, log)
	if err != nil {
		log.Error("Error on connecting to the Auth-service:", err)
		return
	}
	defer func(userServiceConn *grpc.ClientConn) {
		err := userServiceConn.Close()
		if err != nil {
			log.Warn("Unable to close connection Auth-service:", err)
		}
	}(userServiceConn)

	authClient := authv1.NewAuthClient(userServiceConn)

	//
	//rdb := redis.ConnectToRedis(cfg.RedisClient.Host, cfg.RedisClient.Password, cfg.RedisClient.Port, log)
	//
	//mailerSend := mailer.MailerSend{
	//	APIKey: cfg.Hcaptcha,
	//}

	router := routers.Routes(log, authClient)
	log.Info("starting server:", slog.String("address", cfg.Address))
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
