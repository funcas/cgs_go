package model

/*
 *	在线生成工具 https://oktools.net/json2go
 */

// http 连接器配置实体
type HttpConnectorVO struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	QueryParams string `json:"queryParams"`
	ContentType string `json:"contentType"`
	UserAgent   string `json:"userAgent"`
	Timeout     int    `json:"timeout"`
}
