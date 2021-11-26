package role

import (
	"bytes"
	"fmt"
	"log"

	sb "storybot"

	dgo "github.com/bwmarrin/discordgo"
)

func Handler(p *sb.Params) string {
	o := p.IntData.Options[0]
	switch o.Name {
	case "add":
		return add(p)
	case "remove":
		return remove(p)
	case "list":
		return list(p)
	}
	return "No such option for role command."
}

func add(p *sb.Params) string {
	m, err := metaFromData(p)
	if err != nil {
		log.Println(err)
		return "Unknown role."
	}
	if hasRole(p.Member, m.id) {
		return "You already have this role."
	}
	e := p.Environment
	err = p.Session.GuildMemberRoleAdd(e.GuildId, p.Member.User.ID, m.id)
	if err != nil {
		log.Println(err)
		return "Unable to assign role."
	}
	return fmt.Sprintf("Added you to `%s` role.", m.name)
}

func remove(p *sb.Params) string {
	m, err := metaFromData(p)
	if err != nil {
		log.Println(err)
		return "Unknown role."
	}
	if !hasRole(p.Member, m.id) {
		return "Cannot remove role you don't have."
	}
	e := p.Environment
	err = p.Session.GuildMemberRoleRemove(e.GuildId, p.Member.User.ID, m.id)
	if err != nil {
		log.Println(err)
		return "Unable to remove role."
	}
	return fmt.Sprintf("Removed you from `%s` role.", m.name)
}

func list(p *sb.Params) string {
	// Not yet sure whether this is bytes or code points so
	// we'll be pessimistic and assume it's bytes.
	maxMsgLen := 2000
	buf := bytes.Buffer{}
	o := p.IntData.Options[0] // list option
	choice := o.Options[0]    // which option
	which := choice.StringValue()
	for _, c := range p.ChoiceList[p.IntData.Name] {
		roleId := c.Choice.Value.(string)
		if buf.Len() > maxMsgLen {
			return "List of roles would exceed maximum message length."
		}
		if which == "self" && !hasRole(p.Member, roleId) {
			continue
		}
		buf.WriteString(fmt.Sprintf("**%s**\n", c.Name))
		buf.WriteString(fmt.Sprintf("*%s*\n\n", c.Description))
	}
	if buf.Len() == 0 {
		switch which {
		case "self":
			return "You currently have no roles. Use `/role list all` to see names and descriptions of each role."
		case "all":
			return "There are currently no roles that can be assigned."
		}
	}
	return buf.String()
}

func hasRole(m *dgo.Member, id string) bool {
	for i := range m.Roles {
		if m.Roles[i] == id {
			return true
		}
	}
	return false
}

func metaFromData(p *sb.Params) (meta *roleMeta, err error) {
	o := p.IntData.Options[0]
	rr, err := p.Session.GuildRoles(p.Environment.GuildId)
	if err != nil {
		return nil, err
	}
	choice := o.Options[0]
	roleId := choice.StringValue()
	role, err := sb.RoleById(rr, roleId)
	if err != nil {
		return nil, err
	}
	return &roleMeta{
		id:   roleId,
		name: role.Name,
	}, nil
}

type roleMeta struct {
	id   string
	name string
}
