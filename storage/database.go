// Package storage maneja la conexión y operaciones CRUD contra MySQL.
// Implementa el patrón repositorio para separar la lógica de datos
// de la lógica de negocio.
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

// Database encapsula la conexión a MySQL
type Database struct {
	conn *sql.DB
}

// NewDatabase abre y verifica la conexión a MySQL.
// Usa zona horaria de Ecuador (America/Guayaquil) para las fechas.
func NewDatabase() (*Database, error) {
	dsn := "root:@tcp(localhost:3306)/burger_system?parseTime=true&loc=America%2FGuayaquil"

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error abriendo conexión: %v", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("error conectando a MySQL: %v", err)
	}

	log.Println("Conectado a MySQL exitosamente")
	return &Database{conn: conn}, nil
}

// Close cierra la conexión a la base de datos
func (db *Database) Close() {
	db.conn.Close()
}

// ==================== OPERACIONES DE HAMBURGUESAS ====================

// GetHamburguesas retorna todas las hamburguesas de la BD
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

		if err := rows.Scan(&id, &nombre, &descripcion, &precio, &categoria, &ingredientesStr, &disponible); err != nil {
			return nil, err
		}

		h := &menu.Hamburguesa{
			Nombre:       nombre,
			Descripcion:  descripcion,
			Precio:       precio,
			Categoria:    categoria,
			Ingredientes: strings.Split(ingredientesStr, ","),
			Disponible:   disponible,
		}
		h.SetID(id)
		hamburguesas = append(hamburguesas, h)
	}
	return hamburguesas, nil
}

// GetHamburguesaByID busca una hamburguesa por ID. Retorna nil si no existe.
func (db *Database) GetHamburguesaByID(id string) (*menu.Hamburguesa, error) {
	row := db.conn.QueryRow(
		"SELECT id, nombre, descripcion, precio, categoria, ingredientes, disponible FROM hamburguesas WHERE id = ?", id,
	)

	var dbID, nombre, descripcion, categoria, ingredientesStr string
	var precio float64
	var disponible bool

	err := row.Scan(&dbID, &nombre, &descripcion, &precio, &categoria, &ingredientesStr, &disponible)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	h := &menu.Hamburguesa{
		Nombre:       nombre,
		Descripcion:  descripcion,
		Precio:       precio,
		Categoria:    categoria,
		Ingredientes: strings.Split(ingredientesStr, ","),
		Disponible:   disponible,
	}
	h.SetID(dbID)
	return h, nil
}

// CreateHamburguesa inserta una nueva hamburguesa en la BD
func (db *Database) CreateHamburguesa(h *menu.Hamburguesa) error {
	_, err := db.conn.Exec(
		"INSERT INTO hamburguesas (id, nombre, descripcion, precio, categoria, ingredientes, disponible) VALUES (?, ?, ?, ?, ?, ?, ?)",
		h.GetID(), h.Nombre, h.Descripcion, h.Precio, h.Categoria,
		strings.Join(h.Ingredientes, ","), h.Disponible,
	)
	return err
}

// UpdateHamburguesa actualiza nombre, descripción y precio.
// Usa SELECT COUNT para verificar existencia porque MySQL reporta
// rowsAffected=0 cuando los datos enviados son idénticos a los existentes.
func (db *Database) UpdateHamburguesa(id, nombre, descripcion string, precio float64) error {
	var count int
	if err := db.conn.QueryRow("SELECT COUNT(*) FROM hamburguesas WHERE id = ?", id).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("hamburguesa no encontrada")
	}

	_, err := db.conn.Exec(
		"UPDATE hamburguesas SET nombre = ?, descripcion = ?, precio = ? WHERE id = ?",
		nombre, descripcion, precio, id,
	)
	return err
}

// DeleteHamburguesa elimina una hamburguesa si no tiene órdenes asociadas.
// MySQL lanza error de foreign key si hay items_orden referenciando esta hamburguesa.
func (db *Database) DeleteHamburguesa(id string) error {
	result, err := db.conn.Exec("DELETE FROM hamburguesas WHERE id = ?", id)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return fmt.Errorf("no se puede eliminar: esta hamburguesa tiene órdenes asociadas. Primero elimine las órdenes que la contienen")
		}
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("hamburguesa no encontrada")
	}
	return nil
}

// ==================== OPERACIONES DE ÓRDENES ====================

// GetAllOrdenes retorna todas las órdenes ordenadas por fecha descendente
func (db *Database) GetAllOrdenes() ([]*orders.Orden, error) {
	rows, err := db.conn.Query(
		"SELECT id, numero_mesa, tipo_orden, estado, subtotal, total, fecha_hora FROM ordenes ORDER BY fecha_hora DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordenes []*orders.Orden
	for rows.Next() {
		var o orders.Orden
		if err := rows.Scan(&o.ID, &o.NumeroMesa, &o.TipoOrden, &o.Estado, &o.Subtotal, &o.Total, &o.FechaHora); err != nil {
			return nil, err
		}
		ordenes = append(ordenes, &o)
	}
	return ordenes, nil
}

// CreateOrden inserta una orden y todos sus items en la BD
func (db *Database) CreateOrden(o *orders.Orden) error {
	_, err := db.conn.Exec(
		"INSERT INTO ordenes (id, numero_mesa, tipo_orden, estado, subtotal, total) VALUES (?, ?, ?, ?, ?, ?)",
		o.ID, o.NumeroMesa, o.TipoOrden, o.Estado, o.Subtotal, o.Total,
	)
	if err != nil {
		return err
	}

	// Insertar cada item de la orden
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

// GetOrdenByID retorna una orden con sus items. Retorna nil si no existe.
func (db *Database) GetOrdenByID(id string) (*orders.Orden, error) {
	row := db.conn.QueryRow(
		"SELECT id, numero_mesa, tipo_orden, estado, subtotal, total, fecha_hora FROM ordenes WHERE id = ?", id,
	)

	var orden orders.Orden
	err := row.Scan(&orden.ID, &orden.NumeroMesa, &orden.TipoOrden, &orden.Estado, &orden.Subtotal, &orden.Total, &orden.FechaHora)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Cargar los items asociados a esta orden
	rows, err := db.conn.Query(
		"SELECT hamburguesa_id, nombre, cantidad, precio_unit, notas FROM items_orden WHERE orden_id = ?", id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item orders.ItemOrden
		var notas sql.NullString
		if err := rows.Scan(&item.HamburguesaID, &item.Nombre, &item.Cantidad, &item.PrecioUnit, &notas); err != nil {
			return nil, err
		}
		if notas.Valid {
			item.Notas = notas.String
		}
		orden.Items = append(orden.Items, item)
	}
	return &orden, nil
}

// UpdateOrdenEstado cambia el estado de una orden.
// Usa SELECT COUNT para evitar falso negativo si el estado no cambia.
func (db *Database) UpdateOrdenEstado(id, estado string) error {
	var count int
	if err := db.conn.QueryRow("SELECT COUNT(*) FROM ordenes WHERE id = ?", id).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("orden no encontrada")
	}

	_, err := db.conn.Exec("UPDATE ordenes SET estado = ? WHERE id = ?", estado, id)
	return err
}

// DeleteOrden elimina una orden y sus items asociados.
// Primero elimina los items (tabla hija) para respetar las foreign keys.
func (db *Database) DeleteOrden(id string) error {
	_, err := db.conn.Exec("DELETE FROM items_orden WHERE orden_id = ?", id)
	if err != nil {
		return err
	}

	result, err := db.conn.Exec("DELETE FROM ordenes WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("orden no encontrada")
	}
	return nil
}

// ==================== GENERADORES DE ID ====================

// GetNextHamburguesaID genera el siguiente ID secuencial (BURG-001, BURG-002...)
func (db *Database) GetNextHamburguesaID() (string, error) {
	var maxID sql.NullString
	db.conn.QueryRow("SELECT MAX(id) FROM hamburguesas").Scan(&maxID)

	if !maxID.Valid {
		return "BURG-001", nil
	}

	var num int
	fmt.Sscanf(maxID.String, "BURG-%d", &num)
	return fmt.Sprintf("BURG-%03d", num+1), nil
}

// GetNextOrdenID genera el siguiente ID secuencial (ORD-001, ORD-002...)
func (db *Database) GetNextOrdenID() (string, error) {
	var maxID sql.NullString
	db.conn.QueryRow("SELECT MAX(id) FROM ordenes").Scan(&maxID)

	if !maxID.Valid {
		return "ORD-001", nil
	}

	var num int
	fmt.Sscanf(maxID.String, "ORD-%d", &num)
	return fmt.Sprintf("ORD-%03d", num+1), nil
}
