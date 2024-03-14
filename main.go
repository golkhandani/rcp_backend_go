package main

import (
	"github.com/golkhandani/shopWise/app"
	"github.com/golkhandani/shopWise/configs"
)

func main() {

	db := configs.GetDB()
	defer configs.CloseDB()

	app.SetupServerApp(db)
}
