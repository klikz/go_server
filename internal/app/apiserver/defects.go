package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"premier_api/internal/app/model"
)

func (s *server) handleGetDefectsTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlegetdefecttypes", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		result, err := s.store.User().GetDefectsTypes()

		if err != nil {
			fmt.Println("handleGetDefetcsTypes err: ", err)
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
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handledeletedefecttypes", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handledeletedefecttypes error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		_, err := s.store.User().DeleteDefectsTypes(req.ID)

		if err != nil {

			fmt.Println("handledeletedefecttypes err: ", err)
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
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handleadddefecttypes", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleAddDefetctTypes error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Name == "" || req.Line == 0 {
			s.error(w, r, http.StatusBadRequest, errors.New("grong data"))
			return
		}
		_, err := s.store.User().AddDefectsTypes(req.Line, req.Name)

		if err != nil {

			fmt.Println("handleAddDefetctTypes err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		sendData := &respondData{
			Result: "ok",
		}
		s.respond(w, r, http.StatusCreated, sendData)
	}
}
func (s *server) handleAddDefets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handleadddefects", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &model.OtkAddDefect{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleAddDefets error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		_, err := s.store.User().AddDefects(req, u.Name)

		if err != nil {
			fmt.Println("handleAddDefets err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		sendData := &respondData{
			Result: "ok",
		}
		s.respond(w, r, http.StatusCreated, sendData)
	}
}
func (s *server) handleGetRemont() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlegetremont", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		result, err := s.store.User().GetRemont()
		if err != nil {
			fmt.Println("handleGetRemont err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleGetRemontToday() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlegetremonttoday", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		result, err := s.store.User().GetRemontToday()
		if err != nil {
			fmt.Println("handleGetRemontToday err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleGetRemontByDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlegetremontbydate", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &model.Request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleGetRemontByDate error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Date1 == "" || req.Date2 == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("wrong data"))
			return
		}
		result, err := s.store.User().GetRemontByDate(req.Date1, req.Date2)
		if err != nil {
			fmt.Println("handleGetRemontByDate err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleUpdateRemont() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		fmt.Println("remont update: ", u)
		if u.Check {
			_, err := s.store.User().CheckRole("handleupdateremont", u.Role)
			if err != nil {
				fmt.Println("handleUpdateRemont wrong role")
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &model.TokenPerson{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleUpdateRemont error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		_, err := s.store.User().UpdateRemont(u.Name, req.Item)
		if err != nil {
			fmt.Println("handleUpdateRemont err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		sendData := &respondData{
			Result: "ok",
		}
		s.respond(w, r, http.StatusOK, sendData)
	}
}
