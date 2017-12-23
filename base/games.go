package base

import (
  "strings"
  "math/rand"
  "strconv"
  "time"
)

var participants []User
var rafState = false

func (b *Bot) Raffle(C string, U User) {
  if C == "!raffle" {
    var dur int
    var points int
    if 2 < len(strings.SplitAfter(U.message, " ")) {
      points, _ = strconv.Atoi(strings.TrimSpace(strings.SplitAfter(U.message, " ")[1]))
      dur, _ = strconv.Atoi(strings.TrimSpace(strings.SplitAfter(U.message, " ")[2]))
    } else if 1 < len(strings.SplitAfter(U.message, " ")) {
      points, _ = strconv.Atoi(strings.TrimSpace(strings.SplitAfter(U.message, " ")[1]))
      dur = 30
    } else {
      points = 3000
      dur = 30
    }
    b.StartRaffle(float64(dur), points, false, U)
  }
  if C == "!multiraffle" {
    var dur int
    var points int
    
    for i := 0; i < 20; i++ {
      U.username = strconv.Itoa(i)
      participants = append(participants, U)
    }
    
    if 2 < len(strings.SplitAfter(U.message, " ")) {
      points, _ = strconv.Atoi(strings.TrimSpace(strings.SplitAfter(U.message, " ")[1]))
      dur, _ = strconv.Atoi(strings.TrimSpace(strings.SplitAfter(U.message, " ")[2]))
    } else if 1 < len(strings.SplitAfter(U.message, " ")) {
      points, _ = strconv.Atoi(strings.TrimSpace(strings.SplitAfter(U.message, " ")[1]))
      dur = 30
    } else {
      points = 3000
      dur = 30
    }
    b.StartRaffle(float64(dur), points, true, U)
  }
  if C == "!join" && rafState == true {
    for _, usr := range participants {
      if usr == U {
        return
      }
    }
    participants = append(participants, U)
  }
}

func (b *Bot) Roulette(C string, U User) {
  if C == "!roulette" {
    msgSplit := strings.SplitAfter(U.message, " ")
    i := 1
    if  (i >= 1 && i < len(strings.SplitAfter(U.message, " "))) {
      StringRoulPoints := strings.TrimSpace(msgSplit[1])
      StringOldPoints := Query("SELECT points FROM user WHERE name = '"+U.username+"'")
      roulPoints, _ := strconv.Atoi(StringRoulPoints)
      oldPoints, _ := strconv.Atoi(StringOldPoints)
      if StringOldPoints != "" {
        if roulPoints > 0 {
          if roulPoints <= oldPoints {
            ran := rand.Float64()
            if ran > 0.5 {
              nPoints := oldPoints + roulPoints
              snPoints := strconv.Itoa(nPoints)
              Update("user", "points = '"+snPoints+"'", "name", "'"+U.username+"'")
              b.SendMsg(U.displayName + " won the roulette for " + StringRoulPoints + " points and now has " + snPoints + " points! PogChamp")
            } else {
              nPoints := oldPoints - roulPoints
              snPoints := strconv.Itoa(nPoints)
              Update("user", "points = '"+snPoints+"'", "name", "'"+U.username+"'")
              b.SendMsg(U.displayName + " lost the roulette for " + StringRoulPoints + " points and now has " + snPoints + " points! FeelsBadMan")
            }
          } else {
            b.SendMsg(U.displayName + ", you do not have enough points to do a roulette DansGame")
          }
        } else {
          b.SendMsg(U.displayName + ", you can't do negative roulettes")
        }
      } else {
        b.SendMsg(U.displayName + " is not registered in the database, so is not able to do a roulette.")
      }
    } else {
      b.SendMsg("Invalid roulette command")
    }
  }
}

func (B *Bot) Slot(C string, U User) {
  if C == "!slot" {
    emotes := [3][3]string{{"Kappa", "Keepo", "PogChamp"},{"SeemsGood", "DansGame", "4Head"},{"DatSheffy", "LUL", "cmonBruh"}}
    a := int(rand.Float64()*3)
    b := int(rand.Float64()*3)
    c := int(rand.Float64()*3)
    set := int(rand.Float64()*3)
    
    if a == b && b == c {
      Update("user", "points = points + 100", "name", "'"+U.username+"'")
      B.SendMsg(U.displayName+", | "+emotes[set][a]+" | "+emotes[set][b]+" | "+emotes[set][c]+" | -> 3 in a row! You win 100 points PogChamp")
    } else if a == b || b == c {
      Update("user", "points = points + 50", "name", "'"+U.username+"'")
      B.SendMsg(U.displayName+", | "+emotes[set][a]+" | "+emotes[set][b]+" | "+emotes[set][c]+" | -> Pretty close, you win 50 points SeemsGood")
    } else if a == c {
      Update("user", "points = points - 50", "name", "'"+U.username+"'")
      B.SendMsg(U.displayName+", | "+emotes[set][a]+" | "+emotes[set][b]+" | "+emotes[set][c]+" | -> This isn't that good, you lose 50 points FeelsBadMan")
    } else {
      Update("user", "points = points - 100", "name", "'"+U.username+"'")
      B.SendMsg(U.displayName+", | "+emotes[set][a]+" | "+emotes[set][b]+" | "+emotes[set][c]+" | -> Nothing is the same, what are you doing? You lose 100 points LUL")
    }
  }
}

func (b *Bot) Pickpocket(C string, U User) {
  if C == "!stoppp" {
    Update("user", "pickP = 0", "name", "'"+U.username+"'")
    b.SendMsg("You can no longer steal points from " + U.displayName)
  } else if C == "!resumepp" {
    Update("user", "pickP = 1", "name", "'"+U.username+"'")
    b.SendMsg("You can now start stealing points from " + U.displayName)
  } else if C == "!pickpocket" || C == "!pp" {
    msgSplit := strings.SplitAfter(U.message, " ")
    i := 1
    if  (i >= 1 && i < len(strings.SplitAfter(U.message, " "))) {
      target := strings.ToLower(msgSplit[1])
      uPoints := Query("SELECT points FROM user WHERE name = '"+U.username+"'")
      tPoints := Query("SELECT points FROM user WHERE name = '"+target+"'")
      fuPoints, _ := strconv.ParseFloat(uPoints, 64)
      itPoints, _ := strconv.Atoi(tPoints)
      if uPoints != "" && tPoints != "" {
        stealP := int(rand.Float64()*100)
        x := int(rand.Float64()*100)
        if U.sub == "1" {
          if x > 20 {
            if stealP > itPoints {
              Update("user", "points = points + "+tPoints, "name", "'"+U.username+"'")
              Update("user", "points = points - "+tPoints, "name", "'"+target+"'")
              b.SendMsg(U.displayName + " stole all of "+msgSplit[1]+"'s "+tPoints+" points TriHard")
            } else {
              Update("user", "points = points + "+strconv.Itoa(stealP), "name", "'"+U.username+"'")
              Update("user", "points = points - "+strconv.Itoa(stealP), "name", "'"+target+"'")
              b.SendMsg(U.displayName + " stole "+strconv.Itoa(stealP)+" points from "+msgSplit[1]+" TriHard")
            }
          } else if x < 5 {
            lPoints := strconv.Itoa(int(fuPoints * 0.1))
            Update("user", "points = points - "+lPoints, "name", "'"+U.username+"'")
            b.SendMsg(U.displayName + " got caught trying to steal points from "+msgSplit[1]+" and loses "+lPoints+" points")
          } else {
            b.SendMsg(U.displayName + " failed to steal points from "+msgSplit[1])
          }
        } else {
          if x > 40 {
            if stealP > itPoints {
              Update("user", "points = points + "+tPoints, "name", "'"+U.username+"'")
              Update("user", "points = points - "+tPoints, "name", "'"+target+"'")
              b.SendMsg(U.displayName + " stole all of "+msgSplit[1]+"'s "+tPoints+" points TriHard")
            } else {
              Update("user", "points = points + "+strconv.Itoa(stealP), "name", "'"+U.username+"'")
              Update("user", "points = points - "+strconv.Itoa(stealP), "name", "'"+target+"'")
              b.SendMsg(U.displayName + " stole "+strconv.Itoa(stealP)+" points from "+msgSplit[1]+" TriHard")
            }
          } else if x < 10 {
            lPoints := strconv.Itoa(int(fuPoints * 0.1))
            Update("user", "points = points - "+lPoints, "name", "'"+U.username+"'")
            b.SendMsg(U.displayName + " got caught trying to steal points from "+msgSplit[1]+" and loses "+lPoints+" points")
          } else {
            b.SendMsg(U.displayName + " failed to steal points from "+msgSplit[1])
          }
        }
      } else {
        b.SendMsg(U.displayName + ", something went wrong while performing the pickpocket command")
      }
    } else {
      b.SendMsg("Invalid pickpocket command")
    }
  }
}

func (b *Bot) StartRaffle(dur float64, points int, multi bool, U User) {
  var m string
  rafState = true
  if multi == true { m = "multi" } else { m = "" }
  
  b.SendMsg("Started "+m+"raffle for "+strconv.Itoa(points)+" points! Type !join to join the "+m+"raffle PogChamp")
  time.AfterFunc(time.Duration(int(dur*0.25)) * time.Second, func() {b.SendMsg("Hurry up! You still have "+strconv.Itoa(int(dur*0.75))+" seconds left to join the "+m+"raffle for "+strconv.Itoa(points)+" points!")})
  time.AfterFunc(time.Duration(int(dur*0.5)) * time.Second, func() {b.SendMsg("Hurry up! You still have "+strconv.Itoa(int(dur*0.5))+" seconds left to join the "+m+"raffle for "+strconv.Itoa(points)+" points!")})
  time.AfterFunc(time.Duration(int(dur*0.75)) * time.Second, func() {b.SendMsg("Hurry up! You still have "+strconv.Itoa(int(dur*0.25))+" seconds left to join the "+m+"raffle for "+strconv.Itoa(points)+" points!")})
  
  time.AfterFunc(time.Duration(dur) * time.Second, func() {
    rafState = false
    if multi == true {
      winners := ""
      winP := int((rand.Float64() * float64(len(participants)) / 100) * ((rand.Float64() * 25) + 5)) + 1
      wPoints := points / winP
      for i := 0; i < winP; i++ {
        j := int(rand.Float64() * float64(len(participants)))
        winners = winners + " " + participants[j].displayName
        Update("user", "points = points + '"+strconv.Itoa(wPoints)+"'", "name", "'"+participants[j].username+"'")
        for i, v := range participants {
          if v == participants[j] {
            participants = append(participants[:i], participants[i+1:]...)
          }
        }
      }
      b.SendMsg("The raffle has finished! "+strconv.Itoa(winP)+" users have won "+strconv.Itoa(wPoints)+" points each! The winners are:"+winners+" PogChamp")
    } else {
      if len(participants) == 0 {b.SendMsg("No one joined the raffle DansGame"); return}
      ran := int(rand.Float64() * float64(len(participants)))
      if ran == len(participants) {ran = ran-1}
      winner := participants[ran]
      Update("user", "points = points + '"+strconv.Itoa(points)+"'", "name", "'"+winner.username+"'")
      b.SendMsg("The raffle has finished! "+winner.displayName+" has won "+strconv.Itoa(points)+" points PogChamp")
    }
    participants = participants[:0]
  })
}