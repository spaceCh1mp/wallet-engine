package wallet

// FundOrWithdrawWalletRequestBody defines the schema for a fund or withdraw request
type FundOrWithdrawWalletRequestBody struct {
	Amount int `json:"amount"`
}

// ToggleWalletStatusRequestBody defines the schema for an activate/deactivate wallet request
type ToggleWalletStatusRequestBody struct {
	Active *bool `json:"active"`
}
