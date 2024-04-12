package controllers

import (
	"banners/internal/database/models"
	"banners/internal/utils"
	"banners/internal/utils/cache"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var bannerCache cache.BannerCache

func NewBannerController(BannerCache cache.BannerCache) {
	bannerCache = BannerCache
}

var CreateBanner = func(w http.ResponseWriter, r *http.Request) {
	var banner models.Banner
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Respond(w, utils.MessageError("Некорректные данные"))
		return
	}

	if banner.Content == nil || banner.FeatureID == nil || banner.Tags == nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Respond(w, utils.MessageError("Некорректные данные"))
		return
	}

	resp, err := banner.Create()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.MessageError("Внутренняя ошибка сервера"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.Respond(w, resp)
}

var GetBanners = func(w http.ResponseWriter, r *http.Request) {

	//Default -1 -> не учитывать в SQL запросе
	feature_id := utils.ParseOrDefaultInt(r.URL.Query().Get("feature_id"), -1)
	tag_id := utils.ParseOrDefaultInt(r.URL.Query().Get("tag_id"), -1)
	limit := utils.ParseOrDefaultInt(r.URL.Query().Get("tag_id"), 10)
	offset := utils.ParseOrDefaultInt(r.URL.Query().Get("offset"), 0)

	data, err := models.GetBanners(feature_id, tag_id, false, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.MessageError("Внутренняя ошибка сервера"))
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.RespondBanners(w, data)
}

var GetUserBanner = func(w http.ResponseWriter, r *http.Request) {
	feature_id, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Respond(w, utils.MessageError("Некорректные данные"))
		return
	}
	tag_id, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Respond(w, utils.MessageError("Некорректные данные"))
		return
	}
	use_last_revision := utils.ParseOrDefaultBool(r.URL.Query().Get("use_last_revision"), false)

	key := fmt.Sprintf("feature_id=%d&tag_id=%d", feature_id, tag_id)

	var banner = bannerCache.Get(key)

	if banner == nil || use_last_revision {
		data, err := models.GetBanners(feature_id, tag_id, true, 1, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.Respond(w, utils.MessageError("Внутренняя ошибка сервера"))
			return
		}

		if len(*data) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		banner = &(*data)[0]
		bannerCache.Set(key, banner)
	}

	w.WriteHeader(http.StatusOK)
	resp := map[string]interface{}{"content": banner.Content}
	utils.Respond(w, &resp)
}

var UpdateBanner = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	banner_id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Respond(w, utils.MessageError("Некорректные данные"))
		return
	}
	var banner models.Banner
	err = json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Respond(w, utils.MessageError("Некорректные данные"))
		return
	}
	banner.ID = banner_id

	exist, err := models.Exist(banner.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.MessageError("Внутренняя ошибка сервера"))
		return
	}

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = banner.Save()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.MessageError("Внутренняя ошибка сервера"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

var DeleteBanner = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	banner_id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Respond(w, utils.MessageError("Некорректные данные"))
		return
	}

	var banner models.Banner
	banner.ID = banner_id

	exist, err := models.Exist(banner.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.MessageError("Внутренняя ошибка сервера"))
		return
	}

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = banner.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Respond(w, utils.MessageError("Внутренняя ошибка сервера"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
