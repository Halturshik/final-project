package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Halturshik/final-project/GO/db"
	"github.com/Halturshik/final-project/GO/pkg/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	dbConnection, err := db.ConnectDB(dbPath)
	if err != nil {
		log.Fatalf("Ошибка при подключении к БД: %v", err)
	}
	defer dbConnection.Close()

	store := db.NewStore(dbConnection)
	apiServer := api.NewAPI(store)

	r := chi.NewRouter()

	apiServer.Init(r)

	fs := http.FileServer(http.Dir("web"))
	r.Handle("/*", http.StripPrefix("/", fs))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	http.ListenAndServe(":"+port, r)

}
