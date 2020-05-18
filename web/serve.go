package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"scullion/config"
)

func Serve(port int, cfg config.Config) {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/config", configWrapper(cfg))
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func configWrapper(cfg config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "OK")
}
