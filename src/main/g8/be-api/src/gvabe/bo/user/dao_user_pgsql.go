package user

import (
	"fmt"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
	_ "github.com/go-sql-driver/mysql"
)

// NewUserDaoPgsql is helper method to create PostgreSQL-implementation of UserDao
func NewUserDaoPgsql(sqlc *prom.SqlConnect, tableName string) UserDao {
	dao := &UserDaoPgsql{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: mapPgsqlFieldToColNameUser},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: mapPgsqlColNameToFieldUser},
		ColumnsListMap:              map[string][]string{tableName: colsPgsqlUser},
	})
	dao.SetSqlFlavor(prom.FlavorPgSql)
	return dao
}

func InitPgsqlTableUser(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s JSONB, PRIMARY KEY (%s))",
		tableName, colPgsqlUserUsername, colPgsqlUserData, colPgsqlUserUsername)
	if _, err := sqlc.GetDB().Exec(sql); err != nil {
		panic(err)
	}
}

const (
	colPgsqlUserUsername = "uname"
	colPgsqlUserData     = "udata"
)

var (
	colsPgsqlUser              = []string{colPgsqlUserUsername, colPgsqlUserData}
	mapPgsqlFieldToColNameUser = map[string]interface{}{
		fieldUserUsername: colPgsqlUserUsername,
		fieldUserData:     colPgsqlUserData,
	}
	mapPgsqlColNameToFieldUser = map[string]interface{}{
		colPgsqlUserUsername: fieldUserUsername,
		colPgsqlUserData:     fieldUserData,
	}
)

type UserDaoPgsql struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *UserDaoPgsql) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{colPgsqlUserUsername: bo.GboGetAttrUnsafe(fieldUserUsername, reddo.TypeString)}
}

// ‭toBo transforms godal.IGenericBo to business object.
func (dao *UserDaoPgsql) toBo(gbo godal.IGenericBo) *User {
	if gbo == nil {
		return nil
	}
	username := gbo.GboGetAttrUnsafe(fieldUserUsername, reddo.TypeString).(string)
	data := gbo.GboGetAttrUnsafe(fieldUserData, reddo.TypeString).(string)
	return NewUserBo(username, data)
}

// toGbo transforms business object to godal.IGenericBo
func (dao *UserDaoPgsql) toGbo(bo *User) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldUserUsername, bo.GetUsername())
	gbo.GboSetAttr(fieldUserData, bo.GetData())
	return gbo
}

// Delete implements UserDao.Delete
func (dao *UserDaoPgsql) Delete(bo *User) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements UserDao.Create
func (dao *UserDaoPgsql) Create(bo *User) (bool, error) {
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements UserDao.Get
func (dao *UserDaoPgsql) Get(username string) (*User, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{colPgsqlUserUsername: username})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements UserDao.GetN
func (dao *UserDaoPgsql) GetN(fromOffset, maxNumRows int) ([]*User, error) {
	// order ascending by "id" column
	ordering := (&sql.GenericSorting{Flavor: dao.GetSqlFlavor()}).Add(colPgsqlUserUsername)
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, ordering, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*User, 0)
	for _, gbo := range gboList {
		bo := dao.toBo(gbo)
		result = append(result, bo)
	}
	return result, nil
}

// GetAll implements UserDao.GetAll
func (dao *UserDaoPgsql) GetAll() ([]*User, error) {
	return dao.GetN(0, 0)
}

// Update implements UserDao.Update
func (dao *UserDaoPgsql) Update(bo *User) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}
