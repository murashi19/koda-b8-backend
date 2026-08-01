package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/murashi19/koda-b8-backend/internal/di"
	"github.com/murashi19/koda-b8-backend/internal/middleware"
)

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	container, err := di.NewContainer()
	if err != nil {
		log.Fatal(err)
	}
	defer container.Close()

	auth := container.AuthHandler()
	user := container.UserHandler()

	router.POST("/auth/register", auth.Register)
	router.POST("/auth/login", auth.Login)

	authorized := router.Group("/")
	authorized.Use(middleware.Auth(container.Config()))

	authorized.GET("/me", user.Me)
	authorized.PATCH("/me", user.UpdateMyProfile)
	authorized.PATCH("/me/avatar", user.UpdateAvatar)

	users := router.Group("/users")
	users.Use(middleware.Auth(container.Config()))

	users.GET("", user.GetAll)
	users.GET("/:id", user.GetByID)
	users.PUT("/:id", user.Update)
	users.DELETE("/:id", user.Delete)

	log.Fatal(router.Run(":8081"))

}
