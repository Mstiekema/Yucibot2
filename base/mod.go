package base

import (
  "net/url"
  "strings"
  "time"
)

var allowedUser string

func (b *Bot) Links(U User) {  
  msgSplit := make([]string, len(strings.SplitAfter(U.message, " ")))
  var c string
  var u string
  
  if  1 < len(strings.SplitAfter(U.message, " ")) {
    msgSplit = strings.Split(U.message, " ")
    c = strings.TrimSpace(strings.Split(U.message, " ")[0])
    u = strings.TrimSpace(strings.Split(U.message, " ")[1])
  } else {
    msgSplit[0] = U.message
    c = U.message
  }
  
  if  strings.ToLower(allowedUser) != U.username { if U.sub != "true" { if U.mod != "true" {
    for i := 0; i < len(msgSplit); i++ {
      _, err := url.ParseRequestURI(msgSplit[i])
      if err == nil {
        b.SendMsg(`.timeout `+U.username+` 30 Only subs are allowed to post links`)
        return
      }
    }
  }}}
  
  if c == "!permit" && (U.mod == "true" || U.username == strings.ToLower(b.Channel)) {
    allowedUser = u
    b.SendMsg(u+" is now allowed to post links for 30 seconds!")
    time.AfterFunc(30 * time.Second, func() {
      allowedUser = ""
    })  
  }
}