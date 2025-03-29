package main

import (
	"context"
	"database/sql"
	"gator/internal/database"
	"gator/internal/rss"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"time"
)

import (
	"fmt"
	"gator/internal/config"
	"gator/internal/state"
	"os"
)

func HandlerLogin(s *state.State, cmd state.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("not sufficient arguments. The login handler expects a single argument, the username")
	}
	_, err := s.Queries.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user not found in the database. ")
	}
	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't set the user")
	}
	fmt.Println("User" + " " + cmd.Args[0] + " set successfully")
	return nil
}

func HandlerRegister(s *state.State, cmd state.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("not sufficient arguments. The register handler expects a single argument, the name of the user to be created")
	}
	_, err := s.Queries.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: cmd.Args[0]})
	if err != nil {
		return fmt.Errorf("couldn't create the user %v\n", err)
	}
	fmt.Println("User" + " " + cmd.Args[0] + " created successfully")
	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't set the user\n")
	}

	return nil
}

func HandlerReset(s *state.State, cmd state.Command) error {
	err := s.Queries.DeleteAll(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting all users\n.")
	}
	return nil
}

func HandlerUsers(s *state.State, cmd state.Command) error {
	usrs, err := s.Queries.GetUsers(context.Background())
	gatorConfig := config.Read()
	loggedUserName := gatorConfig.CurrentUserName
	if err != nil {
		return fmt.Errorf("error getting all users\n.")
	}
	for _, usr := range usrs {
		if usr == loggedUserName {
			fmt.Printf("* %s (current)\n", usr)
		} else {
			fmt.Printf("* %s\n", usr)
		}
	}
	return nil
}

func HandlerAgg(s *state.State, cmd state.Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("* %v feed\n", feed)
	return nil
}

func HandlerAddFeed(s *state.State, cmd state.Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("not sufficient arguments. The add feed handler expects two arguments - name and url")
	}
	usrName := s.Config.CurrentUserName
	usr, err := s.Queries.GetUser(context.Background(), usrName)
	if err != nil {
		return err
	}
	id := usr.ID
	setFeed, err := s.Queries.SetFeed(context.Background(), database.SetFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    id,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", setFeed)
	return nil
}

func HandlerFeeds(s *state.State, cmd state.Command) error {
	feeds, err := s.Queries.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%s\t%s\t%s\n", feed.Feedname, feed.Feedurl, feed.Username)
	}
	return nil
}

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
	commands.Register("login", HandlerLogin)
	commands.Register("register", HandlerRegister)
	commands.Register("reset", HandlerReset)
	commands.Register("users", HandlerUsers)
	commands.Register("agg", HandlerAgg)
	commands.Register("addfeed", HandlerAddFeed)
	commands.Register("feeds", HandlerFeeds)
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
