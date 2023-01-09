# Discord Virtual Coffees

CLI tool written in Go that bulk creates/deletes voice channels on a Discord server. Each voice channel can have 2 users max, so it operates like a virtual coffee table.  Very bare bones and tailored to my use case.

To create 5 voice channels:

```bash
go run main.go -method=create -number=5
```

To delete 3 voice channels:

```bash
go run main.go -method=delete -number=3
```

To make it work on your server:

1. Create a `.env` file that will store the `BOT_TOKEN` and `SERVER_ID` value
2. Create and add a bot to your Discord server
3. Get the bot's token and the server ID and add it to the `.env` file
4. Update the `TablesCategoryId` constant in the `main.go` file to the channel category inside which you wish to create a channel. You can get the ID by fetching all the channels on your Discord server.
