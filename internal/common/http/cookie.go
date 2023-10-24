package http

import (
	"net/http"
	"time"
)

const (
	AuthTag = "Authentication"
)

func InitCookie(name, value string, expire time.Time, path string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expire,
		Path:     path,
		HttpOnly: true,
	}
}
