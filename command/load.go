package command

import (
	"fmt"
	"os"
	"strings"

	sb "storybot"

	"github.com/BurntSushi/toml"
)

func loadAndUnmarshal(path string, data interface{}) error {
	bb, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to load file: %w", err)
	}
	if err := toml.Unmarshal(bb, data); err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}
	return nil
}

type cmdParts struct {
	Command []cmd
	Insert  []insert
}

func (cp cmdParts) dgoCommands() (cc []sb.Cmd) {
	for _, c := range cp.Command {
		cc = append(cc, c.Cmd)
	}
	return cc
}

type cmd struct {
	sb.Cmd
	AnyoneCanUse bool
}

type insert struct {
	Root    string
	At      []string
	Choices []sb.ChoiceItem
}

func (i insert) dgoChoices() (ch []sb.Choice) {
	for _, c := range i.Choices {
		ch = append(ch, c.Choice)
	}
	return ch
}

func merge(cf cmdParts) error {
	for _, ins := range cf.Insert {
		ch := ins.dgoChoices()
		for _, ia := range ins.At {
			path := strings.Split(ia, ".")
			if len(path) < 2 {
				return fmt.Errorf("path must be at least two segments")
			}
			if ok := insertIntoCommands(cf.Command, ch, path); !ok {
				return fmt.Errorf("could not find command %q", ia)
			}
		}
	}
	return nil
}

func insertIntoCommands(cc []cmd, ch []sb.Choice, path []string) (ok bool) {
	for _, c := range cc {
		if path[0] == c.Name {
			return insertIntoSubcommands(c.Options, ch, path[1:])
		}
	}
	return false
}

func insertIntoSubcommands(oo []sb.Opt, ch []sb.Choice, path []string) (ok bool) {
	for i, o := range oo {
		if path[0] == o.Name {
			if len(path) == 1 {
				oo[i].Choices = ch
				return true
			}
			return insertIntoSubcommands(o.Options, ch, path[1:])
		}
	}
	return false
}
