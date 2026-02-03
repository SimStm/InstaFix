package handlers

import (
	"encoding/json"
	scraper "instafix/handlers/scraper"
	"net/http"
)

// APIResponse representa o formato JSON de resposta
type APIResponse struct {
	Success  bool            `json:"success"`
	PostID   string          `json:"post_id"`
	Username string          `json:"username"`
	Caption  string          `json:"caption"`
	Medias   []scraper.Media `json:"medias"`
	Error    string          `json:"error,omitempty"`
}

// API retorna os dados do post em JSON sem verificação de User-Agent
func API(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obter postID do query parameter
	postID := r.URL.Query().Get("postid")
	if postID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "postid parameter is required",
		})
		return
	}

	// Obter dados do scraper (já tem cache!)
	item, err := scraper.GetData(postID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Verificar se encontrou dados
	if len(item.Medias) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "no media found",
		})
		return
	}

	// Retornar dados em JSON
	json.NewEncoder(w).Encode(APIResponse{
		Success:  true,
		PostID:   item.PostID,
		Username: item.Username,
		Caption:  item.Caption,
		Medias:   item.Medias,
	})
}
