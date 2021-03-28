package main

/*
 * File for firebase(database) functionality
 * author Martin Iversen
 * 28.03.2020
 * version 0.8
 */
import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
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

func DeleteWebhook(id string) error {
	_, err := client.Collection(Collection).Doc(id).Delete(ctx)
	if err != nil {
		return errors.New("Error occurred when trying to delete webhook. Webhook ID: " + id)
	}
	return nil
}

/**
This method will get all the webhooks stored in the database.
Source: https://stackoverflow.com/a/61429531
*/
func GetAll() ([]*firestore.DocumentSnapshot, error) {
	var docs []*firestore.DocumentSnapshot
	iter := client.Collection(Collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func update(id string, data interface{}) error {
	_, err := client.Collection(Collection).Doc(id).Set(ctx, data)
	if err != nil {
		return errors.New("Error while adding webhook to database: " + err.Error())
	}
	return nil

}
