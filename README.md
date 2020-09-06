# Go~with~the~Flow

Set of go scripts to make it easer to run Story constisting of creating accounts, deploying contracts, executing transactions and running scripts on the Flow Blockchain.

Feel free to ask questions to @bjartek or @0xAlchemist in the Flow Discord.

## How to configure the flow.json file
The flow json file for working on devnet should look something like the file you have in the example folder. 

Right now the examples are only creating two accounts/contracts but if you need more feel free to rename the 3-10 keys to suit your needs. Note that the emulator will always create the accounts in that order given that service account. 

The service account is the same as the one used in vscode extension 0.9 so that you will get syntax highlighting in vscode.

## How to configure devnet
Create a ~/.flow-dev.json file that has fully qualified accounts for all your deployed contracts addresses. Make sure that the address is pointed to the latest working devnet node. 

## Examples

1. `cd examples/`
2. make emulator in terminal 1
3. make in terminal 2

## Credits

This project is a rewrite of https://github.com/versus-flow/go-flow-tooling 
