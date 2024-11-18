package requests

import (
	"errors"
	"strings"

	"github.com/Descrout/tosla-go/tosla-go/utils"
)

type Pay3dsRequest struct {
	ThreeDSessionID string `json:"ThreeDSessionId"`
	CardHolderName  string `json:"CardHolderName"`
	CardNo          string `json:"CardNo"`
	ExpireDate      string `json:"ExpireDate"`
	Cvv             string `json:"Cvv"`
}

func (r *Pay3dsRequest) Validate() error {

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
