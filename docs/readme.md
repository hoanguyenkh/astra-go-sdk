# astra-sdk

## Introduction

A Golang SDK for interacting with the Astra Protocol.

## Installation

```go
go get github.com/AstraProtocol/astra-go-sdk
```

## Dependencies

See [go.mod](./../go.mod)


## Import

```go
import (
    github.com/AstraProtocol/astra-go-sdk
)
```

## Tutorials

Following is a series of examples to help you get familiar with the Astra Protocol. The list does not cover all the
capabilities of the SDK, we will try to cover them as much as possible.

* [Introduction](./docs/intro.md)
* [Client](./docs/client.md)
* [Accounts](tutorials/docs/accounts)
    * [Keys](tutorials/docs/accounts/keys.md)
    * [Creating Accounts with HD Wallets](tutorials/docs/accounts/hdwallet_create.md)
    * [Importing Accounts with Mnemonic Strings](tutorials/docs/accounts/hdwallet_import.md)
    * [UTXOs](tutorials/docs/accounts/utxo.md)
        * [Retrieving Output Coins V1](tutorials/docs/accounts/utxo_retrieve.md)
        * [Key Submission](tutorials/docs/accounts/submit_key.md)
        * [UTXO Cache](tutorials/docs/accounts/utxo_cache.md)
        * [Consolidating](tutorials/docs/accounts/consolidate.md)
    * [Account Balances](tutorials/docs/accounts/balances.md)
    * [Account History](tutorials/docs/accounts/tx_history.md)
* [Transactions](tutorials/docs/transactions)
    * [Transaction Parameters](tutorials/docs/transactions/params.md)
    * [Transferring PRV](tutorials/docs/transactions/raw_tx.md)
    * [Transferring Token](tutorials/docs/transactions/raw_tx_token.md)
    * [Initializing Custom Tokens](tutorials/docs/transactions/init_token.md)
    * [Converting UTXOs](tutorials/docs/transactions/convert.md)