// publish 包，该包封装了投稿相关的借口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"Project/common"
	"Project/dao"
	"Project/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// Publish 投稿接口
func Publish(context *gin.Context) {
	var request models.PublishVideoRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = request.Token
	// TODO: JWT auth.

	// Data 需要: `PlayUrl`, `CoverUrl`，其余默认即可
	err := dao.CreateVideoByUserId(request.UserID, request.Data)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, &models.Response{
		StatusCode: common.StatusOK,
		StatusMsg:  "success",
	})
}

// VideoVo 接受从数据库查出来的部分数据的结构体
type VideoVo struct {
	Id            int64  `json:"id" db:"id"`                          // 视频 ID
	UserId        int64  `json:"user_id" db:"user_id"`                // userId
	PlayUrl       string `json:"play_url" db:"play_url"`              // 视频播放地址
	CoverUrl      string `json:"cover_url" db:"cover_url"`            // 视频封面地址
	FavoriteCount int64  `json:"favorite_count" db:"favourite_count"` // 视频点赞总数
	CommentCount  int64  `json:"comment_count" db:"comment_count"`    // 视频评论总数
	IsFavorite    bool   `json:"is_favorite" db:"isfavourite"`        // 是否已点赞
}

var Db *sqlx.DB

//通过token解析后获取用户信息，根据该信息查询数据库中的用户视频列表，封装在一起后返回
//PS:	token解析还没写，现在当做没有token解析这步
func getUserVideoInfoByToken(token string) []models.Video {

	//初始化数据库连接
	database, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
	//若数据库连接失败，报错，返回空
	if err != nil {
		fmt.Println("连接数据库失败：" + err.Error())
		return []models.Video{}
	}
	//连接成功后赋值
	Db = database
	var videos []models.Video //用户视频列表，即最终返回值
	var videoVos []VideoVo    //视频列表信息，用来接收数据库传来的值
	var user []models.User    //用户信息，go的DB操作要求数据库查出来必须用数组类型接收，大概吧，所以就用这个，

	//从数据库查视频信息，用videoVos接收数据
	err = Db.Select(&videoVos, "select * from video where user_id=?", 1)
	//如果查询失败，报错，返回空
	if err != nil {
		fmt.Println("exec failed, ", err)
		return []models.Video{}
	}
	//从数据库查用户信息，这里应该是根据token解析出来的信息来查用户所有信息，还没写token
	err = Db.Select(&user, "select * from user where id=?", 1)
	//如果查询失败，报错，返回空
	if err != nil {
		fmt.Println("exec failed, ", err)
		return []models.Video{}
	}

	//关闭数据库连接，必须卸载err之后的下面，我也不知道为什么，反正文档是这么说的
	defer Db.Close()

	//封装返回值videos，即将视频信息videoVos和用户信息user合并成要求的形式返回
	for _, v := range videoVos {
		t := models.Video{
			ID:            v.Id,
			Author:        user[0],
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
		}
		videos = append(videos, t)
	}

	for _, v := range videos {
		fmt.Println(v)
	}

	//over
	return videos
}

// 真正的pulishList   controller
func PublishList(c *gin.Context) {
	// 声明接收的变量
	var token string = "qqq"

	//获取发布视频信息
	videos := getUserVideoInfoByToken(token)

	//接口返回
	c.JSON(http.StatusOK, models.VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
