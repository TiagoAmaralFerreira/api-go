package main

import (
	"net/http"

	"github.com/TiagoAmaralFerreira/api-go/configs"
	"github.com/TiagoAmaralFerreira/api-go/internal/entity"
	"github.com/TiagoAmaralFerreira/api-go/internal/infra/database"
	"github.com/TiagoAmaralFerreira/api-go/internal/infra/webserver/handlers"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, configs.TokenAuth, configs.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Usuário

	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)

	// Produto
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products", productHandler.GetProducts)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)

	http.ListenAndServe(":8000", r)
}
