# Go with the Flow

Right now this tool required version v0.11.2 of flow-cli. Get using
```sh -ci  "$(curl -fsSL https://storage.googleapis.com/flow-cli/install.sh)" -- v0.11.2```

Set of go scripts to make it easer to run a Story consisting of creating accounts, deploying contracts, executing transactions and running scripts on the Flow Blockchain.

Feel free to ask questions to @bjartek or @0xAlchemist in the Flow Discord.

## How to configure the flow.json file
The flow json file for working on emulator should look something like the file you have in the example folder. 

Right now the examples are only creating one account for all contracts and two user accounts but if you need more feel free to rename the 4-10 keys to suit your needs. Note that the emulator will always create the accounts in that order given that service account. 

The service account is the same as the one used in vscode extension so that you will get syntax highlighting in vscode.

## How to configure devnet
Create a ~/.flow-dev.json file that has fully qualified accounts for all your deployed contracts addresses. Make sure that the address is pointed to the latest working devnet node. 

## How to write contracts
Put contracts in the contract folder and name the file the same as the contract.

## Examples

1. `cd examples/`
2. make emulator in terminal 1
3. make in terminal 2

## Credits

This project is a rewrite of https://github.com/versus-flow/go-flow-tooling 
