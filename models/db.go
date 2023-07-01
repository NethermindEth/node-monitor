package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	_ "time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient    *mongo.Client
	dbContext      context.Context
	database       = "nethermind"
	collectionName = "nodes"
	nodeNotFound   = errors.New("node not found")
)

func InitMongo(uri string) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	dbContext = context.Background()

	if err := client.Ping(dbContext, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	mongoClient = client
}

func getCollection() *mongo.Collection {
	return mongoClient.Database(database).Collection(collectionName)
}

func GetAllData() []Node {
	collection := getCollection()
	cursor, err := collection.Find(dbContext, bson.M{})
	if err != nil {
		log.Println("Failed to fetch data:", err)
		return nil
	}
	defer cursor.Close(dbContext)

	var data []Node
	if err := cursor.All(dbContext, &data); err != nil {
		log.Println("Failed to decode data:", err)
		return nil
	}

	return data
}

func CreateNode(node *Node) error {
	collection := getCollection()
	_, err := getNode(node.Enode)
	if err != nil {
		if errors.Is(err, nodeNotFound) {
			_, err := collection.InsertOne(dbContext, node)
			if err != nil {
				return fmt.Errorf("failed to insert node: %v", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get node: %v", err)
	}

	return nil
}

func getNode(enode string) (*Node, error) {
	collection := getCollection()
	var node Node
	err := collection.FindOne(dbContext, bson.D{{"enode", enode}}).Decode(&node)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nodeNotFound
		}
		return nil, fmt.Errorf("failed to get node: %v", err)
	}

	return &node, nil
}

func AddEntry(enode string, data EthereumNodeData) error {
	collection := getCollection()
	filter := bson.D{{"enode", enode}}

	node, err := getNode(enode)
	if err != nil {
		if errors.Is(err, nodeNotFound) {
			return nil
		}
		return fmt.Errorf("failed to get node: %v", err)
	}

	if node == nil {
		node.Data = make([]EthereumNodeData, 0)
	}
	node.Data = append(node.Data, data)

	update := bson.D{{"$set", node}}

	result, err := collection.UpdateOne(dbContext, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update node: %v", err)
	}

	if result.MatchedCount == 0 {
		return nodeNotFound
	}

	return nil
}
