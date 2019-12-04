package users

import (
	"context"
	"encoding/json"
	utils "github.com/athanatius/godir"
	models "github.com/athanatius/godir/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//CreateUsers insert one to DB
func CreateUsers(res http.ResponseWriter, req *http.Request) {
	//Declare Variable
	var model models.Users
	var model2 models.Users

	// userid := req.Header.Get("user_id")
	// authkey := req.Header.Get("key")
	// uid, _ := primitive.ObjectIDFromHex(userid)

	// if utils.VerifyOwnership(uid, authkey) {

	//Decode Request
	err := json.NewDecoder(req.Body).Decode(&model)

	//Connect DB
	db := utils.ConnectMongoDB()

	//Loop column
	db.Collection("users").FindOne(context.TODO(), bson.M{"username": model.Username}).Decode(&model2)
	if len(model2.Name) != 0 {
		utils.WriteResult(req, res, nil, "User Already Exist!")
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	utils.ErrorHandler(err)
	model.Password = string(password)
	model.FormatTime = time.Now()

	_, err = db.Collection("users").InsertOne(context.TODO(), model)
	utils.ErrorHandler(err)

	//Return Res
	utils.ErrorHandler(err)
	utils.WriteResult(req, res, nil, "User Successfully Created!")
	// } else {
	// 	utils.WriteResult(res, nil, "Access Denied ")
	// 	return
	// }
}

//GetAllUsers return res json Users model
func GetAllUsers(res http.ResponseWriter, req *http.Request) {
	//Declare Variable
	var model models.Users
	var all []models.Users

	switch req.Method {
	case "OPTIONS":
		utils.WriteResult(req, res, nil, "Access Allowed")
		return
	}

	authkey := req.Header.Get("auth_key")
	log.Println(authkey)
	userid := req.Header.Get("user_id")
	uid, _ := primitive.ObjectIDFromHex(userid)

	if utils.VerifyOwnership(uid, authkey) {
		//Connect DB
		db := utils.ConnectMongoDB()
		col, err := db.Collection("users").Find(context.TODO(), bson.M{})

		//Loop column
		for col.Next(context.TODO()) {
			err := col.Decode(&model)
			utils.ErrorHandler(err)
			all = append(all, model)
		}

		if len(all) == 0 {
			utils.WriteResult(req, res, nil, "User Collumn Empty")
			return
		}

		//Return Res
		utils.ErrorHandler(err)
		utils.WriteResult(req, res, all, "Sucessfully Returned "+strconv.Itoa(len(all))+" Users!")
	} else {
		utils.WriteResult(req, res, nil, "Access Denied ")
		return
	}
}

func DeleteUsers(res http.ResponseWriter, req *http.Request) {
	raw_param := mux.Vars(req)
	id := raw_param["id"]
	objid, err := primitive.ObjectIDFromHex(id)
	utils.ErrorHandler(err)

	authkey := req.Header.Get("auth_key")
	if utils.VerifyOwnership(objid, authkey) {
		db := utils.ConnectMongoDB()

		collection := db.Collection("users")
		deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objid})
		utils.ErrorHandler(err)

		if deleteResult.DeletedCount == 0 {
			utils.WriteResult(req, res, nil, "User not found")
			return
		}

		utils.WriteResult(req, res, deleteResult.DeletedCount, "User Deleted!")
	} else {
		utils.WriteResult(req, res, nil, "User Not Found ")
		return
	}
}

func Auth(res http.ResponseWriter, req *http.Request) {
	var user models.Users
	var userauth models.Users
	err := json.NewDecoder(req.Body).Decode(&user)
	utils.ErrorHandler(err)
	db := utils.ConnectMongoDB()
	collection := db.Collection("users")
	utils.ErrorHandler(err)
	collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&userauth)

	// password:= hashAndSalt()
	ismatch := comparePasswords(userauth.Password, []byte(user.Password))
	// log.Println(userauth.Password + "	" + user.Password)

	if ismatch == true {

		authkey := ksuid.New()
		userauth.Auth = authkey.String()

		collection.FindOneAndUpdate(context.TODO(), bson.M{"username": user.Username}, bson.D{{Key: "$set", Value: userauth}})
		utils.WriteResult(req, res, bson.M{"Auth": authkey, "RootPath": userauth.RootPath, "Id": userauth.ID}, "Access Allowed")
	} else {
		utils.WriteResult(req, res, nil, "Access Denied")
	}
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool { // Since we'll be getting the hashed password from the DB it
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}
