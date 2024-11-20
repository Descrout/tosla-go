package requests

import (
	"errors"
	"strings"

	"github.com/Descrout/tosla-go/tosla-go/utils"
)

type Non3dsRequest struct {
	CardHolderName   string `json:"cardHolderName"`
	CardNo           string `json:"cardNo"`
	ExpireDate       string `json:"expireDate"`
	Cvv              string `json:"cvv"`
	OrderID          string `json:"orderId"`
	Amount           int    `json:"amount"`
	Currency         int    `json:"currency"`
	InstallmentCount int    `json:"installmentCount"`
	Description      string `json:"description"`
	Echo             string `json:"echo"`
	ExtraParameters  string `json:"extraParameters"`
}

func (r *Non3dsRequest) Validate() error {

	if len(strings.Split(r.CardHolderName, " ")) < 2 {
		return errors.New("card holder must have name and surname")
	}

	if !utils.LuhnCheck(r.CardNo) {
		return errors.New("invalid card no")
	}

	if utils.CheckExpiration(r.ExpireDate) {
		return errors.New("card is expired or has invalid expiry date")
	}

	return nil
}
