package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"premier_api/internal/app/model"
)

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
		Line  int    `json:"line"`
		Token string `json:"string"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle today body: ", r.Body)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println("error in decode json: ", r.Body, "errror: ", err)
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
