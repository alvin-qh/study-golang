package user

type User struct {
	Id    int64    `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email,omitempty"`
	Phone []string `json:"phone,omitempty"`
}

func New(id int64, name, email string, phone []string) *User {
	return &User{Id: id, Name: name, Email: email, Phone: phone}
}

func (u *User) AddPhone(phone string) {
	u.Phone = append(u.Phone, phone)
}
