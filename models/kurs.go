package models

type KursExchange struct {
	Beli float64 `json:"beli" binding:"required"`
	Jual float64 `json:"jual" binding:"required"`
}

type Kurs struct {
	ID         int64        `json:"id"`
	Symbol     string       `json:"symbol" binding:"required"`
	E_rate     KursExchange `json:"e_rate" binding:"required"`
	TT_counter KursExchange `json:"tt_counter" binding:"required"`
	Bank_notes KursExchange `json:"bank_notes" binding:"required"`
	Date       string       `json:"date" binding:"required"`
}
