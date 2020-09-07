#!/usr/bin/env python3

from base64 import b64encode
from os import urandom


def main():
    random_bytes = urandom(16)
    print(b64encode(random_bytes).decode("utf-8"))


if __name__ == "__main__":
    main()
