package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"scullion/config"
)

func Serve(port int, taskDefs []config.TaskDef) {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/config", configWrapper(taskDefs))
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func configWrapper(taskDefs []config.TaskDef) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(taskDefs)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "OK")
}
