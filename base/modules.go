package base

func (b *Bot) Modules(C string, U User) {
  b.UserInfoComms(C, U)
  b.Basic(C, U)
  b.Roulette(C, U)
  b.Slot(C, U)
  b.Pickpocket(C, U)
  b.Songrequest(C, U)
  go b.Clr(C, U)
}