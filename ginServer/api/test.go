package api

import (
	"github.com/gin-gonic/gin"

	"template/ginServer/api/common"
)

func Test(c *gin.Context) {
	common.Ok(c)
}
