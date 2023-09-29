package mail

type Request struct {
	ID       string `json:"id" validate:"required"`
	CURRENCY string `json:"currency" validate:"required"`
}

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Message Crypto `json:"Message,omitempty"`
}

type Crypto struct {
	Credits Credits `json:"data,omitempty"`
}

type Credits struct {
	PRICE    string `json:"amount,omitempty"`
	ID       string `json:"base,omitempty"`
	CURRENCY string `json:"currency,omitempty"`
}

func GenerateRequest(request Request) string {
	return "https://api.coinbase.com/v2/prices/" + request.ID + "-" + request.CURRENCY + "/buy"
}
