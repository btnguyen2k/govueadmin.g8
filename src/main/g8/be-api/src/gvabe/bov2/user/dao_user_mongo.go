package user

import (
	"github.com/btnguyen2k/henge"
	prommongo "github.com/btnguyen2k/prom/mongo"
)

// NewUserDaoMongo is helper method to create MongoDB-implementation of UserDao
//
// Available since template-v0.3.0
func NewUserDaoMongo(mc *prommongo.MongoConnect, collectionName string, txModeOnWrite bool) UserDao {
	dao := &BaseUserDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoMongo(mc, collectionName, txModeOnWrite)
	return dao
}
