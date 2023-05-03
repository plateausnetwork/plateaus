if ! command -v jq > /dev/null 2>&1; then
  echo >&2 "jq not installed."
  exit 1
fi

NODENAME=`echo $RANDOM | md5sum | head -c 10; echo;`
VALIDATORKEY="node-$NODENAME"
CHAINID="plateaus_432-1"
MONIKER="node-$NODENAME"
KEYRING="file"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
LOCALTESTNET="$HOME/.plateausd"
# to trace evm
#TRACE="--trace"
TRACE=""

# Check if the folder exists
if test -d $LOCALTESTNET; then
    echo "the folder $LOCALTESTNET exists."
    exit
fi

PEERS=""

if [ ! PEERS ]; then
  echo "peers cannot be empty"
  exit
fi

GENESISJSON=`curl http://sentry-nodes.rhizom.me:26657/genesis | jq '.result.genesis'`

if [ ! $GENESISJSON ]; then
  echo "genesis cannot be empty"
  exit
fi

## Set client config
plateausd config keyring-backend $KEYRING --home=$LOCALTESTNET
plateausd config chain-id $CHAINID --home=$LOCALTESTNET

plateausd tendermint unsafe-reset-all --home=$LOCALTESTNET

## if $KEY exists it should be deleted
plateausd keys add $VALIDATORKEY --keyring-backend $KEYRING --algo $KEYALGO --home=$LOCALTESTNET

## Set moniker and chain-id for Evmos (Moniker can be anything, chain-id must be an integer)
plateausd init $MONIKER --chain-id $CHAINID --home=$LOCALTESTNET

# setting genesis.json
echo "copying genesis.json"

echo $GENESISJSON > $LOCALTESTNET/config/genesis.json

sleep 1

# setting peers from config.toml
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = \"$PEERS\"/" $LOCALTESTNET/config/config.toml

## Run this to ensure everything worked and that the genesis file is setup correctly
plateausd validate-genesis --home=$LOCALTESTNET

## Start the node (remove the --pruning=nothing flag if historical queries are not needed)
#plateausd start --pruning=default $TRACE --log_level $LOGLEVEL --home=$LOCALTESTNET

echo "plateausd start --pruning=default $TRACE --log_level $LOGLEVEL --home=$LOCALTESTNET"
