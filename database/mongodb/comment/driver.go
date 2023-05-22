package comment

import (
	"FaRyuk/config"
	"FaRyuk/internal/types"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCommentRepository struct {
	client *mongo.Client
	config *config.Config
}

// TODO move all hardcoded strings (client url, db name, collection name) to a config file and inject it to the repo)
func NewMongoCommentRepository(config *config.Config) *MongoCommentRepository {
	return &MongoCommentRepository{
		client: createMongoDBCLient(),
		config: config,
	}
}

// InsertComment : inserts comment in the database
func (repo *MongoCommentRepository) InsertComment(comment *types.Comment, done chan<- error) {
	defer close(done)
	if comment.ID == "" {
		comment.ID = uuid.NewString() // Generate a new UUID for the comment
	}
	collection := repo.client.Database(repo.config.Database.DbName).Collection("comment")
	_, err := collection.InsertOne(context.TODO(), comment)
	done <- err
}

// GetComments : gets all comments
func (repo *MongoCommentRepository) GetComments(comments chan<- *types.CommentsWithErrorType) {
	var results []*types.Comment
	collection := repo.client.Database(repo.config.Database.DbName).Collection("comment")
	findOptions := options.Find()
	cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		comments <- &types.CommentsWithErrorType{Comments: nil, Err: err}
		return
	}

	for cur.Next(context.Background()) {
		var elem *types.Comment
		err := cur.Decode(&elem)
		if err != nil {
			comments <- &types.CommentsWithErrorType{Comments: nil, Err: err}
			return
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		comments <- &types.CommentsWithErrorType{Comments: nil, Err: err}
		return
	}

	cur.Close(context.Background())
	comments <- &types.CommentsWithErrorType{Comments: results, Err: nil}
}

// RemoveCommentByID : removes comment by ID
func (repo *MongoCommentRepository) RemoveCommentByID(id string, done chan<- error) {
	collection := repo.client.Database(repo.config.Database.DbName).Collection("comment")
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		done <- err
		return
	}
	done <- nil
}

// UpdateComment : updates comment
func (repo *MongoCommentRepository) UpdateComment(r *types.Comment, done chan<- error) {
	collection := repo.client.Database(repo.config.Database.DbName).Collection("comment")
	_, err := collection.UpdateOne(context.Background(), bson.M{"id": r.ID}, bson.M{"$set": r})
	done <- err
}

// GetCommentByID : retrieves comment by ID
func (repo *MongoCommentRepository) GetCommentByID(id string, result chan<- *types.CommentWithErrorType) {
	var comment *types.Comment
	collection := repo.client.Database(repo.config.Database.DbName).Collection("comment")
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&comment)
	if err != nil {
		result <- &types.CommentWithErrorType{Comment: &types.Comment{}, Err: err}
	} else {
		result <- &types.CommentWithErrorType{Comment: comment, Err: nil}
	}
}

// GetCommentsByText : search comments by regular expression
func (repo *MongoCommentRepository) GetCommentsByText(search string, result chan *types.CommentsWithErrorType) {
	defer close(result)
	var results []*types.Comment

	collection := repo.client.Database(repo.config.Database.DbName).Collection("comment")
	filter := bson.M{"content": bson.M{"$regex": ".*" + search + ".*"}}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		result <- &types.CommentsWithErrorType{Comments: nil, Err: err}
		return
	}

	for cur.Next(context.Background()) {
		var elem *types.Comment
		err := cur.Decode(&elem)
		if err != nil {
			result <- &types.CommentsWithErrorType{Comments: nil, Err: err}
			return
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		result <- &types.CommentsWithErrorType{Comments: nil, Err: err}
		return
	}
	// cur.Close(context.Background())
	result <- &types.CommentsWithErrorType{Comments: results, Err: nil}

}

// GetCommentsByTextAndOwner : searchs for comments containing a particular text and that could be accessed by the current user
func (repo *MongoCommentRepository) GetCommentsByTextAndOwner(search string, idUser string, result chan *types.CommentsWithErrorType) {
	defer close(result)
	var results []*types.Comment
	collection := repo.client.Database(repo.config.Database.DbName).Collection("comment")
	filter := bson.M{"content": bson.M{"$regex": ".*" + search + ".*"}}

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		result <- &types.CommentsWithErrorType{Comments: nil, Err: err}
		return
	}

	for cur.Next(context.Background()) {
		var elem *types.Comment
		err := cur.Decode(&elem)
		if err != nil {
			result <- &types.CommentsWithErrorType{Comments: nil, Err: err}
			return
		}
		if elem.Owner == idUser {
			results = append(results, elem)
		}
	}

	if err := cur.Err(); err != nil {
		result <- &types.CommentsWithErrorType{Comments: nil, Err: err}
		return
	}

	cur.Close(context.Background())
	result <- &types.CommentsWithErrorType{Comments: results, Err: nil}
}

// GetCommentsByResultID : get all comments for a given result ID
func (repo *MongoCommentRepository) GetCommentsByResultID(idResult string, result chan<- *types.CommentsWithErrorType) {
	defer close(result)
	var results []*types.Comment

	collection := repo.client.Database(repo.config.Database.DbName).Collection("comment")

	cur, err := collection.Find(context.Background(), bson.M{"idResult": idResult})

	if err != nil {
		result <- &types.CommentsWithErrorType{Comments: nil, Err: err}
		return
	}

	for cur.Next(context.TODO()) {
		var comment *types.Comment
		err := cur.Decode(&comment)
		if err != nil {
			result <- &types.CommentsWithErrorType{Comments: nil, Err: err}
			return
		}
		results = append(results, comment)
	}

	if err := cur.Err(); err != nil {
		result <- &types.CommentsWithErrorType{Comments: nil, Err: err}
		return
	}

	cur.Close(context.TODO())
	result <- &types.CommentsWithErrorType{Comments: results, Err: err}

}

func createMongoDBCLient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// CloseConnection : closes connection with mongo db
func (repo *MongoCommentRepository) CloseConnection() {
	err := repo.client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
		return
	}
}
