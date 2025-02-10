package main

import (
	"backend_spring_2024/internal/models"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type Banner struct {
	TagIDs    []int           `json:"tag_ids"`
	FeatureID int             `json:"feature_id"`
	IsActive  bool            `json:"is_active"`
	Content   json.RawMessage `json:"content"`
}

type BannerID struct {
	BannerID int `json:"banner_id"`
}

func setDescriptionStatusCode(description string, statusCode int, w http.ResponseWriter) {
	w.Header().Set("Description", description)
	w.WriteHeader(statusCode)
}

func (app *Application) getUserBanner(w http.ResponseWriter, r *http.Request) {
	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil || tagID < 1 {
		setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
		return
	}
	featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil || featureID < 1 {
		setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
		return
	}
	useLastRevision := r.URL.Query().Get("use_last_revision") == "true"
	isAdmin := IsRole(r.Header.Get("token"), "admin", app.secretKeyJWT)
	banner, err := app.banners.Get(tagID, featureID, useLastRevision, isAdmin)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			setDescriptionStatusCode("The banner was not found", http.StatusNotFound, w)
		} else {
			setDescriptionStatusCode("Internal server error", http.StatusInternalServerError, w)
		}
		return
	}
	if err = json.NewEncoder(w).Encode(banner); err != nil {
		setDescriptionStatusCode("Internal server error", http.StatusInternalServerError, w)
		return
	}
}

func (app *Application) getBanner(w http.ResponseWriter, r *http.Request) {
	var err error
	var featureID, tagID, limit, offSet int
	if featureIDstr := r.URL.Query().Get("feature_id"); featureIDstr != "" {
		featureID, err = strconv.Atoi(featureIDstr)
		if err != nil {
			setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
			return
		}
	}
	if tagIDstr := r.URL.Query().Get("tag_id"); tagIDstr != "" {
		tagID, err = strconv.Atoi(tagIDstr)
		if err != nil {
			setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
			return
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
			return
		}
	}
	if offSetStr := r.URL.Query().Get("offset"); offSetStr != "" {
		offSet, err = strconv.Atoi(offSetStr)
		if err != nil {
			setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
			return
		}
	}
	banners, err := app.banners.GetBanners(tagID, featureID, limit, offSet)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			setDescriptionStatusCode("The banner was not found", http.StatusNotFound, w)
		} else {
			setDescriptionStatusCode("Internal server error", http.StatusInternalServerError, w)
		}
		return
	}
	if err = json.NewEncoder(w).Encode(banners); err != nil {
		setDescriptionStatusCode("Internal server error", http.StatusInternalServerError, w)
		return
	}
}

func (app *Application) postBanner(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
		return
	}

	var banner Banner
	if err := json.NewDecoder(r.Body).Decode(&banner); err != nil {
		setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
		return
	}
	bannerID, err := app.banners.PostBanner(banner.TagIDs, banner.FeatureID, banner.Content, banner.IsActive)
	if err != nil {
		setDescriptionStatusCode("Internal server error", http.StatusInternalServerError, w)
		return
	}
	setDescriptionStatusCode("Created", http.StatusCreated, w)
	json.NewEncoder(w).Encode(BannerID{BannerID: bannerID})
}

func (app *Application) patchBannerByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
		return
	}

	var banner Banner
	if err = json.NewDecoder(r.Body).Decode(&banner); err != nil {
		setDescriptionStatusCode("Incorrect data", http.StatusBadRequest, w)
		return
	}

	if err = app.banners.PatchBanner(id, banner.FeatureID, banner.TagIDs, banner.Content, banner.IsActive); err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			setDescriptionStatusCode("The banner for the tag was not found", http.StatusNotFound, w)
		} else {
			setDescriptionStatusCode("Internal server error", http.StatusInternalServerError, w)
		}
		return
	}
	setDescriptionStatusCode("OK", http.StatusOK, w)
}

func (app *Application) deleteBannerByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err = app.banners.DeleteBanner(id); err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			setDescriptionStatusCode("The banner for the tag was not found", http.StatusNotFound, w)
		} else {
			setDescriptionStatusCode("Internal server error", http.StatusInternalServerError, w)
		}
		return
	}
	setDescriptionStatusCode("The banner was successfully deleted", http.StatusNoContent, w)
}
