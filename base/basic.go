package base

import (
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strconv"
  "strings"
  "time"
  "github.com/spf13/viper"
)

type StreamInfo struct {
	Data []struct {
		Title string `json:"title"`
		ViewerCount int `json:"viewer_count"`
	} `json:"data"`
}

func (b *Bot) UserInfoComms(C string, U User) {
  if C == "!points" {
    var res = Query("SELECT points FROM user WHERE name = '"+U.username+"'")
    if res != "" {b.SendWhisper("You currently have " + res + " points", U.username)}
  } else if C == "!userpoints" {
    var res = Query("SELECT points FROM user WHERE name = '"+U.username+"'")
    if res != "" {b.SendMsg(U.displayName + " currently has " + res + " points")}
  } else if C == "!lines" {
    var res = Query("SELECT num_lines FROM user WHERE name = '"+U.username+"'")
    if res != "" {b.SendWhisper("You have currently written " + res + " lines", U.username)}
  } else if C == "!rq" {
    line := Query("SELECT log FROM chatlogs WHERE userId = '"+U.userId+"' ORDER BY RAND() LIMIT 1")
    if line != "" {b.SendMsg(U.displayName+": "+line)}
  }
}

func (b *Bot) Basic(C string, U User) {
  viper.SetConfigFile("./config.toml")
  viper.ReadInConfig()
  if C == "!test" {
    b.SendMsg("This is a test message created by " + U.displayName + " FeelsGoodMan")
  } else if C == "!followage" {
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
  } else if C == "!followsince" {
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
  } else if C == "!viewers" {
    info := getStreamInfo(b.Channel)
    if len(info.Data) != 0 {
      b.SendMsg(b.Channel+" currently has "+strconv.Itoa(info.Data[0].ViewerCount)+" viewers")
    } else {
      b.SendMsg(b.Channel+" is currently offline")
    }
  } else if C == "!title" {
    info := getStreamInfo(b.Channel)
    if len(info.Data) != 0 {
      b.SendMsg("Current title: "+info.Data[0].Title)
    } else {
      b.SendMsg(b.Channel+" is currently offline")
    }
  } else if C == "!uptime" {
    resp, _ := http.Get("http://api.yucibot.nl/user/uptime/"+b.Channel)
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
}

func getStreamInfo(chnl string) (result StreamInfo) {
  client := &http.Client{Timeout: time.Second * 10,}
  clientid := viper.GetString("twitch.clientId")
  req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_login="+chnl, nil)
  req.Header.Add("Client-ID", clientid)
  resp, _ := client.Do(req)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  info := StreamInfo{}
  json.Unmarshal(body, &info)
  return info
}

func (b *Bot) getFollowAge(user, chnl string) {
  resp, _ := http.Get("http://api.yucibot.nl/followage/"+user+"/"+chnl)
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
  resp, _ := http.Get("http://api.yucibot.nl/followsince/"+user+"/"+chnl)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  text := string(body[:])
  if strings.Contains(text, "</html>") {
    return
  } else {
    b.SendMsg(text)
  }
}