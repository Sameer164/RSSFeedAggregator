# Terminal RSS Reader

A command-line RSS reader that allows you to follow RSS feeds and browse posts right in your terminal.

## Features

- Follow and manage RSS feeds
- Automatically fetch and store posts from your feeds
- Browse recent posts from all your followed feeds

## Installation

1. Clone this repository
2. Build the application: `make build`
3. Install Postgresql15. Make sure the service is up (`brew services start postgresql@15`)
4. Create gator database. `psql postgres` -> `CREATE DATABASE gator` -> `\c gator` -> `exit`
5. Go over to `sql/schema`. Run `goose postgres "postgres://<username>:@localhost:5432/gator" up`


## Running the Application

### Managing Users
- `./gator login <username>` - Logs in with the username provided if the username exists in the database
- `./gator register <username>` - Creates the user with the username in the database and automatically logs in
- `./gator users` - Displays all the users that can log in with a (*current) to the right of the currently logged in user
- 

### Following Feeds

- `./gator addfeed <feed_name> <url>` - Creates the feed and makes the currently logged in user follow the feed
- `./gator feeds` - Shows all the feeds from all users (so one can follow others' feeds as well)
- `./gator follow <feed_url>` - Makes the logged in user follow the feed with the feed_url. Feeds are unique by url
- `./gator unfollow <feed_url>` - Makes the logged in user unfollow the feed
- `./gator following` - Gives the names of the feeds followed by the logged in user

### Storing & Browsing Feeds

- `./gator browse [limit]` - View recent posts from all your feeds
    - Optional: Specify a limit (default: 2)
- `./gator agg <time>` - Fetches and Saves posts from the feeds you follow every <time> interval. Older feeds fetched first.
