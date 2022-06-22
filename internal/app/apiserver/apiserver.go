package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"

	"premier_api/internal/app/store/sqlstore"

	_ "github.com/lib/pq" // ...
	"github.com/rs/cors"
)

// Start ...
func Start(config *Config) error {

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	srv := newServer(store)

	var c = cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowCredentials: true,
		AllowedMethods: []string{"POST", "PUT", "GET", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"premier_session"},
		MaxAge:         300,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(srv)

	fmt.Println("Server started")

	return http.ListenAndServe(config.BindAddr, handler)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("newDB postgres err: ", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println("newDB ping err: ", err)
		return nil, err
	}

	return db, nil
}
