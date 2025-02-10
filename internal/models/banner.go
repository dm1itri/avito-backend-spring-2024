package models

import (
	"database/sql"
	_ "embed"
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

//go:embed queries/get_user_banner.sql
var getUserBanner string

//go:embed queries/get_user_banner_admin.sql
var getUserBannerAdmin string

//go:embed queries/get_banners.sql
var getBanners string

//go:embed queries/post_banner.sql
var postBanner string

//go:embed queries/patch_banner.sql
var patchBanner string

//go:embed queries/delete_banner.sql
var deleteBanner string

func (model *BannerModel) Get(tagID, featureID int, useLastRevision, isAdmin bool) (row json.RawMessage, err error) {
	if isAdmin {
		err = model.DB.QueryRow(getUserBannerAdmin, tagID, featureID).Scan(&row)
	} else {
		err = model.DB.QueryRow(getUserBanner, tagID, featureID).Scan(&row)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return
}

func (model *BannerModel) GetBanners(tagID, featureID, limit, offSet int) ([]json.RawMessage, error) {
	query := getBanners
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
	query += " " + strings.Join(keys, " AND ")
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

	query := strings.Split(postBanner, ";")
	err = model.DB.QueryRow(query[0], content, isActive).Scan(&ID)
	if err != nil {
		return
	}
	valueStrings := make([]string, 0, len(tagIDs))
	valueArgs := make([]any, 0, len(tagIDs)*3)
	for i := range tagIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, ID, tagIDs[i], featureID)
	}
	if _, err = model.DB.Exec(query[1]+strings.Join(valueStrings, ","), valueArgs...); err != nil {
		return
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

	query := strings.Split(patchBanner, ";")
	if _, err = model.DB.Exec(query[0], content, isActive, ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return
	}
	if _, err = model.DB.Exec(query[1], ID); err != nil {
		return
	}
	valueStrings := make([]string, 0, len(tagIDs))
	valueArgs := make([]any, 0, len(tagIDs)*3)
	for i := range tagIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, ID, tagIDs[i], featureID)
	}
	if _, err = model.DB.Exec(query[2]+" "+strings.Join(valueStrings, ","), valueArgs...); err != nil {
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

	query := strings.Split(deleteBanner, ";")
	if _, err = model.DB.Exec(query[0], ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}
	if _, err = model.DB.Exec(query[1], ID); err != nil {
		return err
	}
	return tx.Commit()
}
