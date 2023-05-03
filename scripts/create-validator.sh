if ! command -v jq > /dev/null 2>&1; then
  echo >&2 "jq not installed."
  exit 1
fi

FROM_NAME="node-6baf715cfc"
EXTERNAL_KEY_PATH="/home/yurigiovani/review/priv-key.txt"

if [ ! $FROM_NAME ]; then
    echo "from name cannot be empty"
    exit
fi

if [ ! $EXTERNAL_KEY_PATH ]; then
    echo "external key path cannot be empty"
    exit
fi

plateausd tx staking create-validator --amount=1xrx --from="$FROM_NAME" --pubkey=`plateausd tendermint show-validator` --commission-rate="0.10" --commission-max-rate="0.20" --commission-max-change-rate="0.01" --min-self-delegation="1" --external-key-path="$EXTERNAL_KEY_PATH"