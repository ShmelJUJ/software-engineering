# usage  python3 user/scripts/generate_algorand_testnet_account.py --encrypt true/false
from algosdk import account, mnemonic
import argparse
import base64

parser = argparse.ArgumentParser()

parser.add_argument("--encrypt", type=bool, help="should encrypt private key", required=True)

args = parser.parse_args()

private_key, address = account.generate_account()
print(f"address: {address}")

print(f"private key: {private_key}")
if args.encrypt:
    print(
        f"mnemonic: {base64.b64encode(mnemonic.from_private_key(private_key).encode())}"
    )
else:
    print(f"mnemonic: {mnemonic.from_private_key(private_key)}")
