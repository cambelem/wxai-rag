package api

import (
	"encoding/json"
	"net/http"
	"wxai-rag/internal/watsonxai"

	"github.com/gorilla/mux"
)

type GenerateTextRequest struct {
	ModelID    string                 `json:"model_id"`
	Input      string                 `json:"input"`
	Parameters map[string]interface{} `json:"parameters"`
}

type GenerateTextResponse struct {
	GeneratedText string `json:"generated_text"`
}

// func SetupRouter(esClient *elasticsearch.Client, wxClient *watsonxai.Client) *mux.Router {
func SetupRouter(wxClient *watsonxai.Client) *mux.Router {

	router := mux.NewRouter()

	// router.HandleFunc("/search", searchHandler(esClient)).Methods("GET")
	router.HandleFunc("/generate-text", generateHandler(wxClient)).Methods("POST")

	return router
}

// func searchHandler(esClient *elasticsearch.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Implement search logic using esClient
// 	}
// }

func generateHandler(wxClient *watsonxai.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req GenerateTextRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.ModelID == "" || req.Input == "" || req.Parameters == nil {
			http.Error(w, "Missing required fields: model_id, input, and parameters", http.StatusBadRequest)
			return
		}

		payload := watsonxai.TextGenerationPayload{
			ModelID:    req.ModelID,
			Input:      req.Input,
			Parameters: req.Parameters,
			ProjectID:  wxClient.ProjectID,
		}

		result, err := wxClient.TextGeneration(payload)
		if err != nil {
			http.Error(w, "Failed to generate text: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := GenerateTextResponse{GeneratedText: result}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
