package main

import (
	"sync"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Wait        *sync.WaitGroup
	MailerChan  chan Message // 在後台發送郵件
	ErrorChan   chan error
	DoneChan    chan bool
}

// 描述郵件伺服器
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string // 附件
	Data        any
	DataMap     map[string]any
	Template    string
}

// a function to listen for messages on the MailerChan

func (m *Mail) sendMail(msg Message, errorChan chan error) {
	if msg.Template == "" {
		msg.Template = "mail"
	}

	// 在msg中自定義寄件者位置
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// pass data to template
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	// build html mail
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		errorChan <- err
	}

	// build plain text mail
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		errorChan <- err
	}

	// 創建服務端
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// smtp 創建客戶端
	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}

	// 設置電子郵件
	email := mail.NewMSG()
	// 寄件者資訊
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	// 設置正文
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	//  檢查是否有附件
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// 發送郵件
	err = email.Send(smtpClient)
	if err != nil {
		errorChan <- err
	}
}

// 建構 html 消息並建構純文本消息
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	return "", nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {

	return "", nil
}

func (m *Mail) getEncryption(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
