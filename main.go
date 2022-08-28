package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env")
	}

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(VoiceStateUpdate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening session", err)
		return
	}

	fmt.Println("Bot is running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = dg.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

type UserState struct {
	User        *discordgo.User
	CurrentVCID string
}

var (
	usermap = map[string]*UserState{}
)

func VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	_, ok := usermap[v.UserID]

	// 新規ユーザー
	if !ok {
		usermap[v.UserID] = new(UserState)
		user, err := s.User(v.UserID)
		if err != nil {
			fmt.Println(err)
			return
		}
		usermap[v.UserID].User = user
	}

	oldVCID := usermap[v.UserID].CurrentVCID    // 前回のボイスチャンネルのID
	usermap[v.UserID].CurrentVCID = v.ChannelID // 今回のボイスチャンネルのIDをセット

	channel, err := s.Channel(v.ChannelID)
	if err != nil {
		fmt.Println(err)
		return
	}

	nickname := v.Member.Nick
	if nickname == "" {
		nickname = v.Member.User.Username
	}
	channelName := channel.Name
	textChannelID := os.Getenv("CHANNEL_ID")

	fmt.Println(oldVCID)
	fmt.Println(v.ChannelID)

	if oldVCID == "" {
		text := "__" + nickname + "__ が " + "**" + channelName + "** に入ったよん"
		_, err := s.ChannelMessageSend(textChannelID, text)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
