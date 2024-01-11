package gateway

import (
	"gateWay/internal/config"
	connection "gateWay/internal/grpc"
	"gateWay/internal/mailer"
	"gateWay/internal/storage/redis"
	authv1 "gateWay/pkg/proto/gen/go/auth"
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

	rdb := redis.ConnectToRedis(cfg.RedisClient.Host, cfg.RedisClient.Password, cfg.RedisClient.Port, log)

	mailerSend := mailer.MailerSend{
		APIKey: cfg.Hcaptcha,
	}
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
