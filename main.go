package main

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/urlHandlers"
	"forum/validateData"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.Engine()
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	staticFiles := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFiles))

	mux.HandleFunc("/register", urlHandlers.HandleRegister)
	mux.HandleFunc("/login", urlHandlers.HandleLogin)
	mux.HandleFunc("/logout", urlHandlers.HandleLogout)
	mux.HandleFunc("/", urlHandlers.HandleForum)
	mux.HandleFunc("/post", urlHandlers.HandlePost)
	mux.HandleFunc("/postcontent", urlHandlers.HandlePostContent)

	fmt.Println("Server hosted at: http://localhost:" + "8000")
	fmt.Println("To Kill Server press Ctrl+C")

	err := server.ListenAndServe()
	validateData.CheckErr(err)
}
