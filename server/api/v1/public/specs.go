package public

import (
	"net/http"
	"oneclickvirt/constant"
	"oneclickvirt/model/common"

	"github.com/gin-gonic/gin"
)

// GetSpecs 获取资源规格配置（公开接口）
// @Summary 获取资源规格配置
// @Description 获取CPU、内存、磁盘、带宽的预定义规格选项
// @Tags 公开接口
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Router /v1/public/specs [get]
func GetSpecs(c *gin.Context) {
	specs := gin.H{
		"cpu":       constant.PredefinedCPUSpecs,
		"memory":    constant.PredefinedMemorySpecs,
		"disk":      constant.PredefinedDiskSpecs,
		"bandwidth": constant.PredefinedBandwidthSpecs,
	}

	c.JSON(http.StatusOK, common.Success(specs))
}
