// Package menu define el modelo de datos para los productos del restaurante.
// Implementa encapsulación (campos privados id y fechaCreado) y
// las interfaces Informable, Identificable y Validable.
package menu

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Hamburguesa representa un producto del menú.
// Los campos id y fechaCreado son privados (encapsulación),
// accesibles solo mediante getters/setters.
type Hamburguesa struct {
	id           string // Privado: acceso solo por GetID/SetID
	Nombre       string
	Descripcion  string
	Precio       float64
	Categoria    string // Clasica, Premium, Vegetariana
	Ingredientes []string
	Disponible   bool
	fechaCreado  time.Time // Privado: solo lectura
}

// NuevaHamburguesa crea una hamburguesa validando precio y categoría.
// Retorna error si los datos no son válidos.
func NuevaHamburguesa(nombre, descripcion string, precio float64, categoria string, ingredientes []string) (*Hamburguesa, error) {
	if precio <= 0 {
		return nil, errors.New("el precio debe ser mayor a 0")
	}

	// Validar que la categoría sea una de las permitidas
	categoriasValidas := []string{"Clasica", "Premium", "Vegetariana"}
	categoriaValida := false
	for _, cat := range categoriasValidas {
		if categoria == cat {
			categoriaValida = true
			break
		}
	}
	if !categoriaValida {
		return nil, errors.New("categoría inválida. Usa: Clasica, Premium o Vegetariana")
	}

	return &Hamburguesa{
		id:           generarID(),
		Nombre:       nombre,
		Descripcion:  descripcion,
		Precio:       precio,
		Categoria:    categoria,
		Ingredientes: ingredientes,
		Disponible:   true,
		fechaCreado:  time.Now(),
	}, nil
}

// Contador global para generación de IDs temporales
var contadorID = 0

// generarID genera un ID único con formato BURG-001, BURG-002, etc.
func generarID() string {
	contadorID++
	return "BURG-" + fmt.Sprintf("%03d", contadorID)
}

// MarshalJSON serializa la hamburguesa a JSON incluyendo el campo privado id.
// Sin este método, json.Marshal ignoraría el campo id por ser privado.
func (h *Hamburguesa) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID           string   `json:"id"`
		Nombre       string   `json:"nombre"`
		Descripcion  string   `json:"descripcion"`
		Precio       float64  `json:"precio"`
		Categoria    string   `json:"categoria"`
		Ingredientes []string `json:"ingredientes"`
		Disponible   bool     `json:"disponible"`
	}{
		ID:           h.id,
		Nombre:       h.Nombre,
		Descripcion:  h.Descripcion,
		Precio:       h.Precio,
		Categoria:    h.Categoria,
		Ingredientes: h.Ingredientes,
		Disponible:   h.Disponible,
	})
}

// GetID retorna el ID de la hamburguesa (solo lectura)
func (h *Hamburguesa) GetID() string {
	return h.id
}

// SetID establece el ID (usado al cargar desde BD)
func (h *Hamburguesa) SetID(id string) {
	h.id = id
}

// GetFechaCreado retorna la fecha de creación (solo lectura)
func (h *Hamburguesa) GetFechaCreado() time.Time {
	return h.fechaCreado
}

// ActualizarPrecio cambia el precio validando que sea positivo
func (h *Hamburguesa) ActualizarPrecio(nuevoPrecio float64) error {
	if nuevoPrecio <= 0 {
		return errors.New("el precio debe ser mayor a 0")
	}
	h.Precio = nuevoPrecio
	return nil
}

// MarcarDisponibilidad cambia el estado de disponibilidad
func (h *Hamburguesa) MarcarDisponibilidad(disponible bool) {
	h.Disponible = disponible
}

// ObtenerInfo retorna información formateada de la hamburguesa.
// Implementa la interface Informable.
func (h *Hamburguesa) ObtenerInfo() string {
	estado := "Disponible"
	if !h.Disponible {
		estado = "No disponible"
	}

	return fmt.Sprintf(
		"[%s] %s - $%.2f (%s)\n  %s\n  Ingredientes: %v",
		h.id, h.Nombre, h.Precio, estado, h.Descripcion, h.Ingredientes,
	)
}

// Validar verifica que la hamburguesa tenga datos válidos.
// Implementa la interface Validable.
func (h *Hamburguesa) Validar() error {
	if h.Nombre == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if h.Precio <= 0 {
		return errors.New("el precio debe ser mayor a 0")
	}
	return nil
}
