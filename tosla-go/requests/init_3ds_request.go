package requests

type Init3dsRequest struct {
	OrderID          string `json:"orderId"`
	CallbackURL      string `json:"callbackUrl"`
	Description      string `json:"description"`
	Echo             string `json:"echo"`
	ExtraParameters  string `json:"extraParameters"`
	Amount           int    `json:"amount"`
	Currency         int    `json:"currency"`
	InstallmentCount int    `json:"installmentCount"`
}
