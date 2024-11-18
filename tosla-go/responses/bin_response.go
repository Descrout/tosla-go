package responses

type BinResponse struct {
	CardPrefix      int                        `json:"CardPrefix"`
	BankID          int                        `json:"BankId"`
	BankCode        string                     `json:"BankCode"`
	BankName        string                     `json:"BankName"`
	CardName        string                     `json:"CardName"`
	CardClass       string                     `json:"CardClass"`
	CardType        string                     `json:"CardType"`
	Country         string                     `json:"Country"`
	BankCommission  int                        `json:"BankCommission"`
	InstallmentInfo map[string]InstallmentRate `json:"InstallmentInfo"`
	Code            int                        `json:"Code"`
	Message         string                     `json:"Message"`
}

type InstallmentRate struct {
	Rate     float64 `json:"Rate"`
	Constant int     `json:"Constant"`
}
