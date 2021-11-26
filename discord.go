package storybot

import (
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
)

type (
	// Type aliases for long-ass names.
	Cmd     = *dgo.ApplicationCommand
	Opt     = *dgo.ApplicationCommandOption
	Choice  = *dgo.ApplicationCommandOptionChoice
	IntData = dgo.ApplicationCommandInteractionData
)

type ChoiceList = map[string][]ChoiceItem

type ChoiceItem struct {
	Choice
	Description string
}

func RoleById(rr []*dgo.Role, id string) (r *dgo.Role, err error) {
	for _, r := range rr {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, fmt.Errorf("no such role with ID %s", id)
}
