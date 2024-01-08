package gateway

import (
	"gateWay/internal/config"
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

	//userServiceConn, err := connectToService(os.Getenv("USER_GRPC_HOST"), os.Getenv("USER_GRPC_PORT"))
	//if err != nil {
	//	log.Fatalln("Error on connecting to the user-service:", err)
	//	return
	//}
	//defer userServiceConn.Close()
	//
	//usersClient := user.NewAuthServiceClient(userServiceConn)
	//
	//rdb := connectToRedis(cfg.RedisConfig{
	//	Host:     os.Getenv("REDIS_HOST"),
	//	Port:     os.Getenv("REDIS_PORT"),
	//	Password: os.Getenv("REDIS_PASSWORD"),
	//})
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
