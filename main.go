package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `bson:"created_At" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}

type UserSaveData struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type Users []*User

func main() {
	router := gin.Default()

	router.GET("/", status)

	api := router.Group("/api")
	api.GET("/users", getUsers)
	api.GET("/users/:userId", getUser)
	api.POST("/users/create", createUser)
	api.DELETE("/users/:userId", deleteUser)
	api.PATCH("/users/:userId", editUser)

	router.Run(":3007")
}

func status(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "up"})
}

func getUsers(ctx *gin.Context) {
	var users Users

	var collection = GetCollection("users")

	cur, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	for cur.Next(ctx) {

		var user User
		err = cur.Decode(&user)

		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		users = append(users, &user)
	}

	ctx.JSON(http.StatusOK, users)
}

func getUser(ctx *gin.Context) {
	var user User
	var collection = GetCollection("users")
	userId := ctx.Param("userId")
	oid, _ := primitive.ObjectIDFromHex(userId)

	err := collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user)

	if err != nil {
		ctx.String(http.StatusNotFound, fmt.Sprintf("user %s not found", userId))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func createUser(ctx *gin.Context) {
	var userData UserSaveData
	var collection = GetCollection("users")

	ctx.BindJSON(&userData)

	oid := primitive.NewObjectID()

	if userData.Name == "" && userData.Email == "" {
		ctx.String(http.StatusBadRequest, "user and email is required")
		return
	}

	user := User{
		ID:        oid,
		Name:      userData.Name,
		Email:     userData.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := collection.InsertOne(ctx, user)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func deleteUser(ctx *gin.Context) {
	var collection = GetCollection("users")
	userId := ctx.Param("userId")
	oid, _ := primitive.ObjectIDFromHex(userId)

	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": oid})

	if err != nil {
		ctx.String(http.StatusBadRequest, "error")
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func editUser(ctx *gin.Context) {
	var userData UserSaveData
	var collection = GetCollection("users")
	userId := ctx.Param("userId")
	oid, _ := primitive.ObjectIDFromHex(userId)

	ctx.BindJSON(&userData)

	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"name":       userData.Name,
			"email":      userData.Email,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		ctx.String(http.StatusBadRequest, "error")
		return
	}

	ctx.JSON(http.StatusOK, result)
}
