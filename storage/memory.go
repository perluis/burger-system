package storage

import (
	"sync"

	"github.com/perluis/burger-system/menu"
)

// Store maneja el almacenamiento en memoria
type Store struct {
	hamburguesas []*menu.Hamburguesa
	mu           sync.RWMutex // Para acceso concurrente seguro
}

// NewStore crea un nuevo store con datos iniciales
func NewStore() *Store {
	store := &Store{
		hamburguesas: []*menu.Hamburguesa{},
	}

	// Cargar datos iniciales
	store.cargarDatosIniciales()

	return store
}

// cargarDatosIniciales crea hamburguesas de ejemplo
func (s *Store) cargarDatosIniciales() {
	ingredientes1 := []string{"carne", "queso", "lechuga", "tomate", "pan"}
	h1, _ := menu.NuevaHamburguesa(
		"Hamburguesa Clásica",
		"Deliciosa hamburguesa tradicional",
		8.99,
		"Clasica",
		ingredientes1,
	)
	s.hamburguesas = append(s.hamburguesas, h1)

	ingredientes2 := []string{"carne", "queso cheddar", "tocino", "bbq", "cebolla", "pan"}
	h2, _ := menu.NuevaHamburguesa(
		"BBQ Premium",
		"Con tocino y salsa BBQ casera",
		12.99,
		"Premium",
		ingredientes2,
	)
	s.hamburguesas = append(s.hamburguesas, h2)

	ingredientes3 := []string{"portobello", "queso", "lechuga", "tomate", "aguacate", "pan integral"}
	h3, _ := menu.NuevaHamburguesa(
		"Veggie Deluxe",
		"Hamburguesa vegetariana gourmet",
		10.99,
		"Vegetariana",
		ingredientes3,
	)
	s.hamburguesas = append(s.hamburguesas, h3)
}

// GetHamburguesas retorna todas las hamburguesas
func (s *Store) GetHamburguesas() []*menu.Hamburguesa {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hamburguesas
}

// GetHamburguesaByID busca una hamburguesa por ID
func (s *Store) GetHamburguesaByID(id string) *menu.Hamburguesa {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, h := range s.hamburguesas {
		if h.GetID() == id {
			return h
		}
	}
	return nil
}
