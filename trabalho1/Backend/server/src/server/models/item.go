package models

type Item struct {
	IdItem      	int64		`json:"idItem"`
	ItemType 	ItemType	`json:"itemType"`
	Price       	float64		`json:"price"`
	Description 	string		`json:"description"`
}