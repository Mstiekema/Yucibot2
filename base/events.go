package base

import (
  "strconv"
)

func (b *Bot) Subs(noticeType, username, displayName, months, plan string) {
  if plan == "1000" { plan = "$4,99" } else if plan == "2000" { plan = "$9,99" } else if plan == "3000" { plan = "$24,99" }
  if noticeType == "sub" {
    b.SendMsg("Thanks "+displayName+" for subscribing using a "+plan+" sub! PogChamp <3")
  }
  if noticeType == "resub" {
    b.SendMsg("Thanks "+displayName+" for subscribing "+months+" months in a row using a "+plan+" sub! PogChamp <3")
  }
  nMonths, _ := strconv.Atoi(months)
  b.SendWhisper("Thanks for the sub! You've received "+strconv.Itoa(500*nMonths)+" points as a reward!", username)
  Update("user", "points = points + "+strconv.Itoa(500*nMonths), "name", "'"+username+"'")
  
  b.StartRaffle(float64(30), (1111*nMonths), true)
}