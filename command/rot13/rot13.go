package rot13

import (
	"strings"

	dgo "github.com/bwmarrin/discordgo"
)

func Handler(data dgo.ApplicationCommandInteractionData) string {
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	content := ""
	clamp := func(i int) int {
		i += 13
		if i > 25 {
			i -= 26
		}
		return i
	}
	for _, r := range data.Options[0].StringValue() {
		if i := strings.IndexRune(lower, r); i >= 0 {
			content += string(lower[clamp(i)])
			continue
		}
		if i := strings.IndexRune(upper, r); i >= 0 {
			content += string(lower[clamp(i)])
			continue
		}
		content += string(r)
	}
	return content
}
