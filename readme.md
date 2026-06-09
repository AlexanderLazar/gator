# Gator, A Minimally-Functional CLI RSS Aggregator

## Requires Go and PostgreSQL

## To install:
Run go install from the root directory

## Configuration:
By default, configuration files are expected at "$HOME/.gatorconfig.json"; these have the form

{
	"db_url": <postgres_URL>,
	"current_user_name": <username>
}

## Usage:
To use the program, run "gator <COMMAND> <ARGS>" from the command line.
The commands are as follows:
  - login <username>
    - Changes to <username>
	- register <username>
	  - Adds <username> as a user
	- users
	  - Lists all registered users 
	- agg <duration>
	  - Starts a process that aggregates from subscribed RSS feeds; duration is the length of time before the aggregator checks again.
	- addfeed <title> <url>
	  - If no user is following a feed, follows the feed with title <title> at <url>
	- feeds
	  - Lists all feeds that the any user is following
	- follow <url>
	  - If another user is already following an RSS feed at <url>, follows the feed.
	- following
	  - Lists all feeds that the currently logged-in user is following.
	- unfollow <url>
	  - Allows the currently logged-in user to unfollow the feed at <url>. If the feed has no users following it, it is removed from the database.
	- browse <limit>
	  - Displays the most recent posts from the logged-in user's followed RSS feeds, up to <limit>. If <limit> is omitted, shows the two most recent posts.
