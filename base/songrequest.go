package base

import (
  "regexp"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "time"
  "strings"
  "strconv"
  "fmt"
)

type VidByWords struct {
  Items []struct {
    Vid struct {
      Id string `json:"videoId"`
    } `json:"id"`
  } `json:"items"`
}

type VidInfo struct {
  Items []struct {
    Snippet struct {
      Title string `json:"title"`
      Thumbnails struct {
        Default struct {
          Url string `json:"url"`
        } `json:"default"`
      } `json:"thumbnails"`
    } `json:"snippet"`
    ContentDetails struct {
      Time string `json:"duration"`
    } `json:"contentDetails"`
  } `json:"items"`
}

func (b *Bot) Songrequest(C string, U User) {
  if C == "!sr" || C == "!songrequest" {
    msgSplit := strings.SplitAfter(U.message, " ")
    i := 1
    if  (i >= 1 && i < len(strings.SplitAfter(U.message, " "))) {
      regex, _ := regexp.Compile(`(?:https?:\/{2})?(?:w{3}\.)?youtu(?:be)?\.(?:com|be)(?:\/watch\?v=|\/)([^\s&]+)`)
      id := strings.SplitAfter(regex.FindString(msgSplit[1]) , "?v=")
      if len(id) >= 2 {
        b.getLinkInfo(id[1], C, U)
      } else {
        if len(msgSplit[1]) == 11 {
          b.getLinkInfo(msgSplit[1], C, U)
        } else {
          songWordsArr := append(msgSplit[1:])
          song := strings.Join(songWordsArr, "")
          song = strings.Replace(song, " ", "%20", -1)
          var netClient = &http.Client{Timeout: time.Second * 10,}
          resp, _ := netClient.Get("https://www.googleapis.com/youtube/v3/search?part=id&q="+song+"&key="+b.ytApiKey)
          defer resp.Body.Close()
          
          body, _ := ioutil.ReadAll(resp.Body)
          vidWords := VidByWords{}
          err := json.Unmarshal(body, &vidWords)
          if err != nil {
            fmt.Println(err)
            return
          }
          b.getLinkInfo(vidWords.Items[0].Vid.Id, C, U)
        }
      }
    }
  }
}

type SongInfo struct {
  Title string
  Username string
}

func (b *Bot) getLinkInfo(L, C string, U User) {
  var netClient = &http.Client{Timeout: time.Second * 10,}
  resp, _ := netClient.Get("https://www.googleapis.com/youtube/v3/videos?id=" + L + "&key=" + b.ytApiKey + "%20&part=snippet,contentDetails,statistics,status")
  defer resp.Body.Close()
  
  body, _ := ioutil.ReadAll(resp.Body)
  vid := VidInfo{}
  err := json.Unmarshal(body, &vid)
  if err != nil {
    fmt.Println(err)
    return
  }

  tReg := regexp.MustCompile(`PT?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)
  matches := tReg.FindStringSubmatch(vid.Items[0].ContentDetails.Time)
  h, _ := strconv.Atoi(strings.TrimSuffix(matches[1], "H"))
  m, _ := strconv.Atoi(strings.TrimSuffix(matches[2], "M"))
  s, _ := strconv.Atoi(strings.TrimSuffix(matches[3], "S"))
  itime := (h*3600)+(m*60)+s
  time := strconv.Itoa(itime)
  
  var db = Conn()
  res, err := db.Query("select name, songid from songrequest where DATE(time) = CURDATE() AND playState = 0")
  if err != nil {
		panic(err.Error())
	}
  defer res.Close()
  
  var names []string
  var songids []string
  
  for res.Next() {
    var name string
    var songid string
    err = res.Scan(&name, &songid)
    names = append(names, name)
    songids = append(songids, songid)
  }
  
  if checkInSlice(L, songids) > 0 {
    b.SendMsg(U.displayName+", this song is already in the queue.")
  } else if checkInSlice(U.displayName, names) > 2 {
    b.SendMsg(U.displayName+", you already have 3 songs in the queue, please wait ")
  } else {
    Insert("songrequest (title, thumb, name, length, songid)", "('"+vid.Items[0].Snippet.Title+"', '" +vid.Items[0].Snippet.Thumbnails.Default.Url+"', '"+U.displayName+"', '"+time+"', '"+L+"')")
    b.SendMsg("Added "+vid.Items[0].Snippet.Title+" to the queue. Song requested by "+U.displayName)
  }
}

func checkInSlice(a string, list []string) int {
  i := 0
  for _, b := range list {
    if b == a {
      i++
    }
  }
  return i
}