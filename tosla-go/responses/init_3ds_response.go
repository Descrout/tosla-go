package responses

type Init3dsResponse struct {
	Code            int    `json:"Code"`
	Message         string `json:"Message"`
	ThreeDSessionID string `json:"ThreeDSessionId"`
	TransactionID   string `json:"TransactionId"`
}
