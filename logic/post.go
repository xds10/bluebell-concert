package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1.生成post_id
	p.ID = snowflake.GenID()
	// 2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	return err
	// 3.返回
}

func GetPostByID(pid int64) (*models.ApiPostDetail, error) {
	// 查询并组合接口数据
	data := &models.ApiPostDetail{}
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed", zap.Error(err))
		return nil, err
	}
	//根据作者id查询
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(pid) failed", zap.Error(err))
		return nil, err
	}

	//根据社区id查询社区详情
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Error(err))
		return nil, err
	}
	data.AuthorName = user.Username
	data.CommunityDetail = community
	data.Post = post
	return data, nil
}
func GetPostList(offset, limit int64) ([]*models.ApiPostDetail, error) {
	// 查询并组合接口数据
	posts, err := mysql.GetPostList(offset, limit)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed", zap.Error(err))
		return nil, err
	}
	data := make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(pid) failed", zap.Error(err))
			continue
		}

		//根据社区id查询社区详情
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			CommunityDetail: community,
			Post:            post,
		}
		data = append(data, postDetail)
	}

	return data, nil
}
