package test

import (
	"context"
	"github.com/onqlavetest1/config"
	"github.com/onqlavetest1/persistence"
	"github.com/onqlavetest1/persistence/redisFactory"
	"log"
	"reflect"
	"testing"
)

func SubTesting(test *testing.M) {
	persistence.InitHandlerFactory()
	//Unreachable Address Testing
	factory := persistence.GetFactory().WithREDIS(config.REDISHandlerFactoryConfigs{Addr: "1.2.3.4:231"}).Subscriber("test").Create()

	err := factory.Start(context.Background())

	if !reflect.DeepEqual(err, redisFactory.ERR_CONNECT_SERVER) {
		log.Fatalln("Test case fail err: ", err, " expected ", redisFactory.ERR_CONNECT_SERVER)
	}

	//Tried to stop the sub before running it
	factory = persistence.GetFactory().WithREDIS(config.REDISHandlerFactoryConfigs{}).Subscriber("test").Create()

	err = factory.Stop()

	if !reflect.DeepEqual(err, redisFactory.ERR_SUB_WAITING) {
		log.Fatalln("Test case fail err: ", err, " expected ", redisFactory.ERR_SUB_WAITING)
	}

	err = factory.Start(context.Background())

	if !reflect.DeepEqual(err, redisFactory.ERR_SUB_GENERAL) {
		log.Fatalln("Test case fail err: ", err, " expected ", redisFactory.ERR_SUB_GENERAL)
	}

	return
}
