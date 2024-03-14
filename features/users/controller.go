package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golkhandani/shopWise/features/auth"
)

type UserController struct {
	UserService  IUserService
	GuardService auth.IGuard
	Router       fiber.Router
}

func (ctl UserController) Register() {

	jwtGaurd := ctl.GuardService.JWTGuard()
	userRoute := ctl.Router.Group("users")

	userRoute.Get("/me", jwtGaurd, ctl.UserService.GetUserProfile)
	userRoute.Post("/me", jwtGaurd, ctl.UserService.CreateUserProfile)
}
