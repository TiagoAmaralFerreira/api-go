package main

import (
	"fmt"
	"net/http"

	"github.com/TiagoAmaralFerreira/api-go/configs"
	"github.com/TiagoAmaralFerreira/api-go/internal/entity"
	"github.com/TiagoAmaralFerreira/api-go/internal/infra/database"
	"github.com/TiagoAmaralFerreira/api-go/internal/infra/webserver/handlers"

	_ "github.com/TiagoAmaralFerreira/api-go/docs"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Finance API
// @version         1.0
// @description     API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Tiago Amaral

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	configs := configs.NewConfig()

	dbHost := configs.DB_HOST
	dbUser := configs.DB_USER
	dbPassword := configs.DB_PASSWORD
	dbName := configs.DB_NAME
	dbPort := configs.DB_PORT

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("jwtExperiesIn", configs.JWT_EXPIRES_IN))

	// Usuário

	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	// Produto
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	http.ListenAndServe(":8000", r)
}
