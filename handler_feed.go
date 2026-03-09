package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Overhustler/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("Error not the correct amount of arguements")
	}

	feedDb, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    feedDb.UserID, // The user from your config/state
		FeedID:    feedDb.ID,     // The ID of the feed you just created
	})
	if err != nil {
		return err
	}

	fmt.Printf("%+v", feedDb)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())

	if err != nil {
		return err
	}
	for i := range feeds {
		fmt.Printf("%+v", feeds[i])
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed, // Replace with actual feed ID
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s %s", follow.FeedName, follow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	followed, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return err
	}
	fmt.Printf("%v", followed)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.db.DeleteFeedFollowsByUserIDandFeedID(context.Background(), database.DeleteFeedFollowsByUserIDandFeedIDParams{UserID: user.ID, FeedID: feed})
	if err != nil {
		return err
	}
	println("feedfollows record deleted successfully")
	return nil
}
func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		if val, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = val
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})

	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("%v\n", post.Title)
		fmt.Printf("%v\n", post.Url)
		if post.Description != "" {
			fmt.Printf("%v\n", post.Description)
		}
		fmt.Printf("%v\n", post.CreatedAt)
	}
	return nil
}
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)

		if err != nil {
			return err
		}

		handler(s, cmd, user)
		return nil
	}
}
