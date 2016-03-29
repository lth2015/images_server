package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Version
const (
	VERSION = "v2"
	ROOT = "./storage/"
)

// For create/delete Folders
type PathMessage struct {
	Path    string `json:string "Path"`
	Message string `json:string "Message"`
}

type Directories struct {
	Dirs []string `json:string "Dirs"`
}

type Files struct {
	Files []string `json:string "Files"`
}

func MakeDir(dir string, pm *PathMessage) {
	if _, err := os.Stat(dir); err != nil {
		if ok := os.MkdirAll(dir, 0755); ok != nil {
			log.Fatalf("Make directory error: directory=%s, error=%s", dir, err.Error())
			pm.Message = ok.Error()
		} else {
			log.Printf("Make directory: directory=%s", dir)
		}
	} else {
		pm.Message = "Direcotry already exists"
	}
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

func GetAccounts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var dirs Directories
	files, err := ioutil.ReadDir(ROOT)
	if err != nil  {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			dirs.Dirs = append(dirs.Dirs, file.Name())
		}
	}

	log.Printf("List Account: account=%v", dirs.Dirs)
	WriteResponse(w, &dirs)
}

func GetContainers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	account = ROOT + account

	var dirs Directories
	files, err := ioutil.ReadDir(account)
	if err != nil  {
		log.Fatal(err)
		pm := PathMessage{Path: account, Message: err.Error()}
		WriteResponse(w, &pm)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			dirs.Dirs = append(dirs.Dirs, file.Name())
		}
	}

	log.Printf("List Containers: containers=%v",dirs.Dirs)
	WriteResponse(w, &dirs)
}

func PostAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	account = ROOT + account

	pm := PathMessage{Path: account, Message: "Folder was created!"}
	MakeDir(account, &pm)
	WriteResponse(w, &pm)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	account = ROOT + account

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

	WriteResponse(w, &pm)
}

func PostContainer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	container = ROOT + account + "/" + container

	pm := PathMessage{Path: container, Message: "Folder was created!"}

	MakeDir(container, &pm)
	WriteResponse(w, &pm)
}

func DeleteContainer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	container = ROOT + account + "/" + container

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

func GetBuckets(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	bucket := "ROOT" + account + "/" + container + "/"

	var dirs Directories
	files, err := ioutil.ReadDir(bucket)
	if err != nil  {
		log.Fatal(err)
		pm := PathMessage{Path: account, Message: err.Error()}
		WriteResponse(w, &pm)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			dirs.Dirs = append(dirs.Dirs, file.Name())
		}
	}

	log.Printf("List Buckets: buckets=%v", dirs.Dirs)
	WriteResponse(w, &dirs)

}

func PutBuckets(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	container = ROOT + account + "/" + container


	pm := PathMessage{Path: container, Message: "File was uploaded successfully"}
	MakeDir(container, &pm)

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
		if err != nil{
			log.Fatalf("Open file error: error=%s", err.Error())
			pm.Message = err.Error()
			break
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create(container + "/" + files[i].Filename)
		defer dst.Close()
		if err != nil {
			log.Fatalf("Create file error: error=%s", err.Error())
			pm.Message = err.Error()
			break
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			log.Fatalf("Copy file error: error=%s", err.Error())
			pm.Message = err.Error()
			break
		}

		log.Printf("Upload file successfully: file=%s", dst.Name())
	}

	pm.Message = "Files are already uploaded!"
	WriteResponse(w, &pm)
}

func PostBucket(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	bucket := params.ByName("bucket")
	bucket = ROOT + account + "/" + container + "/" + bucket
	container = ROOT + account + "/" + container

	pm := PathMessage{Path: container, Message: "File was uploaded successfully"}
	MakeDir(container, &pm)

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
		log.Fatalf("Open file error: error=%s", err.Error())
		pm.Message = err.Error()
	}

	// Create destination file making sure the path is writeable.
	dst, err := os.Create(bucket)
	defer dst.Close()
	if err != nil {
		log.Fatalf("Create file error: error=%s", err.Error())
		pm.Message = err.Error()
	}
	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalf("Copy file error: error=%s", err.Error())
		pm.Message = err.Error()
	}

	pm.Message = "File is already uploaded!"
	WriteResponse(w, &pm)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	container := params.ByName("container")
	bucket := params.ByName("bucket")
	bucket = ROOT + account + "/" + container + "/" + bucket

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
	bucket = ROOT + account + "/" + container + "/" + bucket
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
	router.GET("/api/:version/accounts", GetAccounts)

	router.POST("/api/:version/accounts/:account", PostAccount)
	router.DELETE("/api/:version/accounts/:account", DeleteAccount)

	router.GET("/api/:version/accounts/:account/containers", GetContainers)
	router.PUT("/api/:version/accounts/:account/containers/:container", PutBuckets)
	router.POST("/api/:version/accounts/:account/containers/:container", PostContainer)
	router.DELETE("/api/:version/accounts/:account/containers/:container", DeleteContainer)

	router.GET("/api/:version/accounts/:account/containers/:container", GetBuckets)
 	router.GET("/api/:version/accounts/:account/containers/:container/buckets/:bucket", GetBucket)
	router.POST("/api/:version/accounts/:account/containers/:container/buckets/:bucket", PostBucket)
	router.DELETE("/api/:version/accounts/:account/containers/:container/buckets/:bucket", DeleteBucket)
	log.Fatal(http.ListenAndServe(":8088", router))
}
