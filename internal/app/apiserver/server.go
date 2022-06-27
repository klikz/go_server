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
	//for reports read only permission
	// checkToken := s.router.PathPrefix("/").Subrouter()
	reports := s.router.PathPrefix("/").Subrouter()
	remont := s.router.PathPrefix("/").Subrouter()
	authOtkAddDefect := s.router.PathPrefix("/").Subrouter()
	checktoken := s.router.PathPrefix("/").Subrouter()

	checktoken.Use(s.CheckToken)

	reports.Use(s.authReport)
	remont.Use(s.authRemont)
	authOtkAddDefect.Use(s.authOtkAddDefect)

	reports.HandleFunc("/report/bydate", s.handleGetByDate()).Methods("POST")
	reports.HandleFunc("/report/bydate/models", s.handleByDateModels()).Methods("POST")
	reports.HandleFunc("/report/bydate/models/serial", s.handleGetByDateSerial()).Methods("POST")
	reports.HandleFunc("/report/remont", s.handleGetRemont()).Methods("POST")
	reports.HandleFunc("/report/remont/bydate", s.handleGetRemontByDate()).Methods("POST")

	remont.HandleFunc("/report/remont/update", s.handleUpdateRemont()).Methods("POST")

	//routes for only production processes
	checktoken.HandleFunc("/production/last", s.handlegetLast()).Methods("POST")
	checktoken.HandleFunc("/production/status", s.handlegeGetStatus()).Methods("POST")
	checktoken.HandleFunc("/production/today", s.handleToday()).Methods("POST")
	checktoken.HandleFunc("/production/today/models", s.handleTodayModels()).Methods("POST")
	checktoken.HandleFunc("/production/sector/balance", s.handleSectorBalance()).Methods("POST")
	checktoken.HandleFunc("/production/serial/input", s.handleSerialInput()).Methods("POST")

	checktoken.HandleFunc("/production/packing/last", s.handlegePackingtLast()).Methods("POST")
	checktoken.HandleFunc("/production/packing/today", s.handlegePackingToday()).Methods("POST")
	checktoken.HandleFunc("/production/packing/today/serial", s.handlegePackingTodaySerial()).Methods("POST")
	checktoken.HandleFunc("/production/packing/today/models", s.handlegePackingTodayModels()).Methods("POST")
	checktoken.HandleFunc("/production/packing/serial/input", s.handlePackingSerialInput()).Methods("POST")

	checktoken.HandleFunc("/production/lines", s.handleGetLines()).Methods("POST")
	checktoken.HandleFunc("/production/defects/types", s.handleGetDefectsTypes()).Methods("POST")
	checktoken.HandleFunc("/production/defects/types/delete", s.handleDeleteDefectsTypes()).Methods("POST")
	checktoken.HandleFunc("/production/defects/types/add", s.handleAddDefetctTypes()).Methods("POST")
	checktoken.HandleFunc("/production/defects/add", s.handleAddDefets()).Methods("POST")

	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

	//route for registering
	register := s.router.PathPrefix("/").Subrouter()
	// register.Use(s.authRegister)
	register.HandleFunc("/register", s.handleRegister()).Methods("POST")

	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")

	// s.router.HandleFunc("/components", s.handleComponents()).Methods("GET")
	components := s.router.PathPrefix("/components").Subrouter()
	components.Use(s.authComponents)
	components.HandleFunc("", s.handleComponents()).Methods("POST")
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
