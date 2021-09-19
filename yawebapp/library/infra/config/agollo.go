package config

import (
	"fmt"

	"github.com/shima-park/agollo"
)

type ApolloConf struct {
	UseApollo      bool     `toml:"use_apollo"`
	AppID          string   `toml:"app_id"`
	Cluster        string   `toml:"cluster"`
	NamespaceNames []string `toml:"namespace_names"`
	Host           string   `toml:"host"`
}

func NewAgollo(conf ApolloConf) (agollo.Agollo, error) {
	a, err := agollo.New(conf.Host, conf.AppID, agollo.PreloadNamespaces(conf.NamespaceNames...), agollo.AutoFetchOnCacheMiss(), agollo.FailTolerantOnBackupExists())
	if err != nil {
		return nil, err
	}

	go func() {
		errorCh := a.Start() // Start后会启动goroutine监听变化，并更新agollo对象内的配置cache
		watchCh := a.Watch()
		for {
			select {
			case err := <-errorCh:
				// handle error
				fmt.Println("listening agollo error: ", err)
			case resp := <-watchCh:
				fmt.Println(
					"Namespace:", resp.Namespace,
					"OldValue:", resp.OldValue,
					"NewValue:", resp.NewValue,
					"Error:", resp.Error,
				)
			}
		}
	}()

	return a, nil
}
