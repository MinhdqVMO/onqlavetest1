package redisFactory

import (
	"context"
	"github.com/onqlavetest1/config"
	test2 "github.com/onqlavetest1/test"
	"log"
	"reflect"
	"testing"
)

func TestPub(test *testing.T) {
	command := test2.TestCommand{}
	//Unreachable Address Testing
	factory := NewRedisHandlerFactory(config.REDISHandlerFactoryConfigs{Addr: "1.2.3.4:231"}).Publisher().Create()

	err := factory.PublishCommand(context.Background(), &command)

	if !reflect.DeepEqual(err, ERR_CONNECT_SERVER) {
		log.Fatalln("Test case fail err: ", err, " expected ", ERR_CONNECT_SERVER)
	}

	// Client not initialized
	factory_test_2 := redisHandler{}

	err = factory_test_2.PublishCommand(context.Background(), &command)

	if !reflect.DeepEqual(err, ERR_INITIALIZE_CLIENT) {
		log.Fatalln("Test case fail err: ", err, " expected ", ERR_INITIALIZE_CLIENT)
	}

	//Can't publish to topic
	factory = NewRedisHandlerFactory(config.REDISHandlerFactoryConfigs{Addr: "1.2.3.4:231"}).Publisher().Create()

	err = factory.PublishCommand(context.Background(), &command)

	if !reflect.DeepEqual(err, ERR_PUBLISH_GENERAL) {
		log.Fatalln("Test case fail err: ", err, " expected ", ERR_PUBLISH_GENERAL)
	}

	return
}
