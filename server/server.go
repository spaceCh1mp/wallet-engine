package server

import (
	"fmt"
	"net/http"
	"wallet-engine/config"
	wallet2 "wallet-engine/controllers/wallet"
	"wallet-engine/handlers"
	"wallet-engine/handlers/wallet"
	"wallet-engine/mongo"

	"github.com/gorilla/mux"
)

func Init() error {
	router := mux.NewRouter()

	// initialise dependencies
	mongoConn := mongo.Connect(config.DBURL, config.DBName)
	if err := mongo.EnforceValidation(mongoConn); err != nil {
		return fmt.Errorf("could not set validator %v", err)
	}
	if err := mongo.CreateIndexes(mongoConn, indexes); err != nil {
		return fmt.Errorf("could not create indexes %v", err)
	}

	walletController := wallet2.NewWalletController(mongoConn)
	walletHandler := wallet.NewWalletHandler(walletController)

	// declare routes
	router.Handle("/{id}", handlers.AppHandler(walletHandler.CreateWallet)).Methods(http.MethodPost)
	router.Handle("/{id}/fund",
		walletHandler.ValidateUserHasActiveWallet(
			walletHandler.FundWallet,
		),
	).Methods(http.MethodPost)
	router.Handle("/{id}/withdraw",
		walletHandler.ValidateUserHasActiveWallet(
			walletHandler.WithdrawFromWallet,
		),
	).Methods(http.MethodPost)
	router.Handle(
		"/{id}/activate",
		walletHandler.ValidateUserWalletIsNotActive(
			walletHandler.ActivateWallet,
		),
	).Methods(http.MethodPatch)
	router.Handle(
		"/{id}/deactivate",
		walletHandler.ValidateUserHasActiveWallet(
			walletHandler.DeactivateWallet,
		),
	).Methods(http.MethodPatch)

	return http.ListenAndServe(config.PORT, router)
}
