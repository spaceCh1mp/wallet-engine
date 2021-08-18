package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wallet-engine/context"
	"wallet-engine/controllers/wallet"
	"wallet-engine/errors"
	"wallet-engine/handlers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletHandler interface {
	CreateWallet(w http.ResponseWriter, r *http.Request) error
	FundWallet(w http.ResponseWriter, r *http.Request) error
	WithdrawFromWallet(w http.ResponseWriter, r *http.Request) error
	ActivateWallet(w http.ResponseWriter, r *http.Request) error
	DeactivateWallet(w http.ResponseWriter, r *http.Request) error

	ValidateUserHasActiveWallet(next handlers.AppHandler) http.Handler
	ValidateUserWalletIsNotActive(next handlers.AppHandler) http.Handler
}

func NewWalletHandler(walletController wallet.WalletController) WalletHandler {
	return &walletHandler{
		walletController: walletController,
	}
}

type walletHandler struct {
	walletController wallet.WalletController
}

func (wh walletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	newWallet, err := wh.walletController.CreateWallet(userID)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(newWallet)
}

func (wh walletHandler) FundWallet(w http.ResponseWriter, r *http.Request) error {
	ctxWallet := context.GetWalletData(r)
	fundRequest := &FundOrWithdrawWalletRequestBody{}
	err := handlers.Decode(r.Body, fundRequest)
	if err != nil {
		return err
	}

	fmt.Println(fundRequest.Amount)
	if fundRequest.Amount < 1000 || fundRequest.Amount > 100000000 {
		return errors.BadRequest(fmt.Errorf("you can only fund between NGN10 to NGN1,000,000"))
	}

	wall, err := wh.walletController.FundWallet(ctxWallet.ID, uint64(fundRequest.Amount))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(wall)
}

func (wh walletHandler) WithdrawFromWallet(w http.ResponseWriter, r *http.Request) error {
	ctxWallet := context.GetWalletData(r)
	withdrawRequest := &FundOrWithdrawWalletRequestBody{}
	err := handlers.Decode(r.Body, withdrawRequest)
	if err != nil {
		return err
	}

	fmt.Println(withdrawRequest.Amount)
	fmt.Println(ctxWallet.Balance)
	if withdrawRequest.Amount < 1000 {
		return errors.BadRequest(fmt.Errorf("you can't withdraw %d, minimum amount is NGN100000", withdrawRequest.Amount))
	}
	if withdrawRequest.Amount > int(ctxWallet.Balance) {
		return errors.BadRequest(fmt.Errorf("you do not have sufficient funds"))
	}

	wall, err := wh.walletController.WithdrawFromWallet(ctxWallet.ID, int64(withdrawRequest.Amount))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(wall)
}

func (wh walletHandler) ActivateWallet(w http.ResponseWriter, r *http.Request) error {
	ctxWallet := context.GetWalletData(r)

	wall, err := wh.walletController.ToggleWalletActiveStatus(ctxWallet.ID, true)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(wall)
}

func (wh walletHandler) DeactivateWallet(w http.ResponseWriter, r *http.Request) error {
	ctxWallet := context.GetWalletData(r)

	wall, err := wh.walletController.ToggleWalletActiveStatus(ctxWallet.ID, false)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(wall)
}

func (wh walletHandler) ValidateUserHasActiveWallet(next handlers.AppHandler) http.Handler {
	return handlers.AppHandler(func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["id"]
		userID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		wallet, err := wh.walletController.GetUserWallet(userID)
		if err != nil {
			return err
		}

		r = context.SaveWalletData(r, wallet)
		if wallet.IsActive != nil {
			if *wallet.IsActive == true {
				next.ServeHTTP(w, r)
				return nil
			}

			return errors.BadRequest(fmt.Errorf("user is deactivated"))
		}

		return errors.BadRequest(fmt.Errorf("user does not have wallet"))
	})
}

func (wh walletHandler) ValidateUserWalletIsNotActive(next handlers.AppHandler) http.Handler {
	return handlers.AppHandler(func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["id"]
		userID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		wallet, err := wh.walletController.GetUserWallet(userID)
		if err != nil {
			return err
		}

		r = context.SaveWalletData(r, wallet)
		if wallet.IsActive != nil {
			if *wallet.IsActive == false {
				next.ServeHTTP(w, r)
				return nil
			}

			return errors.BadRequest(fmt.Errorf("user is still active"))
		}

		return errors.BadRequest(fmt.Errorf("user does not have wallet"))
	})
}
