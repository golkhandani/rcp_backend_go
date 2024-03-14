package auth

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golkhandani/shopWise/utils"
)

type IAuthService interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	// Logout(c *fiber.Ctx) error
	// DeleteAccout(c *fiber.Ctx) error
}

type Service struct {
	AuthRepo Repo
}

func (as Service) Register(c *fiber.Ctx) error {
	registerRequest := new(RegisterRequest)
	if err := c.BodyParser(registerRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}
	validate := validator.New()
	if err := validate.Struct(registerRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}

	log.Println(registerRequest.Email)

	if err := as.AuthRepo.CheckIsUserRegistered(registerRequest.Email); err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, 10500))
	}

	insertResult, err := as.AuthRepo.CreateAuth(registerRequest.Email, registerRequest.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.ErrResult(err, 10500))
	}

	log.Println(insertResult)

	data := RegisterResponse{
		Email: insertResult.Email,
	}
	return c.Status(http.StatusOK).JSON(utils.OKResult(data, http.StatusOK))
}

func (as Service) Login(c *fiber.Ctx) error {
	loginRequest := new(LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}
	validate := validator.New()
	if err := validate.Struct(loginRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}

	auth, err := as.AuthRepo.GetAuthByEmailPass(loginRequest.Email, loginRequest.Password)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, 10500))

	}

	signedAccessToken, exp, err := utils.CreateAccessToken(auth.ID.Hex())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.ErrResult(err, http.StatusInternalServerError))
	}

	// Generate encoded token and send it as response.
	signedRefreshToken, _, err := utils.CreateRefreshToken(auth.ID.Hex())

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.ErrResult(err, http.StatusInternalServerError))
	}

	data := LoginResponse{
		AuthToken: Token{
			AccessToken:  signedAccessToken,
			RefreshToken: signedRefreshToken,
			ExpiredAt:    exp,
			TokenType:    "Bearer",
		},
	}
	return c.Status(http.StatusOK).JSON(utils.OKResult(data, http.StatusOK))
}

func (as Service) Refresh(c *fiber.Ctx) error {
	authorization := string(c.Request().Header.Peek("authorization"))
	if !strings.Contains(authorization, "Bearer") {
		return c.Status(http.StatusUnauthorized).JSON(utils.ErrResult(fiber.ErrUnauthorized, http.StatusBadRequest))
	}
	accessToken := strings.Split(authorization, "Bearer ")[1]
	refreshRequest := new(RefreshRequest)
	if err := c.BodyParser(refreshRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}
	validate := validator.New()
	if err := validate.Struct(refreshRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}

	accessTokenData, _, err := utils.DecodeToken(accessToken)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}

	refreshTokenData, refreshExp, err := utils.DecodeToken(refreshRequest.RefreshToken)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}
	if refreshExp < time.Now().Unix() {
		return c.Status(http.StatusUnauthorized).JSON(utils.ErrResult(fiber.ErrUnauthorized, http.StatusBadRequest))
	}

	accessAuthID := accessTokenData[utils.JWTUserIDKey].(string)
	refreshAuthID := refreshTokenData[utils.JWTUserIDKey].(string)

	if accessAuthID != refreshAuthID {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))

	}

	signedAccessToken, exp, err := utils.CreateAccessToken(accessAuthID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.ErrResult(err, http.StatusInternalServerError))
	}

	// Generate encoded token and send it as response.
	signedRefreshToken, _, err := utils.CreateRefreshToken(accessAuthID)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.ErrResult(err, http.StatusInternalServerError))
	}

	data := LoginResponse{
		AuthToken: Token{
			AccessToken:  signedAccessToken,
			RefreshToken: signedRefreshToken,
			ExpiredAt:    exp,
			TokenType:    "Bearer",
		},
	}

	return c.Status(http.StatusOK).JSON(utils.OKResult(data, http.StatusBadRequest))
}
