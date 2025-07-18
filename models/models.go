package models

import "time"

const (
	// Constants for game configuration
	CriticalChance = 10 // 1 in 10 chance (10%)
	BlockChance    = 10 // 1 in 10 chance (10%)

	// BattleLogs
	CriticalHit = "CRITICAL HIT"
	Blocked     = "BLOCKED"
	Victorious  = "VICTORIOUS"

	// Use to specify to UI what LogType a battleLog should be displayed with
	LogTypeTitle       LogType = "title"
	LogTypeHero        LogType = "hero"
	LogTypeEnemy       LogType = "enemy"
	LogTypeInfo        LogType = "info"
	LogTypeDefault     LogType = "default"
	LogTypeSelected    LogType = "selected"
	LogTypeCritical    LogType = "critical"
	LogTypeBlock       LogType = "block"
	LogWaitTimeDefault         = time.Millisecond * 800
	LogWaitTimeShort           = time.Millisecond * 500
	LogWaitTimeLong            = time.Millisecond * 1000
)

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
	BlockChance  int // Maybe we could also have DodgeChance? 🤔
	LifeSteal    int
	CritDamage   int
	Regeneration int
	LifeOnKill   int
	Description  string
	Upgrades     []Upgrade
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
	CurrentEnemy int
	Upgrades     []Upgrade
	GameOver     bool
}

// BattleState contain a snapshot of the state of the battle at a specific turn and the associated logs to display
type BattleState struct {
	Hero      *Player
	Enemy     *Player
	BattleLog *BattleLog
}

// BattleLogs are used to show the outcome of each turn
type BattleLog struct {
	LogMessage  string        // The actual message
	LogType     LogType       // How the log should be displayed by UI
	LogWaitTime time.Duration // Recommended time for UI to wait until display the next log (in milliseconds)
}

type LogType string

// Upgrade represents a possible improvement for the hero
type Upgrade struct {
	Name        string
	Description string
	Effect      func(*Player)
}
