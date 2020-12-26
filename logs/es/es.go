package es

import (
	"github.com/elastic/go-elasticsearch/v6"
)

//记录到es的日志
type esLogger struct {
	*elasticsearch.Client
	DSN   string `json:"dsn"`
	Level int    `json:"level"`
	//formatter logs.LogFormatter
	Formatter string `json:"formatter"`

	//indexNaming IndexNaming
}
