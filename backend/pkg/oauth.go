package backend

// https://www.sohamkamani.com/blog/golang/2019-01-01-jwt-authentication/

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// https://dev.to/douglasmakey/oauth2-example-with-go-3n8a
// https://medium.com/@pliutau/getting-started-with-oauth2-in-go-2c9fae55d187

type Claims struct {
	OAuthState string `json:"oas,omitempty"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

var conf = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
	//		Scopes:       []string{"SCOPE1", "SCOPE2"},
	Endpoint:    github.Endpoint,
	RedirectURL: os.Getenv("GITHUB_OAUTH_REDIRECT_URL"),
}

var token = ""

func OauthLogin(w http.ResponseWriter, r *http.Request) {

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		OAuthState: "random string", // TODO randomize
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   os.Getenv("COOKIES_SECURE") != "false",
	})

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	// url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	url := conf.AuthCodeURL(claims.OAuthState)
	//	http.Redirect(w, r, url, http.StatusMovedPermanently)
	http.Redirect(w, r, url, http.StatusFound)
}

func OauthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO validate state
	// http://www.thread-safe.com/2014/05/the-correct-use-of-state-parameter-in.html

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	//	var code string
	//	if _, err := fmt.Scan(&code); err != nil {
	//		log.Fatal(err)
	//	}
	code := r.FormValue("code")

	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	_ = client

	// store github token in backgorund https://auth0.com/learn/refresh-tokens/

	fmt.Printf("token: %s AccessToken: %s Expiry: %s RefreshToken: %s TokenType: %s Type(): %s Valid(): %t\n", tok, tok.AccessToken, tok.Expiry, tok.RefreshToken, tok.TokenType, tok.Type(), tok.Valid())
	//		log.Printf("Extra(): %v\n", tok.Extra())

	token = tok.AccessToken

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(1 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		// HttpOnly: true,
		// Secure:   true,
		Secure: os.Getenv("COOKIES_SECURE") != "false",
	})

	http.Redirect(w, r, "http://127.0.0.1:8000", http.StatusFound)
}

// TODO implement /refresh https://www.sohamkamani.com/blog/golang/2019-01-01-jwt-authentication/
