package data

import (
	"database/sql"
	"time"
)

type PaymentModel struct {
	DB *sql.DB
}

type Payment struct {
	ID      int64 `json:"id"`
	OrderID int64 `json:"order_id"`
	UserID  int64 `json:"user_id"`

	StripePaymentOrderID string `json:"stripe_payment_intent_id"`

	AmountCents int    `json:"amount_cents"`
	Currency    string `json:"currency"`

	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	PaymentMethod string `json:"payment_method"`
	Captured      bool   `json:"captured"`

	FailureReason *string `json:"failure_reason"`

	Metadata string `json:"metadata"`
	Version  int    `json:"-"`
}
