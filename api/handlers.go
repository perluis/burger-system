// Package api contiene los handlers HTTP que procesan las peticiones
// de la API REST. Cada handler valida la entrada, interactúa con la
// base de datos y retorna respuestas JSON.
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/perluis/burger-system/menu"
	"github.com/perluis/burger-system/orders"
	"github.com/perluis/burger-system/storage"
)

// Server contiene la conexión a BD y expone los handlers HTTP
type Server struct {
	db *storage.Database
}

// NewServer crea una instancia del servidor con la BD inyectada
func NewServer(db *storage.Database) *Server {
	return &Server{db: db}
}

// ==================== HANDLERS DE HAMBURGUESAS ====================

// GetHamburguesas retorna la lista completa de hamburguesas en formato JSON
func (s *Server) GetHamburguesas(w http.ResponseWriter, r *http.Request) {
	hamburguesas, err := s.db.GetHamburguesas()
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error obteniendo hamburguesas: "+err.Error())
		return
	}
	responderJSON(w, http.StatusOK, hamburguesas)
}

// GetHamburguesaByID retorna una hamburguesa específica por su ID
func (s *Server) GetHamburguesaByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	hamburguesa, err := s.db.GetHamburguesaByID(id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error en base de datos: "+err.Error())
		return
	}
	if hamburguesa == nil {
		responderError(w, http.StatusNotFound, "Hamburguesa no encontrada")
		return
	}
	responderJSON(w, http.StatusOK, hamburguesa)
}

// CreateHamburguesa crea una nueva hamburguesa validando los datos de entrada
func (s *Server) CreateHamburguesa(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nombre       string   `json:"nombre"`
		Descripcion  string   `json:"descripcion"`
		Precio       float64  `json:"precio"`
		Categoria    string   `json:"categoria"`
		Ingredientes []string `json:"ingredientes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// NuevaHamburguesa valida precio > 0 y categoría válida
	hamburguesa, err := menu.NuevaHamburguesa(
		input.Nombre, input.Descripcion, input.Precio,
		input.Categoria, input.Ingredientes,
	)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generar ID secuencial desde la BD
	nextID, _ := s.db.GetNextHamburguesaID()
	hamburguesa.SetID(nextID)

	if err := s.db.CreateHamburguesa(hamburguesa); err != nil {
		responderError(w, http.StatusInternalServerError, "Error guardando en base de datos: "+err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, hamburguesa)
}

// UpdateHamburguesa actualiza nombre, descripción y precio de una hamburguesa
func (s *Server) UpdateHamburguesa(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var input struct {
		Nombre      string  `json:"nombre"`
		Descripcion string  `json:"descripcion"`
		Precio      float64 `json:"precio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := s.db.UpdateHamburguesa(id, input.Nombre, input.Descripcion, input.Precio); err != nil {
		responderError(w, http.StatusNotFound, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"message": "Hamburguesa actualizada exitosamente"})
}

// DeleteHamburguesa elimina una hamburguesa. Si tiene órdenes asociadas,
// informa al usuario que debe eliminar las órdenes primero.
func (s *Server) DeleteHamburguesa(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := s.db.DeleteHamburguesa(id)
	if err != nil {
		// Diferenciar error de foreign key vs error genérico
		if strings.Contains(err.Error(), "no disponible") || strings.Contains(err.Error(), "órdenes asociadas") {
			responderError(w, http.StatusConflict, err.Error())
			return
		}
		responderError(w, http.StatusNotFound, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"message": "Hamburguesa eliminada exitosamente"})
}

// ==================== HANDLERS DE ÓRDENES ====================

// GetAllOrdenes retorna todas las órdenes ordenadas por fecha descendente
func (s *Server) GetAllOrdenes(w http.ResponseWriter, r *http.Request) {
	ordenes, err := s.db.GetAllOrdenes()
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error obteniendo órdenes: "+err.Error())
		return
	}
	responderJSON(w, http.StatusOK, ordenes)
}

// CreateOrden crea una orden con sus items, validando que cada hamburguesa exista
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

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	orden, err := orders.NuevaOrden(input.TipoOrden, input.NumeroMesa)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generar ID secuencial desde la BD
	nextID, _ := s.db.GetNextOrdenID()
	orden.ID = nextID

	// Agregar cada item verificando que la hamburguesa exista en BD
	for _, item := range input.Items {
		hamburguesa, err := s.db.GetHamburguesaByID(item.HamburguesaID)
		if err != nil || hamburguesa == nil {
			responderError(w, http.StatusBadRequest, "Hamburguesa no encontrada: "+item.HamburguesaID)
			return
		}

		if err := orden.AgregarItem(
			hamburguesa.GetID(), hamburguesa.Nombre,
			item.Cantidad, hamburguesa.Precio, item.Notas,
		); err != nil {
			responderError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	if err := s.db.CreateOrden(orden); err != nil {
		responderError(w, http.StatusInternalServerError, "Error guardando orden: "+err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, orden)
}

// GetOrdenByID retorna una orden específica con todos sus items
func (s *Server) GetOrdenByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	orden, err := s.db.GetOrdenByID(id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error en base de datos: "+err.Error())
		return
	}
	if orden == nil {
		responderError(w, http.StatusNotFound, "Orden no encontrada")
		return
	}
	responderJSON(w, http.StatusOK, orden)
}

// UpdateOrdenEstado cambia el estado de una orden
// Estados válidos: PENDIENTE, EN_COCINA, LISTA, ENTREGADA, CANCELADA
func (s *Server) UpdateOrdenEstado(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var input struct {
		Estado string `json:"estado"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := s.db.UpdateOrdenEstado(id, input.Estado); err != nil {
		responderError(w, http.StatusNotFound, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"message": "Estado de orden actualizado exitosamente"})
}

// DeleteOrden elimina una orden y todos sus items asociados
func (s *Server) DeleteOrden(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := s.db.DeleteOrden(id); err != nil {
		responderError(w, http.StatusNotFound, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"message": "Orden eliminada exitosamente"})
}

// ==================== FUNCIONES AUXILIARES ====================

// responderJSON serializa datos a JSON y los envía con el código HTTP indicado
func responderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// responderError envía un mensaje de error en formato JSON
func responderError(w http.ResponseWriter, status int, mensaje string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": mensaje})
}
