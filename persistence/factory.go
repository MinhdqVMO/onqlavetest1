package persistence

import (
	"github.com/onqlavetest1/config"
	"github.com/onqlavetest1/persistence/redisFactory"
)

var factory *handlerFactory

type handlerFactory struct {
}

type HandlerFactory interface {
	WithREDIS(configs config.REDISHandlerFactoryConfigs) WithREDISHandler
}

func InitHandlerFactory() {
	if factory != nil {
		return
	}
	factory = &handlerFactory{}
}

func GetFactory() *handlerFactory {
	return factory
}

func (h handlerFactory) WithREDIS(configs config.REDISHandlerFactoryConfigs) WithREDISHandler {
	if configs.MaxConcurrency == 0 {
		configs.MaxConcurrency = config.DEFAULT_CONCURRENCY
	}
	if configs.Network == "" {
		configs.Network = config.DEFAULT_NETWORK
	}
	return redisFactory.NewRedisHandlerFactory(configs)
}
