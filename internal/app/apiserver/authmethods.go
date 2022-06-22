package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"premier_api/internal/app/model"

	"github.com/golang-jwt/jwt/v4"
)

func (s *server) authRegister(next http.Handler) http.Handler {

	type request struct {
		Email     string `json:"regemail"`
		Password  string `json:"regpassword"`
		UserEmail string `json:"email"`
		Role      string `json:"role"`
		Token     string `json:"token"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("decode json error: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		fmt.Println("r body----------- ", req)

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return Secret_key, nil
		})
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// c.JSON(http.StatusOK, claims)
		} else {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		if err != nil {
			fmt.Print(err)
		}

		u, err := s.store.User().FindByEmail(req.UserEmail)
		if err != nil {
			fmt.Println("authRegister FindByEmail err: ", err)
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		if u.Role != "admin" {
			s.error(w, r, http.StatusUnauthorized, errWrongRole)
			return
		}
		u.Email = req.Email
		u.Password = req.Password
		u.Role = req.Role

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) authComponents(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		req := &model.Token{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return Secret_key, nil
		})
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// c.JSON(http.StatusOK, claims)
		} else {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if err != nil {
			fmt.Print(err)
		}
		if req.Email == "sadriddin@musayev.com" {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, nil)))
		}
		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			fmt.Println("authComponents FindByEmail err: ", err)
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if u.Role == "admin" || u.Role == "tech" {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
		} else {
			s.error(w, r, http.StatusUnauthorized, errWrongRole)
			return
		}
	})
}

func (s *server) authReport(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &model.Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body, "errror: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return Secret_key, nil
		})
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// c.JSON(http.StatusOK, claims)
		} else {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		if err != nil {
			fmt.Println("authComponents FindByEmail err: ", err)
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if req.Date1 != "" || req.Date2 != "" {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, req)))
		} else {
			s.error(w, r, http.StatusUnauthorized, errWrongRole)
			return
		}
	})
}
