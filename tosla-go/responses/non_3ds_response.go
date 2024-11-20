package responses

type Non3dsResponse struct {
	Code                int    `json:"Code"`
	Message             string `json:"Message"`
	OrderID             string `json:"OrderId"`
	BankResponseCode    string `json:"BankResponseCode"`
	BankResponseMessage string `json:"BankResponseMessage"`
	AuthCode            string `json:"AuthCode"`
	HostReferenceNumber string `json:"HostReferenceNumber"`
	TransactionID       string `json:"TransactionId"`
	CardHolderName      string `json:"CardHolderName"`
}
