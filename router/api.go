package router

import (
	"fmt"
	"ginapp/conf"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func panicApi(c *gin.Context) {
	filepath := c.DefaultQuery("path", "tmp/a.txt")
	buf, err := os.ReadFile(filepath)
	if err != nil {
		err = errors.Wrap(err, "bad file argument!")
		panic(err)
		// c.String(http.StatusOK, err.Error())
	}
	res := string(buf)
	c.String(http.StatusOK, res)

}
func confApi(c *gin.Context) {
	c.JSON(http.StatusOK, conf.GetConf())
}

type HTTPError struct {
	Error string
}

type User struct {
	ID int `json:"id" example:"1" format:"int64"`
	// 姓名
	Name string `json:"name" form:"name" example:"Alex"`
	// 班名
	ClassName string `json:"class_name" form:"class_name" example:"class1"`
	// 国家
	Country string
	// time rfc3339
	Time time.Time `json:"time" form:"time"`

	// 这是city 说明
	City string `form:"city" example:"Beijing"`
	// form 接收多个值
	Citys []string `form:"city" example:"Bj,Tj"`
	// 这是extra扩展字段
	Extra interface{} `json:"extra"`
}

/*Param Type: query path header body formData
*Refer to: https://github.com/swaggo/swag/blob/master/README.md#param-type
 */

// @Summary      获取用户
// @Description  获取用户详情
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 Cookie header string  	false "token"     default(token=xxx)
// @Param        id    	path      int  	true  "Account ID" Enums(1, 2, 3)
// @Param        trace_id  	query string  true  "trace-id" default(tracd_id_10)
// @Param        queryStr  query  User  true  "query struct" default("name=Alex&page=1")
// @Success      200  {object}  User
// @Header 		 200 {string} Token "qwerty"
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {string}  string "500 error"
// @Router       /user/{id} [get]
func GetUser(c *gin.Context) {
	user := User{}
	if err := c.ShouldBind(&user); err != nil {
		fmt.Println("bind error:", err)
	}
	var res User
	if user.ID == 1 {
		res = User{
			Name: "Alex",
		}
	} else {
		res = User{
			Name: "Other user",
		}
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 添加用户(body)
// @Tags user
// @Accept 		json
// @Produce		json
// @Param body-params body User true "Add user"
// @Success 200 {object} User
// @Failure 400  {object}  HTTPError
// @Router /user [post]
func AddUser(c *gin.Context) {
	user := User{}
	if err := c.ShouldBind(&user); err != nil {
		fmt.Println("bind error:", err)
	}
	c.JSON(http.StatusOK, user)
}

// multipart/form-data
// @Tags user
// @Accept		multipart/form-data
// @Param myFormData   formData   User   true   "task form data"
// @Success 200 {object} User
// @Failure 400  {object}  HTTPError
// @Router /user/upload [post]
func UploadMultipartForm(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

// @Summary Get用户List(query)
// @Tags user
// @Param body-params query User true "user list query"
// @Success 200 {object} User
// @Failure 400  {object}  HTTPError
// @Router /users [get]
func GetUserList(c *gin.Context) {
	req := User{}
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("bind error:", err)
	}
	c.JSON(http.StatusOK, []User{{ID: 1, Name: "Alex"}, {ID: 2, Name: "Tom"}})
}
