package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"premier_api/internal/app/model"
)

func (s *server) handleGetLast() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlegetlast", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleGetLast error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().GetLast(req.Line)
		if err != nil {
			fmt.Println("handleGetLast err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handlegeGetStatus() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlegetstatus", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handlegeGetStatus error decode: ", u.Body, "error: ", err)
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
func (s *server) handleToday() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handletoday", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleToday error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().GetToday(req.Line)

		if err != nil {
			fmt.Println("handleToday err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleTodayModels() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handletodaymodels", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

		}

		req := &request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleTodayModels error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().GetTodayModels(req.Line)

		if err != nil {
			fmt.Println("handleTodayModels err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		fmt.Println("handleTodayModels: ", result)
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleSectorBalance() http.HandlerFunc {
	type request struct {
		Line int `json:"line"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlesectorbalance", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleSectorBalance error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().GetSectorBalance(req.Line)

		if err != nil {
			fmt.Println("handleSectorBalance err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handlegePackingtLast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlepackinglast", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		result, err := s.store.User().GetPackingLast()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handlegePackingToday() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlepackingtoday", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		result, err := s.store.User().GetPackingToday()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handlegePackingTodaySerial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlepackingtodayserial", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		result, err := s.store.User().GetPackingTodaySerial()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handlegePackingTodayModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlepackingtodaymodels", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		result, err := s.store.User().GetPackingTodayModels()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleGetLines() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlegetlines", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		result, err := s.store.User().GetLines()
		if err != nil {
			fmt.Println("handleGetLines err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleGetByDateSerial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlegetbydateserial", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &model.Request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleToday error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Date1 == "" || req.Date2 == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("wrong date"))
			return
		}
		result, err := s.store.User().GetByDateSerial(req.Date1, req.Date2)
		if err != nil {
			fmt.Println("handleToday err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
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
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		req := &request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleSerialInput error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		result, err := s.store.User().SerialInput(req.Line, req.Serial)
		if err != nil {
			fmt.Println("handleSerialInput err: ", err)
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
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		req := &request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handlePackingSerialInput error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Packing == req.Serial || req.Packing == "" {
			fmt.Println("serial xato")
			s.error(w, r, http.StatusBadRequest, errors.New("serial xato"))
			return
		}
		result, err := s.store.User().PackingSerialInput(req.Serial, req.Packing)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return

		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleByDateModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlebydatemodels", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &model.Request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleByDateModels error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Date1 == "" || req.Date2 == "" || req.Line == 0 {
			s.error(w, r, http.StatusBadRequest, errors.New("wrong data"))
			return
		}
		result, err := s.store.User().GetByDateModels(req.Date1, req.Date2, req.Line)
		if err != nil {
			fmt.Println("handleByDateModels err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleGetByDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlebydate", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &model.Request{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleGetByDate error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Date1 == "" || req.Date2 == "" || req.Line == 0 {
			s.error(w, r, http.StatusBadRequest, errors.New("wrong data"))
			return
		}
		result, err := s.store.User().GetByDate(req.Date1, req.Date2, req.Line)
		if err != nil {
			fmt.Println("handleGetByDate err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, result)
	}
}
func (s *server) handleComponents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlecomponents", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &model.ReqBody{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleComponents error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		compontnts, err := s.store.User().ComponentsAll()
		if err != nil {
			fmt.Println("handleComponents err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, compontnts)
	}
}
func (s *server) handleGetInfoBySerial() http.HandlerFunc {
	type Serial struct {
		Serial string `json"serial"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var u = r.Context().Value(ctxKeyUser).(*model.ReqBody)
		if u.Check {
			_, err := s.store.User().CheckRole("handlegetinfobyserial", u.Role)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		req := &Serial{}
		if err := json.Unmarshal([]byte(u.Body), &req); err != nil {
			fmt.Println("handleGetInfoBySerial error decode: ", u.Body, "error: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		compontnts, err := s.store.User().GetInfoBySerial(req.Serial)
		if err != nil {
			fmt.Println("handleGetInfoBySerial err: ", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, compontnts)
	}
}
