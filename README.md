[![Coverage Status](https://coveralls.io/repos/github/bjartek/go-with-the-flow/badge.svg?branch=main)](https://coveralls.io/github/bjartek/go-with-the-flow?branch=main) [![ci](https://github.com/bjartek/go-with-the-flow/actions/workflows/test.yml/badge.svg)](https://github.com/bjartek/go-with-the-flow/actions/workflows/test.yml)

# Go with the Flow

Set of go scripts to make it easer to run a Story consisting of creating accounts, deploying contracts, executing transactions and running scripts on the Flow Blockchain.

Feel free to ask questions to @bjartek in the Flow Discord.

v2 of GoWithTheFlow removed a lot of the code in favor of `flowkit` in the flow-cli. Some of the code from here was
contributed by me into flow-cli like the goroutine based event fetcher.

Breaking changes between v1 and v2:
 - v1 had a config section for discord webhooks. That has been removed since the flow-cli will remove extra config things in flow.json. Store the webhook url in an env variable and use it as argument when creating the DiscordWebhook struct.

Special thanks to @sideninja for helping me get my changes into flow-cli.

## Main features
 - Create a single go file that will start emulator, deploy contracts, create accounts and run scripts and transactions. see `examples/demo/main.go` 
 - Fetch events, store progress in a file and send results to Discord. see `examples/event/main.go`
 - Support inline scripts if you do not want to sture everything in a file when testing 
 - Supports writing tests against transactions and scripts with some limitations on how to implement them. 
 - Asserts to make it easier to use the library in writing tests see `examples/transaction_test.go` for examples

## Gotchas
 - When specifying extra accounts that are created on emulator they are created in alphabetical order, the addresses the emulator assign is always fixed. 
 - tldr; Name your stakeholder acounts in alphabetical order

## Examples

In order to run the demo example you only have to run `make` in the example folder of this project. 
The emulator will be run in memory. 

## Integrate go-with-the-flow into an existing project

You may want to ingetrate go-with-the-flow into an existing project. If this is the case, and if you already have a flow.json at the root of your project, note that go-with-the-flow references accounts from that file by prepending the network to the account name. See the [examples flow.json](https://github.com/bjartek/go-with-the-flow/blob/main/examples/flow.json) for an example. This is done so that it is easy to run a storyline against emulator, tesnet and mainnet. This can be disabled with the `DoNotPrependNetworkToAccountNames` method applied to `gwtf.NewGoWithTheFlowInMemoryEmulator()`. To use go-with-the-flow with your existing flow.json file you need to prepend any given account with either `emulator-`, `testnet-` or `mainnet-` in the flow.json where you define the accounts.

Once you have your flow.json ready, the rest is simple:
- in the root of your project, run `go mod init [NAME]` where '[NAME]' can simply be the github path to your project, ex. "github.com/myGithubHandle/project-name". This will create a 'go.mod' file.
- from the root of your project run `go get github.com/bjartek/go-with-the-flow/v2/gwtf`. This will create a 'go.sum' file and update your 'go.mod' file.
- make a new folder to house your go-with-the-flow scripts either in the root of your project, or in a subfolder. If using a sub-folder, make sure your Makefile (the next step) tests point to it correctly.
- add a `Makefile` to your root directory (see the [examples](https://github.com/bjartek/go-with-the-flow/blob/main/examples) for a reference)

## Credits

This project is a rewrite of https://github.com/versus-flow/go-flow-tooling 
