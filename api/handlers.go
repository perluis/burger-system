package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/perluis/burger-system/menu"
	"github.com/perluis/burger-system/orders"
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

// CreateOrden maneja POST /api/ordenes
func (s *Server) CreateOrden(w http.ResponseWriter, r *http.Request) {
	// Estructura para recibir los datos
	var input struct {
		TipoOrden  string `json:"tipoOrden"`
		NumeroMesa int    `json:"numeroMesa"`
		Items      []struct {
			HamburguesaID string `json:"hamburguesaID"`
			Cantidad      int    `json:"cantidad"`
			Notas         string `json:"notas"`
		} `json:"items"`
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

	// Crear la orden
	orden, err := orders.NuevaOrden(input.TipoOrden, input.NumeroMesa)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Agregar items a la orden
	for _, item := range input.Items {
		// Buscar la hamburguesa para obtener nombre y precio
		hamburguesa := s.store.GetHamburguesaByID(item.HamburguesaID)
		if hamburguesa == nil {
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

	// Guardar en el store
	s.store.AddOrden(orden)

	// Responder
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orden)
}

// GetOrdenByID maneja GET /api/ordenes/{id}
func (s *Server) GetOrdenByID(w http.ResponseWriter, r *http.Request) {
	// Obtener ID de la URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Buscar la orden
	orden := s.store.GetOrdenByID(id)

	if orden == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Orden no encontrada",
		})
		return
	}

	// Responder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orden)
}

// UpdateOrdenEstado maneja PUT /api/ordenes/{id}/estado
func (s *Server) UpdateOrdenEstado(w http.ResponseWriter, r *http.Request) {
	// Obtener ID de la URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Estructura para recibir el nuevo estado
	var input struct {
		Estado string `json:"estado"`
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

	// Actualizar estado
	success := s.store.UpdateOrdenEstado(id, input.Estado)

	if !success {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Orden no encontrada o estado inválido",
		})
		return
	}

	// Responder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Estado de orden actualizado exitosamente",
	})
}
