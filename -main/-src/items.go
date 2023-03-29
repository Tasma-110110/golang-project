package items

import (
	"encoding/json"
	"net/http"
)

type PurchasingItem struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func createPurchasingItem(w http.ResponseWriter, r *http.Request) {

	var item PurchasingItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := getDB()
	_, err = db.Exec("INSERT INTO purchasing_items (name, description, price, quantity) VALUES (?, ?, ?, ?)", item.Name, item.Description, item.Price, item.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
