package state

import (
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
)

type State struct {
	Queries *database.Queries
	Config  *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, handler func(*State, Command) error) {
	c.Commands[name] = handler
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("this command does not exist. provide a different command")
	}
	err := handler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}
