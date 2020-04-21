package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func findUserVoiceState(session *discordgo.Session, userid string) (*discordgo.VoiceState, error) {
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				return vs, nil
			}
		}
	}
	return nil, errors.New("Could not find user's voice state")
}

func findOnlineMembers(s *discordgo.Session, guildID string) ([]string, error) {
	var onlineMem []string
	g, _ := s.Guild(guildID)
	for _, pres := range g.Presences {
		if pres.Status == discordgo.StatusOnline {
			id := pres.User.ID
			onlineMem = append(onlineMem, id)
		}
	}
	return onlineMem, nil
}

func findSharedVoiceChannelmembers(s *discordgo.Session, onlineMemberIDs []string, authorVoiceChannel string) ([]string, error) {
	var voiceChanMembers []string
	for _, id := range onlineMemberIDs {
		onlineUserVoiceState, err := findUserVoiceState(s, id)
		if err != nil {
			continue
		}
		if authorVoiceChannel == onlineUserVoiceState.ChannelID {
			voiceChanMembers = append(voiceChanMembers, onlineUserVoiceState.UserID)
		}
	}

	return voiceChanMembers, nil
}

func shuffleArray(a []string) []string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a
}

func findNickNames(s *discordgo.Session, guildID string, memberIDs []string) ([]string, error) {

	var nickNames []string
	g, _ := s.Guild(guildID)
	for _, onlineMem := range memberIDs {
		for _, gmem := range g.Members {
			if onlineMem == gmem.User.ID {
				nick := gmem.Nick
				if nick == "" {
					user, _ := s.User(onlineMem)
					nick = user.Username
				}
				nickNames = append(nickNames, nick)
				break
			}
		}
	}
	return nickNames, nil
}

func main() {
	discordToken, exists := os.LookupEnv("DISCORD_TOKEN")

	if !exists {
		panic(fmt.Errorf("Could'nt read discord token"))
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
	fmt.Println("Closed down")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!mkteam" {
		vs, _ := findUserVoiceState(s, m.Author.ID)
		oms, _ := findOnlineMembers(s, vs.GuildID)
		voiceChannelMembers, _ := findSharedVoiceChannelmembers(s, oms, vs.ChannelID)
		voiceChannelMembersShuffled := shuffleArray(voiceChannelMembers)
		shuffledNickNames, _ := findNickNames(s, vs.GuildID, voiceChannelMembersShuffled)

		var sb1 strings.Builder
		var sb2 strings.Builder
		sb1.WriteString("**Team 1:**")
		sb2.WriteString("**Team 2:**")

		for i, nickName := range shuffledNickNames {
			if i%2 == 0 {
				sb1.WriteString("\n\t")
				sb1.WriteString(nickName)
			} else {
				sb2.WriteString("\n\t")
				sb2.WriteString(nickName)
			}
		}
		s.ChannelMessageSend(m.ChannelID, sb1.String())
		s.ChannelMessageSend(m.ChannelID, sb2.String())
	}
}
