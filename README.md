# Valorant Discord Stats

Go implementation to post Valorant stats onto Discord  

[Invite this bot to your server](https://discord.com/oauth2/authorize?client_id=396807688039694346&scope=bot&permissions=68608)   

## Prerequisites

* Installed Go  
* Downloaded blitz.gg application  
* Take note of your nametag from blitz.gg
  * Example: https://blitz.gg/valorant/profile/fompei-na1 -> fompei-na1 is the nametag

## Start Discord Bot

Run `go run discord.go -t <token>` to start the Discord bot

Follow this [guide](https://www.writebots.com/discord-bot-token/) up until step 4 to get the token needed to run the bot 

## Install onto EC2 Instance

`sudo yum update`  
`sudo yum install -y git golang`  
`git clone https://github.com/aarlin/valorant-discord-stats.git`  
`sudo amazon-linux-extras install docker`  
`sudo service docker start`  
`sudo usermod -a -G docker ec2-user`  

## Build Docker Image

1. Run `docker build -t valorant-discord-stats .`  
2. Run `docker images` and locate the Image ID for valorant-discord-stats  
3. Run `docker run -it <Image ID>` OR `docker run -d <Image ID>`  
4. Use `docker ps` for currently running processes  

## Commands

| Command             | Description                                       |
|---------------------|---------------------------------------------------|
| !commands           | See list of commands you can query                |
| !career \<nametag>   | See hit percentages for your total career         |
| !last20 \<nametag>   | See hit percentages from the last 20 games played |
| !lastgame \<nametag> | See hit percentages from the last game played     |

## References

https://developers.google.com/sheets/api/quickstart/go 
