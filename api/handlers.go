package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/perluis/burger-system/menu"
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

// CreateHamburguesa maneja POST /api/hamburguesas
func (s *Server) CreateHamburguesa(w http.ResponseWriter, r *http.Request) {
	// Estructura para recibir los datos del JSON
	var input struct {
		Nombre       string   `json:"nombre"`
		Descripcion  string   `json:"descripcion"`
		Precio       float64  `json:"precio"`
		Categoria    string   `json:"categoria"`
		Ingredientes []string `json:"ingredientes"`
	}

	// Decodificar el JSON del body
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	// Crear la hamburguesa usando la función del paquete menu
	hamburguesa, err := menu.NuevaHamburguesa(
		input.Nombre,
		input.Descripcion,
		input.Precio,
		input.Categoria,
		input.Ingredientes,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Agregar al store
	s.store.AddHamburguesa(hamburguesa)

	// Responder con la hamburguesa creada
	w.WriteHeader(http.StatusCreated) // 201
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hamburguesa)
}

// UpdateHamburguesa maneja PUT /api/hamburguesas/{id}
func (s *Server) UpdateHamburguesa(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID de la URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Estructura para recibir los datos
	var input struct {
		Nombre      string  `json:"nombre"`
		Descripcion string  `json:"descripcion"`
		Precio      float64 `json:"precio"`
	}

	// Decodificar JSON
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	// Actualizar en el store
	success := s.store.UpdateHamburguesa(id, input.Nombre, input.Descripcion, input.Precio)

	if !success {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Hamburguesa no encontrada",
		})
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hamburguesa actualizada exitosamente",
	})
}

// DeleteHamburguesa maneja DELETE /api/hamburguesas/{id}
func (s *Server) DeleteHamburguesa(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID de la URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Eliminar del store
	success := s.store.DeleteHamburguesa(id)

	if !success {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Hamburguesa no encontrada",
		})
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hamburguesa eliminada exitosamente",
	})
}
