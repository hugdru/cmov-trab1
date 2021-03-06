package item

import (
	"encoding/json"
	"github.com/pressly/chi"
	"net/http"
	"server/globals/sqlite"
)

func getItems(w http.ResponseWriter, r *http.Request) {
	items, err := sqlite.DB.GetItems()
	if err != nil {
		http.Error(w, "Failed retrieving items", 500)
		return
	}
	itemsSlice, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "Failed converting items to JSON", 500)
		return
	}

	w.Write(buildResponse(itemsSlice, r.URL.Query().Get("hash"), "items"))
}

func getItem(w http.ResponseWriter, r *http.Request) {
	item, err := sqlite.DB.GetItem(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Failed retrieving item", 500)
		return
	}
	itemSlice, err := json.Marshal(item)
	if err != nil {
		http.Error(w, "Failed converting item to JSON", 500)
		return
	}
	w.Write(buildResponse(itemSlice, r.URL.Query().Get("hash"), "item"))
}
