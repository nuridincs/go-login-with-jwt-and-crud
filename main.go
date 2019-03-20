package main

import (
	"fmt"
	"go-task/database"
	"go-task/repository"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Main function
func main() {
	// Init router
	r := mux.NewRouter()
	database.Setup()

	defer database.Disconnect()
	// Route handles & endpoints
	r.HandleFunc("/login", controller.Login).Methods("POST")
	r.HandleFunc("/register", controller.CreateUser).Methods("POST")

	r.HandleFunc("/users", use(controller.GetUsers, basicAuth)).Methods("GET")
	r.HandleFunc("/users/{id}", use(controller.GetUser, basicAuth)).Methods("GET")
	r.HandleFunc("/users/{id}", use(controller.UpdateUser, basicAuth)).Methods("PUT")
	r.HandleFunc("/users/{id}", use(controller.DeleteUser, basicAuth)).Methods("DELETE")

	r.HandleFunc("/list_transaksi", use(controller.ListTransaction, basicAuth)).Methods("POST")
	r.HandleFunc("/save_transaksi", use(controller.SaveTransaction, basicAuth)).Methods("POST")

	r.HandleFunc("/token", controller.Token).Methods("POST")
	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(database.SecretKey()), nil
		})
		//fmt.Printf("unexpected type %T", claims)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}
