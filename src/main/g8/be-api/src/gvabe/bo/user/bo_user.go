// package user contains User business object (BO) and Mysql-DAO implementations.
package user

const (
	fieldUserUsername = "uname"
	fieldUserPassword = "pwd"
	fieldUserName     = "name"
	fieldUserGroupId  = "gid"
)

// User is the business object
type User struct {
	Username string `json:"uname"`
	Password string `json:"pwd"`
	Name     string `json:"name"`
	GroupId  string `json:"gid"`
}

// UserDao defines API to access User storage
type UserDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *User) (bool, error)

	// Create persists a new business object to storage
	Create(bo *User) (bool, error)

	// Get retrieves a business object from storage
	Get(username string) (*User, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*User, error)

	// GetAll retrieves all availables business objects from storage
	GetAll() ([]*User, error)

	// Update modifies an existing business object
	Update(bo *User) (bool, error)
}
