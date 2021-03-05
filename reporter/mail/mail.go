package mail

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

// mailMessage is mail message format.
type mailMessage struct {
	Emails  []string `json:"emails"`
	Text    string   `json:"text"`
	Subject string   `json:"subject"`
}

// AtellaMailConfig is mail Channel configuration.
type AtellaMailConfig struct {
	Address    string   `json:"address"`
	Port       int16    `json:"port"`
	Auth       bool     `json:"auth"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	From       string   `json:"from"`
	To         []string `json:"to"`
	Enabled    bool     `json:"disabled"`
	NetTimeout int
}

// newMailMessage create message for Mail Channel.
func (config *AtellaMailConfig) newMailMessage() *mailMessage {
	local := new(mailMessage)
	local.Emails = config.To
	local.Text = ""
	return local
}

// Send send message (text) via Mail Channel to user's emails , specifying
// in Emails array.
func (config *AtellaMailConfig) Send(text, hostname string) (bool, error) {
	if !config.Enabled {
		return false, nil
	}
	if config.To == nil || len(config.To) < 1 {
		return false, fmt.Errorf("Mail users list are empty")
	}
	msg := config.newMailMessage()
	msg.Subject = fmt.Sprintf("Message from Atella at %s",
		hostname)
	msg.Text = text
	msg.Emails = config.To
	err := config.sendMessage(*msg)
	if err != nil {
		return false, err
	}
	return true, nil
}

// sendMessage send message (text) via Mail Channel to user's emails , specifying
// in Emails array.
func (config *AtellaMailConfig) sendMessage(msg mailMessage) error {
	var (
		emails string = ""
		err    error  = nil
		d      *gomail.Dialer
		conn   gomail.SendCloser
	)

	for i := 0; i < len(msg.Emails); i = i + 1 {
		if msg.Emails[i] != "" {
			if emails == "" {
				emails = msg.Emails[i]
			} else {
				emails = fmt.Sprintf("%s, %s", emails, msg.Emails[i])
			}
		}
	}
	m := gomail.NewMessage()
	m.SetHeader("From", config.From)
	m.SetHeader("To", emails)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Text)

	d = config.dialer()
	conn, err = d.Dial()
	if err != nil {
		return err
	}
	err = gomail.Send(conn, m)
	if err != nil {
		return err
	}
	return nil
}

func (config *AtellaMailConfig) dialer() (d *gomail.Dialer) {
	if config.Username == "" {
		d = &gomail.Dialer{Host: config.Address, Port: int(config.Port)}
	} else {
		d = gomail.NewPlainDialer(config.Address, int(config.Port), config.Username, config.Password)
	}

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d
}
