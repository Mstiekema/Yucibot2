package web

import (
  "net/http"
  "log"
  "github.com/Mstiekema/Yucibot2/web/modules"
)

func MainWeb() {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
  
  http.HandleFunc("/", webmods.Home)
  http.HandleFunc("/test", webmods.Test)
  http.HandleFunc("/user/", webmods.UserPage)
  http.HandleFunc("/songlist/", webmods.Songlist)
  
  err := http.ListenAndServe(":9090", nil)
  if err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}