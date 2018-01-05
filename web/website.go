package web

import (
  "net/http"
  "fmt"

  "github.com/spf13/viper"
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "github.com/markbates/goth"
  "github.com/markbates/goth/gothic"
  "github.com/Mstiekema/Yucibot2/web/modules"
  "github.com/markbates/goth/providers/twitch"
)

func MainWeb() {
  r := mux.NewRouter()
  viper.SetConfigFile("./config.toml")
  err := viper.ReadInConfig(); if err != nil {fmt.Println(err)}
  gothic.Store = sessions.NewCookieStore([]byte(viper.GetString("apiKeys.secretCookieKey")))
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
  goth.UseProviders(twitch.New(viper.GetString("twitch.clientId"), viper.GetString("twitch.clientSecret"), viper.GetString("twitch.loginCallbackUrl")),)
  
  r.HandleFunc("/", webmods.Home)
  r.HandleFunc("/commands", webmods.Commands)
  r.HandleFunc("/stats", webmods.Stats)
  r.HandleFunc("/songlist/", webmods.TodaySonglist)
  r.HandleFunc("/songlist/{date}", webmods.Songlist)
  go r.HandleFunc("/user/{username}", webmods.UserPage)
  r.HandleFunc("/clr", webmods.CLR)
  r.HandleFunc("/auth/{provider}/callback", webmods.Login)
  r.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler)
  r.HandleFunc("/logout", webmods.Logout)
  r.NotFoundHandler = http.HandlerFunc(webmods.Error)
  
  r.HandleFunc("/admin/songlist", webmods.AdminSonglist)
  r.HandleFunc("/admin/modules", webmods.AdminModules)
  r.HandleFunc("/admin/clr", webmods.AdminClr)
  r.HandleFunc("/admin/commands", webmods.AdminCommands)
  
  hub := webmods.NewHub()
  go hub.Run()
  
  r.HandleFunc("/post/getSongs/", webmods.SendSongs)
  r.HandleFunc("/post/adminClr/", func(w http.ResponseWriter, r *http.Request) {webmods.PostAdminClr(hub, w, r)})
  r.HandleFunc("/post/getCLR/", func(w http.ResponseWriter, r *http.Request) {webmods.SendCLR(hub, w, r)})
  
  http.ListenAndServe(":9090", r)
}