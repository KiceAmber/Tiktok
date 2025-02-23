package controller

import (
	"Project/dao"
	"Project/models"
	"Project/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	var q models.CommentActionRequest
	q.Token = c.DefaultQuery("token", "")
	q.VideoID = utils.String2int64(c.DefaultQuery("video_id", ""))
	q.ActionType = int(utils.String2int64(c.DefaultQuery("action_type", "-1")))
	q.CommentText = c.DefaultQuery("comment_text", "")
	q.CommentID = utils.String2int64(c.DefaultQuery("comment_id", "-1"))
	fmt.Println("===============>", q)
	// 从 token 解析 user_id
	myClaims, err := utils.ParseToken(q.Token)
	if err != nil { // token 解析失败
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -2,
			StatusMsg:  "token error!",
		})
		return
	} else { // 如果 token 解析成功，获取 userId
		q.UserID, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}

	var comment *models.Comment
	comment, err = dao.PerformCommentAction(
		q.UserID,
		q.VideoID,
		q.ActionType,
		q.CommentText,
		q.CommentID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommentActionResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		Comment: models.Comment{
			ID:         comment.ID,
			Author:     comment.Author,
			AuthorID:   comment.AuthorID,
			VideoID:    comment.VideoID,
			Content:    comment.Content,
			CreateTime: comment.CreateTime,
			UpdateTime: comment.UpdateTime,
		},
	})
}

// CommentList 评论列表
func CommentList(c *gin.Context) {
	var result models.CommentListResponse
	// 获取请求参数
	token := c.DefaultQuery("token", "")
	video_id, err := strconv.ParseInt(c.DefaultQuery("video_id", "-1"), 10, 64)
	if err != nil { // 获取视频 id 错误
		result.Response.StatusCode = -1
		result.Response.StatusMsg = "git video_id error!"
		c.JSON(http.StatusOK, result)
		return
	}
	// 获取当前用户 id
	var userId int64
	myClaims, err := utils.ParseToken(token)
	if err != nil { // token 解析失败
		result.Response.StatusCode = -2
		result.Response.StatusMsg = "token error!"
		c.JSON(http.StatusOK, result)
		return
	} else { // 如果 token 解析成功，获取 userId
		userId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}
	// 获取评论列表
	result.CommentList, err = dao.GetCommentList(userId, video_id)
	if err != nil { // 如果获取评论列表失败
		result.Response = models.Response{
			StatusCode: -3,
			StatusMsg:  err.Error(),
		}
		c.JSON(http.StatusOK, result)
		return
	}
	result.Response = models.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	c.JSON(http.StatusOK, result)
}
