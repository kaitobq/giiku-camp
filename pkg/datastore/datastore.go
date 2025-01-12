package datastore

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

func New() *datastore.Client {
	ctx := context.Background()

	projectName := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectName == "" {
		panic("GOOGLE_CLOUD_PROJECT environment variable is not set")
	}

	opts := []option.ClientOption{}
	if emulatorHost := os.Getenv("DATASTORE_EMULATOR_HOST"); emulatorHost != "" {
		// エミュレーターを使用する場合のオプションを追加
		opts = append(opts,
			option.WithEndpoint(emulatorHost),
			option.WithoutAuthentication(), // 認証をパス
		)
	}

	dsClient, err := datastore.NewClient(ctx, projectName, opts...)
	if err != nil {
		panic(err)
	}
	return dsClient
}
