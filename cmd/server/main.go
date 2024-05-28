package main

import (
	"github.com/DieegoAlves/API/configs"
	"github.com/DieegoAlves/API/internal/entity"
	"github.com/DieegoAlves/API/internal/infra/database"
	"github.com/DieegoAlves/API/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

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

	//Middleware criado na mão
	//r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(middleware.WithValue("jwt", configs.TokenAuth))
		r.Use(middleware.WithValue("JWTExpiresIn", configs.JWTExpiresIn))
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetAllProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/generate_token", userHandler.GetJWT)

	http.ListenAndServe(":8000", r)
}

//Middleware feito na mão
//func LogRequest(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log.Printf("Request: %s %s", r.Method, r.URL.Path)
//		next.ServeHTTP(w, r)
//	})
//}
