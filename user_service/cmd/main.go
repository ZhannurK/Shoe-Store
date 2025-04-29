package main

import (
	"Shoe-Store/internal/config"
	"Shoe-Store/internal/user/domain"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"Shoe-Store/internal/user/delivery"
	"Shoe-Store/internal/user/repository"
	"Shoe-Store/internal/user/usecase"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("shoe-store.db"), &gorm.Config{})
	db.AutoMigrate(&domain.User{}) // migrate user table

	repo := repository.NewUserRepository(db)
	service := usecase.NewAuthService(repo)
	handler := delivery.NewHandler(service)

	cfg := config.LoadConfig()
	fmt.Println("MongoDB URI:", cfg.MongoURI)
	fmt.Println("Server running at port:", cfg.ServerPort)

	r := gin.Default()
	user := r.Group("/user")
	{
		user.POST("/signup", handler.Signup)
		user.POST("/login", handler.Login)
		user.POST("/confirmpassword", handler.ConfirmPassword)
		user.POST("/changepassword", handler.ChangePassword)
	}

	r.Run(":8080")
}
