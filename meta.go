package storybot

import (
	dgo "github.com/bwmarrin/discordgo"
)

type Environment struct {
	Token string
	AppId string

	DbConn string

	PathSchema      string
	PathCommand     string
	PathCommandData string

	GuildId        string
	RoleIdAdmin    string
	RoleIdMods     string
	RoleIdEveryone string
}

type Params struct {
	Member      *dgo.Member
	Environment *Environment
	ChoiceList  ChoiceList
	IntData     IntData
	Session     *dgo.Session
	Database    *DB
}

type HandlerCallback func(p *Params, i *dgo.Interaction) error
