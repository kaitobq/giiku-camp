package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func New() *firestore.Client {
	ctx := context.Background()

	path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	var app *firebase.App
	var err error

	if path != "" {
		// ローカル環境や開発環境では JSON を使う
		sa := option.WithCredentialsFile(path)
		app, err = firebase.NewApp(ctx, nil, sa)
		if err != nil {
			panic(err)
		}
	} else {
		// 本番環境ではデフォルト認証を使用
		app, err = firebase.NewApp(ctx, nil)
		if err != nil {
			panic(err)
		}
	}

	// Firestore クライアントの作成
	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}

	return client
}
