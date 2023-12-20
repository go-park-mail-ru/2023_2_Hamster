package http

import (
	"net/http"
	"testing"
	"time"
)

func TestInitCookie(t *testing.T) {
	name := "testCookie"
	value := "testValue"
	expire := time.Now().Add(24 * time.Hour)
	path := "/"

	cookie := InitCookie(name, value, expire, path)

	if cookie.Name != name {
		t.Errorf("Expected cookie name to be %s, but got %s", name, cookie.Name)
	}

	if cookie.Value != value {
		t.Errorf("Expected cookie value to be %s, but got %s", value, cookie.Value)
	}

	if !cookie.Expires.Equal(expire) {
		t.Errorf("Expected cookie expiry to be %s, but got %s", expire, cookie.Expires)
	}

	if cookie.Path != path {
		t.Errorf("Expected cookie path to be %s, but got %s", path, cookie.Path)
	}

	if cookie.HttpOnly != true {
		t.Errorf("Expected cookie HttpOnly to be true, but got %v", cookie.HttpOnly)
	}

	if cookie.Secure != true {
		t.Errorf("Expected cookie Secure to be true, but got %v", cookie.Secure)
	}

	if cookie.SameSite != http.SameSiteStrictMode {
		t.Errorf("Expected cookie SameSite to be SameSiteStrictMode, but got %v", cookie.SameSite)
	}
}
