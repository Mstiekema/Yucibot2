package base

import (
  "fmt"
  "strings"
  "math/rand"
  "strconv"
)

func (b *Bot) Roulette(C, usr, M string) {
  if C == "!roulette" {
    msgSplit := strings.SplitAfter(M, " ")
    i := 1
    if  (i >= 1 && i < len(strings.SplitAfter(M, " "))) {
      StringRoulPoints := strings.TrimSpace(msgSplit[1])
      StringOldPoints := Query("user", "points", strings.ToLower(usr))
      roulPoints, _ := strconv.Atoi(StringRoulPoints)
      oldPoints, _ := strconv.Atoi(StringOldPoints)
      if StringOldPoints != "" {
        if roulPoints > 0 {
          if roulPoints <= oldPoints {
            ran := rand.Float64()
            fmt.Println(ran)
            if ran > 0.5 {
              nPoints := oldPoints + roulPoints
              snPoints := strconv.Itoa(nPoints)
              Update("user", "points = '"+snPoints+"'", "'"+usr+"'")
              b.SendMsg(usr + " won the roulette for " + StringRoulPoints + " points and now has " + snPoints + " points! PogChamp")
            } else {
              nPoints := oldPoints - roulPoints
              snPoints := strconv.Itoa(nPoints)
              Update("user", "points = '"+snPoints+"'", "'"+usr+"'")
              b.SendMsg(usr + " lost the roulette for " + StringRoulPoints + " points and now has " + snPoints + " points! FeelsBadMan")
            }
          } else {
            b.SendMsg(usr + ", you do not have enough points to do a roulette DansGame")
          }
        } else {
          b.SendMsg(usr + ", you can't do negative roulettes")
        }
      } else {
        b.SendMsg(usr + " is not registered in the database, so is not able to do a roulette.")
      }
    } else {
      b.SendMsg("Invalid roulette command")
    }
  }
}