package base

import (
  "strings"
  "mvdan.cc/xurls"
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
  
  if strings.ToLower(allowedUser) != U.username { if U.sub != "1" { if U.mod != "1" {
    url := xurls.Relaxed().FindString(U.message)
    if url != "" {
      b.SendMsg(`.timeout `+U.username+` 30 Only subs are allowed to post links`)
      b.SendWhisper(`Only subs are allowed to post links`, U.username)
      return
    }
  }}}
  
  if c == "!permit" && (U.mod == "1" || U.username == strings.ToLower(b.Channel)) {
    allowedUser = u
    b.SendMsg(u+" is now allowed to post links for 30 seconds!")
    time.AfterFunc(30 * time.Second, func() {
      allowedUser = ""
    })  
  }
}

func (b *Bot) ModifyCommands(C string, U User) {
  if U.mod == "1" || U.username == strings.ToLower(b.Channel) {
    split := strings.SplitAfter(U.message, " ")
    if C == "!addpoints" {
      if 2 < len(split) {
        Update("user", "points = points + '"+strings.TrimSpace(split[2])+"'", "name", "'"+strings.ToLower(strings.TrimSpace(split[1]))+"'")
        b.SendWhisper("Succesfully gave "+strings.TrimSpace(split[1])+" "+strings.TrimSpace(split[2])+" points.", U.username)
      }
    } else if C == "!addcom" || C == "!addcommand" {
      if 2 < len(split) {
        commName := split[1]
        if strings.HasPrefix(commName, "!") != true {commName = "!"+commName}
        commResp := strings.Join(append(split[2:]), "")
        Insert("commands (commName, response, cdType, cd)", "('"+commName+"', '"+commResp+"', 'global', '10')")
        b.SendWhisper("Succesfully added "+commName+" to the database.", U.username)
      }
    } else if C == "!editcom" || C == "!editcommand" {
      if 2 < len(split) {
        commName := split[1]
        if strings.HasPrefix(commName, "!") != true {commName = "!"+commName}
        commResp := strings.Join(append(split[2:]), "")
        Update("commands", "response = '"+commResp+"'", "commName", "'"+commName+"'")
        b.SendWhisper("Succesfully edited the "+commName+" command.", U.username)
      }
    } else if C == "!remcom" || C == "!removecommand" {
      if 1 < len(split) {
        commName := strings.TrimSpace(split[1])
        if strings.HasPrefix(commName, "!") != true {commName = "!"+commName}
        Delete("commands", "commName", commName)
        b.SendWhisper("Succesfully removed "+commName+" from the database.", U.username)
      }
    }
  }
}