[![Coverage Status](https://coveralls.io/repos/github/bjartek/go-with-the-flow/badge.svg?branch=main)](https://coveralls.io/github/bjartek/go-with-the-flow?branch=main) [![ci](https://github.com/bjartek/go-with-the-flow/actions/workflows/ci.yml/badge.svg)](https://github.com/bjartek/go-with-the-flow/actions/workflows/ci.yml)

# Go with the Flow

> Tooling to help develop application on the the Flow Blockchain

Set of go scripts to make it easer to run a story consisting of creating accounts,
deploying contracts, executing transactions and running scripts on the Flow Blockchain.
These go scripts also make writing integration tests of your smart contracts much easier.

## Information

### Main features

- Create a single go file that will start emulator, deploy contracts, create accounts and run scripts and transactions. see `examples/demo/main.go`
- Fetch events, store progress in a file and send results to Discord. see `examples/event/main.go`
- Support inline scripts if you do not want to store everything in a file when testing
- Supports writing tests against transactions and scripts with some limitations on how to implement them.
- Asserts to make it easier to use the library in writing tests see `examples/transaction_test.go` for examples

### Gotchas

- When specifying extra accounts that are created on emulator they are created in alphabetical order, the addresses the emulator assign is always fixed.
- tldr; Name your stakeholder accounts in alphabetical order

### Note on v2

v2 of GoWithTheFlow removed a lot of the code in favor of `flowkit` in the flow-cli. Some of the code from here was
contributed by me into flow-cli like the goroutine based event fetcher.

Breaking changes between v1 and v2:

- v1 had a config section for discord webhooks. That has been removed since the flow-cli will remove extra config things in flow.json. Store the webhook url in an env variable and use it as argument when creating the DiscordWebhook struct.

Special thanks to @sideninja for helping me get my changes into flow-cli. and for jayShen that helped with fixing some issues!

## Resources

- Run the demo example in this project with `cd example && make`. The emulator will be run in memory.
- Check [other codebases](https://github.com/bjartek/go-with-the-flow/network/dependents?package_id=UGFja2FnZS0yMjc1NjE0OTAz) that use this project
- Feel free to ask questions to @bjartek in the Flow Discord.

## Credits

This project is a rewrite of https://github.com/versus-flow/go-flow-tooling

## Todo

- fix golangci-rules, disabled a bunch to get the build green
