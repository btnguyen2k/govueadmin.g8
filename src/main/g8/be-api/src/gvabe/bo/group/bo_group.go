// package group contains Group business object (BO) and SQLite-DAO implementations.
package group

const (
	fieldGroupId    = "id"
	fieldGroupName  = "n"
)

// Group is the business object
type Group struct {
	Id    string `json:"id"`
	Name  string `json:"n"`
}

// GroupDao defines API to access Group storage
type GroupDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *Group) (bool, error)

	// Create persists a new business object to storage
	Create(bo *Group) (bool, error)

	// Get retrieves a business object from storage
	Get(id string) (*Group, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*Group, error)

	// GetAll retrieves all availables business objects from storage
	GetAll() ([]*Group, error)

	// Update modifies an existing business object
	Update(bo *Group) (bool, error)
}
