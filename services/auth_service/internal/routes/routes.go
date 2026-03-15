package routes

import (
	"auth_service/internal/handlers"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	pkg "wt/pkg/middleware"
	
)

func Router(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)

	r.Post("/signup", h.UserSignUp)
	r.Post("/login", h.UserLogin) // done
	r.With(pkg.JwtMiddleware).Post("/logout", h.UserLogOut)
	r.Post("/refresh", h.GetNewAccessToken)
	
	return r
}
