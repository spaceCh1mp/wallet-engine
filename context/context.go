package context

import (
	"context"
	"net/http"
	"wallet-engine/models"
)

const (
	walletDataKey = "walletdata"
)

type key struct {
	ctxKey string
}

// SaveWalletData saves wallet data to context
func SaveWalletData(r *http.Request, user *models.Wallet) *http.Request {
	return r.WithContext(save(r.Context(), walletDataKey, user))
}

// GetWalletData returns the wallet data in context
func GetWalletData(r *http.Request) *models.Wallet {
	return r.Context().Value(key{walletDataKey}).(*models.Wallet)
}

func save(ctx context.Context, ctxKey string, value interface{}) context.Context {
	return context.WithValue(ctx, key{ctxKey}, value)
}
