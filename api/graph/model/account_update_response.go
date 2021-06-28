package model

type AccountUpdateResponse struct {
	// The address of the account that was updated.
	Address string `json:"address"`
	// The transactions that caused the update.
	Transactions []Transaction `json:"transactions"`
}
