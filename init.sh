KEY="mykey"
CHAINID="plateaus_432-1"
MONIKER="localtestnet"
KEYRING="os"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# Reinstall daemon
rm -rf ~/.plateausd*
make install

# Set client config
plateausd config keyring-backend $KEYRING
plateausd config chain-id $CHAINID

# if $KEY exists it should be deleted
plateausd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO

# Set moniker and chain-id for Plateaus (Moniker can be anything, chain-id must be an integer)
plateausd init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to axrz
cat $HOME/.plateausd/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="axrz"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
cat $HOME/.plateausd/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="axrz"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
cat $HOME/.plateausd/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="axrz"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
cat $HOME/.plateausd/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="axrz"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
cat $HOME/.plateausd/config/genesis.json | jq '.app_state["inflation"]["params"]["mint_denom"]="axrz"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Set gas limit in genesis
cat $HOME/.plateausd/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Disable base fee in genesis
cat $HOME/.plateausd/config/genesis.json | jq '.app_state.feemarket.params.no_base_fee=true' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Disable incentives
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.incentives.params.enable_incentives=false' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.incentives.params.allocation_limit="0"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.incentives.params.reward_scaler="0"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Disable incentives from inflation
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.inflation.params.enable_inflation=false' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.inflation.params.inflation_distribution.staking_rewards="1"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.inflation.params.inflation_distribution.usage_incentives="0"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.inflation.params.inflation_distribution.community_pool="0"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Disable distributions
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.distribution.params.withdraw_addr_enabled=false' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.distribution.params.community_tax="0"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.distribution.params.base_proposer_reward="0"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
#cat $HOME/.plateausd/config/genesis.json | jq '.app_state.distribution.params.bonus_proposer_reward="0"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Set claims start time
node_address=$(plateausd keys list | grep  "address: " | cut -c12-)
current_date=$(date -u +"%Y-%m-%dT%TZ")
cat $HOME/.plateausd/config/genesis.json | jq -r --arg current_date "$current_date" '.app_state["claims"]["params"]["airdrop_start_time"]=$current_date' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Set claims records for validator account
#amount_to_claim=10000
#cat $HOME/.plateausd/config/genesis.json | jq -r --arg node_address "$node_address" --arg amount_to_claim "$amount_to_claim" '.app_state["claims"]["claims_records"]=[{"initial_claimable_amount":$amount_to_claim, "actions_completed":[false, false, false, false],"address":$node_address}]' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Set claims decay
cat $HOME/.plateausd/config/genesis.json | jq -r --arg current_date "$current_date" '.app_state["claims"]["params"]["duration_of_decay"]="1000000s"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json
cat $HOME/.plateausd/config/genesis.json | jq -r --arg current_date "$current_date" '.app_state["claims"]["params"]["duration_until_decay"]="100000s"' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Claim module account:
# 0xA61808Fe40fEb8B3433778BBC2ecECCAA47c8c47 || xrz1567p2ffg0t2453lr8mdxkuudvdl36mz25kww7t
#cat $HOME/.plateausd/config/genesis.json | jq -r --arg amount_to_claim "$amount_to_claim" '.app_state["bank"]["balances"] += [{"address":"xrz1567p2ffg0t2453lr8mdxkuudvdl36mz25kww7t","coins":[{"denom":"axrz", "amount":$amount_to_claim}]}]' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# disable produce empty block
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.plateausd/config/config.toml
  else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.plateausd/config/config.toml
fi

if [[ $1 == "pending" ]]; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.plateausd/config/config.toml
      sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.plateausd/config/config.toml
      sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.plateausd/config/config.toml
      sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.plateausd/config/config.toml
      sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.plateausd/config/config.toml
      sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.plateausd/config/config.toml
      sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.plateausd/config/config.toml
      sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.plateausd/config/config.toml
      sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.plateausd/config/config.toml
  else
      sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.plateausd/config/config.toml
      sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.plateausd/config/config.toml
      sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.plateausd/config/config.toml
      sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.plateausd/config/config.toml
      sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.plateausd/config/config.toml
      sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.plateausd/config/config.toml
      sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.plateausd/config/config.toml
      sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.plateausd/config/config.toml
      sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.plateausd/config/config.toml
  fi
fi

# Allocate genesis accounts (cosmos formatted addresses)
plateausd add-genesis-account $KEY 3333333333330000000000000axrz --keyring-backend $KEYRING

# Update total supply with claim values
validators_supply=$(cat $HOME/.plateausd/config/genesis.json | jq -r '.app_state["bank"]["supply"][0]["amount"]')
# Bc is required to add this big numbers
#total_supply=$(bc <<< "$amount_to_claim+$validators_supply")
total_supply=3333333333330000000000000
cat $HOME/.plateausd/config/genesis.json | jq -r --arg total_supply "$total_supply" '.app_state["bank"]["supply"][0]["amount"]=$total_supply' > $HOME/.plateausd/config/tmp_genesis.json && mv $HOME/.plateausd/config/tmp_genesis.json $HOME/.plateausd/config/genesis.json

# Sign genesis transaction
plateausd gentx $KEY 1000000000000000000000axrz --keyring-backend $KEYRING --chain-id $CHAINID --gas=0xrz

# Collect genesis tx
plateausd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
plateausd validate-genesis

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
plateausd start --pruning=nothing $TRACE --log_level $LOGLEVEL --json-rpc.api eth,txpool,personal,net,debug,web3
