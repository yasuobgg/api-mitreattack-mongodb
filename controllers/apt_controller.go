package controllers
import (
	//basic libs
    "context"
    "net/http"
    "time"
	"os"
	"log"

	// own modules
	"github.com/yasuobgg/restApiForApt/configs"
    "github.com/yasuobgg/restApiForApt/models"
    "github.com/yasuobgg/restApiForApt/responses"

	// global modules
	"github.com/joho/godotenv"
    "github.com/go-playground/validator/v10"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"

	// web framework
	"github.com/gofiber/fiber/v2"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, os.Getenv("COLNAME"))
var validate = validator.New()

func CreateAPT(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var apt models.APT
    defer cancel()

    //validate the request body
    if err := c.BodyParser(&apt); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.WebResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(&apt); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.WebResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    }

    newAPT := models.APT{
		ID:              apt.ID,
		Name:            apt.Name,
		AssociatedGroup: apt.AssociatedGroup,
		Description:     apt.Description,
		Timestamp:       time.Now().Unix(),
    }

    result, err := userCollection.InsertOne(ctx, newAPT)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.WebResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
    }

    return c.Status(http.StatusCreated).JSON(responses.WebResponse{Status: http.StatusCreated, Message: "success", Data: result})
}

func GetAnAPT(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    aptId := c.Params("id")
    var apt models.APT
    defer cancel()

    err := userCollection.FindOne(ctx, bson.M{"id": aptId}).Decode(&apt)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.WebResponse{Status: http.StatusInternalServerError, Message: "error", Data:  err.Error()})
    }

    return c.Status(http.StatusOK).JSON(responses.WebResponse{Status: http.StatusOK, Message: "success", Data: apt})
}

func EditAnAPT(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    aptId := c.Params("id")
    var apt models.APT
    defer cancel()

    //validate the request body
    if err := c.BodyParser(&apt); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.WebResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(&apt); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.WebResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
    }

    update := bson.M{"name": apt.Name, "associatedgroup": apt.AssociatedGroup, "description": apt.Description, "timestamp": time.Now().Unix()}

    result, err := userCollection.UpdateOne(ctx, bson.M{"id": aptId}, bson.M{"$set": update})
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.WebResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
    }

    //get updated user details
    var updatedAPT models.APT
    if result.MatchedCount == 1 {
        err := userCollection.FindOne(ctx, bson.M{"id": aptId}).Decode(&updatedAPT)
        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(responses.WebResponse{Status: http.StatusInternalServerError, Message: "error", Data:  err.Error()})
        }
    }

    return c.Status(http.StatusOK).JSON(responses.WebResponse{Status: http.StatusOK, Message: "success", Data: updatedAPT})
}

func DeleteAnAPT(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    aptId := c.Params("id")
    defer cancel()

    result, err := userCollection.DeleteOne(ctx, bson.M{"id": aptId})
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.WebResponse{Status: http.StatusInternalServerError, Message: "error", Data:  err.Error()})
    }

    if result.DeletedCount < 1 {
        return c.Status(http.StatusNotFound).JSON(
            responses.WebResponse{Status: http.StatusNotFound, Message: "error", Data: "User with specified ID not found!"},
        )
    }

	return c.SendStatus(http.StatusNoContent)  
}

func GetAllAPTS(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var apts []models.APT
    defer cancel()

    results, err := userCollection.Find(ctx, bson.M{})

    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.WebResponse{Status: http.StatusInternalServerError, Message: "error", Data:  err.Error()})
    }

    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleAPT models.APT
        if err = results.Decode(&singleAPT); err != nil {
            return c.Status(http.StatusInternalServerError).JSON(responses.WebResponse{Status: http.StatusInternalServerError, Message: "error", Data:  err.Error()})
        }

        apts = append(apts, singleAPT)
    }

    return c.Status(http.StatusOK).JSON(
        responses.WebResponse{Status: http.StatusOK, Message: "success", Data: apts},
    )
}

func init(){
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}