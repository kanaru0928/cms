package handler

import (
	"fmt"
	"log"
	"net/http"
)

func DefineListener(port int) {
	mux := http.NewServeMux()

	articleHandler := NewArticleHandler()
	mux.Handle("/articles", articleHandler)

	log.Printf("Starting server on port %d...\n", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
