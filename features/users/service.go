package users

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golkhandani/shopWise/models"
	"github.com/golkhandani/shopWise/utils"
)

type IUserService interface {
	GetUserProfile(c *fiber.Ctx) error
	CreateUserProfile(c *fiber.Ctx) error
}

type UserService struct {
	UserRepo IUserRepo
}

func (us UserService) GetUserProfile(c *fiber.Ctx) error {
	userAuth := c.Locals(utils.LocalUserKey).(models.Auth)
	userProfile, err := us.UserRepo.GetUserByAuthID(userAuth.ID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(utils.ErrResult(err, 201404))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.ErrResult(err, 202500))
	}

	data := UserProfileResponse{
		ID:        userProfile.ID.Hex(),
		Email:     userAuth.Email,
		Username:  userProfile.Username,
		FullName:  userProfile.FullName,
		CreatedAt: userProfile.CreatedAt,
		UpdatedAt: userProfile.UpdatedAt,
	}

	return c.Status(http.StatusOK).JSON(utils.OKResult(data, http.StatusOK))
}

func (us UserService) CreateUserProfile(c *fiber.Ctx) error {
	registerRequest := new(CreateUserProfileRequest)
	if err := c.BodyParser(registerRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}
	validate := validator.New()
	if err := validate.Struct(registerRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, http.StatusBadRequest))
	}

	userAuth := c.Locals(utils.LocalUserKey).(models.Auth)
	createdUserProfile, err := us.UserRepo.CreateUser(userAuth.ID, registerRequest.Username, registerRequest.FullName)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(utils.ErrResult(err, 201404))
		}
		if errors.Is(err, ErrorUserAlreadyHaveProfile) {
			return c.Status(http.StatusBadRequest).JSON(utils.ErrResult(err, 202400))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.ErrResult(err, 203500))
	}

	data := UserProfileResponse{
		ID:        createdUserProfile.ID.Hex(),
		Email:     userAuth.Email,
		Username:  createdUserProfile.Username,
		FullName:  createdUserProfile.FullName,
		CreatedAt: createdUserProfile.CreatedAt,
		UpdatedAt: createdUserProfile.UpdatedAt,
	}

	return c.Status(http.StatusOK).JSON(utils.OKResult(data, http.StatusOK))
}
