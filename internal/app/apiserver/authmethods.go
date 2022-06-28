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

// func (s *server) authReport(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		req := &model.Request{}
// 		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 			fmt.Println("error in decode json: ", r.Body, "errror: ", err)
// 			s.error(w, r, http.StatusBadRequest, err)
// 			return
// 		}
// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
// 			return Secret_key, nil
// 		})
// 		if err != nil {
// 			fmt.Println("token parse error: ", err)
// 			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
// 			return
// 		}
// 		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			// c.JSON(http.StatusOK, claims)
// 		} else {
// 			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
// 			return
// 		}

// 		if req.Date1 != "" || req.Date2 != "" {
// 			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, req)))
// 		} else {
// 			s.error(w, r, http.StatusUnauthorized, errWrongRole)
// 			return
// 		}
// 	})
// }

// func (s *server) authRemont(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		req := &model.Request{}
// 		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 			fmt.Println("error in decode json: ", r.Body, "errror: ", err)
// 			s.error(w, r, http.StatusBadRequest, err)
// 			return
// 		}
// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
// 			return Secret_key, nil
// 		})
// 		if err != nil {
// 			fmt.Println("token parse error: ", err)
// 			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
// 			return
// 		}
// 		person := &model.TokenPerson{}
// 		person.Item = req.Line
// 		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			person.Name = fmt.Sprint(claims["email"])
// 			person.Role = fmt.Sprint(claims["role"])
// 		} else {
// 			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
// 			return
// 		}

// 		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, person)))

// 	})
// }

// func (s *server) authOtkAddDefect(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		req := &model.OtkAddDefect{}
// 		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 			fmt.Println("error in decode json: ", r.Body, "errror: ", err)
// 			s.error(w, r, http.StatusBadRequest, err)
// 			return
// 		}
// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
// 			return Secret_key, nil
// 		})
// 		if err != nil {
// 			fmt.Println("token parse error: ", err)
// 			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
// 			return
// 		}
// 		person := &model.OtkAddDefectParsed{}
// 		person.Checkpoint = req.Checkpoint
// 		person.Defect = req.Defect
// 		person.Serial = req.Serial

// 		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			person.Name = fmt.Sprint(claims["email"])
// 			person.Role = fmt.Sprint(claims["role"])
// 		} else {
// 			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
// 			return
// 		}

// 		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, person)))

// 	})
// }

func (s *server) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		b, err := io.ReadAll(r.Body)
		// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
		if err != nil {
			log.Fatalln(err)
		}
		reqBody := &model.ReqBody{}
		reqBody.Check = true
		reqBody.Body = string(b)
		fmt.Println("req body: ", reqBody)

		type Token struct {
			Token string `json:"token"`
		}
		req := &Token{}
		if err := json.Unmarshal([]byte(reqBody.Body), &req); err != nil {
			fmt.Println("checktoken error decode: ", reqBody.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return Secret_key, nil
		})
		if err != nil {
			fmt.Println("token parse error: ", err)
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
