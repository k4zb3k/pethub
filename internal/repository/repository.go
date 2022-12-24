package repository

import (
	"errors"
	"github.com/k4zb3k/pethub/internal/models"
	"github.com/k4zb3k/pethub/pkg/logging"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

var logger = logging.GetLogger()

type Repository struct {
	Connection *gorm.DB
}

func NewRepository(conn *gorm.DB) *Repository {
	return &Repository{Connection: conn}
}

type Temp struct {
	ID       string
	Name     string
	Password string
}

type Token struct {
	ID     string
	UserId string
	Expire time.Time
}

//================================================

func (r Repository) Registration(name, username string, password []byte, phone string) error {
	sqlQuery := `INSERT INTO users (name, username, password, phone)
				 VALUES (?, ?, ?, ?);`
	tx := r.Connection.Exec(sqlQuery, name, username, password, phone)
	if tx.Error != nil {
		logger.Error(tx.Error)
		return tx.Error
	}

	return nil
}

func (r Repository) Login(username string, password string) (user *models.User, err error) {
	sqlQuery := `select  *from users where username = ?;`

	if err := r.Connection.Raw(sqlQuery, username).Scan(&user).Error; err != nil {
		logger.Infoln("incorrect login or password", err)
		return &models.User{}, err
	}
	if user == nil {
		logger.Infoln("incorrect login or password")
		return &models.User{}, errors.New("incorrect login or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Infoln("invalid login or password", err)
		return &models.User{}, err
	}

	return user, nil
}

func (r Repository) ValidateUsername(username string) (bool, error) {
	var usernameFromDB string
	sqlQuery := `select username from users where username = ?`
	tx := r.Connection.Raw(sqlQuery, username).Scan(&usernameFromDB)
	if tx.Error != nil {
		return true, tx.Error
	}
	if usernameFromDB != "" {
		logger.Infoln("login already registered")
		return true, errors.New("login already registered")
	}
	return false, nil
}

func (r Repository) SetToken(userID int, token string) error {
	sqlQuery := `insert into tokens (user_id, token)
				 values (?, ?);`
	tx := r.Connection.Exec(sqlQuery, userID, token)
	if tx.Error != nil {
		logger.Error("tx error", tx.Error)
		return tx.Error
	}

	return nil
}

func (r Repository) ValidateToken(token string) (int, error) {
	var tokenDB models.Token
	sqlQuery := `select *from tokens where token =?;`

	tx := r.Connection.Raw(sqlQuery, token).Scan(&tokenDB)
	if tx.Error != nil {
		logger.Infoln(tx.Error)
		return 0, tx.Error
	}

	if time.Now().After(tokenDB.Expire) {
		logger.Infoln("token expired")
		return 0, errors.New("token expired")
	}

	return tokenDB.UserID, nil
}

func (r Repository) AddNewAd(UserID int, TypeID int, PetID int, CityID int, title string, description string, photoPath string, reward int) error {
	sqlQuery := `insert into ads (user_id, type_id, pet_id, city_id, title, description, photo_path, reward) values (?,?,?,?,?,?,?,?)`
	tx := r.Connection.Exec(sqlQuery, UserID, TypeID, PetID, CityID, title, description, photoPath, reward)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r Repository) EditAdd(ads *models.Ads) error {
	sqlQuery := `update ads set type_id = ?, pet_id = ?, city_id = ?, title = ?, description = ?, photo_path = ?, reward = ? where id = ? and user_id = ?`
	tx := r.Connection.Exec(sqlQuery, ads.TypeId, ads.PetId, ads.CityId, ads.Title, ads.Description, ads.PhotoPath, ads.Reward, ads.ID, ads.UserID)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r Repository) Paginate(count int, page int) ([]models.Ads, error) {
	var ads []models.Ads

	sqlQuery := `select *from ads limit ? offset ? where is_active = true`

	tx := r.Connection.Raw(sqlQuery, count, (page-1)*count).Scan(&ads)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return ads, nil
}

func (r Repository) Search(q string) ([]models.Ads, error) {
	var ads []models.Ads

	sqlString := "%" + q + "%"

	sqlQuery := `select *from ads where title like ? or description like ? and is_active = true`

	tx := r.Connection.Raw(sqlQuery, sqlString, sqlString).Scan(&ads)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ads, nil
}

func (r Repository) DeleteAd(userId, adsId int) error {
	sqlQuery := `update ads set is_active = false where user_id = ? and id = ?`
	tx := r.Connection.Exec(sqlQuery, userId, adsId)
	if tx.Error != nil {
		log.Println("error in repos delete func")
		return tx.Error
	}
	return nil
}

func (r Repository) Filter(pet int) ([]models.Ads, error) {
	var ads []models.Ads

	sqlString := `select *from ads where pet_id = ? and is_active = true`
	tx := r.Connection.Raw(sqlString, pet).Scan(&ads)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return ads, nil
}

func (r Repository) ValidateFile(filename string) (int, error) {
	var count int
	sqlQuery := `select count(*) as count from ads where photo_path = ?`

	tx := r.Connection.Raw(sqlQuery, filename).Scan(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if count > 0 {
		return count, errors.New("файл уже существует")
	}
	return count, nil
}

func (r Repository) GetMyAds(UserId int) ([]models.Ads, error) {
	var ads []models.Ads
	sqlQuery := `select *from ads where user_id = ?`
	tx := r.Connection.Raw(sqlQuery, UserId).Scan(&ads)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return ads, nil
}
