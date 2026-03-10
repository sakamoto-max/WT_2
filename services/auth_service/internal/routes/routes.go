package routes

import (
	"auth_service/internal/handlers"
	"auth_service/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Router(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	

	r.With(middleware.UserSignUpValidator).Post("/signup", h.UserSignUp)
	r.With(middleware.UserLoginValidator).Post("/login", h.UserLogin) // done
	r.With(middleware.JwtMiddleware).Post("/logout", h.UserLogOut)
	r.With(middleware.NewRefreshTokenValidator).Post("/refresh", h.GetNewAccessToken)

	return r
}
