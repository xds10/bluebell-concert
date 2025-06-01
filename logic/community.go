package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	//查询数据库,找到所有的community
	communityList, err = mysql.GetCommunityList()
	if err != nil {

	}
	return
}

// GetCommunityDetail根据id查询社区详情
func GetCommunityDetail(id int64) (communityList *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetailByID(id)
}
