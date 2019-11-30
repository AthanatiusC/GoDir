package directory

import (
	"context"
	"encoding/json"
	utils "github.com/athanatius/godir"
	models "github.com/athanatius/godir/models"
	// "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"os"

	// "github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	// "time"
)

//CreateUsers insert one to DB
func CreateFolder(res http.ResponseWriter, req *http.Request) {
	// Declare Variable
	var directory models.Directory
	json.NewDecoder(req.Body).Decode(&directory)
	// f, err := os.Create(directory.Path)
	err := os.MkdirAll(directory.Path, 777)
	utils.ErrorHandler(err)
	// defer f.Close()
	utils.WriteResult(res, nil, directory.Path+" Created")
}

//GetAllUsers return res json Users model
func GetDirectory(res http.ResponseWriter, req *http.Request) {
	//Declare Variable
	var model models.Directory
	var file models.Files
	var files []models.Files

	//Decode Request
	err := json.NewDecoder(req.Body).Decode(&model)
	utils.ErrorHandler(err)

	list, err := ioutil.ReadDir(model.Path)
	for _, val := range list {
		file.Size = val.Size()
		file.Name = val.Name()
		file.Path = strings.Join([]string{model.Path, val.Name()}, "/")
		file.LastModified = val.ModTime()
		file.FileMode = val.Mode()
		if val.IsDir() {
			file.Type = "Folder"
		} else {
			format := strings.Split(val.Name(), ".")
			file.Type = format[len(format)-1]
		}
		utils.ErrorHandler(err)
		files = append(files, file)
		// http.DetectContentType()
	}
	if err != nil {
		log.Println(err)
		utils.WriteResult(res, nil, "Directory Not Found!")
		return
	}

	utils.WriteResult(res, files, "Returned "+strconv.Itoa(len(files))+" Object")
}

func DeleteDirectory(res http.ResponseWriter, req *http.Request) {
	var directory models.Directory
	json.NewDecoder(req.Body).Decode(&directory)
	err := os.RemoveAll(directory.Path)
	utils.ErrorHandler(err)
	utils.WriteResult(res, nil, directory.Path+" Deleted")
}

func VerifyOwnership(id primitive.ObjectID, auth_key string) bool {
	var model models.Users

	db := utils.ConnectMongoDB()
	db.Collection("users").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&model)

	if auth_key != "" {
		if auth_key != model.Auth {
			return false
		} else if auth_key == model.Auth {
			return true
		}
	}
	return false
}
