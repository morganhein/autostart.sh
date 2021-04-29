package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

//starts an http server that provides a shell script to download and test with
func main() {
	// Set routing rules
	http.HandleFunc("/shell", shell)

	fmt.Println("Server staring.")
	//Use the default DefaultServeMux.
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func shell(w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprint(w, "error")
		return
	}
	filename := path.Join(cwd, "test_server/test.sh")
	f, err := os.Open(filename)
	if err != nil {
		_, _ = fmt.Fprint(w, "error")
		return
	}
	defer func() {
		_ = f.Close()
	}()
	info, err := os.Stat(filename)
	if err != nil {
		_, _ = fmt.Fprint(w, "error")
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+path.Base(filename))
	http.ServeContent(w, r, f.Name(), info.ModTime(), f)
}
