package routers

import (
	auth "gateWay/internal/handlers/auth/registration"
	authv1 "github.com/DgekoTT/protos/gen/go/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"log/slog"
)

func Routes(log *slog.Logger, client authv1.AuthClient) *chi.Mux {
	const requestMaxAge = 300

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		// TODO: Set to validity domains
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           requestMaxAge,
	})

	r.Use(corsMiddleware.Handler)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", auth.New(log, client))
		r.Get("/info", auth.NewInfo(log))
	})

	//// Rate limiting
	//rate, err := limiter.NewRateFromFormatted("1000-H")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Create a store with the redis client.
	//store, err := storeRedis.NewStoreWithOptions(app.Config.RedisClient, limiter.StoreOptions{
	//	Prefix: "chi_limiter",
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Create a new middleware with the limiter instance.
	//rateLimiter := r.NewMiddleware(limiter.New(store, rate))
	//g.Use(rateLimiter)

	return r
}
