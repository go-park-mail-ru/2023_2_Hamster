package http

import (
	"net/http"
	"time"
)

const (
	AuthCookie = "Authentication"
)

func InitCookie(name, value string, expire time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expire,
		Path:     "/",
		HttpOnly: true,
	}
}
