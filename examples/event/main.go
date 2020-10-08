package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bjartek/go-with-the-flow/gwtf"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func main() {

	webhookID := ""    //the id of the webhhok
	webhookToken := "" //the auth token for the webhook

	g := gwtf.NewGoWithTheFlowDevNet()

	ctx := context.Background()

	c, err := client.New(g.Address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	header, err := c.GetLatestBlockHeader(ctx, true)
	if err != nil {
		panic(err)
	}

	startHeight := header.Height - 10000
	events, err := g.FetchEvents(client.EventRangeQuery{
		Type:        "flow.AccountCreated",
		StartHeight: startHeight,
		EndHeight:   header.Height,
	})

	if err != nil {
		panic(err)
	}

	prettyJSON, err := json.MarshalIndent(&events, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Events emitted: %s\n", string(prettyJSON))

	discord, err := discordgo.New()
	if err != nil {
		panic(err)
	}

	status, err := discord.WebhookExecute(
		webhookID,
		webhookToken,
		false,
		gwtf.EventsToWebhookParams(events))

	if err != nil {
		panic(err)
	}
	spew.Dump(status)
}
