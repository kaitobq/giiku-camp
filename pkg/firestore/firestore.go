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
	if path == "" {
		panic("GOOGLE_APPLICATION_CREDENTIALS is not set")
	}
	sa := option.WithCredentialsFile(path)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		panic(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}

	return client
}
