package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Document struct {
	ID            bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Title         string        `json:"title,omitempty" bson:"title"`
	Filename      string        `json:"filename,omitempty" bson:"filename"`
	S3Path        string        `json:"s3_path,omitempty" bson:"s3_path"`
	Content       string        `json:"content,omitempty" bson:"content"`
	ContentType   string        `json:"content_type,omitempty" bson:"content_type"`
	ThumbnailPath string        `json:"thumbnail_path,omitempty" bson:"thumbnail_path"`
	Timestamp     time.Time     `json:"timestamp,omitempty" bson:"timestamp"`
	Correspondent string        `json:"correspondent,omitempty" bson:"correspondent"`
	Tags          []string      `json:"tags,omitempty" bson:"tags"`
	Hash          string        `json:"hash,omitempty" bson:"hash"`
}
