package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"premier_api/internal/app/store"

	"github.com/gorilla/mux"
)

const (
	sessionName        = "premier_session"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var Secret_key = []byte("secretkeymsd")
var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errWrongRole                = errors.New("wrong role")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

type respondData struct {
	Result string      `json:"result"`
	Err    interface{} `json:"error,omitempty"`
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")

	checktoken := s.router.PathPrefix("/").Subrouter()
	checktoken.Use(s.CheckToken)
	withoutchecktoken := s.router.PathPrefix("/").Subrouter()
	withoutchecktoken.Use((s.WithoutCheckToken))

	// {"line": int, "token": string}
	checktoken.HandleFunc("/api/production/last", s.handleGetLast()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/last", s.handleGetLast()).Methods("POST")

	// {"line": int, "token": string}
	checktoken.HandleFunc("/api/production/status", s.handlegeGetStatus()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/status", s.handlegeGetStatus()).Methods("POST")

	// {"line": int, "token": string}
	checktoken.HandleFunc("/api/production/today", s.handleToday()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/today", s.handleToday()).Methods("POST")

	// {"line": int, "token": string}
	checktoken.HandleFunc("/api/production/today/models", s.handleTodayModels()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/today/models", s.handleTodayModels()).Methods("POST")

	// {"line": int, "token": string}
	checktoken.HandleFunc("/api/production/sector/balance", s.handleSectorBalance()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/sector/balance", s.handleSectorBalance()).Methods("POST")

	// {"line": int, "token": string}
	checktoken.HandleFunc("/api/production/packing/last", s.handlegePackingtLast()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/packing/last", s.handlegePackingtLast()).Methods("POST")

	// {"token": string}
	checktoken.HandleFunc("/api/production/packing/today", s.handlegePackingToday()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/packing/today", s.handlegePackingToday()).Methods("POST")

	// {"token": string}
	checktoken.HandleFunc("/api/production/packing/today/serial", s.handlegePackingTodaySerial()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/packing/today/serial", s.handlegePackingTodaySerial()).Methods("POST")

	// {"token": string}
	checktoken.HandleFunc("/api/production/packing/today/models", s.handlegePackingTodayModels()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/packing/today/models", s.handlegePackingTodayModels()).Methods("POST")

	// {"token": string}
	checktoken.HandleFunc("/api/production/lines", s.handleGetLines()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/lines", s.handlegePackingTodayModels()).Methods("POST")

	// {"token": string}
	checktoken.HandleFunc("/api/production/defects/types", s.handleGetDefectsTypes()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/defects/types", s.handleGetDefectsTypes()).Methods("POST")

	// {"id": int, "token": string}
	checktoken.HandleFunc("/api/production/defects/types/delete", s.handleDeleteDefectsTypes()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/defects/types/delete", s.handleDeleteDefectsTypes()).Methods("POST")

	// {"line": int, "name": string}
	checktoken.HandleFunc("/api/production/defects/types/add", s.handleAddDefetctTypes()).Methods("POST")
	withoutchecktoken.HandleFunc("/production/defects/types/add", s.handleAddDefetctTypes()).Methods("POST")

	//{"serial": string, "checkpoint_id": int, "defect_id": string, "token": string}
	checktoken.HandleFunc("/api/production/defects/add", s.handleAddDefets()).Methods("POST")

	//"date1": string, "date2": string, "line": int, "token": string}
	checktoken.HandleFunc("/api/report/bydate/models/serial", s.handleGetByDateSerial()).Methods("POST")
	withoutchecktoken.HandleFunc("/report/bydate/models/serial", s.handleGetByDateSerial()).Methods("POST")

	//"date1": string, "date2": string, "line": int, "token": string}
	checktoken.HandleFunc("/api/report/bydate", s.handleGetByDate()).Methods("POST")
	withoutchecktoken.HandleFunc("/report/bydate", s.handleGetByDate()).Methods("POST")

	//"date1": string, "date2": string, "line": int, "token": string}
	checktoken.HandleFunc("/api/report/bydate/models", s.handleByDateModels()).Methods("POST")
	withoutchecktoken.HandleFunc("/report/bydate/models", s.handleByDateModels()).Methods("POST")
	//{"token": string}
	checktoken.HandleFunc("/api/report/remont", s.handleGetRemont()).Methods("POST")
	withoutchecktoken.HandleFunc("/report/remont", s.handleGetRemont()).Methods("POST")

	//{"token": string}
	checktoken.HandleFunc("/api/report/remont/today", s.handleGetRemontToday()).Methods("POST")
	withoutchecktoken.HandleFunc("/report/remont/today", s.handleGetRemontToday()).Methods("POST")

	//"date1": string, "date2": string, "token": string}
	checktoken.HandleFunc("/api/report/remont/bydate", s.handleGetRemontByDate()).Methods("POST")
	withoutchecktoken.HandleFunc("/report/remont/bydate", s.handleGetRemontByDate()).Methods("POST")

	//{"id": int, "token": string}
	checktoken.HandleFunc("/api/report/remont/update", s.handleUpdateRemont()).Methods("POST")

	//{"token": string}
	checktoken.HandleFunc("/api/components", s.handleComponents()).Methods("POST")
	withoutchecktoken.HandleFunc("/components", s.handleComponents()).Methods("POST")

	//routes for only production processes
	//{"line": int, "serial": string}
	withoutchecktoken.HandleFunc("/production/serial/input", s.handleSerialInput()).Methods("POST")

	//"serial": string, "packing": string}
	withoutchecktoken.HandleFunc("/production/packing/serial/input", s.handlePackingSerialInput()).Methods("POST")

	//{"serial":string}
	withoutchecktoken.HandleFunc("/checkserial", s.handleGetInfoBySerial()).Methods("POST")

	//data from galileo {model.Galileo}
	withoutchecktoken.HandleFunc("/galileo/input", s.handleGalileo()).Methods("POST")

	//route for registering
	checktoken.HandleFunc("/register", s.handleRegister()).Methods("POST")

	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		w.Header().Set("Role", "admin")
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
