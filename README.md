# go-login-with-jwt-and-crud

Fiturnya:
 - create user
 - login user
 - edit user
 - delete user
 - lihat user
 
 # Instalasi
- git clone https://github.com/nuridincs/go-login-with-jwt-and-crud.git
- go get github.com/dgrijalva/jwt-go
- go get github.com/gorilla/mux
- go get -u github.com/go-sql-driver/mysql
- buat db catatan
- import db catatan.sql
- sesuaikan konfigurasi db di -> database/database.go
- go run main.go

# Testing
- Buka postman
- jalankan localhost:8000/register
  ``` 
  {
	"username" : "dewi",
	"password" : "dewi23",
	"nama" : "dewi"
  }
```
- ketika return success ada respon token,tokennya digunakan ketika ingin akses getuser,update user,delete user.
- jadi di header di tambahkan key : Authorization,value : respon token
