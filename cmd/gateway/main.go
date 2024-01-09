package gateway

import (
	"gateWay/internal/config"
	connection "gateWay/internal/grpc"
	"google.golang.org/grpc"
	"log/slog"
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
	_ = log

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

	usersClient := user.NewAuthServiceClient(userServiceConn)

	rdb := connectToRedis(cfg.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
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
