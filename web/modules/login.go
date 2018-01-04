package webmods

import (
  "fmt"
  "net/http"
  
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
  session.Values["level"] = base.Query("SELECT level FROM user WHERE name = '"+user.Name+"'")
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
  session.Values["level"] = nil
  session.Save(r, w)
  
  w.Header().Set("Location", "/")
  w.WriteHeader(http.StatusTemporaryRedirect)
}