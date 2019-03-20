package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Secret struct (Model)
type Secret struct {
	Key string `json:"secret_key"`
}

// Response struct (Model)
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"access_token"`
}

// Token Get access token
func Token(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var param Secret
	_ = json.NewDecoder(r.Body).Decode(&param)

	url := "http://localhost/api-rajamobil/api/web/index.php/bcakkb/token"

	values := map[string]string{
		"secret_key": param.Key,
	}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var t Response
	err = decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	fmt.Println("response Token:", t.Token)
	json.NewEncoder(w).Encode(t)
}
