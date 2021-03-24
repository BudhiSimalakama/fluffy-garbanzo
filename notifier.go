package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// XMLReturn format that the Roblox model is stored in
type XMLReturn struct {
	XMLName xml.Name `xml:"roblox"`
	Item    struct {
		XMLName    xml.Name `xml:"Item"`
		Properties struct {
			XMLName xml.Name `xml:"Properties"`
			Strings []string `xml:"string"`
		} `xml:"Properties"`
	} `xml:"Item"`
}

// Notification format notifications are stored in
type Notification struct {
	Title   string
	Message string
	Icon    string
	URL     string
	Items   map[string]string
}

// StartNotifier starts notifier checking
func StartNotifier() {
	setLatestKeys()

	for {
		<-time.After(2 * time.Second)
		go checker()
	}
}

func genKey(notif Notification) string {
	return fmt.Sprintf("%s-%s", notif.Title, notif.Message)
}

func setLatestKeys() {
	items := retrieveData()

	latestKeys = []string{}

	for _, notif := range items {
		latestKeys = append(latestKeys, genKey(notif))
	}
}

func retrieveData() []Notification {
	xmlData, _ := HTTPGet("https://assetdelivery.roblox.com/v1/asset?id=317944796")
	var data XMLReturn
	xml.Unmarshal([]byte(xmlData), &data)

	if len(data.Item.Properties.Strings) != 2 {
		time.Sleep(1 * time.Second)
		return retrieveData()
	}

	decoded, _ := base64.StdEncoding.DecodeString(data.Item.Properties.Strings[1])
	var notifs []Notification
	json.NewDecoder(bytes.NewReader(decoded)).Decode(&notifs)

	return notifs
}

func returnEmbed(notif Notification) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color: 15527148,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: notif.Icon,
		},
		Description: fmt.Sprintf("[%s](%s?rbxp=70998839)", notif.Message, strings.Replace(notif.URL, "?rbxp=48103520", "", -1)),
		Title:       notif.Title,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "https://www.rbxleaks.com/discord",
		},
	}
	fields := []*discordgo.MessageEmbedField{}
	for i, field := range notif.Items {
		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s:", i),
				Value:  field,
				Inline: true,
			})
	}
	embed.Fields = fields

	return embed
}

func checker() {
	items := retrieveData()

	for _, notif := range items {
		if !inArray(genKey(notif), latestKeys) {
			latestKeys = append(latestKeys, genKey(notif))
			fmt.Println(fmt.Sprintf("%s: %s [%s]", notif.Title, notif.Message, time.Now().Format(time.RFC1123)))
			embed := returnEmbed(notif)
			Bot.NewNotification(embed)
		}
	}
	if len(latestKeys) > 15 {
		setLatestKeys()
	}
}

func inArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
