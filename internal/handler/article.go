package handler

import "net/http"

type ArticleHandler struct {
	mux *http.ServeMux
}

func NewArticleHandler() *ArticleHandler {
	mux := http.NewServeMux()
	handler := &ArticleHandler{mux: mux}

	handler.mux.HandleFunc("GET /articles", handler.GetArticles)

	return handler
}

func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List of articles"))
}

func (h *ArticleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
