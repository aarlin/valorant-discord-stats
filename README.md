# Valorant Discord Stats

Go implementation to post Valorant stats onto Discord 

## Prerequisites

* Installed Go  
* Downloaded blitz.gg application  
* Take note of your nametag from blitz.gg
  * Example: https://blitz.gg/valorant/profile/fompei-na1 -> fompei-na1 is the nametag

## Start Discord Bot

Run `go run discord.go -t <token>` to start the Discord bot

Follow this [guide](https://www.writebots.com/discord-bot-token/) up until step 4 to get the token needed to run the bot  

## Commands

| Command             | Description                                       |
|---------------------|---------------------------------------------------|
| !commands           | See list of commands you can query                |
| !career <nametag>   | See hit percentages for your total career         |
| !last20 <nametag>   | See hit percentages from the last 20 games played |
| !lastgame <nametag> | See hit percentages from the last game played     |

## References

https://developers.google.com/sheets/api/quickstart/go 
