package event

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	sb "storybot"
)

type event struct {
	Public int64
	Added  string
	Slug   string

	// Fields returned by StoryDevs API.
	Name     string        `json:"name"`
	Summary  string        `json:"summary"`
	Start    int64         `json:"start"`
	Finish   *int64        `json:"finish"`
	Weekly   bool          `json:"weekly"`
	Local    bool          `json:"local"`
	Category sb.DbSliceStr `json:"category"`
	Setting  sb.DbSliceStr `json:"setting"`
}

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

var isEvt = regexp.MustCompile(`^https://storydevs\.com/event/([A-Za-z0-9]{11})$`)

func add(p *sb.Params) string {
	generalErr := "Unable to add event."
	o := p.IntData.Options[0]
	url := o.Options[0].StringValue()
	matches := isEvt.FindStringSubmatch(url)
	if len(matches) == 0 {
		return "Invalid URL."
	}
	slug := matches[1]
	resp, err := http.Get("https://storydevs.com/api/event/" + slug)
	if err != nil {
		log.Println(err)
		return generalErr
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "Couldn't find event."
	}
	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return generalErr
	}
	evt := event{}
	if err := json.Unmarshal(bb, &evt); err != nil {
		log.Println(err)
		return generalErr
	}
	evt.Slug = slug
	evt.Added = p.Member.User.ID

	tx, err := p.Database.Beginx()
	if err != nil {
		log.Println(err)
		return generalErr
	}

	exists, err := tx.Exists(`FROM event WHERE slug = $1`, slug)
	if err != nil {
		fmt.Printf("problem while checking if event exists: %s\n", err)
		return generalErr
	}
	if exists {
		return "Event is already scheduled."
	}

	err = tx.Get(&evt.Public, `
		SELECT
			generate_series FROM GENERATE_SERIES(
				(SELECT MIN(public) FROM event),
				(SELECT MAX(public) FROM event)
			)
		WHERE
			NOT EXISTS(
				SELECT
					public
				FROM
					event
				WHERE
					public = generate_series
			)
	`)
	/*
		There are three potential outcomes to this query:
			- No rows "error"; this can happen if:
				- There's all sequential rows OR
				- there's no rows in the table.
					- In either case, take max and add one.
			- An actual error.
			- There were rows with gaps between public IDs.
				- In this case we use the public ID as populated from the query.
	*/
	switch {
	case errors.Is(err, sql.ErrNoRows):
		err = tx.Get(&evt.Public, `SELECT MAX(public) FROM event`)
		if err != nil {
			fmt.Printf("problem with event max public ID query: %s\n", err)
			return generalErr
		}
		evt.Public++
	case err != nil:
		fmt.Printf("problem with event generate_series query: %s\n", err)
		return generalErr
	}

	err = tx.Insert("event", &evt, []string{
		"public",
		"added",
		"slug",
		"name",
		"summary",
		"start",
		"finish",
		"weekly",
		"local",
		"category",
		"setting",
	})
	if err != nil {
		log.Println(err)
		return generalErr
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return generalErr
	}

	return evt.format("Scheduled event. Members with the `event` role will now be pinged when this event begins.")
}

func remove(p *sb.Params) string {
	return "Not implemented."
}

func list(p *sb.Params) string {
	return "Not implemented."
}
