package config

import "github.com/amanbolat/ca-warehouse-client/api"

type Config struct {
	FmHost         string `split_words:"true" required:"true"`
	FmUser         string `split_words:"true" required:"true"`
	FmPort         string `split_words:"true"`
	FmDatabaseName string `split_words:"true" required:"true"`
	FmPass         string `split_words:"true" required:"true"`
	Printer        string `split_words:"true" required:"true"`
	Debug          bool
	Port           int    `split_words:"true" required:"true"`
	FontPath       string `split_words:"true" required:"true"`
	api.KDNiaoConfig
}
