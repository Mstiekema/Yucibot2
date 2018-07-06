package base

import (
  "os"
  "fmt"
  "time"
  "strconv"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "github.com/spf13/viper"
)

type StreamInfo struct {
	Data []struct {
		Title string `json:"title"`
		ViewerCount int `json:"viewer_count"`
	} `json:"data"`
}

func (b *Bot) UserInfoComms(C string, U User) {
  if C == "!quit" && (U.mod == "1" || U.username == strings.ToLower(b.Channel)) {
    b.SendMsg("Shutting down Yucibot MrDestructoid")
    os.Exit(0)
  } else if C == "!points" {
    exec := func() {
      var res = Query("SELECT points FROM user WHERE name = '"+U.username+"'")
      if res != "" {b.SendWhisper("You currently have " + res + " points", U.username)}
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  } else if C == "!userpoints" {
    exec := func() {
      var res = Query("SELECT points FROM user WHERE name = '"+U.username+"'")
      if res != "" {b.SendMsg(U.displayName + " currently has " + res + " points")}
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  } else if C == "!lines" {
    exec := func() {
      var res = Query("SELECT num_lines FROM user WHERE name = '"+U.username+"'")
      if res != "" {b.SendWhisper("You have currently written " + res + " lines", U.username)}
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  } else if C == "!rq" {
    exec := func() {
      line := Query("SELECT log FROM chatlogs WHERE userId = '"+U.userId+"' ORDER BY RAND() LIMIT 1")
      if line != "" {b.SendMsg(U.displayName+": "+line)}
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  } else if C == "!givepoints" {
    exec := func() {
      split := strings.SplitAfter(U.message, " ")
      if 2 < len(split) {
        giverPoints, _ := strconv.Atoi(Query("SELECT points FROM user WHERE name = '"+U.username+"'"))
        toAddPoints, _ := strconv.Atoi(split[2])
        if giverPoints >= toAddPoints && toAddPoints > 0 {
          Update("user", "points = points + '"+strings.TrimSpace(split[2])+"'", "name", "'"+strings.ToLower(strings.TrimSpace(split[1]))+"'")
          Update("user", "points = points - '"+strings.TrimSpace(split[2])+"'", "name", "'"+U.username+"'")
          b.SendWhisper("Succesfully gave "+strings.TrimSpace(split[1])+" "+strings.TrimSpace(split[2])+" points.", U.username)
        } else {
          b.SendWhisper("You can't give another user more points than you have or steal by giving negative points.", U.username)
        }
      }
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  } else if C == "!vanish" {
    exec := func() {
      b.SendMsg(`.timeout `+U.username+` 1 `+U.displayName+` vanished`)
    }
    b.ExecuteCommand(C, "100", "0", "5", U, exec)
  }
}

func (b *Bot) CustomCommands(C string, U User) {
  var db = Conn()
  var response, commUse, level, points, cd string
  res, err := db.Query(`SELECT response, commUse, level, points, cd FROM commands WHERE commDesc IS NULL AND commName = "`+C+`"`)
  if err != nil {
    fmt.Println(err)
    return
  }
  for res.Next() {
    res.Scan(&response, &commUse, &level, &points, &cd)
    if response != "" {
      exec := func() {
        if strings.Contains(response, "&Count") {
          i, _ := strconv.Atoi(commUse); i++; s := strconv.Itoa(i); Update("commands", `commUse = "`+s+`"`, "commName", "'"+C+"'")
          response = strings.Replace(response, "&Count", s, -1)
        }
        b.SendMsg(response)
      }
      b.ExecuteCommand(C, level, points, cd, U, exec)
    }
  }
  db.Close()
}

var uCooldowns []string
var gCooldowns []string
type fn func()

func (b *Bot) ExecuteCommand(C, level, points, cd string, U User, exec fn) {
  toCd := C+U.username
  for _, b := range gCooldowns { if b == C { return } }
  for _, b := range uCooldowns { if b == toCd { return } }
  if level != "100" {
    if U.mod == "1" { U.sub = "1" }
    if level == "150" && U.sub == "0" {
      b.SendWhisper("You have to be a subscriber to use this command.", U.username)
      return
    }
    if level == "300" && !(U.username == strings.ToLower(b.Channel) || U.mod == "1") {
      b.SendWhisper("You have to be a moderator to use this command.", U.username)
      return
    }
  }
  Points, _ := strconv.Atoi(points)
  if Points != 0 {
    uPoints, _ := strconv.Atoi(Query(`SELECT points FROM user WHERE name = "`+U.username+`"`))
    if uPoints >= Points {
      Update("user", "points = points - "+points, "name", "'"+U.username+"'")
    } else {
      b.SendWhisper("You do not have enough points to use this command.", U.username)
      return
    }
  }
  uCooldowns = append(uCooldowns, toCd)
  gCooldowns = append(gCooldowns, C)
  exec()
  cdi, _ := strconv.Atoi(cd)
  time.AfterFunc(time.Duration(cdi) * time.Second, func() {
    for i, v := range uCooldowns {
      if v == toCd {
        if len(uCooldowns) == 0 {uCooldowns = uCooldowns[:0]}
        uCooldowns = append(uCooldowns[:i], uCooldowns[i+1:]...)
      }
    }
  })
  gCd := 10
  if C == "!sr" || C == "!songrequest" {gCd = 1}
  time.AfterFunc(time.Duration(gCd) * time.Second, func() {
    for i, v := range gCooldowns {
      if v == C {
        if len(gCooldowns) == 0 {gCooldowns = gCooldowns[:0]}
        gCooldowns = append(gCooldowns[:i], gCooldowns[i+1:]...)
      }
    }
  })
}

func (b *Bot) Basic(C string, U User) {
  viper.SetConfigFile("./config.toml")
  viper.ReadInConfig()
  if C == "!followage" {
    exec := func() {
      if 2 == len(strings.SplitAfter(U.message, " ")) {
        user := strings.TrimSpace(strings.SplitAfter(U.message, " ")[1])
        b.getFollowAge(user, b.Channel)
      } else if 2 < len(strings.Split(U.message, " ")) {
        user := strings.TrimSpace(strings.SplitAfter(U.message, " ")[1])
        chnl := strings.TrimSpace(strings.Split(U.message, " ")[2])
        b.getFollowAge(user, chnl)
      } else {
        b.getFollowAge(U.displayName, b.Channel)
      }
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  } else if C == "!followsince" {
    exec := func() {
      if 2 == len(strings.SplitAfter(U.message, " ")) {
        user := strings.TrimSpace(strings.SplitAfter(U.message, " ")[1])
        b.getFollowSince(user, b.Channel)
      } else if 2 < len(strings.Split(U.message, " ")) {
        user := strings.TrimSpace(strings.SplitAfter(U.message, " ")[1])
        chnl := strings.TrimSpace(strings.Split(U.message, " ")[2])
        b.getFollowSince(user, chnl)
      } else {
        b.getFollowSince(U.displayName, b.Channel)
      }
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  } else if C == "!viewers" {
    exec := func() {
      info := b.GetStreamInfo()
      if len(info.Data) != 0 {
        b.SendMsg(b.Channel+" currently has "+strconv.Itoa(info.Data[0].ViewerCount)+" viewers")
      } else {
        b.SendMsg(b.Channel+" is currently offline")
      }
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  } else if C == "!title" {
    exec := func() {
      info := b.GetStreamInfo()
      if len(info.Data) != 0 {
        b.SendMsg("Current title: "+info.Data[0].Title)
      } else {
        b.SendMsg(b.Channel+" is currently offline")
      }
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  } else if C == "!uptime" {
    exec := func() {
      resp, _ := http.Get("http://api.yucibot.com/user/uptime/"+b.Channel)
      defer resp.Body.Close()
      body, _ := ioutil.ReadAll(resp.Body)
      text := string(body[:])
      if strings.Contains(text, "</html>") {
        return
      } else if !strings.Contains(text, "is not live") {
        b.SendMsg(b.Channel+" has been live for "+text)
      } else {
        b.SendMsg(text)
      }
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  }
}

func (b *Bot) GetStreamInfo() (result StreamInfo) {
  client := &http.Client{Timeout: time.Second * 10,}
  clientid := viper.GetString("twitch.clientId")
  req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_login="+b.Channel, nil)
  req.Header.Add("Client-ID", clientid)
  resp, _ := client.Do(req)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  info := StreamInfo{}
  json.Unmarshal(body, &info)
  return info
}

func (b *Bot) getFollowAge(user, chnl string) {
  resp, _ := http.Get("http://api.yucibot.com/followage/"+user+"/"+chnl)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  text := string(body[:])
  if strings.Contains(text, "is not following") || strings.Contains(text, "does not exist") {
    b.SendMsg(text)
  } else if strings.Contains(text, "</html>") {
    return
  } else {
    b.SendMsg(user+" has been following "+chnl+" for "+text)
  }
}

func (b *Bot) getFollowSince(user, chnl string) {
  resp, _ := http.Get("http://api.yucibot.com/followsince/"+user+"/"+chnl)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  text := string(body[:])
  if strings.Contains(text, "</html>") {
    return
  } else {
    b.SendMsg(text)
  }
}
