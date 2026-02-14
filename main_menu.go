package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/perluis/burger-system/menu"
	"github.com/perluis/burger-system/orders"
	"github.com/perluis/burger-system/utils"
)

var (
	hamburguesas []*menu.Hamburguesa
	ordenActual  *orders.Orden
	scanner      = bufio.NewScanner(os.Stdin)
)

func main() {
	inicializarDatos()
	mostrarMenuPrincipal()
}

func inicializarDatos() {
	// Crear algunas hamburguesas por defecto
	ingredientes1 := []string{"carne", "queso", "lechuga", "tomate", "pan"}
	h1, _ := menu.NuevaHamburguesa("Hamburguesa Clásica", "Deliciosa hamburguesa tradicional", 8.99, "Clasica", ingredientes1)
	hamburguesas = append(hamburguesas, h1)

	ingredientes2 := []string{"carne", "queso cheddar", "tocino", "bbq", "cebolla", "pan"}
	h2, _ := menu.NuevaHamburguesa("BBQ Premium", "Con tocino y salsa BBQ casera", 12.99, "Premium", ingredientes2)
	hamburguesas = append(hamburguesas, h2)

	ingredientes3 := []string{"portobello", "queso", "lechuga", "tomate", "aguacate", "pan integral"}
	h3, _ := menu.NuevaHamburguesa("Veggie Deluxe", "Hamburguesa vegetariana gourmet", 10.99, "Vegetariana", ingredientes3)
	hamburguesas = append(hamburguesas, h3)
}

func mostrarMenuPrincipal() {
	for {
		fmt.Println()
		fmt.Println("╔════════════════════════════════════════════╗")
		fmt.Println("║    SISTEMA DE GESTIÓN - RESTAURANTE       ║")
		fmt.Println("╚════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("1. Ver Menú de Hamburguesas")
		fmt.Println("2. Crear Nueva Orden")
		fmt.Println("3. Ver Orden Actual")
		fmt.Println("4. Finalizar Orden")
		fmt.Println("5. Demostración de Interfaces")
		fmt.Println("0. Salir")
		fmt.Println()
		fmt.Print("Seleccione una opción: ")

		opcion := leerLinea()

		switch opcion {
		case "1":
			verMenu()
		case "2":
			crearOrden()
		case "3":
			verOrdenActual()
		case "4":
			finalizarOrden()
		case "5":
			demoInterfaces()
		case "0":
			fmt.Println("\n¡Gracias por usar el sistema! 👋")
			return
		default:
			fmt.Println("\n❌ Opción inválida")
		}
	}
}

func verMenu() {
	fmt.Println("\n╔════════════════════════════════════════════╗")
	fmt.Println("║           MENÚ DE HAMBURGUESAS            ║")
	fmt.Println("╚════════════════════════════════════════════╝\n")

	for i, h := range hamburguesas {
		fmt.Printf("%d. ", i+1)
		utils.MostrarInformacion(h)
		fmt.Println()
	}

	fmt.Println("\nPresione Enter para continuar...")
	leerLinea()
}

func crearOrden() {
	fmt.Println("\n╔════════════════════════════════════════════╗")
	fmt.Println("║           CREAR NUEVA ORDEN               ║")
	fmt.Println("╚════════════════════════════════════════════╝\n")

	fmt.Println("Tipo de orden:")
	fmt.Println("1. Mesa")
	fmt.Println("2. Delivery")
	fmt.Println("3. Para llevar")
	fmt.Print("\nSeleccione: ")

	tipoOpcion := leerLinea()
	var tipoOrden string
	var numeroMesa int

	switch tipoOpcion {
	case "1":
		tipoOrden = "mesa"
		fmt.Print("Número de mesa: ")
		mesaStr := leerLinea()
		numeroMesa, _ = strconv.Atoi(mesaStr)
	case "2":
		tipoOrden = "delivery"
	case "3":
		tipoOrden = "para_llevar"
	default:
		fmt.Println("❌ Opción inválida")
		return
	}

	orden, err := orders.NuevaOrden(tipoOrden, numeroMesa)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	ordenActual = orden
	fmt.Printf("\n✅ Orden %s creada exitosamente\n", orden.ID)

	agregarItems()
}

func agregarItems() {
	for {
		fmt.Println("\n--- Agregar Items ---")
		fmt.Println("Hamburguesas disponibles:")

		for i, h := range hamburguesas {
			fmt.Printf("%d. %s - $%.2f\n", i+1, h.Nombre, h.Precio)
		}

		fmt.Println("0. Terminar de agregar")
		fmt.Print("\nSeleccione hamburguesa: ")

		opcion := leerLinea()

		if opcion == "0" {
			break
		}

		idx, err := strconv.Atoi(opcion)
		if err != nil || idx < 1 || idx > len(hamburguesas) {
			fmt.Println("❌ Opción inválida")
			continue
		}

		h := hamburguesas[idx-1]

		fmt.Print("Cantidad: ")
		cantidadStr := leerLinea()
		cantidad, _ := strconv.Atoi(cantidadStr)

		fmt.Print("Notas especiales (Enter para ninguna): ")
		notas := leerLinea()

		err = ordenActual.AgregarItem(h.GetID(), h.Nombre, cantidad, h.Precio, notas)
		if err != nil {
			fmt.Println("❌ Error:", err)
		} else {
			fmt.Println("✅ Item agregado")
		}
	}
}

func verOrdenActual() {
	if ordenActual == nil {
		fmt.Println("\n❌ No hay orden activa. Cree una orden primero.")
		return
	}

	fmt.Println()
	fmt.Println(ordenActual.ObtenerInfo())
	fmt.Println("\nPresione Enter para continuar...")
	leerLinea()
}

func finalizarOrden() {
	if ordenActual == nil {
		fmt.Println("\n❌ No hay orden activa.")
		return
	}

	fmt.Println("\n╔════════════════════════════════════════════╗")
	fmt.Println("║           FINALIZAR ORDEN                 ║")
	fmt.Println("╚════════════════════════════════════════════╝")

	ordenActual.ActualizarEstado("EN_COCINA")
	fmt.Println("\n✅ Orden enviada a cocina")

	ordenActual.ActualizarEstado("LISTA")
	fmt.Println("✅ Orden lista para entregar")

	ordenActual.ActualizarEstado("ENTREGADA")
	fmt.Println("✅ Orden entregada")

	fmt.Println("\n" + ordenActual.ObtenerInfo())

	ordenActual = nil
	fmt.Println("\nPresione Enter para continuar...")
	leerLinea()
}

func demoInterfaces() {
	fmt.Println("\n╔════════════════════════════════════════════╗")
	fmt.Println("║      DEMOSTRACIÓN DE INTERFACES           ║")
	fmt.Println("╚════════════════════════════════════════════╝\n")

	if len(hamburguesas) > 0 {
		h := hamburguesas[0]
		fmt.Println("🔍 Interface Identificable:")
		utils.MostrarID(h)
		fmt.Println()

		fmt.Println("✓ Interface Validable:")
		utils.ValidarEntidad(h)
		fmt.Println()

		fmt.Println("📄 Interface Informable:")
		utils.MostrarInformacion(h)
	}

	fmt.Println("\nPresione Enter para continuar...")
	leerLinea()
}

func leerLinea() string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
