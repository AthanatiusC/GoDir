package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

//Meta Model
type Meta struct {
	Updated primitive.DateTime
	Created primitive.Timestamp
}

//Directory Model
type Directory struct {
	Name string
	Path string
	File []Files
}

type Files struct {
	Name         string
	Path         string
	Type         string
	Size         int64
	Description  string
	FileMode     os.FileMode
	LastModified time.Time
}

//Users Model
type Users struct {
	ID         primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	Name       string
	Email      string
	Username   string
	Password   string
	Auth       string
	Photo      string
	FormatTime time.Time
	Updated    primitive.Timestamp
	Created    primitive.Timestamp
}
