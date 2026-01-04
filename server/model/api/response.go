package api

// RouteInfo API路由信息
type RouteInfo struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Group  string `json:"group"`
	Name   string `json:"name"`
}
