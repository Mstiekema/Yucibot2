package web

import (
  "net/http"
  "github.com/Mstiekema/Yucibot2/web/modules"
)

func MainWeb() {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
  
  http.HandleFunc("/", webmods.Home)
  http.HandleFunc("/clr", webmods.CLR)
  http.HandleFunc("/test", webmods.Test)
  http.HandleFunc("/user/", webmods.UserPage)
  http.HandleFunc("/songlist/", webmods.Songlist)
  
  http.HandleFunc("/admin/songlist/", webmods.AdminSonglist)
  
  hub := webmods.NewHub()
  go hub.Run()
  
  http.HandleFunc("/post/getSongs/", webmods.SendSongs)
  http.HandleFunc("/post/getCLR/", func(w http.ResponseWriter, r *http.Request) {webmods.SendCLR(hub, w, r)})
  
  http.ListenAndServe(":9090", nil)
}