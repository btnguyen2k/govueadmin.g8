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
	dao.SetSqlFlavor(prom.FlavorMySql)
	return dao
}

func InitSqliteTableUser(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s VARCHAR(64), %s VARCHAR(64), %s VARCHAR(64), PRIMARY KEY (%s))",
		tableName, colSqliteUserUsername, colSqliteUserPassword, colSqliteUserName, colSqliteUserGroupId, colSqliteUserUsername)
	if _, err := sqlc.GetDB().Exec(sql); err != nil {
		panic(err)
	}
}

const (
	colSqliteUserUsername = "uname"
	colSqliteUserPassword = "upwd"
	colSqliteUserName     = "display_name"
	colSqliteUserGroupId  = "gid"
)

var (
	colsSqliteUser              = []string{colSqliteUserUsername, colSqliteUserPassword, colSqliteUserName, colSqliteUserGroupId}
	mapSqliteFieldToColNameUser = map[string]interface{}{
		fieldUserUsername: colSqliteUserUsername,
		fieldUserPassword: colSqliteUserPassword,
		fieldUserName:     colSqliteUserName,
		fieldUserGroupId:  colSqliteUserGroupId,
	}
	mapSqliteColNameToFieldUser = map[string]interface{}{
		colSqliteUserUsername: fieldUserUsername,
		colSqliteUserPassword: fieldUserPassword,
		colSqliteUserName:     fieldUserName,
		colSqliteUserGroupId:  fieldUserGroupId,
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
	bo := &User{
		Username: gbo.GboGetAttrUnsafe(fieldUserUsername, reddo.TypeString).(string),
		Password: gbo.GboGetAttrUnsafe(fieldUserPassword, reddo.TypeString).(string),
		Name:     gbo.GboGetAttrUnsafe(fieldUserName, reddo.TypeString).(string),
		GroupId:  gbo.GboGetAttrUnsafe(fieldUserGroupId, reddo.TypeString).(string),
	}
	return bo
}

// toGbo transforms business object to godal.IGenericBo
func (dao *UserDaoSqlite) toGbo(bo *User) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldUserUsername, bo.Username)
	gbo.GboSetAttr(fieldUserPassword, bo.Password)
	gbo.GboSetAttr(fieldUserName, bo.Name)
	gbo.GboSetAttr(fieldUserGroupId, bo.GroupId)
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
