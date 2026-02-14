package utils

import "fmt"

// MostrarInformacion muestra la información de cualquier entidad Informable
// Demuestra polimorfismo: puede recibir Hamburguesa, Orden, etc.
func MostrarInformacion(entidad Informable) {
	fmt.Println(entidad.ObtenerInfo())
}

// MostrarID muestra el ID de cualquier entidad Identificable
func MostrarID(entidad Identificable) {
	fmt.Printf("ID: %s\n", entidad.GetID())
}

// ValidarEntidad valida cualquier entidad que implemente Validable
// Retorna true si es válida, false si no
func ValidarEntidad(entidad Validable) bool {
	err := entidad.Validar()
	if err != nil {
		fmt.Printf("❌ Validación fallida: %s\n", err)
		return false
	}
	fmt.Println("✅ Validación exitosa")
	return true
}
