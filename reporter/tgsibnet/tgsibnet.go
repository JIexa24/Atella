package tgsibnet

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// AtellaTgSibnetConfig is config for communicate with Sibnet TG bot.
type AtellaTgSibnetConfig struct {
	Address    string
	Port       int16
	Protocol   string
	To         []string
	Enabled    bool
	NetTimeout int
}

// tgSibnetMessage is message format for Sibnet TG bot.
type tgSibnetMessage struct {
	Event     string   `json:"event"`
	Usernames []string `json:"usernames"`
	Text      string   `json:"text"`
}

// tgSibnetPacket is packet format for Sibnet TG bot.
type tgSibnetPacket struct {
	Command string          `json:"command"`
	Message tgSibnetMessage `json:"message"`
}

// newTgSibnetMessage create message for TgSibnet Channel.
func (config *AtellaTgSibnetConfig) newTgSibnetMessage() *tgSibnetMessage {
	local := new(tgSibnetMessage)
	local.Event = "personal"
	local.Usernames = nil
	local.Text = ""
	return local
}

// Send initialize send message (text) via TgSibnet Channel to users,
// specifying in to array in config.
func (config *AtellaTgSibnetConfig) Send(text, hostname string) (bool,
	error) {
	if !config.Enabled {
		return false, nil
	}

	if config.To == nil || len(config.To) < 1 {
		return false, fmt.Errorf("SibnetBot users list are empty")
	}
	msg := config.newTgSibnetMessage()
	msg.Text = fmt.Sprintf("[%s]: %s", hostname, text)
	msg.Usernames = config.To
	result, err := config.sendMessage(*msg)
	if err != nil {
		return false, err
	}
	ret := false
	if result == "ok" {
		ret = true
	}
	return ret, nil
}

// sendMessage send message (text) via TgSibnet Channel to users, specifying in
// to array in config.
func (config *AtellaTgSibnetConfig) sendMessage(msg tgSibnetMessage) (string, error) {
	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", config.Address,
			config.Port),
		time.Duration(config.NetTimeout)*time.Second)
	if err != nil {
		return "", err
	}

	pack := tgSibnetPacket{"sendMessage", msg}
	msgJSON, _ := json.Marshal(pack)

	_, err = conn.Write(msgJSON)
	if err != nil {
		return "", err
	}
	reply, _ := bufio.NewReader(conn).ReadString('\n')
	return reply, nil
}
