package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"io"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func server(w http.ResponseWriter, r *http.Request) {

	// XSS
	io.WriteString(w, r.URL.Query().Get("username"))

	// SQL Injection
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	sql := "SELECT * FROM user WHERE username='" + username + "' AND password='" + password + "'"
	db.Exec(sql)

	// RCE
	userInput := r.URL.Query().Get("command")
	cmd := exec.Command("/bin/sh", "-c", userInput)
	output, _ := cmd.CombinedOutput()
	fmt.Fprintf(w, "Output: %s", output)

	// SSRF
    resp, err := r.URL.Query().Get("url")
    if err != nil {
        fmt.Println("Error fetching URL:", err)
        return
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        return
    }

    fmt.Println(string(body))
}

func main() {
	http.HandleFunc("/", server)
	http.ListenAndServe(":8000", nil)
}
