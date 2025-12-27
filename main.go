package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go"
)

func main() {
	ctx := context.Background()

	valkeyAddr := os.Getenv("VALKEY_ADDR")

	client, err := valkey.NewClient(valkey.ClientOption{
		Username:    "default",
		ClientName:  "app",
		InitAddress: []string{valkeyAddr},
	})
	if err != nil {
		slog.Error(
			"Creating Valkey client failed",
			slog.Any("valkeyAddr", valkeyAddr),
			slog.Any("error", err))
		return
	}

	defer func() {
		client.Close()
		slog.Info("valkey client closed successfully.")
	}()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		slog.Error("Ping failed", slog.Any("error", err))
		return
	}

	slog.Info("Ping succeeded.")

	for {
		t, err := timeFromValkey(ctx, client)
		if err != nil {
			slog.Error("Getting time from valkey failed", slog.Any("error", err))
			return
		}

		slog.Info("Current valkey server time", slog.String("time", t.String()))

		time.Sleep(5 * time.Second)
	}
}

func timeFromValkey(ctx context.Context, client valkey.Client) (*time.Time, error) {
	timeResp := client.Do(
		ctx,
		client.B().Time().Build(),
	)
	if err := timeResp.Error(); err != nil {
		return nil, err
	}

	timeArray, err := timeResp.ToArray()
	if err != nil {
		return nil, err
	} else if len(timeArray) < 2 {
		return nil, errors.New("time response array has less than 2 elements")
	}

	sec, err := timeArray[0].ToString()
	if err != nil {
		return nil, err
	}
	currentUnixSeconds, err := strconv.ParseInt(sec, 10, 64)
	if err != nil {
		return nil, err
	}
	milli, err := timeArray[1].ToString()
	if err != nil {
		return nil, err
	}
	currentUnixMilli, err := strconv.ParseInt(milli, 10, 64)
	if err != nil {
		return nil, err
	}

	t := time.Unix(currentUnixSeconds, currentUnixMilli*1e6)
	return &t, nil
}
