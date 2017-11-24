package base

import (

)

func (b *Bot) Modules(C string, U User) {
  b.UserInfoComms(C, U)
  b.Basic(C, U)
  b.Roulette(C, U)
  b.Slot(C, U)
  b.Pickpocket(C, U)
  b.Songrequest(C, U)
}

func (b *Bot) UserInfoComms(C string, U User) {
  if C == "!points" {
    var res = Query("SELECT points FROM user WHERE name = '"+U.username+"'")
    if res != "" {b.SendWhisper("You currently have " + res + " points", U.username)}
  }
  if C == "!userpoints" {
    var res = Query("SELECT points FROM user WHERE name = '"+U.username+"'")
    if res != "" {b.SendMsg(U.displayName + " currently has " + res + " points")}
  } 
  if C == "!lines" {
    var res = Query("SELECT num_lines FROM user WHERE name = '"+U.username+"'")
    if res != "" {b.SendWhisper("You have currently written " + res + " lines", U.username)}
  }
  if C == "!rq" {
    line := Query("SELECT log FROM chatlogs WHERE userId = '"+U.userId+"' ORDER BY RAND() LIMIT 1")
    if line != "" {b.SendMsg(U.displayName+": "+line)}
  }
}

func (b *Bot) Basic(C string, U User) {
  if C == "!test" {
    b.SendMsg("This is a test message created by " + U.displayName + " FeelsGoodMan")
  }
}