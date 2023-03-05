package models

import (
	"byte_dance_5th/pkg/errortype"
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Message struct {
	Id         int64     `json:"id"`
	UserInfoId int64     `json:"from_user_id"`
	ToUserId   int64     `json:"to_user_id"`
	Content    string    `json:"content"`
	CreateTime int64     `json:"create_time"`
	CreatedAt  time.Time `json:"-"`
}

type MessageDAO struct {
}

var (
	messageDAO  *MessageDAO
	messageOnce sync.Once
)

func NewMessageDAO() *MessageDAO {
	messageOnce.Do(func() {
		messageDAO = new(MessageDAO)
	})
	return messageDAO
}

// CreateMessage 创建聊天记录
func (m *MessageDAO) CreateMessage(msg *Message) error {
	if msg == nil {
		return errors.New("CreateMsgByFromId:" + errortype.PointerIsNilErr)
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(msg).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteMessage 删除消息
func (m *MessageDAO) DeleteMessage(id int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM messages WHERE id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

// QueryMsgById 通过消息ID查询聊天记录
func (m *MessageDAO) QueryMsgById(id int64, msg *Message) error {
	if msg == nil {
		return errors.New("QueryMsgById" + errortype.PointerIsNilErr)
	}
	return DB.Where("id = ?", id).First(msg).Error
}

// QueryMsgListByFromIdAndToId 通过 to_user_id 和 from_user_id 来查询聊天记录列表
func (m *MessageDAO) QueryMsgListByFromIdAndToId(fromId int64, toId int64, messages *[]*Message, preMsgTime int64) error {
	if messages == nil {
		return errors.New("QueryMsgListByFromIdAndToId" + errortype.PointerIsNilErr)
	}
	if err := DB.Model(&Message{}).Where("create_time < ? AND ((to_user_id = ? AND user_info_id = ?) OR (to_user_id = ? AND user_info_id = ?))", preMsgTime, toId, fromId, fromId, toId).Order("created_at DESC").Find(messages).Error; err != nil {
		return err
	}
	return nil
}

// QueryLatestMsgByUid 根据用户uid返回最新的
func (m *MessageDAO) QueryLatestMsgByUid(fromId int64, toId int64) (string, int64, error) {
	var message Message
	//log.Println(fromId, toId)
	if err := DB.Where("(to_user_id = ? AND user_info_id = ?) OR (to_user_id = ? AND user_info_id = ?)", toId, fromId, fromId, toId).Order("created_at DESC").First(&message).Error; err != nil {
		return "", 0, err
	}
	var content string = message.Content
	var msgType int64
	if message.UserInfoId == fromId {
		msgType = 1
	} else {
		msgType = 0
	}
	return content, msgType, nil
}
