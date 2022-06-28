package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"premier_api/internal/app/model"

	"github.com/golang-jwt/jwt/v4"
)

func (s *server) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}
		reqBody := &model.ReqBody{}
		reqBody.Check = true
		reqBody.Body = string(b)
		type Token struct {
			Token string `json:"token"`
		}
		req := &Token{}
		if err := json.Unmarshal([]byte(reqBody.Body), &req); err != nil {
			fmt.Println("CheckToken error decode: ", reqBody.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return Secret_key, nil
		})
		if err != nil {
			fmt.Println("CheckToken token parse error: ", err)
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("action token: ", req.Token)
		} else {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		reqBody.Role = fmt.Sprint(claims["role"])
		reqBody.Name = fmt.Sprint(claims["email"])

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, reqBody)))
	})
}

func (s *server) WithoutCheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}
		reqBody := &model.ReqBody{}
		reqBody.Body = string(b)
		reqBody.Check = false
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, reqBody)))
	})
}
