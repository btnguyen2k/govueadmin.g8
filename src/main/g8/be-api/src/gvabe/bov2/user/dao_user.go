package user

import (
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/henge"
)

const TableUser = "gva_user"
const (
	// UserColMaskUid is name of database column for user's mask-id.
	UserColMaskUid = "zmid"
)

// UserDao defines API to access User storage
//
// Available since template-v0.2.0
type UserDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *User) (bool, error)

	// Create persists a new business object to storage
	Create(bo *User) (bool, error)

	// Get retrieves a business object from storage
	Get(username string) (*User, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int, filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*User, error)

	// GetAll retrieves all available business objects from storage
	GetAll(filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*User, error)

	// Update modifies an existing business object
	Update(bo *User) (bool, error)
}

// BaseUserDaoImpl is a generic implementation of UserDao.
//
// Available since template-v0.3.0
type BaseUserDaoImpl struct {
	henge.UniversalDao
}

// Delete implements UserDao.Delete
func (dao *BaseUserDaoImpl) Delete(user *User) (bool, error) {
	return dao.UniversalDao.Delete(user.sync().UniversalBo)
}

// Create implements UserDao.Create
func (dao *BaseUserDaoImpl) Create(user *User) (bool, error) {
	return dao.UniversalDao.Create(user.sync().UniversalBo)
}

// Get implements UserDao.Get
func (dao *BaseUserDaoImpl) Get(id string) (*User, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewUserFromUbo(ubo), nil
}

// GetN implements UserDao.GetN
func (dao *BaseUserDaoImpl) GetN(fromOffset, maxNumRows int, filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*User, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, filter, sorting)
	if err != nil {
		return nil, err
	}
	result := make([]*User, 0)
	for _, ubo := range uboList {
		app := NewUserFromUbo(ubo)
		result = append(result, app)
	}
	return result, nil
}

// GetAll implements UserDao.GetAll
func (dao *BaseUserDaoImpl) GetAll(filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*User, error) {
	return dao.GetN(0, 0, filter, sorting)
}

// Update implements UserDao.Update
func (dao *BaseUserDaoImpl) Update(user *User) (bool, error) {
	return dao.UniversalDao.Update(user.sync().UniversalBo)
}
