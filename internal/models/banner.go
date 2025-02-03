package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

type Banner struct {
	ID      uint
	content []byte
}

type BannerModel struct {
	DB *sql.DB
}

func (model *BannerModel) Get(tagID, featureID int, useLastRevision, isAdmin bool) (row json.RawMessage, err error) {
	query := `SELECT b.content FROM banners b
    		  JOIN banner_tag_feature btf ON b.id = btf.banner_id AND b.is_active = true
        	  WHERE btf.tag_id = $1 AND btf.feature_id = $2`
	queryAdmin := `SELECT b.content FROM banners b
    		  JOIN banner_tag_feature btf ON b.id = btf.banner_id
        	  WHERE btf.tag_id = $1 AND btf.feature_id = $2`
	err = model.DB.QueryRow(func(isAdmin bool) string {
		if isAdmin {
			return queryAdmin
		}
		return query
	}(isAdmin), tagID, featureID).Scan(&row)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return
}

func (model *BannerModel) GetBanners(tagID, featureID, limit, offSet int) ([]json.RawMessage, error) {
	query := `SELECT b.content FROM banners b
    		  JOIN banner_tag_feature btf ON b.id = btf.banner_id
        	  WHERE `
	keys := make([]string, 0, 2)
	args := make([]any, 0, 4)
	if tagID != 0 {
		args = append(args, tagID)
		keys = append(keys, "btf.tag_id = $1 ")
	}
	if featureID != 0 {
		args = append(args, featureID)
		keys = append(keys, fmt.Sprintf("btf.feature_id = $%d ", len(args)))
	}
	query += strings.Join(keys, " AND ")
	if limit != 0 {
		args = append(args, limit)
		query += fmt.Sprintf("LIMIT $%d ", len(args))
	}
	if offSet != 0 {
		args = append(args, offSet)
		query += fmt.Sprintf("OFFSET $%d", len(args))
	}
	var rowsJSON []json.RawMessage
	rows, err := model.DB.Query(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	for rows.Next() {
		var rowJSON json.RawMessage
		err = rows.Scan(&rowJSON)
		if err != nil {
			log.Fatal("Ошибка сканирования данных:", err)
		}
		rowsJSON = append(rowsJSON, rowJSON)
	}
	return rowsJSON, nil
}

func (model *BannerModel) PostBanner(tagIDs []int, featureID int, content json.RawMessage, isActive bool) (bannerID int, err error) {
	query := `INSERT INTO banners (id, content, is_active) VALUES (DEFAULT, $1, $2) RETURNING id`
	err = model.DB.QueryRow(query, content, isActive).Scan(&bannerID)
	if err != nil {
		return
	}
	query = `INSERT INTO banner_tag_feature (banner_id, tag_id, feature_id) VALUES ($1, $2, $3)`
	for i := range tagIDs {
		_, err = model.DB.Exec(query, bannerID, tagIDs[i], featureID)
		if err != nil {
			return
		}
	}
	return
}

func (model *BannerModel) PatchBanner(bannerID, featureID int, tagIDs []int, content json.RawMessage, isActive bool) (err error) {
	query := `UPDATE banners SET content = $1, is_active = $2 WHERE id = $3`
	if _, err = model.DB.Exec(query, content, isActive, bannerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return
	}
	query = `DELETE FROM banner_tag_feature WHERE banner_id = $1`
	if _, err = model.DB.Exec(query, bannerID); err != nil {
		return
	}
	query = `INSERT INTO banner_tag_feature (banner_id, tag_id, feature_id) VALUES ($1, $2, $3)`
	for i := range tagIDs {
		if _, err = model.DB.Exec(query, bannerID, tagIDs[i], featureID); err != nil {
			return
		}
	}
	return
}
