package main

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"rsoi2/src/gateway-service/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load the env vars: %v", err)
	}

	store := sessions.NewFilesystemStore("", []byte("something-very-secret"))
	gob.Register(map[string]interface{}{})

	r := handlers.Router(store)

	log.Println("server is listening on port: ", port)
	log.Printf("app started")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
