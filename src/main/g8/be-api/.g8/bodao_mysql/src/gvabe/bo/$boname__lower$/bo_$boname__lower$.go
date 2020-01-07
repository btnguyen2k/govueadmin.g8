// package $boname;format="lower"$ contains $boname;format="Camel"$ business object (BO) and Mysql-DAO implementations.
package $boname;format="lower"$

const (
	field$boname;format="Camel"$Id    = "id"
	field$boname;format="Camel"$Name  = "n"
	field$boname;format="Camel"$Value = "v"
)

// $boname;format="Camel"$ is the business object
type $boname;format="Camel"$ struct {
	Id    string `json:"id"`
	Name  string `json:"n"`
	Value int    `json:"v"`
}

// $boname;format="Camel"$Dao defines API to access $boname;format="Camel"$ storage
type $boname;format="Camel"$Dao interface {
	// Delete removes the specified business object from storage
	Delete(bo *$boname;format="Camel"$) (bool, error)

	// Create persists a new business object to storage
	Create(bo *$boname;format="Camel"$) (bool, error)

	// Get retrieves a business object from storage
	Get(id string) (*$boname;format="Camel"$, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*$boname;format="Camel"$, error)

	// GetAll retrieves all availables business objects from storage
	GetAll() ([]*$boname;format="Camel"$, error)

	// Update modifies an existing business object
	Update(bo *$boname;format="Camel"$) (bool, error)
}
