package controllers

import (
	"context"
	"mongoos/configs"
	"mongoos/models"
	"mongoos/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return ErrorHandler(c, http.StatusBadRequest, "error", err)
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return ErrorHandler(c, http.StatusBadRequest, "error", validationErr)
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return ErrorHandler(c, http.StatusInternalServerError, "error", err)
	}
	return SuccessHandler(c, http.StatusCreated, "success", result)

}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		return ErrorHandler(c, http.StatusInternalServerError, "error", err)
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return ErrorHandler(c, http.StatusBadRequest, "error", validationErr)
	}

	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return ErrorHandler(c, http.StatusInternalServerError, "error", err)
	}

	var updatedUser models.User

	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			return ErrorHandler(c, http.StatusInternalServerError, "error", err)
		}

	}
	return SuccessHandler(c, http.StatusCreated, "success", result)
}

func ErrorHandler(c *fiber.Ctx, statusCode int, errMsg string, err error) error {
	return c.Status(statusCode).JSON(responses.UserResponse{
		Status:  statusCode,
		Message: "error",
		Data:    &fiber.Map{"data": err.Error()},
	})
}

func SuccessHandler(c *fiber.Ctx, statusCode int, okMsg, data interface{}) error {
	return c.Status(statusCode).JSON(responses.UserResponse{
		Status:  statusCode,
		Message: "success",
		Data:    &fiber.Map{"data": data},
	})
}
