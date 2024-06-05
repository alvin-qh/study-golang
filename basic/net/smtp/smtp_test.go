package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMTP_NewSMTP(t *testing.T) {
	smtp, err := NewSMTP()
	assert.Nil(t, err)

	assert.Equal(t, "smtp.163.com", smtp.Host)
	assert.Equal(t, 465, smtp.Port)
	assert.Equal(t, "quhao317@163.com", smtp.Sender)
}

func TestSMTP_SendMail(t *testing.T) {
	smtp, err := NewSMTP()
	assert.Nil(t, err)

	err = smtp.Send("mousebaby8080@gmail.com", "测试", "<b>Hello World</b>")
	assert.Nil(t, err)
}
