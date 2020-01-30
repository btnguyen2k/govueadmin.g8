package user

import (
	"fmt"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
	_ "github.com/go-sql-driver/mysql"
)

// NewUserDaoSqlite is helper method to create SQLite-implementation of GroupDao
func NewUserDaoSqlite(sqlc *prom.SqlConnect, tableName string) UserDao {
	dao := &UserDaoSqlite{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: mapSqliteFieldToColNameUser},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: mapSqliteColNameToFieldUser},
		ColumnsListMap:              map[string][]string{tableName: colsSqliteUser},
	})
	dao.SetSqlFlavor(prom.FlavorDefault)
	return dao
}

func InitSqliteTableUser(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s TEXT, PRIMARY KEY (%s))",
		tableName, colSqliteUserUsername, colSqliteUserData, colSqliteUserUsername)
	if _, err := sqlc.GetDB().Exec(sql); err != nil {
		panic(err)
	}
}

const (
	colSqliteUserUsername = "uname"
	colSqliteUserData     = "udata"
)

var (
	colsSqliteUser              = []string{colSqliteUserUsername, colSqliteUserData}
	mapSqliteFieldToColNameUser = map[string]interface{}{
		fieldUserUsername: colSqliteUserUsername,
		fieldUserData:     colSqliteUserData,
	}
	mapSqliteColNameToFieldUser = map[string]interface{}{
		colSqliteUserUsername: fieldUserUsername,
		colSqliteUserData:     fieldUserData,
	}
)

type UserDaoSqlite struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *UserDaoSqlite) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{colSqliteUserUsername: bo.GboGetAttrUnsafe(fieldUserUsername, reddo.TypeString)}
}

// ‭toBo transforms godal.IGenericBo to business object.
func (dao *UserDaoSqlite) toBo(gbo godal.IGenericBo) *User {
	if gbo == nil {
		return nil
	}
	username := gbo.GboGetAttrUnsafe(fieldUserUsername, reddo.TypeString).(string)
	data := gbo.GboGetAttrUnsafe(fieldUserData, reddo.TypeString).(string)
	return NewUserBo(username, data)
}

// toGbo transforms business object to godal.IGenericBo
func (dao *UserDaoSqlite) toGbo(bo *User) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldUserUsername, bo.GetUsername())
	gbo.GboSetAttr(fieldUserData, bo.GetData())
	return gbo
}

// Delete implements UserDao.Delete
func (dao *UserDaoSqlite) Delete(bo *User) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements UserDao.Create
func (dao *UserDaoSqlite) Create(bo *User) (bool, error) {
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements UserDao.Get
func (dao *UserDaoSqlite) Get(username string) (*User, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{colSqliteUserUsername: username})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements UserDao.GetN
func (dao *UserDaoSqlite) GetN(fromOffset, maxNumRows int) ([]*User, error) {
	// order ascending by "id" column
	ordering := (&sql.GenericSorting{Flavor: dao.GetSqlFlavor()}).Add(colSqliteUserUsername)
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
func (dao *UserDaoSqlite) GetAll() ([]*User, error) {
	return dao.GetN(0, 0)
}

// Update implements UserDao.Update
func (dao *UserDaoSqlite) Update(bo *User) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}
