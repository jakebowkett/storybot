package command

import (
	"fmt"

	sb "storybot"

	dgo "github.com/bwmarrin/discordgo"
)

func Create(p *sb.Params) error {
	s := p.Session
	e := p.Environment
	var cf cmdParts
	if err := loadAndUnmarshal(e.PathCommand, &cf); err != nil {
		return err
	}
	if err := loadAndUnmarshal(e.PathCommandData, &cf); err != nil {
		return err
	}
	if err := merge(cf); err != nil {
		return err
	}
	for _, i := range cf.Insert {
		p.ChoiceList[i.Root] = i.Choices
	}

	cc, err := s.ApplicationCommandBulkOverwrite(
		e.AppId,
		e.GuildId,
		cf.dgoCommands(),
	)
	if err != nil {
		return fmt.Errorf("unable to bulk overwrite: %w", err)
	}

	var perm []*dgo.GuildApplicationCommandPermissions
	for _, c := range cf.Command {
		cmdId, ok := cmdIdByName(cc, c.Name)
		if !ok {
			return fmt.Errorf("no command with name %q", c.Name)
		}
		perm = append(perm, &dgo.GuildApplicationCommandPermissions{
			ID: cmdId,
			Permissions: []*dgo.ApplicationCommandPermissions{
				{
					ID:         e.RoleIdEveryone,
					Type:       dgo.ApplicationCommandPermissionTypeRole,
					Permission: c.AnyoneCanUse,
				},
				{
					ID:         e.RoleIdMods,
					Type:       dgo.ApplicationCommandPermissionTypeRole,
					Permission: true,
				},
			},
		})
	}

	err = s.ApplicationCommandPermissionsBatchEdit(e.AppId, e.GuildId, perm)
	if err != nil {
		return fmt.Errorf("unable to batch edit permissions: %w", err)
	}

	return nil
}

func cmdIdByName(cc []sb.Cmd, name string) (id string, ok bool) {
	for _, c := range cc {
		if c.Name == name {
			return c.ID, true
		}
	}
	return "", false
}
