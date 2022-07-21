package redisFactory

import (
	"context"
	"github.com/onqlavetest1/config"
	"log"
	"reflect"
	"testing"
)

func TestSub(test *testing.T) {
	//Unreachable Address Testing
	factory := NewRedisHandlerFactory(config.REDISHandlerFactoryConfigs{Addr: "1.2.3.4:231"}).Subscriber("test").Create()
	err := factory.Start(context.Background())

	if !reflect.DeepEqual(err, ERR_CONNECT_SERVER) {
		log.Fatalln("Test case fail err: ", err, " expected ", ERR_CONNECT_SERVER)
	}

	// Client not initialized
	factory_test_2 := redisHandler{}

	err = factory_test_2.Start(context.Background())

	if !reflect.DeepEqual(err, ERR_INITIALIZE_CLIENT) {
		log.Fatalln("Test case fail err: ", err, " expected ", ERR_INITIALIZE_CLIENT)
	}

	//Tried to stop the sub before running it
	factory = NewRedisHandlerFactory(config.REDISHandlerFactoryConfigs{Addr: "1.2.3.4:231"}).Subscriber("test").Create()

	err = factory.Stop()

	if !reflect.DeepEqual(err, ERR_SUB_WAITING) {
		log.Fatalln("Test case fail err: ", err, " expected ", ERR_SUB_WAITING)
	}

	err = factory.Start(context.Background())

	if !reflect.DeepEqual(err, ERR_SUB_GENERAL) {
		log.Fatalln("Test case fail err: ", err, " expected ", ERR_SUB_GENERAL)
	}

	return
}
