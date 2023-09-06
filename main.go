package main

import (
	"fmt"

	// own modules
    "github.com/yasuobgg/restApiForApt/configs"
    "github.com/yasuobgg/restApiForApt/routes"

	// web framework
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    //run database
    configs.ConnectDB()
	fmt.Println("Connected to MongoDB")

    //routes
    routes.UserRoute(app)

	// app run at port 6000
	fmt.Println("Service is running at port: 6000")
    app.Listen(":6000")
	
}