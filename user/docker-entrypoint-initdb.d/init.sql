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

INSERT INTO public.users
(user_id, first_name, last_name, email)
VALUES
('85e6a060-f914-48d1-b73a-23b7e6c81f46', 'Luffy', 'Monkey d.', 'luffy@example.com'),
('65a8ed73-b6f3-4543-82a6-7ab9ef6e9c7b', 'Nami', '', 'nami@example.com'),
('5eb05a37-cee1-46dd-bafc-c9d32ef95d59', 'Roger', 'Gol d.', 'goldroger@example.com');

INSERT INTO public.algorand_wallets
(wallet_id, user_id, public_key, private_key)
VALUES
('3735b92d-5dcb-4dc0-a8b0-54415d0c52d3', '85e6a060-f914-48d1-b73a-23b7e6c81f46','S654A2A7TYORBPQ4GT4VBRZMVOCLT5EXA34RQISUF6TEHDCHC3GDJF7MTI', 'Y2hvb3NlIHByYWN0aWNlIGtpdGNoZW4gaGVuIGJvcmluZyBrbmVlIHZlbnVlIHJlbmV3IHNldHRsZSByZXNpc3Qgc2hhbGxvdyBraW5nZG9tIGhhd2sgZmF0IHBhdGNoIGNhbXAgdHVuYSBmYW5jeSBwaWNuaWMgaW1wYWN0IGluY2ggcHVscCBvbmlvbiBhYnNvcmIgYWNoaWV2ZQ=='),
('c09e7795-8b25-4639-b7e3-8c2592298eb8', '65a8ed73-b6f3-4543-82a6-7ab9ef6e9c7b', '6KVRYUNUMH4VCTQS55KWQQEXOV6QT7FOX6AZ4RFSIYQGPXN2J6X37OMYS4', 'Z2FwIGdyYXZpdHkgYWJvdmUgdW5hYmxlIGVhc3QgaG9tZSB0b3dhcmQgYmlydGggaHVtYW4gc2FkIHJpYiB2aWxsYWdlIHdoZWF0IGluZGV4IHR1cm4gdm9sY2FubyBzdWJqZWN0IGF1dGhvciBncmFjZSBhcm1vciBidW5kbGUgZXRlcm5hbCBtb2RpZnkgYWJzb3JiIGVnZw=='),
('56526826-61d3-49c7-9496-a0efdab27639','5eb05a37-cee1-46dd-bafc-c9d32ef95d59','SCD7M75KMUCJTIDYHLIW374PTJYKEZYTBQH4TBJRHH5X4VHAVIDAV75Z2Q', 'Zm9ydW0gY291cnNlIHNwb25zb3IgYmlrZSB0dXJ0bGUgcm9hZCBtdXNocm9vbSBiZXlvbmQgc2hvY2sgcmVjYWxsIHJlY2FsbCBneW0gc3RlcmVvIHNlYXJjaCBwaWxvdCBndWlkZSBhcnJhbmdlIGVzdGF0ZSBkaXNoIHNlbnNlIHNlZWsgc2VydmljZSBkd2FyZiBhYnNlbnQgc3VtbWVy');