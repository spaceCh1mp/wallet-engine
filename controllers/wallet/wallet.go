package wallet

import (
	"context"
	"fmt"
	"time"
	"wallet-engine/config"
	"wallet-engine/errors"
	"wallet-engine/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WalletController interface {
	GetUserWallet(userID primitive.ObjectID) (*models.Wallet, error)
	CreateWallet(userID primitive.ObjectID) (*models.Wallet, error)
	FundWallet(id primitive.ObjectID, amount uint64) (*models.Wallet, error)
	WithdrawFromWallet(id primitive.ObjectID, amount int64) (*models.Wallet, error)
	ToggleWalletActiveStatus(id primitive.ObjectID, isActive bool) (*models.Wallet, error)
}

func NewWalletController(db *mongo.Database) WalletController {
	return &walletController{
		collection: db.Collection(config.WalletCollectioName),
	}
}

type walletController struct {
	collection *mongo.Collection
}

func (wc walletController) Bool(value bool) *bool {
	return &value
}

func (wc walletController) GetUserWallet(userID primitive.ObjectID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := wc.collection.FindOne(context.TODO(), bson.M{"userID": userID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound(fmt.Errorf("user not found"))
		}
		return nil, err
	}

	return &wallet, nil
}

func (wc walletController) CreateWallet(userID primitive.ObjectID) (*models.Wallet, error) {
	timestamp := time.Now().UTC()
	newWallet := &models.Wallet{
		UserID:    userID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		IsActive:  wc.Bool(true),
	}
	result, err := wc.collection.InsertOne(context.TODO(), newWallet)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.BadRequest(fmt.Errorf("user already owns wallet"))
		}
		return nil, err
	}

	newWallet.ID = result.InsertedID.(primitive.ObjectID)
	return newWallet, nil
}

func (wc walletController) FundWallet(id primitive.ObjectID, amount uint64) (*models.Wallet, error) {
	return wc.updateWallet(id, bson.M{"$inc": bson.M{"balance": amount}})
}

func (wc walletController) WithdrawFromWallet(id primitive.ObjectID, amount int64) (*models.Wallet, error) {
	return wc.updateWallet(id, bson.M{"$inc": bson.M{"balance": -amount}})
}

func (wc walletController) updateWallet(ID primitive.ObjectID, update bson.M) (*models.Wallet, error) {
	var updatedWallet models.Wallet
	updateOption := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := wc.collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": ID},
		update, updateOption).Decode(&updatedWallet)
	if err != nil {
		return nil, err
	}

	return &updatedWallet, nil
}

func (wc walletController) ToggleWalletActiveStatus(id primitive.ObjectID, isActive bool) (*models.Wallet, error) {
	return wc.updateWallet(id, bson.M{"$set": bson.M{"isActive": isActive}})
}
