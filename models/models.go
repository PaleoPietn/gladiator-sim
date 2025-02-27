package model

// Player represents a gladiator with stats and abilities
type Player struct {
	Name         string
	Health       int
	MaxHealth    int
	AttackMin    int
	AttackMax    int
	Defense      int
	IsHero       bool
	Wins         int
	CritChance   int
	BlockChance  int
	LifeSteal    int
	CritDamage   int
	Regeneration int
	Description  string
}

// BattleResult contains the outcome of an attack
type BattleResult struct {
	Attacker     *Player
	Defender     *Player
	Damage       int
	IsCritical   bool
	IsBlocked    bool
	IsGameOver   bool
	WinnerName   string
	Regeneration int
}

// GameState tracks the overall game progression
type GameState struct {
	CurrentEnemy    int
	UpgradeMode     bool
	Upgrades        []Upgrade
	SelectedUpgrade int
	BattleLog       []string
	GameOver        bool
}

// Upgrade represents a possible improvement for the hero
type Upgrade struct {
	Name        string
	Description string
	Effect      func(*Player)
}

// Constants for game configuration
const (
	CriticalChance = 10 // 1 in 10 chance (10%)
	BlockChance    = 10 // 1 in 10 chance (10%)
	MaxLogEntries  = 20
)

// AddToBattleLog adds a message to the battle log
func (gs *GameState) AddToBattleLog(message string) {
	gs.BattleLog = append(gs.BattleLog, message)
}
