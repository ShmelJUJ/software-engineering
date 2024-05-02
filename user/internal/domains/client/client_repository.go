package client

import (
	"context"
	"encoding/base64"
	"fmt"

	pg "github.com/ShmelJUJ/software-engineering/pkg/postgres"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
)

type ClientRepository struct {
	cluster pg.Postgres
	ctx     context.Context
}

type UserColumn struct {
	client_id  uuid.UUID `db:"user_id"`
	first_name string    `db:"first_name"`
	last_name  string    `db:"last_name"`
	email      string    `db:"email"`
}

type WalletColumn struct {
	wallet_id   uuid.UUID
	user_id     uuid.UUID
	public_key  string
	private_key string
}

func NewClientRepository(url string) (ClientRepository, error) {
	ctx := context.Background()
	pg_cluster, err := pg.New(ctx, url)
	if err != nil {
		return ClientRepository{}, err
	}
	return ClientRepository{
		cluster: *pg_cluster,
		ctx:     ctx,
	}, nil
}

func (repository *ClientRepository) GetClientById(id uuid.UUID) (Client, error) {
	row, err := repository.cluster.Pool.Query(repository.ctx, kGetClientByid, id.String())
	if err != nil {
		if err == pgx.ErrNoRows {
			return Client{}, &UserNotFoundError{client_id: id}
		}
		return Client{}, err
	}
	defer row.Close()
	client_record, err := pgx.CollectOneRow[UserColumn](row, pgx.RowToStructByName[UserColumn])
	if err != nil {
		return Client{}, err
	}
	wallets, err := repository.GetWalletsByUserId(id)
	if err != nil {
		return Client{}, err
	}
	return NewClient(client_record.client_id, client_record.first_name, client_record.last_name, client_record.email, wallets), nil
}

func DecodePrivateKey(private_key string) string {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(private_key)))
	n, err := base64.StdEncoding.Decode(dst, []byte(private_key))
	if err != nil {
		fmt.Println("decode error:", err)
		return ""
	}
	dst = dst[:n]
	return fmt.Sprintf("%q\n", dst)
}

func (repository *ClientRepository) GetWalletsByUserId(id uuid.UUID) ([]wallet, error) {
	rows, err := repository.cluster.Pool.Query(repository.ctx, kGetWalletByUserId, id.String())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &UserNotFoundError{client_id: id}
		}
		return nil, err
	}
	defer rows.Close()

	wallet_records, err := pgx.CollectRows[WalletColumn](rows, pgx.RowToStructByName[WalletColumn])
	wallets := make([]wallet, 0)
	for _, wallet_record := range wallet_records {
		wallets = append(wallets, NewWallet(wallet_record.wallet_id, wallet_record.public_key, DecodePrivateKey(wallet_record.private_key)))
	}
	return wallets, nil
}
