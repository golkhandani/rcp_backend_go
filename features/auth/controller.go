package auth

import (
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	AuthService  IAuthService
	GuardService IGuard
	Router       fiber.Router
}

func (ctl Controller) Register() {
	authRoute := ctl.Router.Group("auth")
	authRoute.Post("/register", ctl.AuthService.Register)
	authRoute.Post("/login", ctl.AuthService.Login)

	// jwtGaurd := ctl.GuardService.JWTGuard()

	authRoute.Post(
		"/refresh",
		/* JWT will be checked inside logic not as middleware */
		ctl.AuthService.Refresh,
	)

	// TODO: to be implemented
	// authRoute.Post("/logout", jwtGaurd, Login(authCollections))
	// authRoute.Delete("/delete-account",jwtGaurd, Login(authCollections))
}
