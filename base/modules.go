package base

func (b *Bot) Modules(C string, U User) {
  b.UserInfoComms(C, U)
  b.Basic(C, U)
  b.Raffle(C, U)
  b.Roulette(C, U)
  b.Slot(C, U)
  b.Pickpocket(C, U)
  b.Songrequest(C, U)
  b.CustomCommands(C, U)
  b.ModifyCommands(C, U)
  go b.Clr(C, U)
}