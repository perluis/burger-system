package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/perluis/burger-system/api"
	"github.com/perluis/burger-system/storage"
)

func main() {
	// Conectar a MySQL
	db, err := storage.NewDatabase()
	if err != nil {
		log.Fatal("❌ Error conectando a MySQL: ", err)
	}
	defer db.Close()

	// Crear servidor con la base de datos
	server := api.NewServer(db)

	// Crear router
	router := mux.NewRouter()

	// ==================== RUTAS API ====================

	// Rutas de hamburguesas
	router.HandleFunc("/api/hamburguesas", server.GetHamburguesas).Methods("GET")
	router.HandleFunc("/api/hamburguesas/{id}", server.GetHamburguesaByID).Methods("GET")
	router.HandleFunc("/api/hamburguesas", server.CreateHamburguesa).Methods("POST")
	router.HandleFunc("/api/hamburguesas/{id}", server.UpdateHamburguesa).Methods("PUT")
	router.HandleFunc("/api/hamburguesas/{id}", server.DeleteHamburguesa).Methods("DELETE")

	// Rutas de órdenes
	router.HandleFunc("/api/ordenes", server.CreateOrden).Methods("POST")
	router.HandleFunc("/api/ordenes/{id}", server.GetOrdenByID).Methods("GET")
	router.HandleFunc("/api/ordenes/{id}/estado", server.UpdateOrdenEstado).Methods("PUT")

	// ==================== RUTAS WEB ====================

	// Página principal
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Error cargando página", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Archivos estáticos (CSS, JS)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Iniciar servidor
	port := "8080"
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║    🍔 BURGER SYSTEM - SERVIDOR ACTIVO      ║")
	fmt.Println("╠════════════════════════════════════════════╣")
	fmt.Println("║  ✅ MySQL conectado                        ║")
	fmt.Printf("║  🌐 Web: http://localhost:%s              ║\n", port)
	fmt.Printf("║  📡 API: http://localhost:%s/api          ║\n", port)
	fmt.Println("╚════════════════════════════════════════════╝")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
