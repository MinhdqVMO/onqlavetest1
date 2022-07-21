package main

import (
	"context"
	"fmt"
	"github.com/onqlavetest1/config"
	"github.com/onqlavetest1/model"
	"github.com/onqlavetest1/persistence"
	"github.com/onqlavetest1/test"
	"log"
	"testing"
	"time"
)

func TestMain(t *testing.M) {
	//to receive and react to commands
	command := test.TestCommand{}
	ctx := context.Background()
	persistence.InitHandlerFactory()
	factory := persistence.GetFactory()
	redisFactory := factory.
		WithREDIS(config.REDISHandlerFactoryConfigs{})

	handler := redisFactory.Subscriber("service").
		RecieveCommand(func() model.DomainCommand { return &command }, func() string { return "" }, func(ctx context.Context, message *model.MessageEnvelope, eventStream chan<- model.DomainEvent) error {
			//we handle the commands here and do whatever is required with the command here
			fmt.Println("Handling Logic")
			return nil
		}).Create()

	errChan := make(chan error, 0)
	defer close(errChan)

	go func() {
		errChan <- handler.Start(ctx)
	}()

	go func() {
		time.Sleep(time.Second * 10)
		errChan <- handler.Stop()
	}()

	time.Sleep(time.Second * 1)

	//Publishing
	for i := 0; i < 10; i++ {
		fmt.Println("Pushing ", i)
		err := redisFactory.Publisher().Create().PublishCommand(ctx, &command)
		if err != nil {
			errChan <- err
			break
		}
	}

	for {
		select {
		case err := <-errChan:
			{
				log.Println(err)
				return
			}
		}
	}
}
