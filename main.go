package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/alexanderlazar/my-rss-aggregator/internal/commands"
	"github.com/alexanderlazar/my-rss-aggregator/internal/config"
	"github.com/alexanderlazar/my-rss-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(err.Error())
	}
	db, err := sql.Open("postgres", conf.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	dbQueries := database.New(db)
	stateptr := &config.State{Configptr: &conf, Db: dbQueries}

	coms := commands.Commands{CommandMap: make(map[string]func(*config.State, commands.Command) error)}

	coms.Register("login", commands.HandlerLogin)
	coms.Register("register", commands.HandlerRegister)
	coms.Register("reset", commands.HandlerResetUsers)
	coms.Register("users", commands.HandlerGetUsers)
	coms.Register("agg", commands.HandlerAgg)
	coms.Register("addfeed", commands.MiddlewareLoggedIn(commands.HandlerAddFeed))
	coms.Register("feeds", commands.HandlerFeeds)
	coms.Register("follow", commands.MiddlewareLoggedIn(commands.HandlerFollow))
	coms.Register("following", commands.MiddlewareLoggedIn(commands.HandlerFollowing))
	coms.Register("unfollow", commands.MiddlewareLoggedIn(commands.HandlerUnfollow))
	coms.Register("browse", commands.MiddlewareLoggedIn(commands.HandlerBrowse))
	inpt := os.Args

	if len(inpt) <= 1 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}
	to_run := inpt[1]
	args := inpt[2:]
	fun := commands.Command{Name: to_run, Args: args}

	err = coms.Run(stateptr, fun)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
