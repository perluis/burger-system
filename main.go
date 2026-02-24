package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/perluis/burger-system/api"
	"github.com/perluis/burger-system/storage"
)

func main() {
	// Crear el store con datos iniciales
	store := storage.NewStore()

	// Crear el servidor con el store
	server := api.NewServer(store)

	// Crear el router
	router := mux.NewRouter()

	// Definir rutas
	router.HandleFunc("/api/hamburguesas", server.GetHamburguesas).Methods("GET")
	router.HandleFunc("/api/hamburguesas/{id}", server.GetHamburguesaByID).Methods("GET")
	router.HandleFunc("/api/hamburguesas", server.CreateHamburguesa).Methods("POST")
	router.HandleFunc("/api/hamburguesas/{id}", server.UpdateHamburguesa).Methods("PUT")
	router.HandleFunc("/api/hamburguesas/{id}", server.DeleteHamburguesa).Methods("DELETE")
	router.HandleFunc("/api/ordenes", server.CreateOrden).Methods("POST")
	router.HandleFunc("/api/ordenes/{id}", server.GetOrdenByID).Methods("GET")
	router.HandleFunc("/api/ordenes/{id}/estado", server.UpdateOrdenEstado).Methods("PUT")

	// Ruta raíz (página de bienvenida)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🍔 API Restaurante de Hamburguesas\n\n")
		fmt.Fprintf(w, "=== ENDPOINTS HAMBURGUESAS ===\n")
		fmt.Fprintf(w, "GET    /api/hamburguesas        - Listar todas\n")
		fmt.Fprintf(w, "GET    /api/hamburguesas/{id}   - Obtener por ID\n")
		fmt.Fprintf(w, "POST   /api/hamburguesas        - Crear nueva\n")
		fmt.Fprintf(w, "PUT    /api/hamburguesas/{id}   - Actualizar\n")
		fmt.Fprintf(w, "DELETE /api/hamburguesas/{id}   - Eliminar\n\n")
		fmt.Fprintf(w, "=== ENDPOINTS ÓRDENES ===\n")
		fmt.Fprintf(w, "POST   /api/ordenes             - Crear orden\n")
		fmt.Fprintf(w, "GET    /api/ordenes/{id}        - Obtener orden\n")
		fmt.Fprintf(w, "PUT    /api/ordenes/{id}/estado - Actualizar estado\n")
	})

	// Iniciar servidor
	port := "8080"
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%s\n", port)
	fmt.Println("📋 Prueba: http://localhost:8080/api/hamburguesas")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
