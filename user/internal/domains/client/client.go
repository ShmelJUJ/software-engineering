package client

import "github.com/google/uuid"

type wallet struct {
	wallet_id   uuid.UUID
	public_key  string
	private_key string
}

type Client struct {
	client_id         uuid.UUID
	client_first_name string
	client_last_name  string
	email             string
	wallets           []wallet
}

func NewClient(client_id uuid.UUID,
	client_first_name string,
	client_last_name string,
	email string, wallets []wallet,
) (client Client) {
	client.client_id = client_id
	client.client_first_name = client_first_name
	client.client_last_name = client_last_name
	client.email = email
	client.wallets = wallets
	return client
}

func (client *Client) GetClientId() uuid.UUID {
	return client.client_id
}

func (client *Client) GetFirstName() string {
	return client.client_first_name
}

func (client *Client) GetLastName() string {
	return client.client_last_name
}

func (client *Client) GetEmail() string {
	return client.email
}

func (client *Client) GetWallets() []wallet {
	return client.wallets
}

func NewWallet(wallet_id uuid.UUID, public_key string, private_key string) wallet {
	return wallet{
		wallet_id:   wallet_id,
		public_key:  public_key,
		private_key: private_key,
	}
}

func (wal *wallet) GetPublicKey() string {
	return wal.public_key
}

func (wal *wallet) GetWalletId() uuid.UUID {
	return wal.wallet_id
}

func (wal *wallet) GetPrivateKey() string {
	return wal.private_key
}
