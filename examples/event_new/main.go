package main

import (
	"fmt"
	"strconv"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {
	var startBlock uint64 = 13978126
	const endBlock uint64 = 14013459
	const address = "access.mainnet.nodes.onflow.org:9000"
	const event = "A.d796ff17107bbff6.Versus.Bid"

	events := gwtf.FetchEvents2(address, []string{event}, startBlock, endBlock)

	for _, ev := range events {
		fields := ev.Fields
		edition, _ := strconv.Unquote(fields["edition"])
		fmt.Printf("%s, %d, %s, %s, %s\n", ev.Time.String(), ev.BlockHeight, edition, fields["bidder"], fields["price"])
	}
}
