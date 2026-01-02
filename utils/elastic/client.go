package elastic

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
)

var (
	client     *elasticsearch.TypedClient
	initOnce   = &sync.Once{}
	defaultCfg = elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "", // 默认用户名，通常从配置文件加载
		Password:  "", // 默认密码，通常从配置文件加载
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
)

// InitClient 初始化 ES 客户端
func InitClient() error {
	var err error
	initOnce.Do(func() {
		cfg := defaultCfg

		// 从配置文件加载 ES 配置
		if viper.IsSet("elasticsearch.addresses") {
			cfg.Addresses = viper.GetStringSlice("elasticsearch.addresses")
		}
		if viper.IsSet("elasticsearch.username") && viper.IsSet("elasticsearch.password") {
			cfg.Username = viper.GetString("elasticsearch.username")
			cfg.Password = viper.GetString("elasticsearch.password")
		}
		if viper.IsSet("elasticsearch.ca_cert") {
			cfg.CACert = []byte(viper.GetString("elasticsearch.ca_cert"))
		}

		client, err = elasticsearch.NewTypedClient(cfg)
		if err != nil {
			log.Printf("Failed to create ES client: %v", err)
			return
		}

		// 测试连接
		_, err = client.Info().Do(context.Background())
		if err != nil {
			log.Printf("Failed to connect to ES: %v", err)
			client = nil
			return
		}

		log.Println("ES client initialized successfully")
	})
	return err
}

// GetClient 获取 ES 客户端实例
func GetClient() (*elasticsearch.TypedClient, error) {
	if client == nil {
		if err := InitClient(); err != nil {
			return nil, err
		}
	}
	return client, nil
}
