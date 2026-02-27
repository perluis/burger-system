// Package orders define el modelo de datos para las órdenes del restaurante.
// Maneja la creación de órdenes, agregado de items y cálculo de totales con IVA.
package orders

import (
	"errors"
	"fmt"
	"time"
)

// Orden representa un pedido del cliente con sus items y totales
type Orden struct {
	ID         string
	NumeroMesa int
	TipoOrden  string // "mesa", "delivery", "para_llevar"
	Items      []ItemOrden
	Estado     string // "PENDIENTE", "EN_COCINA", "LISTA", "ENTREGADA", "CANCELADA"
	Subtotal   float64
	Total      float64
	FechaHora  time.Time
}

// ItemOrden representa un producto dentro de una orden
type ItemOrden struct {
	HamburguesaID string
	Nombre        string
	Cantidad      int
	PrecioUnit    float64
	Notas         string // Ej: "sin cebolla", "extra queso"
}

// Contador para IDs temporales de orden
var contadorOrden = 0

// NuevaOrden crea una orden validando el tipo
func NuevaOrden(tipoOrden string, numeroMesa int) (*Orden, error) {
	tiposValidos := []string{"mesa", "delivery", "para_llevar"}
	tipoValido := false
	for _, tipo := range tiposValidos {
		if tipoOrden == tipo {
			tipoValido = true
			break
		}
	}
	if !tipoValido {
		return nil, errors.New("tipo de orden inválido")
	}

	contadorOrden++
	return &Orden{
		ID:         fmt.Sprintf("ORD-%03d", contadorOrden),
		NumeroMesa: numeroMesa,
		TipoOrden:  tipoOrden,
		Items:      []ItemOrden{},
		Estado:     "PENDIENTE",
		Subtotal:   0,
		Total:      0,
		FechaHora:  time.Now(),
	}, nil
}

// AgregarItem agrega un producto a la orden y recalcula totales
func (o *Orden) AgregarItem(hamburguesaID, nombre string, cantidad int, precioUnit float64, notas string) error {
	if cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a 0")
	}
	if precioUnit <= 0 {
		return errors.New("el precio debe ser mayor a 0")
	}

	o.Items = append(o.Items, ItemOrden{
		HamburguesaID: hamburguesaID,
		Nombre:        nombre,
		Cantidad:      cantidad,
		PrecioUnit:    precioUnit,
		Notas:         notas,
	})
	o.calcularTotales()
	return nil
}

// calcularTotales suma los items y aplica IVA 15% (Ecuador)
func (o *Orden) calcularTotales() {
	var subtotal float64
	for _, item := range o.Items {
		subtotal += float64(item.Cantidad) * item.PrecioUnit
	}
	o.Subtotal = subtotal
	o.Total = subtotal * 1.15 // IVA 15% vigente en Ecuador
}

// ActualizarEstado cambia el estado validando que sea uno permitido
func (o *Orden) ActualizarEstado(nuevoEstado string) error {
	estadosValidos := []string{"PENDIENTE", "EN_COCINA", "LISTA", "ENTREGADA", "CANCELADA"}
	estadoValido := false
	for _, estado := range estadosValidos {
		if nuevoEstado == estado {
			estadoValido = true
			break
		}
	}
	if !estadoValido {
		return errors.New("estado inválido")
	}
	o.Estado = nuevoEstado
	return nil
}

// ObtenerInfo retorna información formateada de la orden.
// Implementa la interface Informable.
func (o *Orden) ObtenerInfo() string {
	info := fmt.Sprintf("=== ORDEN %s ===\n", o.ID)
	info += fmt.Sprintf("Tipo: %s\n", o.TipoOrden)

	if o.TipoOrden == "mesa" {
		info += fmt.Sprintf("Mesa: %d\n", o.NumeroMesa)
	}

	info += fmt.Sprintf("Estado: %s\n", o.Estado)
	info += fmt.Sprintf("Fecha: %s\n\n", o.FechaHora.Format("02/01/2006 15:04"))

	info += "Items:\n"
	for i, item := range o.Items {
		info += fmt.Sprintf("  %d. %dx %s - $%.2f",
			i+1, item.Cantidad, item.Nombre, item.PrecioUnit*float64(item.Cantidad))
		if item.Notas != "" {
			info += fmt.Sprintf(" (%s)", item.Notas)
		}
		info += "\n"
	}

	info += fmt.Sprintf("\nSubtotal: $%.2f\n", o.Subtotal)
	info += fmt.Sprintf("IVA (15%%): $%.2f\n", o.Total-o.Subtotal)
	info += fmt.Sprintf("TOTAL: $%.2f\n", o.Total)

	return info
}
