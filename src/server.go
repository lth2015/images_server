package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"os"
)

// Version
const (
	VERSION = "v1"
)

// For create/delete Folders
type PathMessage struct {
	Path    string `json:string "Path"`
	Message string `json:string "Message"`
}

func ToJson(v interface{}) (string, error) {
	// Json Marshal
	result, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		log.Fatalf("Json Marshal Error: %v", err)
		return "", err
	}
	return string(result), nil
}

func WriteResponse(w http.ResponseWriter, v interface{}) {
	resp, _ := ToJson(v)
	fmt.Fprintf(w, resp)
}

func Healthz(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "ok\n")
}

func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, VERSION+"\n")
}

func PostAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	account = "./storage/" + account

	pm := PathMessage{Path: account, Message: "Folder was created!"}
	if _, err := os.Stat(account); err != nil {
		if ok := os.MkdirAll(account, 0755); ok != nil {
			log.Fatalf("Make account error: accout=%s, error=%s", account, err.Error())
			pm.Message = ok.Error()
		} else {
			log.Printf("Make account: account=%s", account)
		}
	} else {
		pm.Message = "Direcotry already exists"
	}

	WriteResponse(w, &pm)

}

func DeleteAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	account = "./storage/" + account

	pm := PathMessage{Path: account, Message: "Folder was deleted!"}
	if _, err := os.Stat(account); err == nil {
		if ok := os.RemoveAll(account); ok != nil {
			log.Fatalf("Remove account error: account=%s, error=%s", account, err.Error())
			pm.Message = ok.Error()
		} else {
			log.Printf("Delete account: account=%s", account)
		}
	} else {
		pm.Message = "Directory not exists"
	}

	/*
	// Json Marshal
	result, err := json.MarshalIndent(&pm, "", " ")
	if err != nil {
		log.Fatalf("Json Marshal Error: %v", err)
		return
	}
	*/
	WriteResponse(w, &pm)
}

func PostContainer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	container = "./storage/" + account + "/" + container

	pm := PathMessage{Path: container, Message: "Folder was created!"}
	if _, err := os.Stat(container); err != nil {
		if ok := os.MkdirAll(container, 0755); ok != nil {
			log.Fatalf("Make container error: accout=%s, error=%s", container, err.Error())
			pm.Message = ok.Error()
		} else {
			log.Printf("Make container: container=%s", container)
		}
	} else {
		pm.Message = "Direcotry already exists"
	}

	/*
	// Json Marshal
	result, err := json.MarshalIndent(&pm, "", " ")
	if err != nil {
		log.Fatalf("Json Marshal Error: %v", err)
		return
	}
	*/

	WriteResponse(w, &pm)

}

func DeleteContainer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	container = "./storage/" + account + "/" + container

	pm := PathMessage{Path: container, Message: "Folder was deleted!"}
	if _, err := os.Stat(container); err == nil {
		if ok := os.RemoveAll(container); ok != nil {
			log.Fatalf("Remove container error: container=%s, error=%s", container, err.Error())
			pm.Message = ok.Error()
		} else {
			log.Printf("Delete container: container=%s", container)
		}
	} else {
		pm.Message = "Directory not exists"
	}

	WriteResponse(w, &pm)
}

func PutBuckets(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	container = "./storage/" + account + "/" + container

	pm := PathMessage{Path: container, Message: "File was uploaded successfully"}

	// Parse the multipart form in the request
	err := r.ParseMultipartForm(1000000)
	if err != nil {
		pm.Message = err.Error()
		WriteResponse(w, &pm)
		return
	}

	// Get a ref to the parsed multipart form
	m := r.MultipartForm

	// Get the *file
	files := m.File["file"]
	for i, _ := range files {
		// For each *file, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			pm.Message = err.Error()
			break
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create(container + "/" + files[i].Filename)
		defer dst.Close()
		if err != nil {
			pm.Message = err.Error()
			break
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			pm.Message = err.Error()
			break
		}
	}

	WriteResponse(w, &pm)
}

func PostBucket(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	bucket := params.ByName("bucket")
	bucket = "./storage/" + account + "/" + container + "/" + bucket

	pm := PathMessage{Path: container, Message: "File was uploaded successfully"}

	// Parse the multipart form in the request
	err := r.ParseMultipartForm(1000000)
	if err != nil {
		pm.Message = err.Error()
		WriteResponse(w, &pm)
	}

	// Get a ref to the parsed multipart form
	m := r.MultipartForm

	// Get the *file
	files := m.File["file"]
	if 1 != len(files) {
		pm.Message = "You must upload only one file at a time on this method"
		WriteResponse(w, &pm)
		return
	}

	// Get files[0] as the file witch will be uploaded
	file := files[0]
	src, err := file.Open()
	defer src.Close()
	if err != nil {
		pm.Message = err.Error()
	}

	// Create destination file making sure the path is writeable.
	dst, err := os.Create(bucket)
	defer dst.Close()
	if err != nil {
		pm.Message = err.Error()
	}
	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		pm.Message = err.Error()
	}

	WriteResponse(w, &pm)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	bucket := params.ByName("bucket")
	bucket = "./storage/" + account + "/" + container + "/" + bucket

	pm := PathMessage{Path: bucket, Message: "file was deleted!"}
	if _, err := os.Stat(bucket); err == nil {
		if ok := os.Remove(bucket); ok != nil {
			log.Fatalf("Remove bucket error: bucket=%s, error=%s", bucket, err.Error())
			pm.Message = ok.Error()
		} else {
			log.Printf("Delete bucket: bucket=%s", bucket)
		}
	} else {
		pm.Message = "File not exists"
	}

	WriteResponse(w, &pm)
}

func GetBucket(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	bucket := params.ByName("bucket")
	bucket = "./storage/" + account + "/" + container + "/" + bucket
	log.Printf("Get Bucket: bucket=%s", bucket)
	if _, err := os.Stat(bucket); err != nil {
		log.Printf("Get Bucket error: err=%s", err.Error())
		http.Error(w, http.StatusText(404), 404)
		return
	}

	http.ServeFile(w, r, bucket)
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	router := httprouter.New()
	router.GET("/version", Version)
	router.GET("/healthz", Healthz)
	router.POST("/api/:version/:account", PostAccount)
	router.DELETE("/api/:version/:account", DeleteAccount)
	router.POST("/api/:version/:account/:container", PostContainer)
	router.DELETE("/api/:version/:account/:container", DeleteContainer)
	router.GET("/api/:version/:account/:container/:bucket", GetBucket)
	router.PUT("/api/:version/:account/:container", PutBuckets)
	router.POST("/api/:version/:account/:container/:bucket", PostBucket)
	router.DELETE("/api/:version/:account/:container/:bucket", DeleteBucket)
	log.Fatal(http.ListenAndServe(":8088", router))
}
