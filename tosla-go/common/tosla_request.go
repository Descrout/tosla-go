package common

type ToslaRequest struct {
	ClientID string `json:"clientId"`
	APIUser  string `json:"apiUser"`
	Rnd      string `json:"rnd"`
	TimeSpan string `json:"timeSpan"`
	Hash     string `json:"hash"`
}
