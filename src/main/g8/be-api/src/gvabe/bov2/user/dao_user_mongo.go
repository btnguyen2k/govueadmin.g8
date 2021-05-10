package user

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

// NewUserDaoMongo is helper method to create MongoDB-implementation of UserDao
//
// Available since template-v0.3.0
func NewUserDaoMongo(mc *prom.MongoConnect, collectionName string, txModeOnWrite bool) UserDao {
	dao := &BaseUserDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoMongo(mc, collectionName, txModeOnWrite)
	return dao
}

// // UserDaoMongo is MongoDB-implementation of UserDao
// //
// // Available since template-v0.3.0
// type UserDaoMongo struct {
// 	henge.UniversalDao
// }
//
// // GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
// func (dao *UserDaoMongo) GdaoCreateFilter(_ string, gbo godal.IGenericBo) godal.FilterOpt {
// 	return godal.MakeFilter(map[string]interface{}{henge.FieldId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)})
// }
//
// // Delete implements UserDao.Delete
// func (dao *UserDaoMongo) Delete(user *User) (bool, error) {
// 	return dao.UniversalDao.Delete(user.sync().UniversalBo)
// }
//
// // Create implements UserDao.Create
// func (dao *UserDaoMongo) Create(user *User) (bool, error) {
// 	return dao.UniversalDao.Create(user.sync().UniversalBo)
// }
//
// // Get implements UserDao.Get
// func (dao *UserDaoMongo) Get(id string) (*User, error) {
// 	ubo, err := dao.UniversalDao.Get(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return NewUserFromUbo(ubo), nil
// }
//
// // GetN implements UserDao.GetN
// func (dao *UserDaoMongo) GetN(fromOffset, maxNumRows int, filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*User, error) {
// 	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, filter, sorting)
// 	if err != nil {
// 		return nil, err
// 	}
// 	result := make([]*User, 0)
// 	for _, ubo := range uboList {
// 		app := NewUserFromUbo(ubo)
// 		result = append(result, app)
// 	}
// 	return result, nil
// }
//
// // GetAll implements UserDao.GetAll
// func (dao *UserDaoMongo) GetAll(filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*User, error) {
// 	return dao.GetN(0, 0, filter, sorting)
// }
//
// // Update implements UserDao.Update
// func (dao *UserDaoMongo) Update(user *User) (bool, error) {
// 	return dao.UniversalDao.Update(user.sync().UniversalBo)
// }
