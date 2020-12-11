package router

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/mall-api/internal/router/api/admin"
)

func pmsProductRouter(r *gin.RouterGroup) {
	r1 := r.Group("/product")
	r1.POST("/create", admin.PmsProductController.Create)
}