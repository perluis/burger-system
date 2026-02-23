package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/perluis/burger-system/storage"
)

// Server contiene el store y los handlers
type Server struct {
	store *storage.Store
}

// NewServer crea un nuevo servidor
func NewServer(store *storage.Store) *Server {
	return &Server{
		store: store,
	}
}

// GetHamburguesas maneja GET /api/hamburguesas
func (s *Server) GetHamburguesas(w http.ResponseWriter, r *http.Request) {
	// Obtener todas las hamburguesas del store
	hamburguesas := s.store.GetHamburguesas()

	// Convertir a JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hamburguesas)
}

// GetHamburguesaByID maneja GET /api/hamburguesas/{id}
func (s *Server) GetHamburguesaByID(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID de la URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Buscar la hamburguesa
	hamburguesa := s.store.GetHamburguesaByID(id)

	if hamburguesa == nil {
		// No existe
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Hamburguesa no encontrada",
		})
		return
	}

	// Devolver la hamburguesa en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hamburguesa)
}
