package model

// 聊天室对应的用户
type UserRoom struct {
	ID        int64 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RoomID    int64 `json:"room_id" gorm:"room_id"`
	UserID    int64 `json:"user_id" gorm:"user_id"`
	MessageID int64 `json:"message_id" gorm:"message_id"`
	//CreateTime time.Time `json:"create_time" gorm:"create_time"`
	//UpdateTime time.Time `json:"update_time" gorm:"update_time"`
}

func (UserRoom) TableName() string {
	return "user_room"
}
