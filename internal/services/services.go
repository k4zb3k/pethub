package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/k4zb3k/pethub/internal/models"
	"github.com/k4zb3k/pethub/internal/repository"
	"github.com/k4zb3k/pethub/pkg/logging"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

var logger = logging.GetLogger()

type Services struct {
	Repository *repository.Repository
}

func NewServices(rep *repository.Repository) *Services {
	return &Services{Repository: rep}
}

//============================================================

func (s *Services) Register(userInfo *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	err = s.Repository.Registration(userInfo.Name, userInfo.Username, hash, userInfo.Phone)
	if err != nil {
		return err
	}

	return nil
}

func (s *Services) Login(userInfo *models.User) (string, error) {
	userFromDB, err := s.Repository.Login(userInfo.Username, userInfo.Password)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	userInfo.ID = userFromDB.ID

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(userInfo.Password))
	if err != nil {
		return "", errors.New("incorrect password")
	}

	buf := make([]byte, 256)

	_, err = rand.Read(buf)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(buf)

	return token, nil
}

func (s *Services) ValidateLogin(login string) (bool, error) {
	isLoginUsed, err := s.Repository.ValidateUsername(login)
	if err != nil {
		return true, err
	}
	if isLoginUsed {
		return true, errors.New("login already registered")
	}
	return false, nil
}

func (s *Services) ValidateToken(token string) (int, error) {
	userId, err := s.Repository.ValidateToken(token)
	if err != nil {
		log.Println("token is not available", err)
		return 0, err
	}

	return userId, nil
}

func (s *Services) SetToken(token *models.Token) error {
	err := s.Repository.SetToken(token.UserID, token.Token)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *Services) AddNewAd(AdsInfo *models.Ads) error { //post_ad
	//userId, err := s.Repository.ValidateToken()

	err := s.Repository.AddNewAd(AdsInfo.UserID, AdsInfo.TypeId, AdsInfo.PetId, AdsInfo.CityId, AdsInfo.Title, AdsInfo.Description, AdsInfo.PhotoPath, AdsInfo.Reward)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Info("Объявление успешно опубликовано")
	return nil
}

func (s *Services) EditAd(AdsInfo *models.Ads) error {
	err := s.Repository.EditAdd(AdsInfo)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Info("Объявление успешно изменено")
	return nil
}

func (s *Services) DeleteAd(userId, adsId int) error {
	err := s.Repository.DeleteAd(userId, adsId)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Info("Успешно удалено")
	return nil

}

func (s *Services) Paginate(count int, page int) ([]models.Ads, error) {
	ads, err := s.Repository.Paginate(count, page)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (s *Services) Search(q string) ([]models.Ads, error) {
	ads, err := s.Repository.Search(q)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return ads, nil
}

func (s *Services) GetMyAds(userId int) ([]models.Ads, error) {
	ads, err := s.Repository.GetMyAds(userId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return ads, nil
}

func (s *Services) Filter(pet int) ([]models.Ads, error) {
	ads, err := s.Repository.Filter(pet)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (s *Services) ValidateFile(path string) (int, error) {
	count, err := s.Repository.ValidateFile(path)
	if err != nil {
		return 0, err
	}
	return count, nil
}
