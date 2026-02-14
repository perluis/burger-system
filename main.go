package main

import (
	"fmt"

	"github.com/perluis/burger-system/menu"
)

func main() {
	fmt.Println("=== SISTEMA DE GESTIÓN - RESTAURANTE ===")
	fmt.Println()

	ingredientes := []string{"carne", "queso", "lechuga", "tomate", "pan"}
	hamburguesa, err := menu.NuevaHamburguesa(
		"Hamburguesa Clásica",
		"Deliciosa hamburguesa con ingredientes frescos",
		8.99,
		"Clasica",
		ingredientes,
	)

	if err != nil {
		fmt.Println("Error al crear hamburguesa:", err)
		return
	}

	fmt.Println("✅ Hamburguesa creada:")
	fmt.Println(hamburguesa.ObtenerInfo())
	fmt.Println()

	fmt.Println("📝 Actualizando precio a $9.99...")
	err = hamburguesa.ActualizarPrecio(9.99)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("✅ Precio actualizado")
	}
	fmt.Println()

	fmt.Println("❌ Marcando como no disponible...")
	hamburguesa.MarcarDisponibilidad(false)
	fmt.Println()

	fmt.Println("📋 Información actualizada:")
	fmt.Println(hamburguesa.ObtenerInfo())
}
