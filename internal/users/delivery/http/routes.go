package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/middleware"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users"
)

func MapUserRoute(router *chi.Mux, h users.Handlers, mw *middleware.MiddlewareManager) {
	// User routes
	router.Route("/user", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())
			r.Get("/me", h.Me())
			r.Put("/me", h.UpdateMe())
			// r.Patch("/me/updatepass", h.UpdatePasswordMe())
			// Admin routes
			r.Group(func(r chi.Router) {
				r.Use(mw.SuperUser())
				r.Get("/", h.GetMulti())
				r.Post("/", h.Create())
			})
			// Per id routes
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Get())
				// Admin routes
				r.Group(func(r chi.Router) {
					r.Use(mw.SuperUser())
					r.Delete("/", h.Delete())
					r.Put("/", h.Update())
					r.Patch("/updatepass", h.UpdatePassword())
					r.Get("/logoutall", h.LogoutAllAdmin())
				})
			})
		})
	})
}
