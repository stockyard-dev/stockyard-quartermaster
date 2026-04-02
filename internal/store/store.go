package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type InventoryItem struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Category string `json:"category"`
	Location string `json:"location"`
	Quantity int `json:"quantity"`
	PurchaseDate string `json:"purchase_date"`
	PurchasePrice int `json:"purchase_price"`
	Condition string `json:"condition"`
	SerialNumber string `json:"serial_number"`
	Notes string `json:"notes"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"quartermaster.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS inventory(id TEXT PRIMARY KEY,name TEXT NOT NULL,category TEXT DEFAULT '',location TEXT DEFAULT '',quantity INTEGER DEFAULT 1,purchase_date TEXT DEFAULT '',purchase_price INTEGER DEFAULT 0,condition TEXT DEFAULT 'good',serial_number TEXT DEFAULT '',notes TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *InventoryItem)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO inventory(id,name,category,location,quantity,purchase_date,purchase_price,condition,serial_number,notes,created_at)VALUES(?,?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Category,e.Location,e.Quantity,e.PurchaseDate,e.PurchasePrice,e.Condition,e.SerialNumber,e.Notes,e.CreatedAt);return err}
func(d *DB)Get(id string)*InventoryItem{var e InventoryItem;if d.db.QueryRow(`SELECT id,name,category,location,quantity,purchase_date,purchase_price,condition,serial_number,notes,created_at FROM inventory WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Category,&e.Location,&e.Quantity,&e.PurchaseDate,&e.PurchasePrice,&e.Condition,&e.SerialNumber,&e.Notes,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]InventoryItem{rows,_:=d.db.Query(`SELECT id,name,category,location,quantity,purchase_date,purchase_price,condition,serial_number,notes,created_at FROM inventory ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []InventoryItem;for rows.Next(){var e InventoryItem;rows.Scan(&e.ID,&e.Name,&e.Category,&e.Location,&e.Quantity,&e.PurchaseDate,&e.PurchasePrice,&e.Condition,&e.SerialNumber,&e.Notes,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *InventoryItem)error{_,err:=d.db.Exec(`UPDATE inventory SET name=?,category=?,location=?,quantity=?,purchase_date=?,purchase_price=?,condition=?,serial_number=?,notes=? WHERE id=?`,e.Name,e.Category,e.Location,e.Quantity,e.PurchaseDate,e.PurchasePrice,e.Condition,e.SerialNumber,e.Notes,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM inventory WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM inventory`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]InventoryItem{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["category"];ok&&v!=""{where+=" AND category=?";args=append(args,v)}
    if v,ok:=filters["condition"];ok&&v!=""{where+=" AND condition=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,category,location,quantity,purchase_date,purchase_price,condition,serial_number,notes,created_at FROM inventory WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []InventoryItem;for rows.Next(){var e InventoryItem;rows.Scan(&e.ID,&e.Name,&e.Category,&e.Location,&e.Quantity,&e.PurchaseDate,&e.PurchasePrice,&e.Condition,&e.SerialNumber,&e.Notes,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    return m
}
