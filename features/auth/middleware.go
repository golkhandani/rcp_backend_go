package auth

import (
	"net/http"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golkhandani/shopWise/utils"
)

type IGuard interface {
	JWTGuard() func(*fiber.Ctx) error
}

type Guard struct {
	AuthRepo Repo
}

func (g Guard) JWTGuard() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		ContextKey: utils.JWTContextKey,
		SigningKey: utils.JWTSigningKey,
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals(utils.JWTContextKey).(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			uid := claims[utils.JWTUserIDKey].(string)
			user, err := g.AuthRepo.GetAuthByID(uid)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(utils.ErrResult(err, http.StatusUnauthorized))
			}

			c.Locals(utils.LocalUserKey, *user)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(utils.ErrResult(err, http.StatusUnauthorized))
			}
			return c.Next()
		},
	})
}
