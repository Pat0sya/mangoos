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
	var user models.User

	ctx, cancel := CreateContextWithTimeout()
	defer cancel()

	if err := ParseAndValidate(c, &user); err != nil {
		return err
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := InsertOne(ctx, userCollection, newUser)
	if err != nil {
		return ErrorHandler(c, http.StatusInternalServerError, "error inserting user", err)
	}

	return SuccessHandler(c, http.StatusCreated, "user created successfully", result)
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := CreateContextWithTimeout()
	defer cancel()

	userId := c.Params("userId")
	var user models.User

	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return ErrorHandler(c, http.StatusBadRequest, "invalid user ID", err)
	}

	if err := FindOne(ctx, userCollection, bson.M{"id": objId}, &user); err != nil {
		return ErrorHandler(c, http.StatusInternalServerError, "user not found", err)
	}

	return SuccessHandler(c, http.StatusOK, "user found", user)
}

func EditAUser(c *fiber.Ctx) error {
	userId := c.Params("userId")
	var user models.User

	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return ErrorHandler(c, http.StatusBadRequest, "invalid user ID", err)
	}

	ctx, cancel := CreateContextWithTimeout()
	defer cancel()

	if err := ParseAndValidate(c, &user); err != nil {
		return err
	}

	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}
	result, err := UpdateOne(ctx, userCollection, bson.M{"id": objId}, update)
	if err != nil {
		return ErrorHandler(c, http.StatusInternalServerError, "error updating user", err)
	}
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			return ErrorHandler(c, http.StatusInternalServerError, "user not found", err)
		}
	}

	return SuccessHandler(c, http.StatusOK, "user updated successfully", result)
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := CreateContextWithTimeout()
	defer cancel()
	userId := c.Params("userId")
	objId, _ := primitive.ObjectIDFromHex(userId)
	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return ErrorHandler(c, http.StatusBadRequest, "error delete the user ID", err)
	}
	if result.DeletedCount < 1 {
		return ErrorHandler(c, http.StatusNotFound, "user not found", err)
	}
	return SuccessHandler(c, http.StatusOK, "user delete successfully", result)

}

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := CreateContextWithTimeout()
	var users []models.User
	defer cancel()
	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return ErrorHandler(c, http.StatusInternalServerError, "user not found", err)
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return ErrorHandler(c, http.StatusInternalServerError, "users not found", err)
		}
		users = append(users, singleUser)
	}
	return SuccessHandler(c, http.StatusOK, "users found successfully", users)
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

func ParseAndValidate[T any](c *fiber.Ctx, result *T) error {

	if err := c.BodyParser(result); err != nil {
		return ErrorHandler(c, http.StatusBadRequest, "error parsing request body", err)
	}

	// Validate the parsed data
	if validationErr := validate.Struct(result); validationErr != nil {
		return ErrorHandler(c, http.StatusBadRequest, "validation error", validationErr)
	}

	return nil
}

func CreateContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
func InsertOne(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(ctx, document)
}
func UpdateOne(ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return collection.UpdateOne(ctx, filter, bson.M{"$set": update})
}
func FindOne(ctx context.Context, collection *mongo.Collection, filter interface{}, result interface{}) error {
	return collection.FindOne(ctx, filter).Decode(result)
}
