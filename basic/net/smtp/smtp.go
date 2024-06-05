package smtp

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	// 加载 `.env` 环境变量文件
	godotenv.Load()
}

// 定义邮件内容模板
const (
	template = "To: %s\r\n" + // 收件人
		"Subject: %s\r\n" + // 邮件标题
		"From: %s <%s>\r\n" + // 发件人昵称 <发件人地址>
		"MIME-Version: 1.0\r\n" + // MIME 版本
		"Content-Type: text/html; charset=utf-8\r\n" + // 内容类型
		"\r\n" +
		"%s" // 邮件正文
)

// 定义 SMTP 发送结构体
type SMTP struct {
	Host     string    // SMTP 地址
	Port     int       // SMTP 端口号
	Auth     smtp.Auth // 发送认证实例 (明文)
	Sender   string    // 发送人地址
	Nickname string    // 发送人昵称
}

// 实例化结构体
func NewSMTP() (*SMTP, error) {
	// 获取端口号
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, err
	}

	// 创建 SMTP 实例
	return &SMTP{
		Host:     os.Getenv("SERVER"),
		Port:     port,
		Auth:     smtp.PlainAuth("", os.Getenv("ACCOUNT"), os.Getenv("PASSWORD"), os.Getenv("SERVER")),
		Sender:   os.Getenv("SENDER"),
		Nickname: os.Getenv("NICKNAME"),
	}, nil
}

// 发送邮件
func (s *SMTP) Send(to string, subject string, msg string) error {
	// 填充邮件模板
	msg = fmt.Sprintf(template, to, subject, s.Nickname, s.Sender, msg)

	// 发送邮件
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Host, s.Port),
		s.Auth,
		s.Sender,
		[]string{to},
		[]byte(msg),
	)
}
