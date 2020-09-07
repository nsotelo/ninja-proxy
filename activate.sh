#! /usr/bin/env bash

python3 -m venv .venv
source .venv/bin/activate
pip install ./utilities/

echo -e "Run \e[32m\e[1mninja-key\e[0m to get a new key."
echo -e "Run \e[32m\e[1mninja-link\e[0m to generate a temporary link (use -h for more information)"
