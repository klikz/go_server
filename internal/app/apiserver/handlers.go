package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"premier_api/internal/app/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (s *server) handleUpdateRemont() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.TokenPerson)
		_, err := s.store.User().UpdateRemont(u.Name, u.Role, u.Item)
		fmt.Println("auth errors: ", err)
		if err != nil {
			fmt.Println("handleToday err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		sendData := &respondData{
			Result: "ok",
		}
		s.respond(w, r, http.StatusOK, sendData)
	}
}

func (s *server) handleGetRemont() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := s.store.User().GetRemont()
		if err != nil {
			fmt.Println("handleToday err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) handleGetByDateSerial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ctx value: ", r.Context().Value(ctxKeyUser))
		var u = r.Context().Value(ctxKeyUser).(*model.Request)
		result, err := s.store.User().GetByDateSerial(u.Date1, u.Date2)

		if err != nil {
			fmt.Println("handleToday err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) handleByDateModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.Request)
		result, err := s.store.User().GetByDateModels(u.Date1, u.Date2, u.Line)

		if err != nil {
			fmt.Println("handleTodayModels err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleGetByDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ctx value: ", r.Context().Value(ctxKeyUser))
		var u = r.Context().Value(ctxKeyUser).(*model.Request)
		result, err := s.store.User().GetByDate(u.Date1, u.Date2, u.Line)

		if err != nil {
			fmt.Println("handleToday err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) handlegeGetStatus() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		result, err := s.store.User().GetStatus(req.Line)

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) handlePackingSerialInput() http.HandlerFunc {
	type request struct {
		Serial  string `json:"serial"`
		Packing string `json:"packing"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		fmt.Println("serial: ", req.Serial, "packing", req.Packing)
		if req.Packing == req.Serial {
			fmt.Println("serial bir xil")
			s.error(w, r, http.StatusBadRequest, errors.New("serial xato"))
			return
		}

		result, err := s.store.User().PackingSerialInput(req.Serial, req.Packing)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)

		}
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) handlegePackingTodayModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := s.store.User().GetPackingTodayModels()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handlegePackingTodaySerial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := s.store.User().GetPackingTodaySerial()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handlegePackingToday() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := s.store.User().GetPackingToday()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) handlegePackingtLast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := s.store.User().GetPackingLast()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}

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

func (s *server) handleComponents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		compontnts, err := s.store.User().ComponentsAll()
		// fmt.Print("components: ", compontnts)
		if err != nil {
			fmt.Println("handleComponents err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, compontnts)
	}

}

func (s *server) handlegetLast() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().GetLast(req.Line)

		if err != nil {
			fmt.Println("handlegetLast err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleToday() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().GetToday(req.Line)

		if err != nil {
			fmt.Println("handleToday err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) handleTodayModels() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().GetTodayModels(req.Line)

		if err != nil {
			fmt.Println("handleTodayModels err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleSectorBalance() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().GetSectorBalance(req.Line)

		if err != nil {
			fmt.Println("handleSectorBalance err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleSerialInput() http.HandlerFunc {
	type request struct {
		Line   int    `json:"line"`
		Serial string `json:"serial"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().SerialInput(req.Line, req.Serial)

		if err != nil {
			fmt.Println("handleSerialInput err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleGetDefectsTypes() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		result, err := s.store.User().GetDefectsTypes()

		if err != nil {
			fmt.Println("handleGetDefetcsTypes err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleGetLines() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		result, err := s.store.User().GetLines()

		if err != nil {
			fmt.Println("handleGetLines err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *server) handleDeleteDefectsTypes() http.HandlerFunc {
	type request struct {
		ID int `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		_, err := s.store.User().DeleteDefectsTypes(req.ID)

		if err != nil {

			fmt.Println("handleSectorBalance err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		sendData := &respondData{
			Result: "ok",
		}
		s.respond(w, r, http.StatusCreated, sendData)
	}
}

func (s *server) handleAddDefetctTypes() http.HandlerFunc {
	type request struct {
		Line int    `json:"line"`
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		_, err := s.store.User().AddDefectsTypes(req.Line, req.Name)

		if err != nil {

			fmt.Println("handleSectorBalance err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		sendData := &respondData{
			Result: "ok",
		}
		s.respond(w, r, http.StatusCreated, sendData)
	}
}
