// Package utils define interfaces que permiten polimorfismo en el sistema.
// Tanto Hamburguesa como Orden implementan estas interfaces,
// lo que permite tratarlas de forma genérica.
package utils

import "time"

// Identificable: cualquier entidad con un ID único
type Identificable interface {
	GetID() string
}

// Actualizable: entidades que registran su fecha de creación
type Actualizable interface {
	GetFechaCreado() time.Time
}

// Informable: entidades que pueden mostrar su información como texto
type Informable interface {
	ObtenerInfo() string
}

// Validable: entidades que pueden validar sus propios datos
type Validable interface {
	Validar() error
}
