package handlers

import (
	"AlexSarva/go-shortener/crypto"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// GenerateCookie generate cookie for the user
func GenerateCookie(userID uuid.UUID) http.Cookie {
	session := crypto.Encrypt(userID, crypto.SecretKey)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "session", Value: session, Expires: expiration, Path: "/"}
	return cookie
}

// GetCookie get cookie from request
func GetCookie(r *http.Request) (uuid.UUID, error) {
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		log.Println(cookieErr)
		return uuid.UUID{}, ErrNotValidCookie
	}
	userID, cookieDecryptErr := crypto.Decrypt(cookie.Value, crypto.SecretKey)
	if cookieDecryptErr != nil {
		return uuid.UUID{}, cookieDecryptErr
	}
	return userID, nil

}

// CookieHandler middleware adds cookie for new user
func CookieHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, userIDErr := GetCookie(r)
		if userIDErr != nil {
			log.Println(userIDErr)
			userCookie := GenerateCookie(uuid.New())
			log.Println(userCookie)
			r.AddCookie(&userCookie)
			http.SetCookie(w, &userCookie)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
