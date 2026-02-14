package utils

import "time"

// Identificable representa cualquier entidad que tiene un ID único
type Identificable interface {
	GetID() string
}

// Actualizable representa entidades que pueden ser actualizadas
type Actualizable interface {
	GetFechaCreado() time.Time
}

// Informable representa entidades que pueden mostrar su información
type Informable interface {
	ObtenerInfo() string
}

// Validable representa entidades que pueden ser validadas
type Validable interface {
	Validar() error
}
