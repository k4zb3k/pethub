package server

import (
	"encoding/json"
	"fmt"
	"github.com/k4zb3k/pethub/internal/models"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (s *Server) Registration(w http.ResponseWriter, r *http.Request) {
	var user models.User

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
	}

	err = json.Unmarshal(bytes, &user)
	if err != nil {
		logger.Error(err)
		return
	}

	isLoginUsed, err := s.Services.ValidateLogin(user.Username)
	if err != nil {
		BadRequest(w, err)
		return
	}
	if isLoginUsed {
		BadRequest(w, err)
		return
	}

	err = s.Services.Register(&user)
	if err != nil {
		logger.Error(err)
		InternalServerError(w, err)
		return
	}
	_, err = w.Write([]byte("Successful registration"))
	if err != nil {
		logger.Error(err)
		//InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var user *models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		logger.Error(err)
		return
	}

	// body, err := json.NewEncoder(r.Body).Encode(user)

	token, err := s.Services.Login(user)
	if err != nil {
		logger.Error(err)

		message := fmt.Sprintf("%v", err)
		ErrResponse(w, 400, message)

		return
	}

	newToken := models.Token{
		UserID: user.ID,
		Token:  token,
	}

	if err = s.Services.SetToken(&newToken); err != nil {
		return
	}

	_, err = w.Write([]byte(token))
	if err != nil {
		return
	}

	w.WriteHeader(202)
}

func (s *Server) AddNewAd(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// получение userId из context

	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	file, handler, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile {
			ErrResponse(w, 400, "Запрос не содержит файла")
			logger.Error(err)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	defer file.Close()

	count, err := s.Services.ValidateFile("files/images/" + handler.Filename)
	if err != nil {
		logger.Error(err)
		return
	}

	if count > 0 {
		http.Error(w, "Файл уже существует", http.StatusBadRequest)
		logger.Error("Файл уже существует")
		return
	}

	if handler.Size > 5242880 { //если размер файла больше 5 мб
		http.Error(w, "Файл больше 5 мб", http.StatusBadRequest)
		logger.Info("Файл больше 5 мб")
		return
	}

	f, err := os.OpenFile("./files/images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error(err)
		return
	}

	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		return
	}

	ads := r.FormValue("adsInfo")

	var AdsInfo models.Ads
	AdsInfo.UserID = userID

	err = json.Unmarshal([]byte(ads), &AdsInfo)
	if err != nil {
		logger.Error(err)
		return
	}

	photoPath := "files/images/" + handler.Filename
	AdsInfo.PhotoPath = photoPath

	if AdsInfo.TypeId <= 0 || AdsInfo.PetId <= 0 || AdsInfo.CityId <= 0 {
		ErrResponse(w, 400, "введите корректные индентификаторы")
		logger.Infoln("введен некорректные индентификаторы")
		return
	}

	err = s.Services.AddNewAd(&AdsInfo)
	if err != nil {
		return
	}

	//log.Println(photoPath)

	w.Write([]byte("Объявление успешно опубликовано"))
}

func (s *Server) EditAd(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// получение userId из context

	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	file, handler, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile {
			http.Error(w, "Запрос не содержит файла", http.StatusBadRequest)
			logger.Error(err)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	defer file.Close()

	count, err := s.Services.ValidateFile("files/images/" + handler.Filename)
	if err != nil {
		logger.Error(err)
		return
	}

	if count > 0 {
		http.Error(w, "Файл уже существует", http.StatusBadRequest)
		logger.Error("Файл уже существует")
		return
	}

	if handler.Size > 5242880 { //если размер файла больше 5 мб
		http.Error(w, "Файл больше 5 мб", http.StatusBadRequest)
		logger.Info("Файл больше 5 мб")
		return
	}

	f, err := os.OpenFile("./files/images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error(err)
		return
	}

	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		return
	}

	ads := r.FormValue("adsInfo")

	var AdsInfo models.Ads
	AdsInfo.UserID = userID

	err = json.Unmarshal([]byte(ads), &AdsInfo)
	if err != nil {
		logger.Error(err)
		return
	}

	photoPath := "files/images/" + handler.Filename
	AdsInfo.PhotoPath = photoPath

	if AdsInfo.TypeId <= 0 || AdsInfo.PetId <= 0 || AdsInfo.CityId <= 0 {
		ErrResponse(w, 400, "введите корректные индентификаторы")
		logger.Infoln("введен некорректные индентификаторы")
		return
	}

	err = s.Services.EditAd(&AdsInfo)
	if err != nil {
		return
	}
	//w.Write([]byte(body))
	w.WriteHeader(200)
}

func (s *Server) DeleteAd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)
	AdsId, err := strconv.Atoi(r.URL.Query().Get("ads-id"))
	if err != nil {
		logger.Error(err)
		return
	}
	err = s.Services.DeleteAd(userID, AdsId)
	if err != nil {
		logger.Error(err)
		return
	}
	w.Write([]byte("Успешно удалено"))
	w.WriteHeader(202)
}

func (s *Server) Paginate(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query()
	count := name.Get("count")
	counts, err := strconv.Atoi(count)
	if err != nil {
		return
	}
	page := name.Get("page")
	pages, err := strconv.Atoi(page)
	if err != nil {
		return
	}
	ads, err := s.Services.Paginate(counts, pages)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(ads)
	if err != nil {
		return
	}
}

func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q := query.Get("q")

	ads, err := s.Services.Search(q)
	if err != nil {
		return
	}
	err = json.NewEncoder(w).Encode(ads)
	if err != nil {
		return
	}

}

func (s *Server) Filter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q, err := strconv.Atoi(query.Get("pet"))
	if err != nil {
		return
	}

	ads, err := s.Services.Filter(q)

	if err != nil {
		return
	}
	err = json.NewEncoder(w).Encode(ads)
	if err != nil {
		return
	}
}

func (s *Server) GetMyAds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)
	ads, err := s.Services.GetMyAds(userID)
	err = json.NewEncoder(w).Encode(ads)
	if err != nil {
		return
	}
}
