package vote

import (
	"fmt"
	"strings"

	sb "storybot"

	dgo "github.com/bwmarrin/discordgo"
)

const (
	maxChoices   = 9
	maxChoiceLen = 64
)

func Handler(p *sb.Params) (string, sb.HandlerCallback) {
	s := p.IntData.Options[0].StringValue()
	ss := strings.Split(s, ",")
	if len(ss) > maxChoices {
		return fmt.Sprintf("Too many choices. Have %d, maximum %d", len(ss), maxChoices), nil
	}
	for i := range ss {
		ss[i] = fmt.Sprintf("`%d.` ", i+1) + strings.TrimSpace(ss[i])
		if len(ss[i]) > maxChoiceLen {
			return fmt.Sprintf("Choice %d is too long. It has %d characters, maximum is %d", i+1, len(ss[i]), maxChoiceLen), nil
		}
	}
	return strings.Join(ss, "\n"), callback(len(ss))
}

func callback(n int) sb.HandlerCallback {
	return func(p *sb.Params, i *dgo.Interaction) error {
		s := p.Session
		e := p.Environment
		msg, err := s.InteractionResponse(e.AppId, i)
		if err != nil {
			return fmt.Errorf("unable to acquire interaction response: %w", err)
		}
		for i := 0; i < n; i++ {
			err := s.MessageReactionAdd(msg.ChannelID, msg.ID, num[i])
			if err != nil {
				return fmt.Errorf("unable to add reaction: %w", err)
			}
		}
		return nil
	}
}

var num = []string{
	"1️⃣",
	"2️⃣",
	"3️⃣",
	"4️⃣",
	"5️⃣",
	"6️⃣",
	"7️⃣",
	"8️⃣",
	"9️⃣",
}
