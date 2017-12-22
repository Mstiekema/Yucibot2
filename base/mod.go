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

func (b *Bot) ModifyCommands(C string, U User) {
  if U.mod == "true" || U.username == strings.ToLower(b.Channel) {
    split := strings.SplitAfter(U.message, " ")
    if C == "!addcom" || C == "!addcommand" {
      if 2 < len(strings.SplitAfter(U.message, " ")) {
        commName := split[1]
        if strings.HasPrefix(commName, "!") != true {commName = "!"+commName}
        commResp := strings.Join(append(split[2:]), "")
        Insert("commands (commName, response, cdType, cd)", "('"+commName+"', '"+commResp+"', 'global', '10')")
        b.SendWhisper("Succesfully added "+commName+" to the database.", U.username)
      }
    } else if C == "!editcom" || C == "!editcommand" {
      if 2 < len(strings.SplitAfter(U.message, " ")) {
        commName := split[1]
        if strings.HasPrefix(commName, "!") != true {commName = "!"+commName}
        commResp := strings.Join(append(split[2:]), "")
        Update("commands", "response = '"+commResp+"'", "commName", "'"+commName+"'")
        b.SendWhisper("Succesfully edited the "+commName+" command.", U.username)
      }
    } else if C == "!remcom" || C == "!removecommand" {
      if 1 < len(strings.SplitAfter(U.message, " ")) {
        commName := strings.TrimSpace(split[1])
        if strings.HasPrefix(commName, "!") != true {commName = "!"+commName}
        Delete("commands", "commName", commName)
        b.SendWhisper("Succesfully removed "+commName+" from the database.", U.username)
      }
    }
  }
}