package main

import (
	"errors"
	"math/rand"
	"strings"
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
		if pres.Status == discordgo.StatusOnline && pres.User.Bot == false {
			id := pres.User.ID
			onlineMem = append(onlineMem, id)
		}
	}

	if len(onlineMem) == 0 {
		return nil, errors.New("No Members found online")
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

func createEmbedOutput(team1 string, team2 string) *discordgo.MessageEmbed {

	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x000000, // Black
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Team 1:",
				Value:  team1,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Team 2:",
				Value:  team2,
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     "Team Selection",
	}
}

// MessageCreate This function will be called (due to AddHandler above) every time a new message is created on any channel that the autenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!mktest" {
		s.ChannelMessageSend(m.ChannelID, "Nothing currently in test")
	}

	if m.Content == "!mkteam" {
		vs, _ := findUserVoiceState(s, m.Author.ID)
		oms, err := findOnlineMembers(s, vs.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		voiceChannelMembers, _ := findSharedVoiceChannelmembers(s, oms, vs.ChannelID)
		voiceChannelMembersShuffled := shuffleArray(voiceChannelMembers)
		shuffledNickNames, _ := findNickNames(s, vs.GuildID, voiceChannelMembersShuffled)

		var sb1 strings.Builder
		var sb2 strings.Builder

		for i, nickName := range shuffledNickNames {
			if i%2 == 0 {
				sb1.WriteString(nickName)
				sb1.WriteString("\n")
			} else {
				sb2.WriteString(nickName)
				sb2.WriteString("\n")
			}
		}
		if sb1.String() == "" {
			sb1.WriteString("Empty")
		}

		if sb2.String() == "" {
			sb2.WriteString("Empty")
		}

		s.ChannelMessageSendEmbed(m.ChannelID, createEmbedOutput(sb1.String(), sb2.String()))
	}
}
