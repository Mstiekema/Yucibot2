package base

import (
  "fmt"
  "strings"
  "math/rand"
  "strconv"
)

func (b *Bot) Roulette(C string, U User) {
  if C == "!roulette" {
    msgSplit := strings.SplitAfter(U.message, " ")
    i := 1
    if  (i >= 1 && i < len(strings.SplitAfter(U.message, " "))) {
      StringRoulPoints := strings.TrimSpace(msgSplit[1])
      StringOldPoints := Query("user", "points", U.username)
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
              Update("user", "points = '"+snPoints+"'", "'"+U.username+"'")
              b.SendMsg(U.displayName + " won the roulette for " + StringRoulPoints + " points and now has " + snPoints + " points! PogChamp")
            } else {
              nPoints := oldPoints - roulPoints
              snPoints := strconv.Itoa(nPoints)
              Update("user", "points = '"+snPoints+"'", "'"+U.username+"'")
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
      Update("user", "points = points + 100", "'"+U.username+"'")
      B.SendMsg(U.displayName+", | "+emotes[set][a]+" | "+emotes[set][b]+" | "+emotes[set][c]+" | -> 3 in a row! You win 100 points PogChamp")
    } else if a == b || b == c {
      Update("user", "points = points + 50", "'"+U.username+"'")
      B.SendMsg(U.displayName+", | "+emotes[set][a]+" | "+emotes[set][b]+" | "+emotes[set][c]+" | -> Pretty close, you win 50 points SeemsGood")
    } else if a == c {
      Update("user", "points = points - 50", "'"+U.username+"'")
      B.SendMsg(U.displayName+", | "+emotes[set][a]+" | "+emotes[set][b]+" | "+emotes[set][c]+" | -> This isn't that good, you lose 50 points FeelsBadMan")
    } else {
      Update("user", "points = points - 100", "'"+U.username+"'")
      B.SendMsg(U.displayName+", | "+emotes[set][a]+" | "+emotes[set][b]+" | "+emotes[set][c]+" | -> Nothing is the same, what are you doing? You lose 100 points LUL")
    }
  }
}