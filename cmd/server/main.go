package main

import (
	"github.com/DieegoAlves/API/configs"
	_ "github.com/DieegoAlves/API/docs"
	"github.com/DieegoAlves/API/internal/entity"
	"github.com/DieegoAlves/API/internal/infra/database"
	"github.com/DieegoAlves/API/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

// @title				Go Expert API Example
// @version 			1.0
// @description 		Product API with Authentication
// @termsOfService		http://swagger.io/terms/
//
// @contact.name 		Diego Alves Ferreira
// @contact.url			https://www.linkedin.com/in/dieegoalves/
// @contact.email		diegoaf@ucl.br
//
// @license.name		Revolution Softwares
// @license.url			http://www.revolutionsoftwares.com
//
// @host				localhost:8000
// @BasePath			/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs, err := configs.LoadConfig("./cmd/server")
	if err != nil {
		panic(err)

	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Product{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JWTExpiresIn", configs.JWTExpiresIn))

	//Middleware criado na mão
	//r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetAllProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Post("/generate_token", userHandler.GetJWT)
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	err = http.ListenAndServe(":8000", r)
	if err != nil {
		return
	}

}

// Middleware feito na mão
//func LogRequest(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log.Printf("Request: %s %s", r.Method, r.URL.Path)
//		next.ServeHTTP(w, r)
//	})
//}
