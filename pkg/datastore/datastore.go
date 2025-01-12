package datastore

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

func New() *datastore.Client {
	ctx := context.Background()

	path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	projectName := os.Getenv("GOOGLE_CLOUD_PROJECT")
	var dsClient *datastore.Client
	var err error

	if path != "" {
		// ローカル環境や開発環境では JSON を使う
		sa := option.WithCredentialsFile(path)
		dsClient, err = datastore.NewClient(ctx, projectName, sa)
		if err != nil {
			panic(err)
		}
	} else {
		// 本番環境ではデフォルト認証を使用
		dsClient, err = datastore.NewClient(ctx, projectName)
		if err != nil {
			panic(err)
		}
	}

	return dsClient
}
