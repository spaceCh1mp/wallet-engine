package wallet

// FundOrWithdrawWalletRequestBody defines the schema for a fund or withdraw request
type FundOrWithdrawWalletRequestBody struct {
	Amount int `json:"amount"`
}
