package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Overhustler/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Expected a single argument, the username")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("User name was set to %s", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Expected a single argument, the name")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]})
	if err != nil {
		return err
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("%s was created successfully", user.Name)
	fmt.Printf("%v", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		println("Delete unsuccessful")
		return err
	}
	println("All users deleted.")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	names, err := s.db.ListUsers(context.Background())
	if err != nil {
		return err
	}
	for i := range names {
		if names[i] == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)", names[i])
			continue
		}
		fmt.Printf("* %s", names[i])
	}
	return nil
}
