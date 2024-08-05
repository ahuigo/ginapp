package router

import (
	"ginapp/store"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Stock struct {
	Code      string `gorm:"primary_key" `
	Price     uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func getDb() (db *gorm.DB, err error) {
	Db := store.GetPgDB()
	return Db.DB, nil
}

func insertHandler(c *gin.Context) {
	res := "ok"
	_, err := getDb()
	if err != nil {
		res = err.Error()
	}

	c.String(http.StatusOK, res)

}
