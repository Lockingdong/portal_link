package main

import (
	"fmt"
	"net/http"
	"portal_link/pkg"
	"portal_link/pkg/config"
)

func main() {
	config.Init()
	db := pkg.NewPG(config.GetDBConfig().DSN())
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
