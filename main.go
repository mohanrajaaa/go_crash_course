package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

var (
	svc        *s3.S3
	bucketName = "s3://t2tbucket/" // Change this to your bucket name
)

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/files", createFile).Methods("POST")
	router.HandleFunc("/files/{filename}", readFile).Methods("GET")
	router.HandleFunc("/files/{filename}", updateFile).Methods("PUT")
	router.HandleFunc("/files/{filename}", deleteFile).Methods("DELETE")

	// Initialize AWS Session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		fmt.Println("Failed to create AWS session,", err)
		return
	}
	svc = s3.New(sess)

	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", router)
}

func createFile(w http.ResponseWriter, r *http.Request) {
	var file File
	_ = json.NewDecoder(r.Body).Decode(&file)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file.Name),
		Body:   aws.ReadSeekCloser(strings.NewReader(file.Content)),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(file)
}

func readFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer result.Body.Close()

	content, err := ioutil.ReadAll(result.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(string(content))
}

func updateFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	var file File
	_ = json.NewDecoder(r.Body).Decode(&file)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   aws.ReadSeekCloser(strings.NewReader(file.Content)),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(file)
}
func deleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
