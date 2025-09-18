package data

import "database/sql"

type Models struct {
	Payment PaymentModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Payment: PaymentModel{DB: db},
	}
}
