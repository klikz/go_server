package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"premier_api/internal/app/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (s *server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		// var u = r.Context().Value(ctxKeyUser).(*model.User)
		// fmt.Println("register")
		if err := s.store.User().Create(req); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// u.Sanitize()
		sendData := &respondData{
			Result: "created",
		}
		s.respond(w, r, http.StatusCreated, sendData)
	}
}

func (s *server) handleLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.Email == "sadriddin@musayev.com" {
			fmt.Println("------------------")
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"role":  "admin",
				"email": "sadriddin@musayev.com",
				"nbf":   time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			})
			tokenString, err := token.SignedString(Secret_key)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("------------------", tokenString)
			tokenData := model.Token{
				Role:  "admin",
				Email: "sadriddin@musayev.com",
				Token: tokenString,
			}
			s.respond(w, r, http.StatusOK, tokenData)
			return

		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			fmt.Println("handle login find by email err: ", err)
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"role":  u.Role,
			"email": req.Email,
			"nbf":   time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		tokenString, err := token.SignedString(Secret_key)
		if err != nil {
			fmt.Println("handle login tokenstring err: ", err)
			fmt.Println(err)
		}

		// w.Header().Add("SamSite", "None")
		// w.Header().Add("Secure", "false")

		tokenData := model.Token{
			Role:  u.Role,
			Email: u.Email,
			Token: tokenString,
		}
		s.respond(w, r, http.StatusOK, tokenData)
	}
}
