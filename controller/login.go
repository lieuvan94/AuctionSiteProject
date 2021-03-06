package controller

import (
	"AuctionSiteProject/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//	@Tags Logging Controller
// 	@Summary Tạo tài khoản mới
//  @Description Tạo tài khoản mới bằng JSON form, trả về session token nếu thành công.
//  @Param RegisterInfo body model.RegisterForm true "Resgister Information"
//  @Success 200 {body} string "Session token"
//	@Failure 400 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /signup [POST]
func RegisterJSON(c *gin.Context) {
	db := GetDBInstance().Db
	var newUser model.RegisterForm

	//check if the json form is valid
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Not a valid JSON!",
		})
		return
	}

	//Check for empty field and password length
	if newUser.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Vui lòng nhập tên",
		})
		return
	}
	if len(newUser.UserPassword) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Mật khẩu phải có tối thiểu 4 kí tự",
		})
		return
	}

	//Fetch userdata from database to check for existing username
	if exist, err := checkUserByID(newUser.UserID); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Error while fetching user data",
		})
		return
	} else if exist == true {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Tên truy cập đã có người sử dụng!",
		})
		return
	}

	//Encrypt the password
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.UserPassword), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Encrypt passsword error",
		})
		return
	}
	passwordHash := string(hash)

	//Filling information in struct
	newUserform := model.UserCommon{
		UserID:          newUser.UserID,
		UserName:        "",
		UserPhone:       "",
		UserEmail:       newUser.UserEmail,
		UserGender:      0,
		UserAddress:     "",
		UserPassword:    passwordHash,
		UserAccessLevel: 1,
		UserCreateat:    time.Now(),
	}
	UserSessionToken, _ := tokenGenerate(newUser.UserID)

	//Save account info to database
	errInsertDb := db.Table("user_common").Create(newUserform).Error
	if errInsertDb != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errInsertDb,
			"message": "Error: Cannot save new user",
		})
		return
	}

	//Generating success repond
	var rsp = model.SignupLoginResponse{
		ResponseTime: time.Now().String(),
		Code:         0,
		Message:      "Đăng kí thành công",
		Data:         newUserform,
		SessionToken: UserSessionToken,
	}

	c.JSON(http.StatusOK, rsp)
	return
}

//	@Tags Logging Controller
//	@Summary Đăng nhập bằng JSON
//  @Description Đăng nhập bằng JSON form, trả về session token nếu thành công.
//  @Param LoginForm body model.LoginForm true "Login Information"
//  @Param password body string true "password"
//  @Success 200 {body} string "Session token"
//	@Failure 400 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /login [POST]
func LoginJSON(c *gin.Context) {
	//db := GetDBInstance().Db
	var userLogin model.LoginForm

	err := c.BindJSON(&userLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Not a valid JSON!",
		})
		return
	}
	//Validing Login Info
	if exist, err := validLoginInfo(userLogin); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Error while fetching user data",
		})
		return
	} else if exist == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Tên truy cập hoặc mật khẩu không đúng!",
		})
		return
	}
	//Generate token
	var token string
	if token, err = tokenGenerate(userLogin.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Error,cannot create login session!",
		})
		return
	}

	//return Session token
	c.JSON(200, gin.H{
		"sessiontoken": token,
	})
	return
}
