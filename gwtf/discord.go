package gwtf

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)



// SendEventsToWebhook Sends events to the webhook with the given name from flow.json
func (dw DiscordWebhook) SendEventsToWebhook(events []*FormatedEvent) (*discordgo.Message, error) {

	discord, err := discordgo.New()
	if err != nil {
		return nil, err
	}

	status, err := discord.WebhookExecute(
		dw.ID,
		dw.Token,
		dw.Wait,
		EventsToWebhookParams(events))

	if err != nil {
		return nil, err
	}
	return status, nil
}

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
