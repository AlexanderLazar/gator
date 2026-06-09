package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/alexanderlazar/my-rss-aggregator/internal/config"
	"github.com/alexanderlazar/my-rss-aggregator/internal/database"
	"github.com/alexanderlazar/my-rss-aggregator/internal/rss"
	"github.com/google/uuid"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandMap map[string]func(*config.State, Command) error
}

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("login handler expects exactly one argument")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	err = s.Configptr.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	return nil
}

func HandlerRegister(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("login handler expects exactly one argument")
	}
	dbUser, err := s.Db.CreateUser(context.Background(),
		database.CreateUserParams{ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0]})
	if err != nil {
		return err
	}
	s.Configptr.SetUser(dbUser.Name)
	fmt.Println("user created")
	fmt.Println("Name: ", dbUser.Name)
	fmt.Println("ID: ", dbUser.ID)
	fmt.Println("Created: ", dbUser.CreatedAt)
	fmt.Println("Updated: ", dbUser.UpdatedAt)
	return nil
}

func HandlerAddFeed(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("feed handler expects both a name and a url")
	}
	// uid, err := s.Db.GetUser(context.Background(), s.Configptr.User)
	// if err != nil {
	// 	return err
	// }
	dbFeed, err := s.Db.AddFeed(context.Background(),
		database.AddFeedParams{Name: cmd.Args[0],
			Url:    cmd.Args[1],
			UserID: user.ID,
		})
	if err != nil {
		return err
	}
	feedfollow := database.CreateFeedFollowParams{
		CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: dbFeed.UserID, FeedID: dbFeed.ID}
	_, err = s.Db.CreateFeedFollow(context.Background(), feedfollow)
	if err != nil {
		return err
	}
	fmt.Println("feed added")
	fmt.Println("Name: ", dbFeed.Name)
	fmt.Println("URL: ", dbFeed.Url)
	fmt.Println("User ID: ", dbFeed.UserID)
	return nil
}

func HandlerFeeds(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return errors.New("feed list expects no arguments")
	}
	feedlist, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, item := range feedlist {
		fmt.Println("Name: ", item.Name)
		fmt.Println("URL: ", item.Url)
		// fmt.Println("Username: ", item.UserName)
	}
	return nil
}

func HandlerFollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("follow expects the URL of the feed to follow")
	}
	url := cmd.Args[0]
	// usr, err := s.Db.GetUser(context.Background(), s.Configptr.User)
	// if err != nil {
	// 	return err
	// }
	feedid, err := s.Db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}
	feedfollow := database.CreateFeedFollowParams{
		CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feedid}
	followRow, err := s.Db.CreateFeedFollow(context.Background(), feedfollow)
	if err != nil {
		return err
	}
	fmt.Println("Feed: ", followRow.FeedName)
	fmt.Println("User: ", followRow.UserName)
	return nil
}

func HandlerUnfollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("unfollow expects exactly one feed to unfollow")
	}
	feedid, err := s.Db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	err = s.Db.Unfollow(context.Background(),
		database.UnfollowParams{UserID: user.ID, FeedID: feedid})
	return err
}

func HandlerFollowing(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) > 0 {
		return errors.New("following expects no arguments")
	}
	following, err := s.Db.GetFeedFollowsForUser(context.Background(), s.Configptr.User)
	if err != nil {
		return err
	}
	fmt.Printf("User: %v is following:\n", user.Name)
	for _, follow := range following {
		fmt.Println(follow.FeedName)
	}
	return nil
}

func HandlerBrowse(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) > 1 {
		return errors.New("browse takes at most one argument")
	}
	var err error
	limit := 2
	if len(cmd.Args) == 1 {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
	}
	params := database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)}
	posts, err := s.Db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println(len(posts))
	for _, post := range posts {
		fmt.Println(post.Title)
		fmt.Println(post.Description)
	}
	return nil
}

func HandlerResetUsers(s *config.State, cmd Command) error {
	err := s.Db.ResetUsers(context.Background())
	return err
}

func HandlerGetUsers(s *config.State, cmd Command) error {
	usrs, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, usr := range usrs {
		if usr == s.Configptr.User {
			curr := "(current)"
			fmt.Println("* ", usr, curr)
		} else {
			fmt.Println("* ", usr)
		}
	}
	return nil
}

func HandlerAgg(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("agg expexcts only a duration")
	}
	dur, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("Collecting feeds every ", dur)
	ticker := time.NewTicker(dur)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	// return nil
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	cfun, ok := c.CommandMap[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return cfun(s, cmd)
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) error {
	c.CommandMap[name] = f
	return nil
}
func MiddlewareLoggedIn(handler func(s *config.State, cmd Command, user database.User) error) func(*config.State, Command) error {
	return func(s *config.State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Configptr.User)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *config.State) error {
	next, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {

		return err
	}
	marknext := database.MarkFeedFetchedParams{Url: next.Url, LastFetchedAt: sql.NullTime{Time: next.LastFetchedAt.Time, Valid: true}}
	err = s.Db.MarkFeedFetched(context.Background(), marknext)
	if err != nil {

		return err
	}
	feed, err := rss.FetchFeed(context.Background(), next.Url)
	if err != nil {

		return err
	}
	for _, item := range feed.Channel.Item {
		if item.Title == "" {
			continue
		}
		pblTime, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {

			pblTime, err = time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {

				continue
			}
		}
		// feedid, err := s.Db.GetFeed(context.Background(), next.Url)
		// if err != nil {
		// 	log.Printf(err.Error())
		// 	return err
		// }
		post := database.CreatePostParams{ID: uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pblTime,
			FeedID:      next.ID,
		}
		err = s.Db.CreatePost(context.Background(), post)

		if err != nil {
			continue
		}
	}
	return nil
}
