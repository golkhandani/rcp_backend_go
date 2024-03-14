package app

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golkhandani/shopWise/configs"
	"github.com/golkhandani/shopWise/features/auth"
	"github.com/golkhandani/shopWise/features/users"
	"go.mongodb.org/mongo-driver/mongo"
)

const _LogFormat = "${pid} ${locals:requestid} ${status} - ${method} ${path}\n"

func SetupServerApp(db *mongo.Database) *fiber.App {

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: false,
	},
	)
	app.Use(requestid.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: _LogFormat,
	}))

	api := app.Group("/api")

	// SETUP REPOS

	authRepo := auth.Repo{
		AuthCollection: db.Collection("auth"),
	}
	userRepo := users.UserRepo{
		UserCollection: db.Collection("users"),
	}

	// SETUP SERVICES
	userService := users.UserService{
		UserRepo: userRepo,
	}

	authService := auth.Service{
		AuthRepo: authRepo,
	}

	gaurdService := auth.Guard{
		AuthRepo: authRepo,
	}

	// SETUP CONTROLLERS
	authContoller := auth.Controller{
		AuthService:  authService,
		GuardService: gaurdService,
		Router:       api,
	}

	userContoller := users.UserController{
		UserService:  userService,
		GuardService: gaurdService,
		Router:       api,
	}

	// REGISTER CONTROLLERS
	authContoller.Register()
	userContoller.Register()

	serverErr := app.Listen(fmt.Sprintf(":%d", configs.Env.Port))
	if serverErr != nil {
		log.Fatal(serverErr)
	}
	return app
}
