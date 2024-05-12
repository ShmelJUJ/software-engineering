package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShmelJUJ/software-engineering/user/gen"
	"github.com/ShmelJUJ/software-engineering/user/internal/domains/client"
	"github.com/go-faster/sdk/zctx"
	"github.com/ogen-go/ogen/conv"
	"go.uber.org/zap"
)

const pgURL = "postgres://user_service:user_root@user_postgres:5432/user_db?sslmode=disable"

// Compile-time check for Handler.
var _ user.Handler = (*Handler)(nil)

type Handler struct {
	user.UnimplementedHandler // automatically implement all methods
}

func (h Handler) GetClientById(ctx context.Context, params user.GetClientByIdParams) (user.GetClientByIdRes, error) {
	zctx.From(ctx).Info("GetClient", zap.Any("params", params))
	repository, err := client.NewClientRepository(pgURL, ctx) //TODO наебашить конфиг для подключения к бд
	var not_found *client.UserNotFoundError
	if err != nil {

		return &user.GetClientByIdBadRequest{}, err
	}
	converted_client_id, err := conv.ToUUID(params.ClientID)
	if err != nil {
		zctx.From(ctx).Error(err.Error())
		return &user.GetClientByIdBadRequest{}, err
	}
	user_from_db, err := repository.GetClientById(converted_client_id)
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
	zctx.From(ctx).Info("GetWallet", zap.Any("params", params))
	repository, err := client.NewClientRepository(pgURL, ctx) //TODO наебашить конфиг для подключения к бд
	var not_found *client.UserNotFoundError
	if err != nil {

		return &user.GetWalletByIdBadRequest{}, err
	}
	converted_client_id, err := conv.ToUUID(params.ClientID)
	if err != nil {
		zctx.From(ctx).Error(err.Error())
		return &user.GetWalletByIdBadRequest{Code: "parse_error", Message: "cant parse client_id to uuid"}, err
	}
	user_from_db, err := repository.GetClientById(converted_client_id)
	if err != nil {
		if errors.As(err, &not_found) {
			return &user.GetWalletByIdNotFound{Code: "user_not_found", Message: not_found.Error()}, nil
		}
		return &user.GetWalletByIdBadRequest{}, err
	}
	for _, wallet := range user_from_db.GetWallets() {
		if wallet.GetWalletId().String() == params.WalletID {
			return &user.Wallet{
				PublicKey:  wallet.GetPublicKey(),
				PrivateKey: wallet.GetPrivateKey(),
			}, nil
		}
	}
	return &user.GetWalletByIdNotFound{Code: "wallet_not_found", Message: fmt.Sprintf("wallet with id %s not found", params.WalletID)}, nil
}

func (h Handler) GetAuthToken(ctx context.Context, request user.OptAuthRequest) (r user.GetAuthTokenRes, _ error) {
	if !request.Set {
		zctx.From(ctx).Error("Empty Get Auth Token Request")
		return &user.GetAuthTokenBadRequest{}, nil
	}
	zctx.From(ctx).Info("GetToken", zap.Any("params", request.Value.Email))
	repository, err := client.NewClientRepository(pgURL, ctx) //TODO наебашить конфиг для подключения к бд
	var not_found *client.UserNotFoundError
	if err != nil {

		return &user.GetAuthTokenBadRequest{}, err
	}
	user_from_db, err := repository.GetClientByEmail(request.Value.Email)
	if err != nil {
		if errors.As(err, &not_found) {
			return &user.GetAuthTokenNotFound{Code: "user_not_found", Message: not_found.Error()}, nil
		}
		return &user.GetAuthTokenBadRequest{}, err
	}
	auth_token, err := user_from_db.GetAuthToken(request.Value.Password)
	var wrong_password_err *client.WrongPasswordError
	if err != nil {
		if errors.As(err, &wrong_password_err) {
			return &user.GetAuthTokenForbidden{Code: "wrong_password", Message: wrong_password_err.Error()}, wrong_password_err
		}
		return &user.GetAuthTokenBadRequest{}, err
	}
	return &user.AuthResponse{
		AuthToken: auth_token,
	}, nil
}
