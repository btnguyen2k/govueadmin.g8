package $boname;format="lower"$

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/mongo"
	"github.com/btnguyen2k/prom"
	"main/src/goapi"
)

// newMongoConnect is helper method to create connection pool to MongoDB server
func newMongoConnect() *prom.MongoConnect {
	// read from config or use default value
	url := goapi.AppConfig.GetString("dao$boname__lower$.mongodb.url", "mongodb://test:test@localhost:27017/test")
	db := goapi.AppConfig.GetString("dao$boname__lower$.mongodb.db", "test")
	timeoutMs := goapi.AppConfig.GetInt32("dao$boname__lower$.mongodb.timeout", 10000)
	if mongoConnect, err := prom.NewMongoConnect(url, db, int(timeoutMs)); err != nil {
		panic(err)
	} else if mongoConnect == nil {
		panic("error creating [prom.MongoConnect] instance")
	} else {
		return mongoConnect
	}
}

// new$boname;format="Camel"$DaoMMongo is helper method to create MongoDB-implementation of $boname;format="Camel"$Dao
func new$boname;format="Camel"$DaoMongo(mc *prom.MongoConnect, collectionName string) $boname;format="Camel"$Dao {
	dao := &$boname;format="Camel"$DaoMongo{collectionName: collectionName}
	dao.GenericDaoMongo = mongo.NewGenericDaoMongo(mc, godal.NewAbstractGenericDao(dao))
	return dao
}

const (
	collection$boname;format="Camel"$ = "$boname__lower$"
)

type $boname;format="Camel"$DaoMongo struct {
	*mongo.GenericDaoMongo
	collectionName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *$boname;format="Camel"$DaoMongo) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	id, _ := bo.GboGetAttr(field$boname;format="Camel"$Id, reddo.TypeString)
	return map[string]interface{}{field$boname;format="Camel"$Id: id}
}

// ‭toBo transforms godal.IGenericBo to business object.
func (dao *$boname;format="Camel"$DaoMongo) toBo(gbo godal.IGenericBo) *$boname;format="Camel"$ {
	if gbo == nil {
		return nil
	}
	bo := &$boname;format="Camel"${
		Id:   gbo.GboGetAttrUnsafe(field$boname;format="Camel"$Id, reddo.TypeString).(string),
		Name: gbo.GboGetAttrUnsafe(field$boname;format="Camel"$Name, reddo.TypeString).(string),
		Value: int(gbo.GboGetAttrUnsafe(field$boname;format="Camel"$Value, reddo.TypeInt).(int64)),
	}
	return bo
}

// toGbo transforms business object to godal.IGenericBo
func (dao *$boname;format="Camel"$DaoMongo) toGbo(bo *$boname;format="Camel"$) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(field$boname;format="Camel"$Id, bo.Id)
	gbo.GboSetAttr(field$boname;format="Camel"$Name, bo.Name)
	gbo.GboSetAttr(field$boname;format="Camel"$Value, bo.Value)
	return gbo
}

// Delete implements $boname;format="Camel"$Dao.Delete
func (dao *$boname;format="Camel"$DaoMongo) Delete(bo *$boname;format="Camel"$) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.collectionName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements $boname;format="Camel"$Dao.Create
func (dao *$boname;format="Camel"$DaoMongo) Create(bo *$boname;format="Camel"$) (bool, error) {
	numRows, err := dao.GdaoCreate(dao.collectionName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements $boname;format="Camel"$Dao.Get
func (dao *$boname;format="Camel"$DaoMongo) Get(id string) (*$boname;format="Camel"$, error) {
	gbo, err := dao.GdaoFetchOne(dao.collectionName, map[string]interface{}{field$boname;format="Camel"$Id: id})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements $boname;format="Camel"$Dao.GetN
func (dao *$boname;format="Camel"$DaoMongo) GetN(fromOffset, maxNumRows int) ([]*$boname;format="Camel"$, error) {
	// order ascending by "id" column
	ordering := map[string]int{field$boname;format="Camel"$Id: 1}
	gboList, err := dao.GdaoFetchMany(dao.collectionName, nil, ordering, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*$boname;format="Camel"$, 0)
	for _, gbo := range gboList {
		bo := dao.toBo(gbo)
		result = append(result, bo)
	}
	return result, nil
}

// GetAll implements $boname;format="Camel"$Dao.GetAll
func (dao *$boname;format="Camel"$DaoMongo) GetAll() ([]*$boname;format="Camel"$, error) {
	return dao.GetN(0, 0)
}

// Update implements $boname;format="Camel"$Dao.Update
func (dao *$boname;format="Camel"$DaoMongo) Update(bo *$boname;format="Camel"$) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.collectionName, dao.toGbo(bo))
	return numRows > 0, err
}
