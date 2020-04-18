package web

import (
	"fmt"
	"log"
	"net/http"
)

func Serve(port int) {
	http.HandleFunc("/health", healthHandler)
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
