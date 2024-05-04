package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShmelJUJ/software-engineering/user/internal/domains/client"
	"github.com/ShmelJUJ/software-engineering/user/gen"

	"github.com/go-faster/sdk/zctx"
	"go.uber.org/zap"
)

const pgURL = "postgres://postgres:root@user_postgres:5432/postgres?sslmode=disable"

// Compile-time check for Handler.
var _ user.Handler = (*Handler)(nil)

type Handler struct {
	user.UnimplementedHandler // automatically implement all methods
}

func (h Handler) GetClientById(ctx context.Context, params user.GetClientByIdParams) (user.GetClientByIdRes, error) {
	zctx.From(ctx).Info("GetClient", zap.Any("params", params))
	repository, err := client.NewClientRepository(pgURL) //TODO наебашить конфиг для подключения к бд
	var not_found *client.UserNotFoundError
	if err != nil {

		return &user.GetClientByIdBadRequest{}, err
	}
	user_from_db, err := repository.GetClientById(params.ClientID)
	if err != nil {
		if errors.As(err, &not_found) {
			return &user.GetClientByIdNotFound{Code: "user_not_found", Message: not_found.Error()}, nil
		}
		return &user.GetClientByIdBadRequest{}, err
	}
	return &user.User{
		ClientID:  user_from_db.GetClientId(),
		FirstName: user.NewOptString(user_from_db.GetFirstName()),
		LastName:  user.NewOptString(user_from_db.GetLastName()),
		Email:     user_from_db.GetEmail(),
		Wallets: func() (wallet_ids []string) {
			for _, wallet := range user_from_db.GetWallets() {
				wallet_ids = append(wallet_ids, wallet.GetPublicKey())
			}
			return wallet_ids
		}(),
	}, nil
}

func (h Handler) GetWalletById(ctx context.Context, params user.GetWalletByIdParams) (r user.GetWalletByIdRes, _ error) {
	zctx.From(ctx).Info("GetClient", zap.Any("params", params))
	repository, err := client.NewClientRepository(pgURL) //TODO наебашить конфиг для подключения к бд
	var not_found *client.UserNotFoundError
	if err != nil {

		return &user.GetWalletByIdBadRequest{}, err
	}
	user_from_db, err := repository.GetClientById(params.ClientID)
	if err != nil {
		if errors.As(err, &not_found) {
			return &user.GetWalletByIdNotFound{Code: "user_not_found", Message: not_found.Error()}, nil
		}
		return &user.GetWalletByIdBadRequest{}, err
	}
	for _, wallet := range user_from_db.GetWallets() {
		if wallet.GetWalletId() == params.WalletID {
			return &user.Wallet{
				PublicKey:  wallet.GetPublicKey(),
				PrivateKey: wallet.GetPrivateKey(),
			}, nil
		}
	}
	return &user.GetWalletByIdNotFound{Code: "wallet_not_found", Message: fmt.Sprintf("wallet with id %s not found", params.WalletID.String())}, nil
}
