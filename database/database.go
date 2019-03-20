package database

import (
	"database/sql"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	db *sql.DB
)

//Setup global connect
func Setup() {
	d, err := sql.Open("mysql", "root:@tcp(localhost:3306)/catatan?charset=utf8")
	if err != nil {
		panic(err)
	}

	db = d
}

//Connect global connect
func Connect() *sql.DB {
	return db
}

//Disconnect global connect
func Disconnect() error {
	return db.Close()
}

//SecretKey global secretkey
func SecretKey() string {
	return "apitest"
}

//GetClaims get claims token
func GetClaims(tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey()), nil
	})
	fmt.Println(err)
	return claims
}
