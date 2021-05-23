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
