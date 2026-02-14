package menu

import (
	"errors"
	"fmt"
	"time"
)

// Hamburguesa representa un producto del menú del restaurante.
// Incluye toda la información necesaria para gestión de inventario,
// pricing y disponibilidad en tiempo real.
type Hamburguesa struct {
	ID           string // Formato: BURG-001, BURG-002, etc.
	Nombre       string
	Descripcion  string
	Precio       float64 // Precio en USD
	Categoria    string  // Valores permitidos: "Clasica", "Premium", "Vegetariana"
	Ingredientes []string
	Disponible   bool // false cuando está temporalmente agotada
	FechaCreado  time.Time
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
		ID:           id,
		Nombre:       nombre,
		Descripcion:  descripcion,
		Precio:       precio,
		Categoria:    categoria,
		Ingredientes: ingredientes,
		Disponible:   true,
		FechaCreado:  time.Now(),
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
		h.ID,
		h.Nombre,
		h.Precio,
		estado,
		h.Descripcion,
		h.Ingredientes,
	)
}
