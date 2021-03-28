package main

/*
 * File for firebase(database) functionality
 * author Martin Iversen
 * 22.03.2020
 * version 0.2
 *
 */
import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
)

var ctx context.Context
var client *firestore.Client

const Collection = "webhook"

func Init() error {
	// Firebase initialisation
	ctx = context.Background()

	// Authenticate
	opt := option.WithCredentialsFile("./assignment-2-13402-firebase-adminsdk-j2q0b-c9eb380f52.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error initializing client: %v", err)
	}
	return nil
}

func AddWebhook(webhook interface{}) (string, error) {
	newEntry, _, err := client.Collection(Collection).Add(ctx, webhook)
	if err != nil {
		return "", errors.New("Error while adding webhook to database: " + err.Error())
	}
	return newEntry.ID, nil
}
