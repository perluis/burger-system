// Punto de entrada del sistema Burger System.
// Configura la conexión a MySQL, define las rutas de la API REST
// y levanta el servidor HTTP con soporte CORS.
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

// corsHandler intercepta las peticiones HTTP antes de que lleguen al router.
// Agrega los headers CORS necesarios para que el frontend pueda
// comunicarse con la API sin restricciones del navegador.
type corsHandler struct {
	router *mux.Router
}

func (c *corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Responder directamente a peticiones preflight (OPTIONS)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	c.router.ServeHTTP(w, r)
}

func main() {
	// Conectar a la base de datos MySQL
	db, err := storage.NewDatabase()
	if err != nil {
		log.Fatal("Error conectando a MySQL: ", err)
	}
	defer db.Close()

	// Crear el servidor con la conexión a BD
	server := api.NewServer(db)
	router := mux.NewRouter()

	// Rutas API - Hamburguesas (CRUD completo)
	router.HandleFunc("/api/hamburguesas", server.GetHamburguesas).Methods("GET")
	router.HandleFunc("/api/hamburguesas/{id}", server.GetHamburguesaByID).Methods("GET")
	router.HandleFunc("/api/hamburguesas", server.CreateHamburguesa).Methods("POST")
	router.HandleFunc("/api/hamburguesas/{id}", server.UpdateHamburguesa).Methods("PUT")
	router.HandleFunc("/api/hamburguesas/{id}", server.DeleteHamburguesa).Methods("DELETE")

	// Rutas API - Órdenes (CRUD + gestión de estado)
	router.HandleFunc("/api/ordenes", server.GetAllOrdenes).Methods("GET")
	router.HandleFunc("/api/ordenes", server.CreateOrden).Methods("POST")
	router.HandleFunc("/api/ordenes/{id}", server.GetOrdenByID).Methods("GET")
	router.HandleFunc("/api/ordenes/{id}/estado", server.UpdateOrdenEstado).Methods("PUT")
	router.HandleFunc("/api/ordenes/{id}", server.DeleteOrden).Methods("DELETE")

	// Página principal - sirve el HTML del frontend
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Error cargando página", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Archivos estáticos (imágenes, CSS)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Iniciar el servidor con el wrapper CORS
	port := "8080"
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║    🍔 BURGER SYSTEM - SERVIDOR ACTIVO      ║")
	fmt.Println("╠════════════════════════════════════════════╣")
	fmt.Println("║  ✅ MySQL conectado                        ║")
	fmt.Printf("║  🌐 Web: http://localhost:%s              ║\n", port)
	fmt.Printf("║  📡 API: http://localhost:%s/api          ║\n", port)
	fmt.Println("╚════════════════════════════════════════════╝")

	log.Fatal(http.ListenAndServe(":"+port, &corsHandler{router: router}))
}
