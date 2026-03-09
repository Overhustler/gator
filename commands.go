package main

import (
	"errors"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if val, ok := c.cmds[cmd.name]; ok {
		return val(s, cmd)
	} else {
		return errors.New("command does not exist")
	}
}

func (c *commands) register(name string, f func(*state, command) error) {

	c.cmds[name] = f
}
