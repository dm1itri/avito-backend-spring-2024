package main

import (
	"backend_spring_2024/internal/models"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type Response struct {
	Description string `json:"description"`
}

type Banner struct {
	TagIDs    []int           `json:"tag_ids"`
	FeatureID int             `json:"feature_id"`
	IsActive  bool            `json:"is_active"`
	Content   json.RawMessage `json:"content"`
}

func (app *Application) getUserBanner(w http.ResponseWriter, r *http.Request) {
	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil || tagID < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil || featureID < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	useLastRevision := r.URL.Query().Get("use_last_revision") == "true"
	isAdmin := app.isRole(r, "admin")
	banner, err := app.banners.Get(tagID, featureID, useLastRevision, isAdmin)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	if err = json.NewEncoder(w).Encode(banner); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func (app *Application) getBanner(w http.ResponseWriter, r *http.Request) {
	var err error
	featureID := -1
	if featureIDstr := r.URL.Query().Get("feature_id"); featureIDstr != "" {
		featureID, err = strconv.Atoi(featureIDstr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	tagID := -1
	if tagIDstr := r.URL.Query().Get("tag_id"); tagIDstr != "" {
		tagID, err = strconv.Atoi(tagIDstr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	limit := -1
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	offSet := -1
	if offSetStr := r.URL.Query().Get("offset"); offSetStr != "" {
		offSet, err = strconv.Atoi(offSetStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	banners, err := app.banners.GetBanners(tagID, featureID, limit, offSet)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	if err = json.NewEncoder(w).Encode(banners); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func (app *Application) postBanner(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var banner Banner
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.banners.PostBanner(banner.TagIDs, banner.FeatureID, banner.Content, banner.IsActive)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *Application) patchBannerByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(Response{Description: "OK"}); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func (app *Application) deleteBannerByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(Response{Description: "OK"}); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
