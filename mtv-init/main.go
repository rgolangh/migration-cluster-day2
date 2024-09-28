package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const savePath = "/tmp/vmware"

func main() {
	http.HandleFunc("/", bootstrapHandler)
	http.HandleFunc("/migrations/bootstrap", bootstrapHandler)
	http.HandleFunc("/upload", uploadHandler)
    err := os.Mkdir(savePath, os.ModePerm)
    if err != nil && !errors.Is(err, os.ErrExist){
        panic(err)
    }
	http.Handle("/vmware/", http.StripPrefix("/vmware/",http.FileServer(http.Dir(savePath))))
	fmt.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", nil)
}

func bootstrapHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("form").Parse(fff)
	if err != nil {
		http.Error(w, "Unable to load form", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Max file size 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	// Retrieve file and other form data
	file, _, err := r.FormFile("vddk")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	target, err := os.Create(path.Join(savePath,"vddk.tar.gz"))
	if err != nil {
		http.Error(w, "Error creating destination file", http.StatusInternalServerError)
		return
	}
	defer target.Close()

	_, err = io.Copy(target, file)
	if err != nil {
		http.Error(w, "Error writing destination file", http.StatusInternalServerError)
		return
	}

	envFile, err := os.Create(path.Join(savePath,"env"))
	if err != nil {
		http.Error(w, "Error creating destination file", http.StatusInternalServerError)
		return
	}
	defer envFile.Close()

    j := fmt.Sprintf("url=%s\nusername=%s\npassword=%s\n",
		r.FormValue("url"),
		r.FormValue("username"),
		r.FormValue("password"))
	_, err = io.Copy(envFile, strings.NewReader(j))
	if err != nil {
		http.Error(w, "Error writing destination env file", http.StatusInternalServerError)
		return
	}

    err = os.WriteFile(path.Join(savePath, "done"), nil, os.ModePerm)
    if err != nil {
		http.Error(w, "Error writing done file", http.StatusInternalServerError)
        return 
    }
	// For now, just return a simple confirmatio
	fmt.Fprintf(w, "<html><body>File and vmware credentials recieved and avaiable under <a href=\"/vmware\" />/vmware</a></body></html>\n")
}

const fff = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <link rel="stylesheet" href="fonts.css" />
  <!-- Include latest PatternFly CSS via CDN -->
  <link rel="stylesheet" href="https://unpkg.com/@patternfly/patternfly/patternfly.css" crossorigin="anonymous" >
  <link rel="stylesheet" href="style.css" />
  <title>PatternFly Basic CodeSandbox Example</title>
</head>
<body>
  <header class="pf-v5-c-masthead" id="basic-masthead">
  <span class="pf-v5-c-masthead__toggle">
    <button
      class="pf-v5-c-button pf-m-plain"
      type="button"
      aria-label="Global navigation"
    >
      <i class="fas fa-bars" aria-hidden="true"></i>
    </button>
  </span>
  <div class="pf-v5-c-masthead__main">
    <a class="pf-v5-c-masthead__brand" href="#">Migration Preparations</a>
  </div>
  <div class="pf-v5-c-masthead__content">
  </div>
  </header>
   <!-- Main content -->
   <main class="pf-c-page__main">
       <div class="pf-c-form">
           <h1 class="pf-c-title pf-m-2xl"></h1>
           <form action="/upload" method="post" enctype="multipart/form-data" class="pf-c-form">
               <div class="pf-c-form__group">
                   <label for="vddk" class="pf-c-form__label">VDDK File</label>
                   <input type="file" name="vddk" id="vddk" class="pf-c-form-control" required />
               </div>
               <div class="pf-c-form__group">
                   <label for="username" class="pf-c-form__label">vSphere Username</label>
                   <input type="text" name="username" id="username" class="pf-c-form-control" required />
               </div>
               <div class="pf-c-form__group">
                   <label for="password" class="pf-c-form__label">vSphere Password</label>
                   <input type="password" name="password" id="password" class="pf-c-form-control" required />
               </div>
               <div class="pf-c-form__group">
                   <label for="url" class="pf-c-form__label">vSphere URL</label>
                   <input type="url" name="url" id="url" class="pf-c-form-control" required />
               </div>
               <button type="submit" class="pf-c-button pf-m-primary">Submit</button>
           </form>
       </div>
   </main>
</body>
</html>
`
