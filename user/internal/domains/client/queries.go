package client

const (
	kGetClientByid = `SELECT user_id, first_name, last_name, email
						FROM public.users
						WHERE user_id = $1::UUID;`

	kGetWalletByUserId = `SELECT wallet_id, user_id, public_key, private_key
							FROM public.wallets
							WHERE user_id = $1::UUID`
)
