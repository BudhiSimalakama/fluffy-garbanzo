# RBXNotifier

Source code for the RBXNotifier bot.

Due to the discord rate limit and the amount of servers the bot is in the bot has gotten extremely slow with notifications. I decided to open source this so people will be able to run this for their own server and have faster notifications.

Notifications are sourced from Roblox+

## Running the bot
### Prerequisites
To run this bot you need to create one on the [Discord developer portal](https://discordapp.com/developers/applications)

You must also have [Go](https://golang.org/) installed.

I suggest hosting this on an actual server but it is possible to just run from your desktop.

### Configuration
Rename the conf.json.example file to conf.json.

Copy the bot token from the Discord developer portal _not the client token_, you can find it on the bot tab. Put this token in conf.json after "Bot " on Discord Token.

Ex. `"DiscordToken": "Bot asdasdasdasdasd.asdf.asdfasdfasdf"`

Unless you are hosting this bot for more than about 1000 servers leave MinShardCount at 1.

Open command prompt in the directory and run `go build`, this will compile the code into an application.

To run the bot just open the application while the conf.json file is in the same directory as it. If it closes instantly it is likely that the config file is incorrect.