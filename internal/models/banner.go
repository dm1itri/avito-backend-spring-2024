package models

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type Banner struct {
	ID      uint
	content []byte
}

type BannerModel struct {
	DB *sql.DB
}

func (model *BannerModel) Get(tagID, featureID int, useLastRevision, isAdmin bool) (json.RawMessage, error) {
	query := `SELECT b.content FROM banners b
    		  JOIN banner_tag_feature btf ON b.id = btf.banner_id AND b.is_active = true
        	  WHERE btf.tag_id = $1 AND btf.feature_id = $2`
	queryAdmin := `SELECT b.content FROM banners b
    		  JOIN banner_tag_feature btf ON b.id = btf.banner_id
        	  WHERE btf.tag_id = $1 AND btf.feature_id = $2`
	var row json.RawMessage
	err := model.DB.QueryRow(func(isAdmin bool) string {
		if isAdmin {
			return queryAdmin
		}
		return query
	}(isAdmin), tagID, featureID).Scan(&row)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return row, nil
}
