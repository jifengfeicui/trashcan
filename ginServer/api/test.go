package api

import (
	"github.com/gin-gonic/gin"

	"trashcan/ginServer/api/common"
)

func Test(c *gin.Context) {
	common.Ok(c)
}
