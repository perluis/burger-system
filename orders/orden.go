package orders

import (
	"errors"
	"fmt"
	"time"
)

// Orden representa un pedido del cliente
type Orden struct {
	ID         string
	NumeroMesa int
	TipoOrden  string // "mesa", "delivery", "para_llevar"
	Items      []ItemOrden
	Estado     string // "PENDIENTE", "EN_COCINA", "LISTA", "ENTREGADA"
	Subtotal   float64
	Total      float64
	FechaHora  time.Time
}

// ItemOrden representa un item dentro de una orden
type ItemOrden struct {
	HamburguesaID string
	Nombre        string
	Cantidad      int
	PrecioUnit    float64
	Notas         string // "sin cebolla", "extra queso", etc.
}

// Contador para IDs de orden
var contadorOrden = 0

// NuevaOrden crea una nueva orden
func NuevaOrden(tipoOrden string, numeroMesa int) (*Orden, error) {
	// Validar tipo de orden
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
	id := fmt.Sprintf("ORD-%03d", contadorOrden)

	return &Orden{
		ID:         id,
		NumeroMesa: numeroMesa,
		TipoOrden:  tipoOrden,
		Items:      []ItemOrden{},
		Estado:     "PENDIENTE",
		Subtotal:   0,
		Total:      0,
		FechaHora:  time.Now(),
	}, nil
}

// AgregarItem agrega un item a la orden
func (o *Orden) AgregarItem(hamburguesaID, nombre string, cantidad int, precioUnit float64, notas string) error {
	if cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a 0")
	}

	if precioUnit <= 0 {
		return errors.New("el precio debe ser mayor a 0")
	}

	item := ItemOrden{
		HamburguesaID: hamburguesaID,
		Nombre:        nombre,
		Cantidad:      cantidad,
		PrecioUnit:    precioUnit,
		Notas:         notas,
	}

	o.Items = append(o.Items, item)
	o.calcularTotales()

	return nil
}

// calcularTotales calcula el subtotal y total con IVA
func (o *Orden) calcularTotales() {
	var subtotal float64

	for _, item := range o.Items {
		subtotal += float64(item.Cantidad) * item.PrecioUnit
	}

	o.Subtotal = subtotal
	// IVA 12% (Ecuador)
	o.Total = subtotal * 1.12
}

// ActualizarEstado cambia el estado de la orden
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

// ObtenerInfo retorna información de la orden
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
	info += fmt.Sprintf("IVA (12%%): $%.2f\n", o.Total-o.Subtotal)
	info += fmt.Sprintf("TOTAL: $%.2f\n", o.Total)

	return info
}
