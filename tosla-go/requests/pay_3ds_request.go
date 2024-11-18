package requests

type Pay3dsRequest struct {
	ThreeDSessionID string `json:"ThreeDSessionId"`
	CardHolderName  string `json:"CardHolderName"`
	CardNo          string `json:"CardNo"`
	ExpireDate      string `json:"ExpireDate"`
	Cvv             string `json:"Cvv"`
}
