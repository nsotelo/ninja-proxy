#! /usr/bin/env python3

from Crypto.Cipher import AES
import argparse
import base64
from datetime import datetime
from datetime import timedelta
import os

KEY = b"Sixteen byte key"
EXPIRY = datetime.utcnow() + timedelta(minutes=5)


def encode(data):
    return base64.urlsafe_b64encode(data).decode("utf8").replace("\n", "")


def encrypt(data: str, expiry: datetime):
    cipher = AES.new(KEY, AES.MODE_GCM)
    nonce = cipher.nonce
    ciphertext, tag = cipher.encrypt_and_digest(f"{data};{expiry.isoformat()}".encode("utf-8"))
    return encode(nonce), encode(ciphertext + tag)


def test():
    url = "http://username:password@example.com:12345"
    nonceRaw, encryptedRaw = encrypt(url, EXPIRY)
    encrypted = base64.urlsafe_b64decode(encryptedRaw)
    ciphertext, tag = encrypted[:-16], encrypted[-16:]
    cipher = AES.new(KEY, AES.MODE_GCM, nonce=base64.urlsafe_b64decode(nonceRaw))
    plaintext = cipher.decrypt(ciphertext)
    plain_url, expiry = plaintext.decode("utf8").split(";")
    cipher.verify(tag)
    return url == plain_url and expiry == EXPIRY.isoformat()


def main():
    parser = argparse.ArgumentParser(description="Generate an encrypted link.")
    parser.add_argument(
        "lifetime", type=int, nargs="1", help="How many seconds the link should be open for. (Starting now)"
    )
    parser.add_argument(
        "URL", type=str, nargs="1", help="The URL to encode.",
    )
    args = parser.parse_args()
    expiry = datetime.utcnow() + timedelta(seconds=args.lifetime)
    username, password = encrypt(args.URL, expiry)
    print(f"{username}:{password}@localhost:7777")


if __name__ == "__main__":
    main()
