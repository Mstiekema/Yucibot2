package web

import (
  "strings"
  "html/template"
  "fmt"
  "net/http"
  "log"
  "github.com/Mstiekema/Yucibot2/base"
)

type User struct {
	Username string
  Points string
  Lines string
}

func MainWeb() {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
  
  http.HandleFunc("/", home)
  http.HandleFunc("/test", test)
  http.HandleFunc("/user/", userPage)
  
  err := http.ListenAndServe(":9090", nil)
  if err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}

func home(w http.ResponseWriter, r *http.Request){
  // Testing some stuff, some other stuff will be here soon
  channel := "Mstiekema"
  t, err := template.ParseFiles("./web/templates/home.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.Execute(w, channel)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}

func userPage(w http.ResponseWriter, r *http.Request) {
  usr := strings.Replace(strings.SplitAfter(r.URL.Path, "/")[2], "/", "", 2)
  points := base.Query("user", "points", usr)
  lines := base.Query("user", "num_lines", usr)
  u := &User{Username: usr, Points: points, Lines: lines}
  
  t, err := template.ParseFiles("./web/templates/user.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.Execute(w, u)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}

func test(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "This is a test page xD")
}