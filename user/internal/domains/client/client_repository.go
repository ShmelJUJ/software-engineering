package client

import (
	"context"
	"encoding/base64"
	"fmt"

	pg "github.com/ShmelJUJ/software-engineering/pkg/postgres"
	"github.com/go-faster/sdk/zctx"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
)

type ClientRepository struct {
	cluster pg.Postgres
	ctx     context.Context
}

type UserColumns struct {
	UserId    uuid.UUID `db:"user_id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"user_password"`
}

type WalletColumn struct {
	WalletId   uuid.UUID `db:"wallet_id"`
	UserId     uuid.UUID `db:"user_id"`
	PublicKey  string    `db:"public_key"`
	PrivateKey string    `db:"private_key"`
}

func NewClientRepository(url string, ctx context.Context) (ClientRepository, error) {
	lg := zctx.From(ctx)
	pg_cluster, err := pg.New(ctx, url)
	if err != nil {
		lg.Error(err.Error())
		return ClientRepository{}, err
	}
	return ClientRepository{
		cluster: *pg_cluster,
		ctx:     ctx,
	}, nil
}

func (repository *ClientRepository) GetClientById(id uuid.UUID) (Client, error) {
	row, err := repository.cluster.Pool.Query(repository.ctx, kGetClientByid, id.String())
	lg := zctx.From(repository.ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Client{}, &UserNotFoundError{client_id: id.String()}
		}
		lg.Error(err.Error())
		return Client{}, err
	}
	defer row.Close()
	client_record, err := pgx.CollectOneRow[UserColumns](row, pgx.RowToStructByName[UserColumns])
	if err != nil {
		lg.Error(err.Error())

		return Client{}, err
	}
	wallets, err := repository.GetWalletsByUserId(id)
	if err != nil {
		lg.Error(err.Error())

		return Client{}, err
	}

	return NewClient(client_record.UserId, client_record.FirstName, client_record.LastName, client_record.Email, DecodePrivateData(client_record.Password), wallets), nil
}

func (repository *ClientRepository) GetClientByEmail(email string) (Client, error) {
	row, err := repository.cluster.Pool.Query(repository.ctx, kGetClientByEmail, email)
	lg := zctx.From(repository.ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Client{}, &UserNotFoundError{client_id: email}
		}
		lg.Error(err.Error())
		return Client{}, err
	}
	defer row.Close()
	client_record, err := pgx.CollectOneRow[UserColumns](row, pgx.RowToStructByName[UserColumns])
	if err != nil {
		lg.Error(err.Error())

		return Client{}, err
	}
	wallets, err := repository.GetWalletsByUserId(client_record.UserId)
	if err != nil {
		lg.Error(err.Error())

		return Client{}, err
	}
	return NewClient(client_record.UserId, client_record.FirstName, client_record.LastName, client_record.Email, DecodePrivateData(client_record.Password), wallets), nil
}

func DecodePrivateData(private_data string) string {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(private_data)))
	n, err := base64.StdEncoding.Decode(dst, []byte(private_data))
	if err != nil {
		fmt.Println("decode error:", err)
		return ""
	}
	dst = dst[:n]
	return fmt.Sprintf("%q", dst)
}

func (repository *ClientRepository) GetWalletsByUserId(id uuid.UUID) ([]wallet, error) {
	rows, err := repository.cluster.Pool.Query(repository.ctx, kGetWalletByUserId, id.String())
	lg := zctx.From(repository.ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &UserNotFoundError{client_id: id.String()}
		}
		lg.Error(err.Error())

		return nil, err
	}
	defer rows.Close()

	wallet_records, err := pgx.CollectRows[WalletColumn](rows, pgx.RowToStructByName[WalletColumn])
	wallets := make([]wallet, 0)
	for _, wallet_record := range wallet_records {
		wallets = append(wallets, NewWallet(wallet_record.WalletId, wallet_record.PublicKey, DecodePrivateData(wallet_record.PrivateKey)))
	}
	return wallets, nil
}
