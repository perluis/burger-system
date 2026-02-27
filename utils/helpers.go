// Package utils proporciona funciones utilitarias que operan sobre interfaces.
// Demuestran polimorfismo: una misma función puede recibir cualquier tipo
// que implemente la interface requerida (Hamburguesa, Orden, etc).
package utils

import "fmt"

// MostrarInformacion imprime la info de cualquier entidad Informable
func MostrarInformacion(entidad Informable) {
	fmt.Println(entidad.ObtenerInfo())
}

// MostrarID imprime el ID de cualquier entidad Identificable
func MostrarID(entidad Identificable) {
	fmt.Printf("ID: %s\n", entidad.GetID())
}

// ValidarEntidad valida cualquier entidad Validable y retorna el resultado
func ValidarEntidad(entidad Validable) bool {
	err := entidad.Validar()
	if err != nil {
		fmt.Printf("❌ Validación fallida: %s\n", err)
		return false
	}
	fmt.Println("✅ Validación exitosa")
	return true
}
