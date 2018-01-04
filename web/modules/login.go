package webmods

import (
  "fmt"
  "net/http"
  "strconv"
  
  "github.com/markbates/goth/gothic"
  "github.com/Mstiekema/Yucibot2/base"
)

func Login(w http.ResponseWriter, r *http.Request) {
  user, err := gothic.CompleteUserAuth(w, r)
  if err != nil {
    fmt.Fprintln(w, err)
    return
  }
  
  session, _ := gothic.Store.Get(r, "loginSession")
  session.Values["userId"] = user.UserID
  session.Values["username"] = user.Name
  session.Values["displayName"] = user.NickName
  session.Values["profile_pic"] = user.AvatarURL
  session.Values["points"] = base.Query("SELECT points FROM user WHERE name = '"+user.Name+"'")
  session.Values["lines"] = base.Query("SELECT num_lines FROM user WHERE name = '"+user.Name+"'")
  session.Values["level"], _ = strconv.Atoi(base.Query("SELECT level FROM user WHERE name = '"+user.Name+"'"))
  session.Save(r, w)
  
  w.Header().Set("Location", "/user/"+user.Name)
  w.WriteHeader(http.StatusTemporaryRedirect)
}

func Logout(w http.ResponseWriter, r *http.Request) {
  gothic.Logout(w, r)
  
  session, _ := gothic.Store.Get(r, "loginSession")
  session.Values["userId"] = nil
  session.Values["username"] = nil
  session.Values["displayName"] = nil
  session.Values["profile_pic"] = nil  
  session.Values["points"] = nil
  session.Values["lines"] = nil
  session.Values["level"] = 100
  session.Save(r, w)
  
  w.Header().Set("Location", "/")
  w.WriteHeader(http.StatusTemporaryRedirect)
}