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

func (model *BannerModel) PostBanner(tagIDs []int, featureID int, content json.RawMessage, isActive bool) (ID int, err error) {
	tx, err := model.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO banners (id, content, is_active) VALUES (DEFAULT, $1, $2) RETURNING id`
	err = model.DB.QueryRow(query, content, isActive).Scan(&ID)
	if err != nil {
		return
	}
	query = `INSERT INTO banner_tag_feature (banner_id, tag_id, feature_id) VALUES ($1, $2, $3)`
	for i := range tagIDs {
		_, err = model.DB.Exec(query, ID, tagIDs[i], featureID)
		if err != nil {
			return
		}
	}

	return ID, tx.Commit()
}

func (model *BannerModel) PatchBanner(ID, featureID int, tagIDs []int, content json.RawMessage, isActive bool) (err error) {
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `UPDATE banners SET content = $1, is_active = $2 WHERE id = $3`
	if _, err = model.DB.Exec(query, content, isActive, ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return
	}
	query = `DELETE FROM banner_tag_feature WHERE banner_id = $1`
	if _, err = model.DB.Exec(query, ID); err != nil {
		return
	}
	query = `INSERT INTO banner_tag_feature (banner_id, tag_id, feature_id) VALUES `
	valueStrings := make([]string, 0, len(tagIDs))
	valueArgs := make([]any, 0, len(tagIDs)*3)
	for i := range tagIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, ID, tagIDs[i], featureID)
	}
	if _, err = model.DB.Exec(query+strings.Join(valueStrings, ","), valueArgs...); err != nil {
		return
	}

	return tx.Commit()
}

func (model *BannerModel) DeleteBanner(ID int) error {
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `DELETE FROM banners WHERE id = $1`
	if _, err = model.DB.Exec(query, ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}
	query = `DELETE FROM banner_tag_feature WHERE banner_id = $1`
	if _, err = model.DB.Exec(query, ID); err != nil {
		return err
	}

	return tx.Commit()
}
