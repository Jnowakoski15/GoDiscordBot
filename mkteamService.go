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

// MessageCreate This function will be called (due to AddHandler above) every time a new message is created on any channel that the autenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
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
