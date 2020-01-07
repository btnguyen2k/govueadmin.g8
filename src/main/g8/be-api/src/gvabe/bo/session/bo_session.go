// package session contains Session business object (BO) and SQLite-DAO implementations.
package session

import "time"

const (
	fieldSessionId       = "id"
	fieldSessionData     = "data"
	fieldSessionExpiry   = "exp"
	fieldSessionParentId = "pid"
)

// Session is the business object
type Session struct {
	Id       string    `json:"id"`
	Data     string    `json:"data"`
	Expiry   time.Time `json:"exp"`
	ParentId string    `json:"pid"`
}

// SessionDao defines API to access Group storage
type SessionDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *Session) (bool, error)

	// Create persists a new business object to storage
	Create(bo *Session) (bool, error)

	// Get retrieves a business object from storage
	Get(id string) (*Session, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*Session, error)

	// GetAll retrieves all availables business objects from storage
	GetAll() ([]*Session, error)

	// Update modifies an existing business object
	Update(bo *Session) (bool, error)
}
