package user

import (
	"github.com/btnguyen2k/henge"
	promdynamodb "github.com/btnguyen2k/prom/dynamodb"
)

// NewUserDaoDynamodb is helper method to create AWS DynamoDB-implementation of UserDao
//
// Available since template-v0.3.0
func NewUserDaoDynamodb(adc *promdynamodb.AwsDynamodbConnect, tableName string) UserDao {
	dao := &BaseUserDaoImpl{}
	spec := &henge.DynamodbDaoSpec{UidxAttrs: [][]string{{UserFieldMaskId}}}
	dao.UniversalDao = henge.NewUniversalDaoDynamodb(adc, tableName, spec)
	return dao
}
