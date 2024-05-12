package client

const (
	kGetClientByid = `SELECT user_id, first_name, last_name, email, user_password
						FROM public.users
						WHERE user_id = $1::UUID;`
	kGetClientByEmail = `SELECT user_id, first_name, last_name, email, user_password
						FROM public.users
						WHERE email = $1::TEXT;`

	kGetWalletByUserId = `SELECT wallet_id, user_id, public_key, private_key
							FROM public.algorand_wallets
							WHERE user_id = $1::UUID`
)
