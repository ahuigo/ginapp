package router

import (
	// swag: https://github.com/swaggo/swag
	// gin-swagger: https://github.com/swaggo/swag/tree/master/example/celler
	"ginapp/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwagger(r *gin.Engine) {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "ahuigo API"
	docs.SwaggerInfo.Description = `
# markdown
1. item1
2. itme2
`
	docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "m:" + fmt.Sprint(*port)
	r.Use(func(c *gin.Context) {
		docs.SwaggerInfo.Host = c.Request.Host
	})

	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
