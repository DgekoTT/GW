package config

import (
	"gateWay/internal/mailer"
	"gateWay/lib/jwt"
	"github.com/kataras/hcaptcha"
	"github.com/redis/go-redis/v9"
	"time"
	authv1 "youTeam/protos/gen/go/auth"
)

type options struct {
	JWTAuthTTL         time.Duration `yaml:"token_ttl" env-required:"true"`
	JWTVerificationTTL time.Duration `yaml:"token_verification_ttl" env-required:"true"`
	AppURL             string        // Application API URL
	Env                string        `yaml:"env" env-default:"local"`
}

type Config struct {
	options     options
	token       *jwt.TokenMaker
	authClient  authv1.AuthClient
	redisClient *redis.Client
	mailer      mailer.MailerSend
	hcaptcha    *hcaptcha.Client
}
