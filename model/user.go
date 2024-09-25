package model

type User struct {
	UserID   int64  `json:"user_id" gorm:"column:user_id"`
	Username string `json:"username" gorm:"column:user_name"`
	Password string `json:"password" gorm:"column:password"`
}

func (User) TableName() string {
	return "user"
}
