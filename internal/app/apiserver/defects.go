package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"premier_api/internal/app/model"
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

func (s *server) handleGetRemontByDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.Request)
		result, err := s.store.User().GetRemontByDate(u.Date1, u.Date2)
		if err != nil {
			fmt.Println("handleToday err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
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

func (s *server) handleAddDefets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.OtkAddDefectParsed)

		_, err := s.store.User().AddDefects(u) //serial, name string, checkpoint, defect int

		if err != nil {

			fmt.Println("handleSectorBalance err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
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
