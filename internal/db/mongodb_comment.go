package db

import (
	"FaRyuk/internal/types"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertComment : inserts comment in the database
func (db *Handler) InsertComment(r *types.Comment) error {
	collection := db.client.Database("faryuk").Collection("comment")
	_, err := collection.InsertOne(context.TODO(), r)
	return err
}

// GetComments : gets all comments
func (db *Handler) GetCommentsWithChannel(comments chan<- types.CommentsWithErrorType) {
	var results []types.Comment
	collection := db.client.Database("faryuk").Collection("comment")
	findOptions := options.Find()
	cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		comments <- types.CommentsWithErrorType{Comments: make([]types.Comment, 0), Err: err}
		return
	}

	for cur.Next(context.Background()) {
		var elem types.Comment
		err := cur.Decode(&elem)
		if err != nil {
			comments <- types.CommentsWithErrorType{Comments: make([]types.Comment, 0), Err: err}
			return
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		comments <- types.CommentsWithErrorType{Comments: make([]types.Comment, 0), Err: err}
		return
	}

	cur.Close(context.Background())
	comments <- types.CommentsWithErrorType{Comments: results, Err: nil}
}

// RemoveCommentByID : removes comment by ID
func (db *Handler) RemoveCommentByIDWithChannel(id string, done chan<- error) {
	collection := db.client.Database("faryuk").Collection("comment")
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		done <- err
		return
	}
	done <- nil
}

// UpdateComment : updates comment
func (db *Handler) UpdateCommentWithChannel(r *types.Comment, done chan<- bool) {
	collection := db.client.Database("faryuk").Collection("comment")
	_, err := collection.UpdateOne(context.Background(), bson.M{"id": r.ID}, bson.M{"$set": r})
	done <- err == nil
}

// GetCommentByID : retrieves comment by ID
func (db *Handler) GetCommentByIDWithChannel(id string, result chan<- types.CommentWithErrorType) {
	var comment types.Comment
	collection := db.client.Database("faryuk").Collection("comment")
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&comment)
	if err != nil {
		result <- types.CommentWithErrorType{Comment: types.Comment{}, Err: err}
	} else {
		result <- types.CommentWithErrorType{Comment: comment, Err: nil}
	}
}

// GetCommentsByText : search comments by regular expression
func (db *Handler) GetCommentsByText(search string, result chan types.CommentsWithErrorType) {
	defer close(result)
	var results []types.Comment

	collection := db.client.Database("faryuk").Collection("comment")
	filter := bson.M{"content": bson.M{"$regex": ".*" + search + ".*"}}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		result <- types.CommentsWithErrorType{Comments: []types.Comment{}, Err: err}
		return
	}

	for cur.Next(context.Background()) {
		var elem types.Comment
		err := cur.Decode(&elem)
		if err != nil {
			result <- types.CommentsWithErrorType{Comments: make([]types.Comment, 0), Err: err}
			return
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		result <- types.CommentsWithErrorType{Comments: make([]types.Comment, 0), Err: err}
		return
	}
	// cur.Close(context.Background())
	result <- types.CommentsWithErrorType{Comments: results, Err: nil}

}

// GetCommentsByTextAndOwner : searchs for comments containing a particular text and
//
//	that could be accessed by the current user
func (db *Handler) GetCommentsByTextAndOwner(search string, idUser string, result chan types.CommentsWithErrorType) {
	defer close(result)
	var results []types.Comment
	collection := db.client.Database("faryuk").Collection("comment")
	filter := bson.M{"content": bson.M{"$regex": ".*" + search + ".*"}}

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		result <- types.CommentsWithErrorType{Comments: []types.Comment{}, Err: err}
		return
	}

	for cur.Next(context.Background()) {
		var elem types.Comment
		err := cur.Decode(&elem)
		if err != nil {
			result <- types.CommentsWithErrorType{Comments: []types.Comment{}, Err: err}
			return
		}
		if elem.Owner == idUser {
			results = append(results, elem)
		}
	}

	if err := cur.Err(); err != nil {
		result <- types.CommentsWithErrorType{Comments: []types.Comment{}, Err: err}
		return
	}

	cur.Close(context.Background())
	result <- types.CommentsWithErrorType{Comments: results, Err: nil}
}

// GetCommentsByResult : get all comments for a given result
func (db *Handler) GetCommentsByResult(idResult string, result chan<- types.CommentsWithErrorType) {
	defer close(result)
	var results []types.Comment

	collection := db.client.Database("faryuk").Collection("comment")
	findOptions := options.Find()

	cur, err := collection.Find(context.TODO(), findOptions)
	if err != nil {
		result <- types.CommentsWithErrorType{Comments: []types.Comment{}, Err: err}
		return
	}

	for cur.Next(context.TODO()) {
		var elem types.Comment
		err := cur.Decode(&elem)
		if err != nil {
			result <- types.CommentsWithErrorType{Comments: []types.Comment{}, Err: err}
			return
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		result <- types.CommentsWithErrorType{Comments: []types.Comment{}, Err: err}
		return
	}

	cur.Close(context.TODO())
	result <- types.CommentsWithErrorType{Comments: []types.Comment{}, Err: err}

}
