package main

import (
	"database/sql"
	"gator/internal/database"
	_ "github.com/lib/pq"
)

import (
	"fmt"
	"gator/internal/config"
	"gator/internal/state"
	"os"
)

func main() {
	gatorConfig := config.Read()
	db, err := sql.Open("postgres", gatorConfig.DbUrl)
	if err != nil {
		fmt.Println("couldn't connect to db")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s := &state.State{Config: gatorConfig, Queries: dbQueries}

	commands := &state.Commands{Commands: make(map[string]func(*state.State, state.Command) error)}
	commands.Register("login", middlewareLoggedIn(HandlerLogin))
	commands.Register("register", HandlerRegister)
	commands.Register("reset", HandlerReset)
	commands.Register("users", HandlerUsers)
	commands.Register("agg", HandlerAgg)
	commands.Register("addfeed", middlewareLoggedIn(HandlerAddFeed))
	commands.Register("feeds", HandlerFeeds)
	commands.Register("follow", middlewareLoggedIn(HandlerFollow))
	commands.Register("following", middlewareLoggedIn(HandlerFollowing))
	commands.Register("unfollow", middlewareLoggedIn(HandlerUnfollow))
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
	command := state.Command{Name: argsWithoutProg[0], Args: argsWithoutProg[1:]}
	err = commands.Run(s, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
