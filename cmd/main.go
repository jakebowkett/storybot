package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	sb "storybot"
	"storybot/command"

	"github.com/BurntSushi/toml"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	env = flag.String("e", "", "TOML file containing environment variables.")
	cmd = flag.Bool("c", false, "Create commands. Overwrites previous commands.")
)

func main() {

	flag.Parse()

	e, err := loadEnv(*env)
	if err != nil {
		log.Fatalf("unable to load environment: %s", err)
	}

	db, err := setupDB(e.DbConn, e.PathSchema)
	if err != nil {
		log.Fatalf("unable to set up database: %s", err)
	}

	s, err := dgo.New("Bot " + e.Token)
	if err != nil {
		log.Fatalf("unable to create discordgo session: %s", err)
	}
	defer s.Close()

	p := &sb.Params{
		Session:     s,
		Environment: e,
		ChoiceList:  make(sb.ChoiceList),
		Database:    db,
	}
	/*
		ChoiceList is populated by command.Create - we need an
		alternate way of populating it for cases where the -c
		flag isn't passed.
	*/
	if *cmd {
		if err = command.Create(p); err != nil {
			log.Fatalf("unable to create commands: %s", err)
		}
	}

	s.Identify.Intents = dgo.IntentsGuildMessages

	_ = s.AddHandler(command.Handler(p))

	if err := s.Open(); err != nil {
		log.Fatalf("unable to open session: %s", err)
	}

	println("StoryBot online.")

	end := make(chan os.Signal, 1)
	signal.Notify(end, os.Interrupt)
	<-end
}

func loadEnv(path string) (e *sb.Environment, err error) {
	bb, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load file: %w", err)
	}
	e = &sb.Environment{}
	if err := toml.Unmarshal(bb, &e); err != nil {
		return nil, fmt.Errorf("unable to unmarshal: %w", err)
	}
	return e, nil
}

func setupDB(connStr, schema string) (*sb.DB, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	bb, err := os.ReadFile(schema)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(bb)); err != nil {
		return nil, err
	}
	return &sb.DB{DB: db}, nil
}
