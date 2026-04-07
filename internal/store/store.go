package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// InventoryItem represents a single tracked inventory item.
// PurchasePrice is stored as cents (integer) to avoid float math errors.
type InventoryItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Category      string `json:"category"`
	Location      string `json:"location"`
	Quantity      int    `json:"quantity"`
	PurchaseDate  string `json:"purchase_date"`
	PurchasePrice int    `json:"purchase_price"`
	Condition     string `json:"condition"`
	SerialNumber  string `json:"serial_number"`
	Notes         string `json:"notes"`
	CreatedAt     string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "quartermaster.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS inventory(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		category TEXT DEFAULT '',
		location TEXT DEFAULT '',
		quantity INTEGER DEFAULT 1,
		purchase_date TEXT DEFAULT '',
		purchase_price INTEGER DEFAULT 0,
		condition TEXT DEFAULT 'good',
		serial_number TEXT DEFAULT '',
		notes TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *InventoryItem) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Condition == "" {
		e.Condition = "good"
	}
	_, err := d.db.Exec(
		`INSERT INTO inventory(id, name, category, location, quantity, purchase_date, purchase_price, condition, serial_number, notes, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Name, e.Category, e.Location, e.Quantity, e.PurchaseDate, e.PurchasePrice, e.Condition, e.SerialNumber, e.Notes, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *InventoryItem {
	var e InventoryItem
	err := d.db.QueryRow(
		`SELECT id, name, category, location, quantity, purchase_date, purchase_price, condition, serial_number, notes, created_at
		 FROM inventory WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Name, &e.Category, &e.Location, &e.Quantity, &e.PurchaseDate, &e.PurchasePrice, &e.Condition, &e.SerialNumber, &e.Notes, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []InventoryItem {
	rows, _ := d.db.Query(
		`SELECT id, name, category, location, quantity, purchase_date, purchase_price, condition, serial_number, notes, created_at
		 FROM inventory ORDER BY name ASC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []InventoryItem
	for rows.Next() {
		var e InventoryItem
		rows.Scan(&e.ID, &e.Name, &e.Category, &e.Location, &e.Quantity, &e.PurchaseDate, &e.PurchasePrice, &e.Condition, &e.SerialNumber, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *InventoryItem) error {
	_, err := d.db.Exec(
		`UPDATE inventory SET name=?, category=?, location=?, quantity=?, purchase_date=?, purchase_price=?, condition=?, serial_number=?, notes=?
		 WHERE id=?`,
		e.Name, e.Category, e.Location, e.Quantity, e.PurchaseDate, e.PurchasePrice, e.Condition, e.SerialNumber, e.Notes, e.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM inventory WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM inventory`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []InventoryItem {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (name LIKE ? OR category LIKE ? OR location LIKE ? OR serial_number LIKE ?)"
		args = append(args, "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}
	if v, ok := filters["category"]; ok && v != "" {
		where += " AND category=?"
		args = append(args, v)
	}
	if v, ok := filters["condition"]; ok && v != "" {
		where += " AND condition=?"
		args = append(args, v)
	}
	if v, ok := filters["location"]; ok && v != "" {
		where += " AND location=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, name, category, location, quantity, purchase_date, purchase_price, condition, serial_number, notes, created_at
		 FROM inventory WHERE `+where+` ORDER BY name ASC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []InventoryItem
	for rows.Next() {
		var e InventoryItem
		rows.Scan(&e.ID, &e.Name, &e.Category, &e.Location, &e.Quantity, &e.PurchaseDate, &e.PurchasePrice, &e.Condition, &e.SerialNumber, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns total item count, total quantity, total value (in cents),
// and breakdowns by category, condition, and location.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":        d.Count(),
		"total_qty":    0,
		"total_value":  0,
		"by_category":  map[string]int{},
		"by_condition": map[string]int{},
		"by_location":  map[string]int{},
	}

	var totalQty, totalValue int
	d.db.QueryRow(`SELECT COALESCE(SUM(quantity), 0), COALESCE(SUM(quantity * purchase_price), 0) FROM inventory`).Scan(&totalQty, &totalValue)
	m["total_qty"] = totalQty
	m["total_value"] = totalValue

	if rows, _ := d.db.Query(`SELECT category, COUNT(*) FROM inventory WHERE category != '' GROUP BY category`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var k string
			var c int
			rows.Scan(&k, &c)
			by[k] = c
		}
		m["by_category"] = by
	}

	if rows, _ := d.db.Query(`SELECT condition, COUNT(*) FROM inventory WHERE condition != '' GROUP BY condition`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var k string
			var c int
			rows.Scan(&k, &c)
			by[k] = c
		}
		m["by_condition"] = by
	}

	if rows, _ := d.db.Query(`SELECT location, COUNT(*) FROM inventory WHERE location != '' GROUP BY location`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var k string
			var c int
			rows.Scan(&k, &c)
			by[k] = c
		}
		m["by_location"] = by
	}

	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
