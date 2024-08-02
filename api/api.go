package api

import (
	"net/http"
	"wxai-rag/internal/watsonxai"
	"wxai-rag/pkg/elasticsearch"

	"github.com/gorilla/mux"
)

func SetupRouter(esClient *elasticsearch.Client, wxClient *watsonxai.Client) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/search", searchHandler(esClient)).Methods("GET")
	router.HandleFunc("/generate", generateHandler(wxClient)).Methods("POST")

	return router
}

func searchHandler(esClient *elasticsearch.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement search logic using esClient
	}
}

func generateHandler(wxClient *watsonxai.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement text generation logic using wxClient
	}
}
