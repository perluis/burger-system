package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/perluis/burger-system/menu"
	"github.com/perluis/burger-system/orders"
)

// Database maneja la conexión a MySQL
type Database struct {
	conn *sql.DB
}

// NewDatabase crea una nueva conexión a la base de datos
func NewDatabase() (*Database, error) {
	// Conexión: usuario:contraseña@tcp(host:puerto)/base_de_datos
	dsn := "root:@tcp(localhost:3306)/burger_system?parseTime=true"

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error abriendo conexión: %v", err)
	}

	// Verificar conexión
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("error conectando a MySQL: %v", err)
	}

	log.Println("✅ Conectado a MySQL exitosamente")

	return &Database{conn: conn}, nil
}

// Close cierra la conexión
func (db *Database) Close() {
	db.conn.Close()
}

// ==================== HAMBURGUESAS ====================

// GetHamburguesas obtiene todas las hamburguesas
func (db *Database) GetHamburguesas() ([]*menu.Hamburguesa, error) {
	rows, err := db.conn.Query("SELECT id, nombre, descripcion, precio, categoria, ingredientes, disponible FROM hamburguesas")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hamburguesas []*menu.Hamburguesa

	for rows.Next() {
		var id, nombre, descripcion, categoria, ingredientesStr string
		var precio float64
		var disponible bool

		err := rows.Scan(&id, &nombre, &descripcion, &precio, &categoria, &ingredientesStr, &disponible)
		if err != nil {
			return nil, err
		}

		ingredientes := strings.Split(ingredientesStr, ",")

		h := &menu.Hamburguesa{
			Nombre:       nombre,
			Descripcion:  descripcion,
			Precio:       precio,
			Categoria:    categoria,
			Ingredientes: ingredientes,
			Disponible:   disponible,
		}
		h.SetID(id)

		hamburguesas = append(hamburguesas, h)
	}

	return hamburguesas, nil
}

// GetHamburguesaByID obtiene una hamburguesa por ID
func (db *Database) GetHamburguesaByID(id string) (*menu.Hamburguesa, error) {
	row := db.conn.QueryRow("SELECT id, nombre, descripcion, precio, categoria, ingredientes, disponible FROM hamburguesas WHERE id = ?", id)

	var nombre, descripcion, categoria, ingredientesStr string
	var dbID string
	var precio float64
	var disponible bool

	err := row.Scan(&dbID, &nombre, &descripcion, &precio, &categoria, &ingredientesStr, &disponible)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ingredientes := strings.Split(ingredientesStr, ",")

	h := &menu.Hamburguesa{
		Nombre:       nombre,
		Descripcion:  descripcion,
		Precio:       precio,
		Categoria:    categoria,
		Ingredientes: ingredientes,
		Disponible:   disponible,
	}
	h.SetID(dbID)

	return h, nil
}

// CreateHamburguesa inserta una nueva hamburguesa
func (db *Database) CreateHamburguesa(h *menu.Hamburguesa) error {
	ingredientesStr := strings.Join(h.Ingredientes, ",")

	_, err := db.conn.Exec(
		"INSERT INTO hamburguesas (id, nombre, descripcion, precio, categoria, ingredientes, disponible) VALUES (?, ?, ?, ?, ?, ?, ?)",
		h.GetID(), h.Nombre, h.Descripcion, h.Precio, h.Categoria, ingredientesStr, h.Disponible,
	)

	return err
}

// UpdateHamburguesa actualiza una hamburguesa
func (db *Database) UpdateHamburguesa(id, nombre, descripcion string, precio float64) error {
	result, err := db.conn.Exec(
		"UPDATE hamburguesas SET nombre = ?, descripcion = ?, precio = ? WHERE id = ?",
		nombre, descripcion, precio, id,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("hamburguesa no encontrada")
	}

	return nil
}

// DeleteHamburguesa elimina una hamburguesa
func (db *Database) DeleteHamburguesa(id string) error {
	result, err := db.conn.Exec("DELETE FROM hamburguesas WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("hamburguesa no encontrada")
	}

	return nil
}

// ==================== ÓRDENES ====================

// CreateOrden inserta una nueva orden
func (db *Database) CreateOrden(o *orders.Orden) error {
	// Insertar orden
	_, err := db.conn.Exec(
		"INSERT INTO ordenes (id, numero_mesa, tipo_orden, estado, subtotal, total) VALUES (?, ?, ?, ?, ?, ?)",
		o.ID, o.NumeroMesa, o.TipoOrden, o.Estado, o.Subtotal, o.Total,
	)
	if err != nil {
		return err
	}

	// Insertar items
	for _, item := range o.Items {
		_, err := db.conn.Exec(
			"INSERT INTO items_orden (orden_id, hamburguesa_id, nombre, cantidad, precio_unit, notas) VALUES (?, ?, ?, ?, ?, ?)",
			o.ID, item.HamburguesaID, item.Nombre, item.Cantidad, item.PrecioUnit, item.Notas,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOrdenByID obtiene una orden por ID
func (db *Database) GetOrdenByID(id string) (*orders.Orden, error) {
	// Obtener orden
	row := db.conn.QueryRow("SELECT id, numero_mesa, tipo_orden, estado, subtotal, total, fecha_hora FROM ordenes WHERE id = ?", id)

	var orden orders.Orden
	err := row.Scan(&orden.ID, &orden.NumeroMesa, &orden.TipoOrden, &orden.Estado, &orden.Subtotal, &orden.Total, &orden.FechaHora)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Obtener items
	rows, err := db.conn.Query("SELECT hamburguesa_id, nombre, cantidad, precio_unit, notas FROM items_orden WHERE orden_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item orders.ItemOrden
		var notas sql.NullString
		err := rows.Scan(&item.HamburguesaID, &item.Nombre, &item.Cantidad, &item.PrecioUnit, &notas)
		if err != nil {
			return nil, err
		}
		if notas.Valid {
			item.Notas = notas.String
		}
		orden.Items = append(orden.Items, item)
	}

	return &orden, nil
}

// UpdateOrdenEstado actualiza el estado de una orden
func (db *Database) UpdateOrdenEstado(id, estado string) error {
	result, err := db.conn.Exec("UPDATE ordenes SET estado = ? WHERE id = ?", estado, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("orden no encontrada")
	}

	return nil
}

// GetNextHamburguesaID obtiene el siguiente ID disponible
func (db *Database) GetNextHamburguesaID() (string, error) {
	var maxID sql.NullString
	err := db.conn.QueryRow("SELECT MAX(id) FROM hamburguesas").Scan(&maxID)
	if err != nil {
		return "BURG-001", nil
	}

	if !maxID.Valid {
		return "BURG-001", nil
	}

	// Extraer número del ID (BURG-001 -> 1)
	var num int
	fmt.Sscanf(maxID.String, "BURG-%d", &num)

	return fmt.Sprintf("BURG-%03d", num+1), nil
}

// GetNextOrdenID obtiene el siguiente ID de orden
func (db *Database) GetNextOrdenID() (string, error) {
	var maxID sql.NullString
	err := db.conn.QueryRow("SELECT MAX(id) FROM ordenes").Scan(&maxID)
	if err != nil {
		return "ORD-001", nil
	}

	if !maxID.Valid {
		return "ORD-001", nil
	}

	var num int
	fmt.Sscanf(maxID.String, "ORD-%d", &num)

	return fmt.Sprintf("ORD-%03d", num+1), nil
}
