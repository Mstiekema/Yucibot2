package webmods

import (
  "log"
  "time"
  "strconv"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "html/template"
  "github.com/spf13/viper"
  "github.com/markbates/goth/gothic"
  "github.com/Mstiekema/Yucibot2/base"
)

type StreamInfo struct {
	Data []struct {
		Title string `json:"title"`
	} `json:"data"`
}

func Home(w http.ResponseWriter, r *http.Request){
  client := &http.Client{Timeout: time.Second * 10,}
  viper.SetConfigFile("./config.toml"); viper.ReadInConfig(); clientid := viper.GetString("twitch.clientId"); chnl := viper.GetString("twitch.channel")
  req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_login="+chnl, nil)
  req.Header.Add("Client-ID", clientid); resp, _ := client.Do(req); defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body); i := StreamInfo{}; json.Unmarshal(body, &i)
  var state bool; var title string;
  if len(i.Data) != 0 {state = true} else {state = false;}
  Home := map[string]interface{}{"streamer": chnl, "title": title, "state": state}
  LoadPage(w, r, "./web/templates/home.html", Home)
}

func Commands(w http.ResponseWriter, r *http.Request){
  db := base.Conn()
  res, err := db.Query("SELECT level, commName, COALESCE(commDesc, '') as commDesc, COALESCE(commUse, '') as commUse, COALESCE(response, '') as response, points, cd FROM commands")
  if err != nil {
    panic(err.Error())
  }
  defer res.Close()
  
  var levels, commNames, commDescs, commUses, responses, pointss, cds []string
  for res.Next() {
    var level, commName, commDesc, commUse, response, points, cd string
    err = res.Scan(&level, &commName, &commDesc, &commUse, &response, &points, &cd)
    levels = append(levels, level)
    commNames = append(commNames, commName)
    commDescs = append(commDescs, commDesc)
    commUses = append(commUses, commUse)
    responses = append(responses, response)
    pointss = append(pointss, points)
    cds = append(cds, cd)
  }
  commands := map[string]interface{}{
    "level": levels,
    "commName": commNames,
    "commDesc": commDescs,
    "commUse": commUses,
    "response": responses,
    "points": pointss,
    "cd": cds,
  }
  db.Close()
  LoadPage(w, r, "./web/templates/commands.html", commands)
}

func Stats(w http.ResponseWriter, r *http.Request){
  db := base.Conn()
  var lines, userCount, songRequests, timeoutCount, banCount string
  
  l, _ := db.Query(`SELECT COUNT(*) FROM chatlogs`); for l.Next() {l.Scan(&lines)}
  u, _ := db.Query(`SELECT COUNT(*) FROM user`); for u.Next() {u.Scan(&userCount)}
  s, _ := db.Query(`SELECT COUNT(*) FROM songrequest`); for s.Next() {s.Scan(&songRequests)}
  t, _ := db.Query(`SELECT COUNT(*) FROM adminlogs where type = "timeout"`); for t.Next() {t.Scan(&timeoutCount)}
  b, _ := db.Query(`SELECT COUNT(*) FROM adminlogs where type = "ban"`); for b.Next() {b.Scan(&banCount)}
  
  var nlNames, totalLines, pNames, totalPoints, onlineHours, onNames, offlineHours, offNames []string
  var nlName, totalLine, pName, totalPoint, onlineHour, onName, offlineHour, offName string
  tl, _ := db.Query("SELECT name, num_lines FROM user ORDER BY num_lines DESC LIMIT 25")
  for tl.Next() {tl.Scan(&nlName, &totalLine); totalLines = append(totalLines, totalLine); nlNames = append(nlNames, nlName)}
  tp, _ := db.Query("SELECT name, points FROM user ORDER BY points DESC LIMIT 25")
  for tp.Next() {tp.Scan(&pName, &totalPoint); totalPoints = append(totalPoints, totalPoint); pNames = append(pNames, pName)}
  ton, _ := db.Query("SELECT name, timeOnline FROM user ORDER BY timeOnline DESC LIMIT 25")
  for ton.Next() {ton.Scan(&onName, &onlineHour); onlineHours = append(onlineHours, onlineHour); onNames = append(onNames, onName)}
  tof, _ := db.Query("SELECT name, timeOffline FROM user ORDER BY timeOffline DESC LIMIT 25")
  for tof.Next() {tof.Scan(&offName, &offlineHour); offlineHours = append(offlineHours, offlineHour); offNames = append(offNames, offName)}
  
  Stats := map[string]interface{}{
    "lines": lines,
    "userCount": userCount,
    "songRequests": songRequests,
    "timeoutCount": timeoutCount,
    "banCount": banCount,
    "totalLines": totalLines,
    "nlNames": nlNames,
    "totalPoints": totalPoints,
    "pNames": pNames,
    "onlineHours": onlineHours,
    "onNames": onNames,
    "offlineHours": offlineHours,
    "offNames": offNames,
  }
  db.Close()
  LoadPage(w, r, "./web/templates/stats.html", Stats)
}

func Bets(w http.ResponseWriter, r *http.Request) {
  // Display all bets on this page and who lost / won and how many points they betted and such
  var db = base.Conn()
  var betIds, betInfos, betOptionss, betWinners, betStates, betsIds, betsUsers, betsPointss, betsToWins, allBets []string; 
  var betId, betInfo, betOptions, betWinner, betState, betsId, betsUser, betsPoints, betsToWin, bet string; 
  rs, _ := db.Query(`SELECT betInfo FROM bet`); for rs.Next() {rs.Scan(&bet); allBets = append(allBets, bet);}
  res, _ := db.Query(`SELECT bet.betId, bet.betInfo, bet.betOptions, bet.betWinner, bet.betState, bets.betId, bets.betUser, bets.betPoints, bets.betToWin FROM bets INNER JOIN bet ON bets.betId = bet.betId`)
  for res.Next() {
    res.Scan(&betId, &betInfo, &betOptions, &betWinner, &betState, &betsId, &betsUser, &betsPoints, &betsToWin); 
    betIds = append(betIds, betId);
    betInfos = append(betInfos, betInfo);
    betOptionss = append(betOptionss, betOptions);
    betWinners = append(betWinners, betWinner);
    betStates = append(betStates, betState);
    betsIds = append(betsIds, betsId);
    betsUsers = append(betsUsers, betsUser);
    betsPointss = append(betsPointss, betsPoints);
    betsToWins = append(betsToWins, betsToWin);
  }
  bets := make(map[string][]map[string]string)
  for i := 0; i < len(allBets); i++ {
    bets[allBets[i]] = nil
    for j := 0; j < len(betIds); j++ {
      if allBets[i] == betInfos[j] {
        entry := map[string]string{"id": betsIds[j], "username": betsUsers[j], "points": betsPointss[j], "win": betsToWins[j]}
        bets[allBets[i]] = append(bets[allBets[i]], entry)
      }
    }
  }
  Bets := map[string]interface{}{"testBet": bets, "bets": allBets, "betId": betIds, "betInfo": betInfos, "betOptions": betOptionss, "betWinner": betWinners, "betState": betStates, "betsId": betsIds, "betsUser": betsUsers, "betsPoints": betsPointss, "betsToWin": betsToWins}
  LoadPage(w, r, "./web/templates/bets.html", Bets)
}

func Error(w http.ResponseWriter, r *http.Request) {
  LoadPage(w, r, "./web/templates/404.html", nil)
}

func LoadPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
  parseWatchTime := func (y string) string {
    x, _ := strconv.Atoi(y); days := x / 60 / 24; hours := (x / 60) % 24; minutes := x % 60
    return strconv.Itoa(days)+" days, "+strconv.Itoa(hours)+" hours, "+strconv.Itoa(minutes)+" minutes"
  }
  add := func (x, y int) int { return x + y }
  funcs := template.FuncMap{"add": add, "parseWatchTime": parseWatchTime}
  session, _ := gothic.Store.Get(r, "loginSession")
  if data != nil {
    session.Values["Info"] = data
  }
  t, err := template.New("").Funcs(funcs).ParseFiles(tmpl, "./web/templates/header.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.ExecuteTemplate(w, "base", session)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}