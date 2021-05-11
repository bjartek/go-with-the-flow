package main

import (
	"fmt"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {

	g := gwtf.NewGoWithTheFlowMainNet()

	/* ekaitza
		Start(13978126).
		End(14013459).
		*/
	eventsFetcher := g.EventFetcher().
		Start(14178514).
		End(14213332).
		Event("A.d796ff17107bbff6.Versus.Bid")

	events, err := eventsFetcher.Run()
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%v", events)
	for _, ev := range events {
		fmt.Printf("%d, %v,%s, %s, %s\n", ev.BlockHeight, ev.Time, ev.Fields["edition"], ev.Fields["bidder"], ev.Fields["price"])
	}
	/*
		for _, ev := range events {
			fields := ev.Fields

			price, err := strconv.ParseFloat(fields["price"], 64)
			if err != nil {
				panic(err)
			}
			if val, ok := bid[fields["bidder"]]; ok {
				bid[fields["bidder"]] = val + price
			} else {
				bid[fields["bidder"]] = price
			}
		}
		type kv struct {
			Key   string
			Value float64
		}

		var ss []kv
		for k, v := range bid {
			ss = append(ss, kv{k, v})
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		for _, kv := range ss {
			fmt.Printf("%s, %v\n", kv.Key, kv.Value)
		}
		/*
			//if you add an eventHook in discord to the flow.json in the example folder you can use this to send the events to a discord channel
			msg, err := eventsFetcher.SendEventsToWebhook("gwtf")
			if err != nil {
				panic(err)
			}
			fmt.Println("send message with id", msg.ID)
	*/

}
