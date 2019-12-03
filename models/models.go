package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

type Application struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ServerDir string             `server_root`
	AppName   string             `application_name`
}

//Meta Model
type Meta struct {
	Updated primitive.Timestamp `json:"updated_at"`
	Created primitive.Timestamp `json:"created_at"`
}

//Directory Model
type Directory struct {
	Name string  `json:"directory_name"`
	Path string  `json:"requested_path"`
	File []Files `json:"file_list"`
}

type Files struct {
	Name         string      `json:"file_name"`
	Path         string      `json:"file_path"`
	Type         string      `json:"file_type"`
	Size         int64       `json:"file_size"`
	Description  string      `json:"file_description"`
	FileMode     os.FileMode `json:"file_mode"`
	LastModified time.Time   `json:"last_modified"`
}

//Users Model
type Users struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name"`
	Email      string             `json:"email"`
	Username   string             `json:"username"`
	Password   string             `json:"password"`
	Auth       string             `json:"authentication_key"`
	Photo      string             `json:"picture_url"`
	RootPath   string             `json:"root_path"`
	FormatTime time.Time          `json:"format_time"`
	MetaData   Meta               `json:"meta_data"`
}
