package controller

import (
	"encoding/json"
	"fmt"
	"go-task/database"

	// "go-task/helper"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// User struct (Model)
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nama     string `json:"nama"`
}

// Login Get single user
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	db := database.Connect()
	var each = User{}
	err := db.QueryRow(
		"select id, username, password, nama from user where username = ? AND password = ?",
		user.Username,
		user.Password,
	).Scan(
		&each.ID,
		&each.Username,
		&each.Password,
		&each.Nama,
	)
	if err != nil {
		fmt.Println(err.Error())
		var msg = map[string]string{
			"status":  "404",
			"message": "Data tidak ditemukan",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	sign.Claims = jwt.MapClaims{
		"id":       each.ID,
		"username": each.Username,
		"nama":     each.Nama,
	}
	token, err := sign.SignedString([]byte(database.SecretKey()))
	if err != nil {
		var msg = map[string]string{
			"status":  "401",
			"message": err.Error(),
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	var msg = map[string]interface{}{
		"status":  "200",
		"message": "Success",
		"data":    each,
		"token":   token,
	}
	json.NewEncoder(w).Encode(msg)
}

// GetUsers Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()
	rows, err := db.Query("select id, username, password, nama from user")
	if err != nil {
		fmt.Println(err.Error())
		var msg = map[string]string{
			"status":  "404",
			"message": "Data tidak ditemukan",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	var users []User
	for rows.Next() {
		var each = User{}
		rows.Scan(&each.ID, &each.Username, &each.Password, &each.Nama)
		users = append(users, each)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUser Get single user
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	db := database.Connect()
	var each = User{}
	err := db.QueryRow(
		"select id, username, password, nama from user where id= ? ",
		params["id"],
	).Scan(
		&each.ID,
		&each.Username,
		&each.Password,
		&each.Nama,
	)
	if err != nil {
		fmt.Println(err.Error())
		var msg = map[string]string{
			"status":  "404",
			"message": "Data tidak ditemukan",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	json.NewEncoder(w).Encode(each)
}

// CreateUser Add new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	db := database.Connect()
	_, err := db.Exec(
		"insert into user (username,password,nama) values (?, ?, ?)",
		user.Username,
		user.Password,
		user.Nama,
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var msg = map[string]string{
		"status":  "200",
		"message": "Data berhasil disimpan",
	}

	// from := "email"
	// fromAlias := "Catatan Ku"
	// subject := "Konfirmasi Akun"
	// receiver := "email"
	// dataEmail := map[string]string{
	// 	"username": user.Nama,
	// }
	// re := mail.NewRequest(from, fromAlias, []string{receiver}, subject)
	// go re.Send("view/confirm.html", dataEmail)

	json.NewEncoder(w).Encode(msg)
}

// UpdateUser Update user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	db := database.Connect()
	res, err := db.Exec(
		"update user set username = ?, password = ?, nama = ? where id = ?",
		user.Username,
		user.Password,
		user.Nama,
		params["id"],
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		var msg = map[string]string{
			"status":  "304",
			"message": "Tidak ada data yang diubah",
		}

		json.NewEncoder(w).Encode(msg)
		return
	}

	var msg = map[string]string{
		"status":  "200",
		"message": "Data berhasil diubah",
	}

	json.NewEncoder(w).Encode(msg)
}

// DeleteUser Delete user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	db := database.Connect()
	res, err := db.Exec(
		"delete from user where id = ?",
		params["id"],
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
	}
	if count == 0 {
		var msg = map[string]string{
			"status":  "304",
			"message": "Tidak ada data yang terhapus",
		}

		json.NewEncoder(w).Encode(msg)
		return
	}

	var msg = map[string]string{
		"status":  "200",
		"message": "Data berhasil dihapus",
	}

	json.NewEncoder(w).Encode(msg)
}
