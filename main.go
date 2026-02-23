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

	// Ruta raíz (página de bienvenida)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🍔 API Restaurante de Hamburguesas\n\n")
		fmt.Fprintf(w, "Endpoints disponibles:\n")
		fmt.Fprintf(w, "GET /api/hamburguesas - Listar todas las hamburguesas\n")
		fmt.Fprintf(w, "GET /api/hamburguesas/{id} - Obtener hamburguesa por ID\n")
	})

	// Iniciar servidor
	port := "8080"
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%s\n", port)
	fmt.Println("📋 Prueba: http://localhost:8080/api/hamburguesas")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
