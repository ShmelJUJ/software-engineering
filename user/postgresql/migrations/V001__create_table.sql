CREATE TABLE IF NOT EXISTS public.users (
    user_id UUID NOT NULL PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    email TEXT
);

CREATE TABLE IF NOT EXISTS public.algorand_wallets(
    wallet_id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    public_key TEXT,
    private_key TEXT,
    FOREIGN KEY (user_id) REFERENCES public.users (user_id)
);