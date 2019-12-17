package directory

import (
	"encoding/json"
	utils "github.com/AthanatiusC/godir"
	models "github.com/AthanatiusC/godir/models"
	// "github.com/gorilla/mux"
	"os"

	// "github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "golang.org/x/crypto/bcrypt"
	// "archive/zip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	// "time"
)

//CreateUsers insert one to DB
func CreateFolder(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "OPTIONS":
		utils.WriteResult(req, res, nil, "Access Allowed")
		return
	}
	// Declare Variable
	var directory models.Directory
	json.NewDecoder(req.Body).Decode(&directory)
	if directory.Path == "" {
		utils.WriteResult(req, res, nil, directory.Path+"Path cannot be null")
		return
	}
	// f, err := os.Create(directory.Path)
	err := os.MkdirAll(directory.Path, 777)
	utils.ErrorHandler(err)
	// defer f.Close()
	utils.WriteResult(req, res, nil, directory.Path+" Created")
}

type RenamePayload struct {
	Oldpath string `json:"old_path"`
	Newpath string `json:"new_path"`
}

func RenameFolder(res http.ResponseWriter, req *http.Request) {
	userid := req.Header.Get("user_id")
	authkey := req.Header.Get("auth_key")
	uid, _ := primitive.ObjectIDFromHex(userid)

	switch req.Method {
	case "OPTIONS":
		utils.WriteResult(req, res, nil, "Access Allowed")
		return
	}

	if utils.VerifyOwnership(uid, authkey) {
		var rename RenamePayload
		json.NewDecoder(req.Body).Decode(&rename)
		err := os.Rename(rename.Oldpath, rename.Newpath)
		isErr := utils.ErrorHandler(err)
		if isErr {
			utils.WriteResult(req, res, nil, "Action Failed")
		} else {
			utils.WriteResult(req, res, nil, "Successfully Renamed")
		}
	} else {
		utils.WriteResult(req, res, nil, "Access Denied")
	}
}

//GetAllUsers return res json Users model
func GetDirectory(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "OPTIONS":
		utils.WriteResult(req, res, nil, "Access Allowed")
		return
	}
	//Declare Variable
	var model models.Directory
	var file models.Files
	var files []models.Files

	userid := req.Header.Get("user_id")
	authkey := req.Header.Get("auth_key")
	uid, _ := primitive.ObjectIDFromHex(userid)

	if utils.VerifyOwnership(uid, authkey) {

		//Decode Request
		err := json.NewDecoder(req.Body).Decode(&model)
		// utils.ErrorHandler(err)

		list, err := ioutil.ReadDir(model.Path)
		if err != nil {
			log.Println(err)
			utils.WriteResult(req, res, nil, "Directory Not Found!")
			return
		}
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
			files = append(files, file)
			// http.DetectContentType()
		}
		utils.WriteResult(req, res, files, "Returned "+strconv.Itoa(len(files))+" Object")
	} else {
		utils.WriteResult(req, res, nil, "Access Denied")
	}
}

func DeleteDirectory(res http.ResponseWriter, req *http.Request) {
	userid := req.Header.Get("user_id")
	authkey := req.Header.Get("auth_key")
	uid, _ := primitive.ObjectIDFromHex(userid)

	switch req.Method {
	case "OPTIONS":
		utils.WriteResult(req, res, nil, "Access Allowed")
		return
	}

	if utils.VerifyOwnership(uid, authkey) {
		switch req.Method {
		case "OPTIONS":
			utils.WriteResult(req, res, nil, "Access Allowed")
			return
		}
		var directory models.Directory
		if utils.IsExists(directory.Path) {
			utils.WriteResult(req, res, nil, "File/Folder not found!")
		}
		json.NewDecoder(req.Body).Decode(&directory)
		err := os.RemoveAll(directory.Path)
		utils.ErrorHandler(err)
		utils.WriteResult(req, res, nil, directory.Path+" Deleted")
	} else {
		utils.WriteResult(req, res, nil, "Access Denied ")
		return
	}
}

func DownloadFile(res http.ResponseWriter, req *http.Request) {
	userid := req.Header.Get("user_id")
	authkey := req.Header.Get("auth_key")
	uid, _ := primitive.ObjectIDFromHex(userid)

	switch req.Method {
	case "OPTIONS":
		utils.WriteResult(req, res, nil, "Access Allowed")
		return
	}

	if utils.VerifyOwnership(uid, authkey) {
		if req.Header.Get("Path") == "" {
			utils.WriteResult(req, res, nil, "Empty Path")
			return
		}
		Openfile, err := os.Open(req.Header.Get("Path"))
		utils.ErrorHandler(err)

		defer Openfile.Close() //Close after function return

		Filename := Openfile.Name()

		log.Println("User : " + req.Header.Get("user_id") + " Requested : " + Filename)

		//File is found, create and send the correct headers

		//Get the Content-Type of the file
		//Create a buffer to store the header of the file in
		FileHeader := make([]byte, 512)
		//Copy the headers into the FileHeader buffer
		Openfile.Read(FileHeader)
		//Get content type of file
		FileContentType := http.DetectContentType(FileHeader)

		//Get the file size
		FileStat, _ := Openfile.Stat()                     //Get info from file
		FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

		//Send the headers
		res.Header().Set("Content-Disposition", "attachment; filename="+req.Header.Get("Name"))
		res.Header().Set("Content-Type", FileContentType)
		res.Header().Set("Content-Length", FileSize)

		//Send the file
		//We read 512 bytes from the file already, so we reset the offset back to 0
		Openfile.Seek(0, 0)
		io.Copy(res, Openfile) //'Copy' the file to the client
	} else {
		utils.WriteResult(req, res, nil, "Access Denied ")
		return
	}
}

func UploadFile(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "OPTIONS":
		utils.WriteResult(req, res, nil, "Access Allowed")
		return
	}

	userid := req.Header.Get("user_id")
	authkey := req.Header.Get("auth_key")
	uid, _ := primitive.ObjectIDFromHex(userid)

	if utils.VerifyOwnership(uid, authkey) {
		req.ParseMultipartForm(1000)
		file, handler, err := req.FormFile("Files")
		Path := req.FormValue("Path")
		// Name := req.FormValue("Name")

		utils.ErrorHandler(err)
		defer file.Close()

		log.Println("User uploaded " + handler.Filename)
		log.Printf("File Size: %+v\n", handler.Size)
		log.Printf("MIME Header: %+v\n", handler.Header)

		f, err := os.Create(Path)
		io.Copy(f, file)
		defer f.Close()
		if err != nil {
			log.Println(err)
			return
		}

		utils.WriteResult(req, res, nil, "File Successfully uploaded")
	} else {
		utils.WriteResult(req, res, nil, "Access Denied ")
		return
	}
}