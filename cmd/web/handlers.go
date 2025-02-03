package main

import (
	"backend_spring_2024/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type Response struct {
	Description string `json:"description"`
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
	fmt.Fprintln(w, "getBanners")
}

func (app *Application) postBanner(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "postBanner")
}

func (app *Application) patchBannerByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(Response{Description: "OK"})
}

func (app *Application) deleteBannerByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(Response{Description: "OK"})
}
