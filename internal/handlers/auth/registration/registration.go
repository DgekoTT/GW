package auth

import (
	"context"
	"errors"
	"fmt"
	resp "gateWay/internal/lib/api/response"
	authv1 "github.com/DgekoTT/protos/gen/go/auth"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"regexp"
	"time"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
	userID string `json:"user"`
}

type Registration interface {
	RegisterUser(email, password string) (string, error)
}

func New(log *slog.Logger, client authv1.AuthClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.RegisterUser.new"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {

			log.Error("failed to decode request body", err.Error())

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err = validateEmail(req.Email); err != nil {

			log.Error("wrong email", err.Error())

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		if err = validatePassword(req.Password); err != nil {
			log.Error("incorrect password", err.Error())

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		grpcReq := &authv1.RegisterRequest{
			Email:    req.Email,
			Password: req.Password,
		}
		response, err := client.Register(ctx, grpcReq)
		if err != nil {
			log.Error("Ошибка при регистрации пользователя", err.Error())
			render.JSON(w, r, resp.Error("Ошибка при регистрации"))
			return
		}
		log.Info("url added", slog.String("id", response.UserId))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			userID:   response.UserId,
		})
	}
}

func NewInfo(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.new.info"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		render.JSON(w, r, "yess")
	}
}

func validateEmail(email string) error {
	// Проверка, что Email не пустой
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Проверка формата Email
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func isValidEmail(email string) bool {
	// Регулярное выражение для проверки электронной почты
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Проверка на наличие хотя бы одной заглавной буквы
	hasUpper := regexp.MustCompile(`[A-Z]`)
	if !hasUpper.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Проверка на наличие хотя бы одного специального символа
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+{}:"<>?|[\]\\;',./]`)
	if !hasSpecial.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
