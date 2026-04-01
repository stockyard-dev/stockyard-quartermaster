package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-quartermaster/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){q:=r.URL.Query().Get("q");cat:=r.URL.Query().Get("category");list,_:=s.db.List(q,cat);if list==nil{list=[]store.Item{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var item store.Item;json.NewDecoder(r.Body).Decode(&item);if item.Name==""{writeError(w,400,"name required");return};s.db.Create(&item);writeJSON(w,201,item)}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
