# Phoenix Manager
Tower of Fantasy party finder using Discord Bot API written in GO and using [discordgo](https://github.com/bwmarrin/discordgo) API

## Build your own
### Requirement
* GO 
> I used 1.19.1 so it might not be usable in some version but you can try :)
* Your own discord bot
> Token and Application ID because slash command require Application ID to register command
### Building
1. Download the source
2. Create file name "config.env"
3. Create ENV variable name `BotToken` and `AppID` and may look like this
```
BotToken=ThisIsMyBotToken
AppID=5634655231564
```
and of course use your own

4. Run using `go run .` or `go build` then run .exe and you are done!

## Know Issue
* You can join the same party and when you press "Leave Party", it will panic

## Consider to Implement
* Party that only available to specific Guild only

## Q & A
- Why named "Phoenix Manger"?
> Because my CREW name is "PHOENIX"
- Why GO?
> TBH I like Python because it's easy to write and understand but it DOES NOT compile down to .exe that can execute quickly by just simply double click.
> Of course there's a way but that's too much work to be done so I choose GO
