package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

func goog(w http.ResponseWriter, req *http.Request) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// minor hack to check if bucket exists (acts as credential check too.)
	attrs, err := client.Bucket("gcs2aws-bucket-test").Attrs(ctx)
	if err != nil {
		fmt.Fprintf(w, "It could be missing credentials, or: \n%s", err)
		fmt.Println(err)
		return
	}
	fmt.Println(attrs)

	rc, err := client.Bucket("gcs2aws-bucket-test").Object("gcs.txt").NewReader(ctx)
	if err != nil {
		fmt.Fprintf(w, "It could be missing credentials, or: \n%s", err)
		fmt.Println(err)
		return
	}
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		fmt.Fprintf(w, "It could be missing credentials, or: \n%s", err)
		fmt.Println(err)
		return
	}
	rc.Close()

	// Print text data to stdout
	fmt.Println("file contents:", string(slurp))
	// Print text data to screen
	fmt.Fprintf(w, "%s", string(slurp))

}

func aws(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	http.HandleFunc("/goog", goog)
	http.HandleFunc("/aws", aws)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))

}
