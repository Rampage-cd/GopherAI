package message

import (
	"GopherAI/common/mysql"
	"GopherAI/model"
)

func CreateMessage(message *model.Message) (*model.Message, error) {
	err := mysql.DB.Create(message).Error
	return message, err
}
