package image

// DownloadImageRequest 下载镜像请求
type DownloadImageRequest struct {
	ImageID      uint   `json:"imageId"`      // 系统镜像ID
	ProviderType string `json:"providerType"` // provider类型
	InstanceType string `json:"instanceType"` // 实例类型
	Architecture string `json:"architecture"` // 架构
}
