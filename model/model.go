package model

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	jsoniter "github.com/json-iterator/go"
)

//Config ...Database login info
type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		Address  string `json:"address"`
	} `json:"database"`
}

//Items ...Used by gorm and json
type Items struct {
	ItemID          int    `gorm:"type:int(11)" json:"itemid"`
	ItemName        string `gorm:"type:varchar(255) json:"itemname"`
	ItemSaleStatus  string `gorm:"type:varchar(30) json:"itemsalestatus"`
	ItemCondition   string `gorm:"type:varchar(30)" json:"itemcondition"`
	CategoriesID    int    `gorm:"type:int(11)" json:"categoriesid"`
	ItemDescription string `gorm:"type:text"`
}

//ItemImage ...Used by gorm and json
type ItemImage struct {
	ItemID    int    `gorm:"type:int(11)" json:"itemid"`
	ImageLink string `gorm:"text" json:"imagelink"`
}

type Result struct {
	Item  []Items
	Image []ItemImage
}

//UserCommon ...Used by gorm and json
type UserCommon struct {
	UserID    string `gorm:"type:varchar(255)" json:"userid"`
	UserName  string `gorm:"type:varchar(100)" json:"name"`
	UserPhone string `gorm:"type:varchar(15)" json:"phone"`
	//UserBirth 			time.Time `gorm:"type:date" json:"birthdate"`
	UserGender       byte   `gorm:"type:char(1)" json:"gender"`
	UserAddress      string `gorm:"type:varchar(255)" json:"address"`
	UserPassword     string `gorm:"type:varchar(255)" json:"password"`
	UserAccessLevel  int    `gorm:"type:int" json:"accesslevel"`
	UserSessionToken string `gorm:"type:TEXT" json:"token"`
}

//BidSession ...Used by gorm and json
type BidSession struct {
	SessionID        int       `gorm:"type:int(11)" json:"sessionid"`
	ItemID           int       `gorm:"type:int(11)" json:"itemid"`
	SellerID         string    `gorm:"type:varchar(255)" json:"sellerid"`
	SessionStatus    string    `gorm:"type:varchar(30)" json:"sessionstatus"`
	SessionStartDate time.Time `gorm:"type:datetime" json:"startdate"`
	SessionEndDate   time.Time `gorm:"type:datetime" json:"enddate"`
}

//UserWishlist ...used by gorm and json
type UserWishlist struct {
	UserID string `gorm:"type:varchar(255)" json:"userid"`
	ItemID int    `gorm:"type:int(11)" json:"itemid"`
}

//SignupLoginResponse ...Respond form
type SignupLoginResponse struct {
	ResponseTime string     `json:"responseTime"`
	Code         int        `json:"code"`
	Message      string     `json:"message"`
	Data         UserCommon `json:"data"`
}

//AuthorizationHeader ...Used to get session token in header
type AuthorizationHeader struct {
	Token string `header:"Authorization"`
}

var (
	SecretKey = "thonking"
)

func DecodeDataFromJsonFile(f *os.File, data interface{}) error {
	jsonParser := jsoniter.NewDecoder(f)
	err := jsonParser.Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func SetupConfig() Config {
	var conf Config

	// Đọc file config.dev.json
	configFile, err := os.Open("config.local.json")
	if err != nil {
		// Nếu không có file config.dev.json thì đọc file config.default.json
		configFile, err = os.Open("config.default.json")
		if err != nil {
			panic(err)
		}
		defer configFile.Close()
	}
	defer configFile.Close()

	// Parse dữ liệu JSON và bind vào conf
	err = DecodeDataFromJsonFile(configFile, &conf)
	if err != nil {
		log.Println("Không đọc được file config.")
		panic(err)
	}

	return conf
}

func ConnectDb(user string, password string, database string, address string) *gorm.DB {
	connectionInfo := fmt.Sprintf(`%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local`, user, password, address, database)

	db, err := gorm.Open("mysql", connectionInfo)
	if err != nil {
		panic(err)
	}
	return db
}
