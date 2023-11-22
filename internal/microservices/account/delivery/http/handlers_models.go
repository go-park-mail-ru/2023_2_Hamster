package http

type CreateAccount struct {
	Balance     float64 `json:"balance"`
	MeanPayment string  `json:"mean_payment"`
}
