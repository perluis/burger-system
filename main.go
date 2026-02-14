package main

import (
	"fmt"

	"github.com/perluis/burger-system/menu"
	"github.com/perluis/burger-system/orders"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║  SISTEMA DE GESTIÓN - RESTAURANTE         ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()

	// ========== CREAR HAMBURGUESAS ==========
	fmt.Println("📋 CREANDO MENÚ...")
	fmt.Println()

	ingredientes1 := []string{"carne", "queso", "lechuga", "tomate", "pan"}
	clasica, err := menu.NuevaHamburguesa(
		"Hamburguesa Clásica",
		"Deliciosa hamburguesa con ingredientes frescos",
		8.99,
		"Clasica",
		ingredientes1,
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("✅", clasica.ObtenerInfo())
	fmt.Println()

	ingredientes2 := []string{"carne", "queso cheddar", "tocino", "bbq", "cebolla caramelizada", "pan"}
	premium, err := menu.NuevaHamburguesa(
		"Hamburguesa BBQ Premium",
		"Con tocino ahumado y salsa BBQ casera",
		12.99,
		"Premium",
		ingredientes2,
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("✅", premium.ObtenerInfo())
	fmt.Println()

	// ========== CREAR ORDEN ==========
	fmt.Println("🛒 CREANDO ORDEN...")
	fmt.Println()

	orden, err := orders.NuevaOrden("mesa", 5)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("✅ Orden creada: %s\n", orden.ID)
	fmt.Println()

	// ========== AGREGAR ITEMS ==========
	fmt.Println("➕ AGREGANDO ITEMS A LA ORDEN...")
	fmt.Println()

	err = orden.AgregarItem(clasica.ID, clasica.Nombre, 2, clasica.Precio, "sin cebolla")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("✅ Agregado: 2x Hamburguesa Clásica")

	err = orden.AgregarItem(premium.ID, premium.Nombre, 1, premium.Precio, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("✅ Agregado: 1x Hamburguesa BBQ Premium")
	fmt.Println()

	// ========== MOSTRAR ORDEN ==========
	fmt.Println(orden.ObtenerInfo())
	fmt.Println()

	// ========== ACTUALIZAR ESTADO ==========
	fmt.Println("🔄 ACTUALIZANDO ESTADO DE LA ORDEN...")
	fmt.Println()

	err = orden.ActualizarEstado("EN_COCINA")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("✅ Estado actualizado a: %s\n", orden.Estado)
	fmt.Println()

	err = orden.ActualizarEstado("LISTA")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("✅ Estado actualizado a: %s\n", orden.Estado)
	fmt.Println()

	err = orden.ActualizarEstado("ENTREGADA")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("✅ Estado actualizado a: %s\n", orden.Estado)
	fmt.Println()

	// ========== RESUMEN FINAL ==========
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║  RESUMEN FINAL                            ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println(orden.ObtenerInfo())
}
