#! /usr/bin/env python3

from Crypto.Cipher import AES
import argparse
import base64
from datetime import datetime
from datetime import timedelta
import json


def _parse_header(arg):
    raw_key, value = arg.split("=")
    return (raw_key.title(), value)


def encode(data):
    return base64.urlsafe_b64encode(data).decode("utf8").replace("\n", "")


def encrypt(data: str, expiry: datetime, headers: dict, key: bytes):
    cipher = AES.new(key, AES.MODE_GCM)
    nonce = cipher.nonce
    ciphertext, tag = cipher.encrypt_and_digest(f"{data};{expiry.isoformat()};{json.dumps(headers)}".encode("utf-8"))
    return encode(nonce), encode(ciphertext + tag)


def main():
    parser = argparse.ArgumentParser(description="Generate an encrypted link.")
    parser.add_argument("lifetime", type=int, help="How many seconds the link should be open for. (Starting now)")
    parser.add_argument("key", type=str, help="Base 64 encoded key. Should be 16, 24 or 32 bytes.")
    parser.add_argument(
        "URL", type=str, help="The URL to encode.",
    )
    parser.add_argument(
        "--headers",
        metavar="KEY=VALUE",
        nargs="*",
        help="Set a number of key-value pairs "
        "(do not put spaces before or after the = sign). "
        "If a value contains spaces, you should define "
        "it with double quotes: "
        'foo="this is a sentence". Note that '
        "values are always treated as strings.",
        default=[],
    )
    args = parser.parse_args()
    headers = dict(map(_parse_header, args.headers))
    expiry = datetime.utcnow() + timedelta(seconds=args.lifetime)
    username, password = encrypt(args.URL, expiry, headers, base64.b64decode(args.key))
    print(f"{username}:{password}@localhost:7777")


if __name__ == "__main__":
    main()
