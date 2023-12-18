package controllers

import (
	"backend/configs"
	"backend/models"
	"backend/responses"
	"net/http"
	"time"

	"github.com/andskur/argon2-hashing"
	"github.com/go-playground/validator/v10"
	jtoken "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gofr.dev/pkg/gofr"
)

var UserCollection = configs.GetCollection(configs.DB, "user")
var authValidator = validator.New()

func CreateUser(ctx *gofr.Context) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	if validationErr := authValidator.Struct(user); validationErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: gofr.Map{"error": validationErr.Error()}})
	}

	// Hash the Password
	pwdHash, err := argon2.GenerateFromPassword([]byte(user.Password), argon2.DefaultParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Username: user.Username,
		Password: string(pwdHash),
		FullName: user.FullName,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	return ctx.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: gofr.Map{"user": result}})
}

func GetUser(ctx *gofr.Context) error {
	userId := ctx.Params("userId")
	objId, _ := primitive.ObjectIDFromHex(userId)

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user, gofr.FindOneOptions().SetProjection(bson.M{"password": 0}))

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	return ctx.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: gofr.Map{"user": user}})
}

func EditPassword(ctx *gofr.Context) error {
	var passwordChgReq models.PasswordChangeRequest
	var existingUser models.User

	if err := ctx.BodyParser(&passwordChgReq); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: gofr.Map{"error": "Invalid request payload"}})
	}

	// Check Old Password with user
	user := GetUserDetailsFromToken(ctx)
	userId := user[0]
	objId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	err = userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&existingUser)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	err = argon2.CompareHashAndPassword([]byte(existingUser.Password), []byte(passwordChgReq.OldPassword))

	if err != nil {
		return ctx.Status(http.StatusForbidden).JSON(responses.UserResponse{Status: http.StatusForbidden, Message: "error", Data: gofr.Map{"error": "Incorrect Old Password"}})
	}

	newPassword, hashErr := argon2.GenerateFromPassword([]byte(passwordChgReq.NewPassword), argon2.DefaultParams)

	if hashErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	update := bson.M{"password": string(newPassword)}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	return ctx.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: gofr.Map{"result": result}})
}

func LoginUser(ctx *gofr.Context) error {
	var userData models.LoginRequest
	var queryUser models.User

	if err := ctx.BodyParser(&userData); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	if validationErr := authValidator.Struct(userData); validationErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: gofr.Map{"error": validationErr.Error()}})
	}

	// Get User
	err := userCollection.FindOne(ctx, bson.M{"username": userData.Username}).Decode(&queryUser)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	// Compare Hash
	hashCompareError := argon2.CompareHashAndPassword([]byte(queryUser.Password), []byte(userData.Password))
	if hashCompareError != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "unauthorized", Data: gofr.Map{"error": "Wrong Password"}})
	}

	day := time.Hour * 24
	claims := jtoken.MapClaims{
		"ID":       queryUser.Id,
		"Username": queryUser.Username,
		"FullName": queryUser.FullName,
		"exp":      time.Now().Add(day * 1).Unix(),
	}

	// Generate JWT
	tokenGenerator := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	token, err := tokenGenerator.SignedString([]byte(configs.JWTSecret()))

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: gofr.Map{"error": err.Error()}})
	}

	return ctx.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "authenticated", Data: gofr.Map{"token": token}})
}

func GetUserDetailsFromToken(ctx *gofr.Context) []string {
	user := ctx.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	userId := claims["ID"].(string)
	username := claims["Username"].(string)
	fullname := claims["FullName"].(string)
	return []string{userId, username, fullname}
}

func GetProfile(ctx *gofr.Context) error {
	userData := GetUserDetailsFromToken(ctx)
	return ctx.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: gofr.Map{"user": gofr.Map{"id": userData[0], "username": userData[1], "fullname": userData[2]}}})
}
