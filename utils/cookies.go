package utils

import "net/http"

func GetCookie(name string, r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	val := tokenCookie.Value
	return val, nil
}