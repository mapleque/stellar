package stellar

// Engine ...
type Engine interface {
	// AddCharactor ...
	AddCharactor(opts ...CharactorOption) (Charactor, error)
	// GetCharactor ...
	GetCharactor(id string) (Charactor, bool)
	// RemoveCharactor ...
	RemoveCharactor(id string) bool

	// DoAction ...
	DoAction(actions ...Action) error

	// GetEvents ...
	GetEvents(id int) (events []Event, lastEventId int)

	// Start start the game
	Start() error

	// Stop stop the game
	Stop() error
}

var _ Engine = &engine{}

// EngineOption ...
type EngineOption func(*engine) error

// Charactor ...
type Charactor interface {
	// GetName ...
	GetName() string
	// GetMaxHP ...
	GetMaxHP() int64
	// GetCurHP ...
	GetCurHP() int64
	// IsDead ...
	IsDead() bool
}

var _ Charactor = &charactor{}

// CharactorOption ...
type CharactorOption func(*charactor) error

// Action ...
type Action interface {
	Type() ActionType
	Source() Charactor
	Targets() []Charactor
}

type ActionType string

const (
	ActionTypeAttack ActionType = "Attack"
)

// Event ...
type Event interface {
	Type() EventType
	Message() string
	Charactors() []Charactor
}

type EventType string

const (
	EventTypeAddCharactor EventType = "add charactor"
	EventTypeDoAttack     EventType = "do attack"
	EventTypeBeAttack     EventType = "be attack"
	EventTypeDead         EventType = "dead"
)
