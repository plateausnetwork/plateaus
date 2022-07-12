#!/bin/bash

KEY="mykey"
CHAINID="plateaus_9000-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t plateaus-datadir.XXXXX)

echo "create and add new keys"
./plateausd keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init Plateaus with moniker=$MONIKER and chain-id=$CHAINID"
./plateausd init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./plateausd add-genesis-account \
"$(./plateausd keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000aplateaus,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./plateausd gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./plateausd collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./plateausd validate-genesis --home $DATA_DIR

echo "starting plateaus node $i in background ..."
./plateausd start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started plateaus node"
tail -f /dev/null