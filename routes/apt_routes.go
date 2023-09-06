package routes

import (
    "github.com/yasuobgg/restApiForApt/controllers" 

    "github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	app.Get("/apt/all", controllers.GetAllAPTS)
	app.Get("/apt/:id", controllers.GetAnAPT)
    app.Post("/apt", controllers.CreateAPT) 
	app.Put("/apt/:id", controllers.EditAnAPT)
	app.Delete("/apt/:id", controllers.DeleteAnAPT)
}