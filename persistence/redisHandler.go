package persistence

import (
	"github.com/onqlavetest1/persistence/handler"
)

type WithREDISHandler interface {
	Subscriber(subscriberName string) handler.RecieveAndCreateHandler
	Publisher() handler.CreateHandler
}
