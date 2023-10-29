package http

import (
	"net/http"
	"time"
)

func InitCookie(name, value string, expire time.Time, path string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expire,
		Path:     path,
		HttpOnly: true,                    // No JS permision
		Secure:   true,                    // httpsOnly
		SameSite: http.SameSiteStrictMode, // CSRF protection: No cross-site cookie
	}
}
