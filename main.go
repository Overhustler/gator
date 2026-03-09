package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Overhustler/gator/internal/config"
	"github.com/Overhustler/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	config := config.Read()
	var currentSate state
	currentSate.cfg = &config
	db, err := sql.Open("postgres", currentSate.cfg.DBURL)
	dbQueries := database.New(db)
	currentSate.db = dbQueries

	comms := commands{cmds: make(map[string]func(*state, command) error)}
	comms.register("login", handlerLogin)
	comms.register("register", handlerRegister)
	comms.register("reset", handlerReset)
	comms.register("users", handlerListUsers)
	comms.register("agg", handlerAgg)
	comms.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	comms.register("feeds", handlerFeeds)
	comms.register("following", middlewareLoggedIn(handlerFollowing))
	comms.register("follow", middlewareLoggedIn(handlerFollow))
	comms.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	comms.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("no arguments provided")
	}
	commandName := os.Args[1]
	var commandInput []string
	if len(os.Args) > 2 {
		commandInput = os.Args[2:]
	}
	command := command{
		name: commandName,
		args: commandInput,
	}
	err = comms.run(&currentSate, command)
	if err != nil {
		log.Fatal(err)
	}
}
