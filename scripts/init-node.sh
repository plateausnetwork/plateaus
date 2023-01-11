NODENAME=`echo $RANDOM | md5sum | head -c 10; echo;`
VALIDATORKEY="val-node-$NODENAME"
CHAINID="plateaus_432-1"
MONIKER="node-$NODENAME"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
LOCALTESTNET="`pwd`/localnet-setup/"
# to trace evm
#TRACE="--trace"
TRACE=""

#PEERS="c3947db1b922f3be28e038a486308b0c52f41dc0@127.0.0.1:26656"
PEERS=`curl http://localhost:26657/genesis | jq '.result.genesis.app_state.genutil.gen_txs[0].body.memo'`
GENESISJSON=`curl http://localhost:26657/genesis | jq '.result.genesis'`

if [ ! $PEERS ]; then
  echo "peers cannot be empty"
  exit
fi

if [ ! $GENESISJSON ]; then
  echo "genesis cannot be empty"
  exit
fi

## creating a localnet folder
rm -rf `pwd`/localnet-setup && mkdir $LOCALTESTNET

## Set client config
plateausd config keyring-backend $KEYRING --home=$LOCALTESTNET
plateausd config chain-id $CHAINID --home=$LOCALTESTNET

plateausd tendermint unsafe-reset-all --home=$LOCALTESTNET

## if $KEY exists it should be deleted
plateausd keys add $VALIDATORKEY --keyring-backend $KEYRING --algo $KEYALGO --home=$LOCALTESTNET

## Set moniker and chain-id for Evmos (Moniker can be anything, chain-id must be an integer)
plateausd init $MONIKER --chain-id $CHAINID --home=$LOCALTESTNET

# getting the genesis.json
GENESISJSON=`curl http://localhost:26657/genesis | jq '.result.genesis'`

# setting genesis.json
echo "copying genesis.json"

echo $GENESISJSON > $LOCALTESTNET/config/genesis.json

sleep 1

# setting peers from config.toml
#peers="2eaff30af3ef7c05b27a966826aadc8ccb86db18@127.0.0.1:26656"
#peers="d50c54142ca612af9146f3c9f6c6702fbde05d97@127.0.0.1:26656"
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = $PEERS/" $LOCALTESTNET/config/config.toml

## Run this to ensure everything worked and that the genesis file is setup correctly
plateausd validate-genesis --home=$LOCALTESTNET

## Start the node (remove the --pruning=nothing flag if historical queries are not needed)
plateausd start --pruning=default $TRACE --log_level $LOGLEVEL --home=$LOCALTESTNET

echo "plateausd start --pruning=default $TRACE --log_level $LOGLEVEL --home=$LOCALTESTNET"
