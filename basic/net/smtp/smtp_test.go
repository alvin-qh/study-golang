package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建 SMTP 实例
func TestSMTP_NewSMTP(t *testing.T) {
	smtp, err := NewSMTP()
	assert.Nil(t, err)

	assert.Equal(t, "smtp.163.com", smtp.Host)
	assert.Equal(t, 465, smtp.Port)
	assert.Equal(t, "quhao317@163.com", smtp.Sender)
}

// 测试发送邮件
func TestSMTP_SendMail(t *testing.T) {
	t.Skipf("Run this test manual")

	smtp, err := NewSMTP()
	assert.Nil(t, err)

	err = smtp.Send("mousebaby8080@gmail.com", "测试", "<p>Hello World!</p><p>This mail was sent from Golang</p>")
	assert.Nil(t, err)
}
