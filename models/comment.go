package models

import (
	"byte_dance_5th/pkg/errortype"
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Comment struct {
	Id         int64     `json:"id"`
	UserInfoId int64     `json:"-"`
	VideoId    int64     `json:"-"`
	User       UserInfo  `json:"user" gorm:"-"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"-"`
	CreateDate string    `json:"create_date" gorm:"-"`
}

type CommentDAO struct {
}

var (
	commentDao  *CommentDAO
	commentOnce sync.Once
)

func NewCommentDao() *CommentDAO {
	commentOnce.Do(func() {
		commentDao = new(CommentDAO)
	})
	return commentDao
}

// CreateAndCntAddOne 创建评论并将视频评论数量加一
func (c *CommentDAO) CreateAndCntAddOne(comment *Comment) error {
	//log.Println("CreateAndCntAddOne")
	if comment == nil {
		return errors.New("CreateAndCntAddOne:" + errortype.PointerIsNilErr)
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			//log.Println("CreateAndCntAddOne：评论未能入库")
			return err
		}
		if err := tx.Exec("UPDATE videos v SET v.comment_count = v.comment_count+1 WHERE v.id=?", comment.VideoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteAndCntSubOne 删除评论并将视频评论数量减一
func (c *CommentDAO) DeleteAndCntSubOne(commentId, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		//从表中删除
		if err := tx.Exec("DELETE FROM comments WHERE id = ?", commentId).Error; err != nil {
			return err
		}
		//将视频评论数ID减一
		if err := tx.Exec("UPDATE videos v SET v.comment_count = v.comment_count-1 WHERE v.id=? AND v.comment_count>0", videoId).Error; err != nil {
			return err
		}
		// 事务成功提交nil
		return nil
	})
}

// QueryCommentById 通过评论ID查询评论
func (c *CommentDAO) QueryCommentById(id int64, comment *Comment) error {
	if comment == nil {
		return errors.New("QueryCommentById" + errortype.PointerIsNilErr)
	}
	return DB.Where("id = ?", id).First(comment).Error
}

// QueryCommentListByVideoId 通过视频ID查询评论
func (c *CommentDAO) QueryCommentListByVideoId(videoId int64, comments *[]*Comment) error {
	if comments == nil {
		return errors.New("QueryCommentListByVideoId" + errortype.PointerIsNilErr)
	}
	if err := DB.Model(&Comment{}).Where("video_id = ?", videoId).Find(comments).Error; err != nil {
		return err
	}
	return nil
}
