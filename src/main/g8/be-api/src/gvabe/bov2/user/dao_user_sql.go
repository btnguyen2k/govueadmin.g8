package user

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/prom"

	"main/src/henge"
)

const TableUser = "gva_user"
const (
	UserCol_MaskUid = "zmid"
)

// NewUserDaoSql is helper method to create SQL-implementation of UserDao
//
// available since template-v0.2.0
func NewUserDaoSql(sqlc *prom.SqlConnect, tableName string) UserDao {
	dao := &UserDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(
		sqlc, tableName,
		map[string]string{UserCol_MaskUid: UserField_MaskId})
	return dao
}

// UserDaoSql is SQL-implementation of UserDao
//
// available since template-v0.2.0
type UserDaoSql struct {
	henge.UniversalDao
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *UserDaoSql) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
	return map[string]interface{}{henge.ColId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
}

// Delete implements UserDao.Delete
func (dao *UserDaoSql) Delete(user *User) (bool, error) {
	return dao.UniversalDao.Delete(user.UniversalBo.Clone())
}

// Create implements UserDao.Create
func (dao *UserDaoSql) Create(user *User) (bool, error) {
	return dao.UniversalDao.Create(user.sync().UniversalBo.Clone())
}

// Get implements UserDao.Get
func (dao *UserDaoSql) Get(id string) (*User, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewUserFromUbo(ubo), nil
}

// GetN implements UserDao.GetN
func (dao *UserDaoSql) GetN(fromOffset, maxNumRows int) ([]*User, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows)
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
func (dao *UserDaoSql) GetAll() ([]*User, error) {
	return dao.GetN(0, 0)
}

// Update implements UserDao.Update
func (dao *UserDaoSql) Update(user *User) (bool, error) {
	return dao.UniversalDao.Update(user.sync().UniversalBo.Clone())
}
