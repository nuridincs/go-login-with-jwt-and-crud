package controller

import (
	"encoding/json"
	"fmt"
	"go-task/database"
	"net/http"
	"strings"
)

// Trx struct (Model)
type Trx struct {
	ID         string `json:"id"`
	Tanggal    string `json:"tanggal"`
	Keterangan string `json:"keterangan"`
	Harga      int    `json:"harga"`
	IDUser     int    `json:"id_user"`
	Tipe       int    `json:"tipe"`
}

// Param struct (Model)
type Param struct {
	Page  int    `json:"page"`
	Tahun string `json:"tahun"`
	Bulan string `json:"bulan"`
}

// ListTransaction Get data transaksi
func ListTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var param Param
	_ = json.NewDecoder(r.Body).Decode(&param)

	tokenString := r.Header.Get("Authorization")
	claims := database.GetClaims(tokenString)

	db := database.Connect()

	limit := 0
	if param.Page != 0 {
		limit = (param.Page - 1) * 10
	}

	condition := []string{}
	condition = append(condition, "id_user = "+claims["id"].(string))
	if param.Tahun != "" {
		if param.Bulan == "" {
			condition = append(condition, "tanggal like '%"+param.Tahun+"%'")
		} else {
			bulan := fmt.Sprintf("%02s", param.Bulan)
			condition = append(condition, "tanggal like '%"+param.Tahun+"-"+bulan+"%'")
		}
	}
	where := strings.Join(condition, " AND ")

	rows, err := db.Query(
		"select id, tanggal, keterangan, harga, id_user, tipe "+
			"from transaksi "+
			"where "+where+
			"order by tanggal desc "+
			"limit ?,10",
		limit,
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
	var data []Trx
	for rows.Next() {
		var each = Trx{}
		rows.Scan(
			&each.ID,
			&each.Tanggal,
			&each.Keterangan,
			&each.Harga,
			&each.IDUser,
			&each.Tipe,
		)
		data = append(data, each)
	}
	var totalRows int
	var totalUangMasuk int
	var totalUangKeluar int
	err = db.QueryRow(
		"select count(id) as total_rows " +
			"from transaksi " +
			"where " + where,
	).Scan(&totalRows)

	err = db.QueryRow(
		"select sum(harga) as total_uang_masuk " +
			"from transaksi " +
			"where " + where +
			"and tipe = 2",
	).Scan(&totalUangMasuk)

	err = db.QueryRow(
		"select sum(harga) as total_uang_keluar " +
			"from transaksi " +
			"where " + where +
			"and tipe = 1",
	).Scan(&totalUangKeluar)

	if err != nil {
	}

	var msg = map[string]interface{}{
		"status":            "200",
		"message":           "Success",
		"data":              data,
		"total_rows":        totalRows,
		"total_uang_masuk":  totalUangMasuk,
		"total_uang_keluar": totalUangKeluar,
	}
	json.NewEncoder(w).Encode(msg)
}

// SaveTransaction Get data transaksi
func SaveTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var param Trx
	_ = json.NewDecoder(r.Body).Decode(&param)

	tokenString := r.Header.Get("Authorization")
	claims := database.GetClaims(tokenString)

	db := database.Connect()
	_, err := db.Exec(
		"insert into transaksi (tanggal, keterangan, harga, id_user, tipe) values (?, ?, ?, ?, ?)",
		param.Tanggal,
		param.Keterangan,
		param.Harga,
		claims["id"],
		param.Tipe,
	)
	if err != nil {
		fmt.Println(err.Error())
		var msg = map[string]string{
			"status":  "500",
			"message": err.Error(),
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	var msg = map[string]string{
		"status":  "200",
		"message": "Data berhasil disimpan",
	}
	json.NewEncoder(w).Encode(msg)
}
