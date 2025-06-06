package user

type User struct {
	ID       int64
	Email    string
	Username string
	Password string
}

func (User) TableName() string {
	return "user_t"
}
