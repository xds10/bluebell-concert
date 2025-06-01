package mysql

import (
	"bluebell/models"
	"database/sql"
	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	if err = db.Select(&data, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			zap.L().Warn("there is no community")
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (*models.CommunityDetail, error) {
	communityDetail := new(models.CommunityDetail)
	sqlStr := `select 
    			community_id,community_name,introduction,create_time 
				from community 
				where community_id = ?`
	err := db.Get(communityDetail, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return communityDetail, err
}
