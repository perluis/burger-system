package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/perluis/burger-system/menu"
	"github.com/perluis/burger-system/orders"
	"github.com/perluis/burger-system/storage"
)

// Server contiene la base de datos y los handlers
type Server struct {
	db *storage.Database
}

// NewServer crea un nuevo servidor con conexión a BD
func NewServer(db *storage.Database) *Server {
	return &Server{
		db: db,
	}
}

// ==================== HANDLERS HAMBURGUESAS ====================

// GetHamburguesas maneja GET /api/hamburguesas
func (s *Server) GetHamburguesas(w http.ResponseWriter, r *http.Request) {
	hamburguesas, err := s.db.GetHamburguesas()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error obteniendo hamburguesas: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hamburguesas)
}

// GetHamburguesaByID maneja GET /api/hamburguesas/{id}
func (s *Server) GetHamburguesaByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	hamburguesa, err := s.db.GetHamburguesaByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error en base de datos: " + err.Error(),
		})
		return
	}

	if hamburguesa == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Hamburguesa no encontrada",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hamburguesa)
}

// CreateHamburguesa maneja POST /api/hamburguesas
func (s *Server) CreateHamburguesa(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nombre       string   `json:"nombre"`
		Descripcion  string   `json:"descripcion"`
		Precio       float64  `json:"precio"`
		Categoria    string   `json:"categoria"`
		Ingredientes []string `json:"ingredientes"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	// Crear hamburguesa con validaciones
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

	// Obtener siguiente ID de la BD
	nextID, _ := s.db.GetNextHamburguesaID()
	hamburguesa.SetID(nextID)

	// Guardar en BD
	err = s.db.CreateHamburguesa(hamburguesa)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error guardando en base de datos: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hamburguesa)
}

// UpdateHamburguesa maneja PUT /api/hamburguesas/{id}
func (s *Server) UpdateHamburguesa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		Nombre      string  `json:"nombre"`
		Descripcion string  `json:"descripcion"`
		Precio      float64 `json:"precio"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	err = s.db.UpdateHamburguesa(id, input.Nombre, input.Descripcion, input.Precio)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hamburguesa actualizada exitosamente",
	})
}

// DeleteHamburguesa maneja DELETE /api/hamburguesas/{id}
func (s *Server) DeleteHamburguesa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := s.db.DeleteHamburguesa(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hamburguesa eliminada exitosamente",
	})
}

// ==================== HANDLERS ÓRDENES ====================

// CreateOrden maneja POST /api/ordenes
func (s *Server) CreateOrden(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TipoOrden  string `json:"tipoOrden"`
		NumeroMesa int    `json:"numeroMesa"`
		Items      []struct {
			HamburguesaID string `json:"hamburguesaID"`
			Cantidad      int    `json:"cantidad"`
			Notas         string `json:"notas"`
		} `json:"items"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	// Crear orden
	orden, err := orders.NuevaOrden(input.TipoOrden, input.NumeroMesa)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Obtener siguiente ID
	nextID, _ := s.db.GetNextOrdenID()
	orden.ID = nextID

	// Agregar items
	for _, item := range input.Items {
		hamburguesa, err := s.db.GetHamburguesaByID(item.HamburguesaID)
		if err != nil || hamburguesa == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Hamburguesa no encontrada: " + item.HamburguesaID,
			})
			return
		}

		err = orden.AgregarItem(
			hamburguesa.GetID(),
			hamburguesa.Nombre,
			item.Cantidad,
			hamburguesa.Precio,
			item.Notas,
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
	}

	// Guardar en BD
	err = s.db.CreateOrden(orden)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error guardando orden: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orden)
}

// GetOrdenByID maneja GET /api/ordenes/{id}
func (s *Server) GetOrdenByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	orden, err := s.db.GetOrdenByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error en base de datos: " + err.Error(),
		})
		return
	}

	if orden == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Orden no encontrada",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orden)
}

// UpdateOrdenEstado maneja PUT /api/ordenes/{id}/estado
func (s *Server) UpdateOrdenEstado(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		Estado string `json:"estado"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	err = s.db.UpdateOrdenEstado(id, input.Estado)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Estado de orden actualizado exitosamente",
	})
}
