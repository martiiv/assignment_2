package main

/*
 * File for firebase(database) functionality
 * Has 5 functions:
 * 					Init(): 			For initializing the database using a key
 *					AddWebhook(): 		For adding a webhook to the database
 *					DeleteWebhook():	For deleting a webhook
 *					GetAll():			For getting all the webhooks in the database
 *										(this function is taken from stackoverflow:https://stackoverflow.com/a/61429531)
 *					Update():			For updating information on a webhook in the database
 * author Martin Iversen
 * 29.03.2020
 * version 1.0
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

//Initializing DB
var ctx context.Context
var client *firestore.Client

const Collection = "webhook" //Defining the name of the collection we will be dealing with

/*
 * Function for initializing the database, will be used when starting the pp
 */
func Init() error {
	// Firebase initialisation
	ctx = context.Background()

	// Authenticate with key file from firebase
	opt := option.WithCredentialsFile("./assignment-2-13402-firebase-adminsdk-j2q0b-c9eb380f52.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error occurred initializing client: %v", err)
	}

	return nil
}

/*
 * Function for adding a webhook to the database
 * Returns the ID an object is given when the database creates
 */
func AddWebhook(webhook interface{}) (string, error) {
	newEntry, _, err := client.Collection(Collection).Add(ctx, webhook) //Adds to the database
	if err != nil {
		return "", errors.New("Error occurred when adding webhook to database: " + err.Error())
	}
	return newEntry.ID, nil //Returns the id of an entry in the database collection
}

/*
 * Function for deleting a webhook from the database
 */
func DeleteWebhook(id string) error {
	_, err := client.Collection(Collection).Doc(id).Delete(ctx) //Deletes from the database
	if err != nil {
		return errors.New("Error occurred when trying to delete webhook. Webhook ID: " + id)
	}
	return nil
}

/*
 * Function for getting all entries in a database
 * Source: https://stackoverflow.com/a/61429531
 * Returns an object containing document snapshots of the entries in the database
 */
func GetAll() ([]*firestore.DocumentSnapshot, error) {
	var docs []*firestore.DocumentSnapshot               //Defining object to be returned
	iter := client.Collection(Collection).Documents(ctx) //Gets all entries in the database
	for {                                                //Iterates through the database
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
	return docs, nil //Returns a list of entries
}

/*
 * Function for updating information on an entry in the database
 */
func update(id string, data interface{}) error {
	_, err := client.Collection(Collection).Doc(id).Set(ctx, data)
	if err != nil {
		return errors.New("Error while adding webhook to database: " + err.Error())
	}
	return nil
}
