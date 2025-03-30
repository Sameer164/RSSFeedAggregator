package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/rss"
	"gator/internal/state"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func middlewareLoggedIn(handler func(s *state.State, cmd state.Command, user database.User) error) func(s *state.State, cmd state.Command) error {
	return func(s *state.State, cmd state.Command) error {
		if s.Config.CurrentUserName == "" {
			return fmt.Errorf("no users are currently logged , there are no users in the database, so register one before you login")
		}
		usr, err := s.Queries.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user %s not found in the database\n", s.Config.CurrentUserName)
		}
		return handler(s, cmd, usr)
	}
}

func HandlerLogin(s *state.State, cmd state.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("not sufficient arguments. The login handler expects a single argument, the username")
	}
	usr, err := s.Queries.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user %s not found in the database\n", cmd.Args[0])
	}
	err = s.Config.SetUser(usr.Name)
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
	err = s.Config.SetUser("")
	if err != nil {
		return fmt.Errorf("couldn't clear the user in the config. May not work properly.\n")
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

func HandlerAgg(s *state.State, cmd state.Command, usr database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("not sufficient arguments. The aggregate handler expects a single argument - the duration to fetch feeds on")
	}
	timeInterval, _ := time.ParseDuration(cmd.Args[0])
	ticker := time.NewTicker(timeInterval)
	for ; ; <-ticker.C {
		feed, err := s.Queries.GetNextFeedToFetch(context.Background(), usr.ID)
		if err != nil {
			return fmt.Errorf("couldn't fetch next feed for user %s\n.", usr.ID)
		}
		_, err = s.Queries.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID: feed.ID, LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true}})
		if err != nil {
			return fmt.Errorf("couldn't mark feed for user %s\n.", usr.ID)
		}

		fetchedFeed, err := rss.FetchFeed(context.Background(), feed.Url)
		if err != nil {
			return fmt.Errorf("couldn't fetch feed\n.")
		}
		for _, post := range fetchedFeed.Channel.Item {
			publishedDateParsed, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", post.PubDate)
			if err != nil {
				return fmt.Errorf("couldn't parse published date\n.")
			}
			publishedDateParsed = publishedDateParsed.UTC()
			_, err = s.Queries.CreatePost(context.Background(), database.CreatePostParams{ID: uuid.New(),
				CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Title: post.Title,
				Url: post.Link, Description: post.Description, PublishedAt: publishedDateParsed, FeedID: feed.ID})
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					continue
				}
				return fmt.Errorf("couldn't create post for user %s\n.", usr.ID)
			}
		}
	}
}

func HandlerAddFeed(s *state.State, cmd state.Command, usr database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("not sufficient arguments. The add feed handler expects two arguments - name and url")
	}
	usrName := s.Config.CurrentUserName
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
	fmt.Println("User" + " " + usrName + " created " + setFeed.Name + "!")
	err = HandlerFollow(s, state.Command{Name: "follow", Args: []string{setFeed.Url}}, usr)
	if err != nil {
		return err
	}
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

func HandlerFollow(s *state.State, cmd state.Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("arguments count doesn't match. The follow handler expects one argument - the url of the feed")
	}
	feed, err := s.Queries.GetFeedFromURL(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("Couldn't find feed " + cmd.Args[0])
		return err
	}
	followInfo, err := s.Queries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		fmt.Println("Couldn't create the follow feed")
		return err
	}
	fmt.Printf("User: %s started following %s\n", followInfo.UserName, followInfo.FeedName)
	return nil
}

func HandlerFollowing(s *state.State, cmd state.Command, user database.User) error {
	if s.Config.CurrentUserName == "" {
		return fmt.Errorf("no users are currently logged in")
	}
	feedNames, err := s.Queries.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("The following feeds are followed by: %s\n", s.Config.CurrentUserName)
	for _, feedName := range feedNames {
		fmt.Printf("* %s\n", feedName)
	}
	return nil
}

func HandlerUnfollow(s *state.State, cmd state.Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("argument count doesn't match. The follow handler expects one argument - the url of the feed")
	}
	feed, err := s.Queries.GetFeedFromURL(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("Couldn't find feed " + cmd.Args[0])
		return err
	}
	_, err = s.Queries.Unfollow(context.Background(), database.UnfollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		fmt.Println("Couldn't unfollow " + cmd.Args[0])
		return err
	}
	return nil
}

func HandlerBrowse(s *state.State, cmd state.Command, user database.User) error {
	var limit int32 = 2
	if len(cmd.Args) == 1 {
		v, err := strconv.ParseInt(cmd.Args[0], 10, 64)
		if err != nil {
			return err
		}
		limit = int32(v)
	}
	posts, err := s.Queries.GetPostsForUser(context.Background(), database.GetPostsForUserParams{ID: user.ID, Limit: limit})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("* %s\n", post.Title)
	}
	return nil
}
