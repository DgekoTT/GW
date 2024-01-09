package connection

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

const BackOffAttempts = 10
const BackOffDuration = 2 * time.Second

func ConnectToService(name string, port int, log *slog.Logger) (*grpc.ClientConn, error) {
	url := fmt.Sprintf("%s:%d", name, port)

	var counts uint8

	for {
		contextDial, cancelDial := context.WithTimeout(context.Background(), time.Second)
		defer cancelDial()

		conn, err := grpc.DialContext(
			contextDial, url, grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)

		if err != nil {
			log.Error(fmt.Sprintf("MicroService '%s' not connected yet \n", name))
			counts++
		} else {
			log.Info(fmt.Sprintf("Connected to '%s'!\n", name))
			return conn, nil
		}

		if counts > BackOffAttempts {
			return nil, err
		}

		log.Info(fmt.Sprintf("Try connect '%s' again after 2 seconds", name))
		time.Sleep(BackOffDuration)

		continue
	}
}
