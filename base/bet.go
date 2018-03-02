package base

import (
  "fmt"
  "strings"
  "strconv"
)

func (b *Bot) Bet(C string, U User) {
  if C == "!bet" {
    split := strings.SplitAfter(U.message, " ")
    if 2 < len(split) {
      betToWin := split[1]
      betPoints := split[2]; ibetPoints, _ := strconv.Atoi(betPoints)
      betId := Query("SELECT betId FROM bet WHERE betState = 0")
      if betId == "" {b.SendWhisper("There is currently no bet in progress.", U.username); return}
      checkBet := Query("SELECT betsId FROM bets where betId = '"+betId+"' AND betUser = '"+U.userId+"'")
      if checkBet != "" {b.SendWhisper("You've already bet on this bet.", U.username); return}
      oldPoints := Query("SELECT points FROM user WHERE name = '"+U.username+"'"); ioldPoints, _ := strconv.Atoi(oldPoints)
      if ibetPoints <= 0 {b.SendWhisper("You can do negative nor 0 point bets. Please try again using !bet CHOICE POINTS", U.username); return}
      if ioldPoints >= ibetPoints {
        Update("user", "points = points - "+betPoints, "name", "'"+U.username+"'")
        Insert("bets (betId, betUser, betPoints, betToWin)", "('"+betId+"', '"+U.userId+"', '"+betPoints+"', '"+betToWin+"')")
        b.SendWhisper("Succesfully added your bet to the database, good luck!", U.username)
      } else {
        b.SendWhisper("You do not have this many points", U.username)
      }
    } else {
      b.SendWhisper("You sent an invalid bet command, please try again.", U.username)
    }
  }
  if C == "!cbet" {
    exec := func() {
      a := Query("SELECT betInfo FROM bet WHERE betState = 0")
      c := Query("SELECT betId FROM bet WHERE betState = 0")
      d := Query("SELECT betOptions FROM bet WHERE betState = 0")
      if a == "" {b.SendWhisper("This isn't a valid bet ID, please try again", U.username); return}
      b.SendMsg("The current bet is #"+c+": "+a+" You can choose from: "+d)
      b.SendMsg("Bet on this bet by using '!bet CHOICE POINTS'! Good luck!")
    }
    b.ExecuteCommand(C, "100", "0", "10", U, exec)
  }
  if C == "!addbet" && (U.mod == "1" || U.username == strings.ToLower(b.Channel)) {
    betState := Query("SELECT betId FROM bet WHERE betState = 0")
    if betState != "" {b.SendWhisper("You have to end the current bet first before you can start another one.", U.username); return}
    s := strings.SplitAfter(U.message, "|")
    fmt.Println(len(s))
    if 1 >= len(s) {b.SendWhisper("You didn't specify the winners / didn't use the correct syntax. Please use '!addbet BET | CHOICES'", U.username); return;}
    betInfo := strings.TrimPrefix(s[0], "!addbet ")
    betOptions := s[1]
    Insert("bet (betInfo, betOptions)", "('"+betInfo+"', '"+betOptions+"')")
    c := Query("SELECT betId FROM bet WHERE betInfo = '"+betInfo+"'")
    b.SendMsg("Started a new bet PogChamp The bet is #"+c+": "+betInfo+" You can choose from: "+betOptions)
    b.SendMsg("Bet on this bet by using '!bet CHOICE POINTS'! Good luck!")
  }
  if C == "!endbet" && (U.mod == "1" || U.username == strings.ToLower(b.Channel)) {
    betId := Query("SELECT betId FROM bet WHERE betState = 0")
    if betId == "" {b.SendWhisper("There is currently no bet in progress.", U.username); return}
    split := strings.SplitAfter(U.message, " ")
    if 1 < len(split) {
      var db = Conn()
      rows, err := db.Query("select betUser, betPoints, betToWin from bets where betId = ?", betId)
      if err != nil {fmt.Println(err); return}
      defer rows.Close()
      var betUsers, betPointss, betToWins []string
      for rows.Next() {
        var betUser, betPoints, betToWin string
        err = rows.Scan(&betUser, &betPoints, &betToWin)
        betUsers = append(betUsers, betUser)
        betPointss = append(betPointss, betPoints)
        betToWins = append(betToWins, betToWin)
      }
      Update("bet", "betWinner = '"+strings.ToLower(strings.TrimSpace(split[1]))+"'", "'betState'", "0")
      Update("bet", "betState = 1", "'betState'", "0")
      for i := 0; i < len(betUsers); i++ {
        if strings.ToLower(strings.TrimSpace(betToWins[i])) == strings.ToLower(strings.TrimSpace(split[1])) {
          wPoints, _ := strconv.Atoi(betPointss[i])
          Update("user", "points = points + "+strconv.Itoa(wPoints*2), "userId", "'"+betUsers[i]+"'")
        }
      }
      b.SendMsg("'"+split[1]+"' is the winning bet! Users who guessed correctly will now be rewarded PogChamp")
      db.Close()
    } else {
      b.SendWhisper("You sent an invalid end bet command, please try again.", U.username)
    }
  }
}