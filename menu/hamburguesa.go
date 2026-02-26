package menu

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Hamburguesa representa un producto del menú del restaurante.
// Utiliza encapsulación para proteger campos críticos como ID y FechaCreado.
type Hamburguesa struct {
	id           string    // Privado - no se puede cambiar desde fuera
	Nombre       string    // Público
	Descripcion  string    // Público
	Precio       float64   // Público - se modifica con ActualizarPrecio()
	Categoria    string    // Público
	Ingredientes []string  // Público
	Disponible   bool      // Público - se modifica con MarcarDisponibilidad()
	fechaCreado  time.Time // Privado - solo lectura
}

// NuevaHamburguesa crea y retorna una nueva hamburguesa con ID autogenerado.
// Valida que el precio sea mayor a 0 y que la categoría sea válida.
// Retorna error si los datos son inválidos.
func NuevaHamburguesa(nombre, descripcion string, precio float64, categoria string, ingredientes []string) (*Hamburguesa, error) {
	// Validar precio positivo
	if precio <= 0 {
		return nil, errors.New("el precio debe ser mayor a 0")
	}

	// Validar categoría
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

	// Generar ID único
	id := generarID()

	// Crear y retornar la hamburguesa
	return &Hamburguesa{
		id:           id, // Minúscula
		Nombre:       nombre,
		Descripcion:  descripcion,
		Precio:       precio,
		Categoria:    categoria,
		Ingredientes: ingredientes,
		Disponible:   true,
		fechaCreado:  time.Now(), // Minúscula
	}, nil

}

// Contador global para generación de IDs
var contadorID = 0

// generarID genera un ID único para hamburguesas.
// Formato: BURG-001, BURG-002, etc.
func generarID() string {
	contadorID++
	return "BURG-" + fmt.Sprintf("%03d", contadorID)
}

// GetID retorna el ID de la hamburguesa (solo lectura)
func (h *Hamburguesa) GetID() string {
	return h.id
}

// SetID establece el ID de la hamburguesa (usado para cargar desde BD)
func (h *Hamburguesa) SetID(id string) {
	h.id = id
}

// GetFechaCreado retorna la fecha de creación (solo lectura)
func (h *Hamburguesa) GetFechaCreado() time.Time {
	return h.fechaCreado
}

// ActualizarPrecio cambia el precio de la hamburguesa.
// Valida que el nuevo precio sea mayor a 0.
func (h *Hamburguesa) ActualizarPrecio(nuevoPrecio float64) error {
	if nuevoPrecio <= 0 {
		return errors.New("el precio debe ser mayor a 0")
	}
	h.Precio = nuevoPrecio
	return nil
}

// MarcarDisponibilidad cambia el estado de disponibilidad.
func (h *Hamburguesa) MarcarDisponibilidad(disponible bool) {
	h.Disponible = disponible
}

// ObtenerInfo retorna una representación en texto de la hamburguesa.
func (h *Hamburguesa) ObtenerInfo() string {
	estado := "Disponible"
	if !h.Disponible {
		estado = "No disponible"
	}

	return fmt.Sprintf(
		"[%s] %s - $%.2f (%s)\n  %s\n  Ingredientes: %v",
		h.id,
		h.Nombre,
		h.Precio,
		estado,
		h.Descripcion,
		h.Ingredientes,
	)
}

// Validar verifica que la hamburguesa tenga datos válidos
// Implementa la interface Validable
func (h *Hamburguesa) Validar() error {
	if h.Nombre == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if h.Precio <= 0 {
		return errors.New("el precio debe ser mayor a 0")
	}
	return nil
}

// MarshalJSON personaliza la serialización JSON para incluir el ID privado.
// Esto es necesario porque el campo 'id' es privado (minúscula) y Go
// no lo incluye automáticamente en la serialización JSON.
// Sin este método, el frontend no recibiría los IDs de las hamburguesas.
func (h *Hamburguesa) MarshalJSON() ([]byte, error) {
	// Estructura auxiliar que define cómo se verá el JSON de salida
	type HamburguesaJSON struct {
		ID           string   `json:"id"`
		Nombre       string   `json:"nombre"`
		Descripcion  string   `json:"descripcion"`
		Precio       float64  `json:"precio"`
		Categoria    string   `json:"categoria"`
		Ingredientes []string `json:"ingredientes"`
		Disponible   bool     `json:"disponible"`
	}

	// Convertir la hamburguesa a la estructura JSON y serializar
	return json.Marshal(HamburguesaJSON{
		ID:           h.id,
		Nombre:       h.Nombre,
		Descripcion:  h.Descripcion,
		Precio:       h.Precio,
		Categoria:    h.Categoria,
		Ingredientes: h.Ingredientes,
		Disponible:   h.Disponible,
	})
}
