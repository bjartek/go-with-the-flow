package gwtf

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

//EventsToWebhookParams convert events to rich webhook
func EventsToWebhookParams(events []*FormatedEvent) *discordgo.WebhookParams {
	embeds := []*discordgo.MessageEmbed{}
	for _, event := range events {

		fields := []*discordgo.MessageEmbedField{}
		for name, value := range event.Fields {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  name,
				Value: value,
			})
		}

		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:  event.Name,
			Type:   discordgo.EmbedTypeRich,
			Fields: fields,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("blockHeight %d @ %s", event.BlockHeight, event.Time),
			},
		})
	}

	return &discordgo.WebhookParams{
		Embeds: embeds,
	}
}
