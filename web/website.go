package web

import (
  "net/http"
  "fmt"
  
  "github.com/Mstiekema/Yucibot2/web/modules"

  "github.com/spf13/viper"
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "github.com/markbates/goth"
  "github.com/markbates/goth/gothic"
  "github.com/markbates/goth/providers/twitch"
)

var store = sessions.NewCookieStore([]byte("supersecretkey"))

func MainWeb() {
  r := mux.NewRouter()
  gothic.Store = store
  viper.SetConfigFile("./config.toml")
  err := viper.ReadInConfig()
  if err != nil {fmt.Println(err)}
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
  goth.UseProviders(twitch.New(viper.GetString("twitch.clientId"), viper.GetString("twitch.clientSecret"), viper.GetString("twitch.loginCallbackUrl")),)
  
  r.HandleFunc("/", webmods.Home)
  r.HandleFunc("/clr", webmods.CLR)
  r.HandleFunc("/user/{username}", webmods.UserPage)
  r.HandleFunc("/songlist/", webmods.TodaySonglist)
  r.HandleFunc("/songlist/{date}", webmods.Songlist)
  r.HandleFunc("/auth/{provider}/callback", webmods.Login)
  r.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler)
  r.HandleFunc("/logout", webmods.Logout)
  r.NotFoundHandler = http.HandlerFunc(webmods.Error)
  
  r.HandleFunc("/admin/songlist", webmods.AdminSonglist)
  
  hub := webmods.NewHub()
  go hub.Run()
  
  r.HandleFunc("/post/getSongs/", webmods.SendSongs)
  r.HandleFunc("/post/getCLR/", func(w http.ResponseWriter, r *http.Request) {webmods.SendCLR(hub, w, r)})
  
  http.ListenAndServe(":9090", r)
}