package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/damoon/michael-collins/backend/pkg"
)

func main() {
	log.Println("Hello, world.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html>Hello, you have requested: %s</br>\n<a href='/oauth-login'>login via github</a>", r.URL.Path)
	})

	http.HandleFunc("/oauth-login", backend.OauthLogin)

	http.HandleFunc("/oauth-callback", backend.OauthCallback)

	http.HandleFunc("/jwt-refresh", backend.JWTRefresh)

	http.HandleFunc("/details", backend.Details)

	http.ListenAndServe(":8080", nil)
}
