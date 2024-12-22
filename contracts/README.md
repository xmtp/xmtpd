# XMTP Contracts

- [XMTP Contracts](#xmtp-contracts)
  - [Messages Contracts](#messages-contracts)
  - [XMTP Node Registry](#xmtp-node-registry)
  - [Usage](#usage)
    - [Prerequisites](#prerequisites)
    - [Install](#install)
    - [Test](#test)
    - [Run static analysis](#run-static-analysis)
  - [Scripts](#scripts)
    - [Messages contracts](#messages-contracts-1)
    - [Node registry](#node-registry)

**⚠️ Experimental:** This software is in early development. Expect frequent changes and unresolved issues.

This repository contains all the smart contracts that underpin the XMTP decentralized network.

## Messages Contracts

The messages contracts manage the blockchain state for `GroupMessages` and `IdentityUpdates` sent by clients to the network.

These contracts ensure transparency and provide a historical record of state changes.

## XMTP Node Registry

The `XMTP Node Registry` maintains a blockchain-based record of all node operators participating in the XMTP network. This registry serves as a source of truth for the network's active node participants, contributing to the network's integrity.

The registry is currently implemented following the [ERC721](https://eips.ethereum.org/EIPS/eip-721) standard.

## Usage

The project is built with the `Foundry` framework, and dependency management is handled using `soldeer`.

Additionally, it uses `slither` for static analysis.

### Prerequisites

[Install foundry](https://book.getfoundry.sh/getting-started/installation)

[Install slither](https://github.com/crytic/slither?tab=readme-ov-file#how-to-install)

### Install

As the project uses `soldeer`, update the dependencies by running:

```shell
forge soldeer update
```

Build the contracts:

```shell
forge build
```

### Test

To run the unit tests:

```shell
forge test
```

### Run static analysis

Run the analysis with `slither`:

```shell
slither .
```

## Scripts

The project includes deployer and upgrade scripts.

### Messages contracts

- Configure the environment by creating an `.env` file, with this content:

```shell
### Main configuration
PRIVATE_KEY=0xYourPrivateKey # Private key of the EOA deploying the contracts

### XMTP deployment configuration
XMTP_GROUP_MESSAGES_ADMIN_ADDRESS=0x12345abcdf # the EOA assuming the admin role in the GroupMessages contract.
XMTP_IDENTITY_UPDATES_ADMIN_ADDRESS=0x12345abcdf # the EOA assuming the admin role in the IdentityUpdates contract.
```

- Run the desired script with:

```shell
forge script --rpc-url <RPC_URL> --broadcast <PATH_TO_SCRIPT>
```

Example:

```shell
forge script --rpc-url http://localhost:7545 --broadcast script/DeployGroupMessages.s.sol
```

The scripts output the deployment and upgrade in the `output` folder.

### Node registry

**⚠️:** The node registry hasn't been fully migrated to forge scripts.

- Deploy with `forge create`:

```shell
forge create --broadcast --legacy --json --rpc-url $DOCKER_RPC_URL --private-key $PRIVATE_KEY "src/Nodes.sol:Nodes"
```
