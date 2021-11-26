package command

import (
	"log"

	sb "storybot"
	"storybot/command/event"
	"storybot/command/pronoun"
	"storybot/command/role"

	dgo "github.com/bwmarrin/discordgo"
)

/*
Return type cannot be an alias due to the way handlers are
validated by discordgo.
*/
func Handler(p *sb.Params) func(s *dgo.Session, i *dgo.InteractionCreate) {

	return func(s *dgo.Session, i *dgo.InteractionCreate) {

		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()

		if i.Type != dgo.InteractionApplicationCommand {
			return
		}

		data := i.ApplicationCommandData()
		p.IntData = data
		p.Member = i.Member
		content := ""
		switch data.Name {
		case "pronoun":
			content = pronoun.Handler(p)
		case "role":
			content = role.Handler(p)
		case "event":
			content = event.Handler(p)
		default:
			content = "No such command."
		}

		resp := dgo.InteractionResponse{
			Type: dgo.InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{Content: content},
		}
		if err := s.InteractionRespond(i.Interaction, &resp); err != nil {
			log.Println(err)
		}
	}
}
