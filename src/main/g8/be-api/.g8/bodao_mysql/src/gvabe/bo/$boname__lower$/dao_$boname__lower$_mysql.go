package $boname;format="lower"$

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
	_ "github.com/go-sql-driver/mysql"
	"main/src/goapi"
	"strings"
	"time"
)

// newMysqlConnection is helper method to create connection pool to MySQL server
func newMysqlConnection() *prom.SqlConnect {
	// read from config or use default value
	driver := "mysql"
	timeoutMs := goapi.AppConfig.GetInt32("dao$boname__lower$.mysql.timeout", 10000)
	timezone := goapi.AppConfig.GetString("dao$boname__lower$.mysql.timezone", "Asia/Ho_Chi_Minh")
	dsn := goapi.AppConfig.GetString("dao$boname__lower$.mysql.dsn",
		"test:test@tcp(localhost:3306)/test?charset=utf8mb4,utf8&parseTime=true&loc="+strings.Replace(timezone, "/", "%2f", -1))
	if sqlConnect, err := prom.NewSqlConnectWithFlavor(driver, dsn, int(timeoutMs), nil, prom.FlavorMySql); err != nil {
		panic(err)
	} else if sqlConnect == nil {
		panic("error creating [prom.SqlConnect] instance")
	} else {
		loc, _ := time.LoadLocation(timezone)
		sqlConnect.SetLocation(loc)
		return sqlConnect
	}
}

// new$boname;format="Camel"$DaoMysql is helper method to create MySQL-implementation of $boname;format="Camel"$Dao
func new$boname;format="Camel"$DaoMysql(sqlc *prom.SqlConnect, tableName string) $boname;format="Camel"$Dao {
	dao := &$boname;format="Camel"$DaoMysql{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: mapMysqlFieldToColName$boname;format="Camel"$},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: mapMysqlColNameToField$boname;format="Camel"$},
		ColumnsListMap:              map[string][]string{tableName: colsMysql$boname;format="Camel"$},
	})
	dao.SetSqlFlavor(prom.FlavorMySql)
	return dao
}

/*
Table schema for $boname;format="Camel"$:

CREATE TABLE IF NOT EXISTS tbl_$boname__lower$ (
	id			VARCHAR(64),
	val_str		VARCHAR(255)		NOT NULL,
	val_int		INT					NOT NULL DEFAULT (0),
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
*/

const (
	tableMysql$boname;format="Camel"$    = "tbl_$boname__lower$"
	colMysql$boname;format="Camel"$Id    = "id"
	colMysql$boname;format="Camel"$Name  = "val_str"
	colMysql$boname;format="Camel"$Value = "val_int"
)

var (
	colsMysql$boname;format="Camel"$              = []string{colMysql$boname;format="Camel"$Id, colMysql$boname;format="Camel"$Name, colMysql$boname;format="Camel"$Value}
	mapMysqlFieldToColName$boname;format="Camel"$ = map[string]interface{}{
		field$boname;format="Camel"$Id   : colMysql$boname;format="Camel"$Id,
		field$boname;format="Camel"$Name : colMysql$boname;format="Camel"$Name,
		field$boname;format="Camel"$Value: colMysql$boname;format="Camel"$Value,
	}
	mapMysqlColNameToField$boname;format="Camel"$ = map[string]interface{}{
		colMysql$boname;format="Camel"$Id   : field$boname;format="Camel"$Id,
		colMysql$boname;format="Camel"$Name : field$boname;format="Camel"$Name,
		colMysql$boname;format="Camel"$Value: field$boname;format="Camel"$Value,
	}
)

type $boname;format="Camel"$DaoMysql struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *$boname;format="Camel"$DaoMysql) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{colMysql$boname;format="Camel"$Id: bo.GboGetAttrUnsafe(field$boname;format="Camel"$Id, reddo.TypeString)}
}

// ‭toBo transforms godal.IGenericBo to business object.
func (dao *$boname;format="Camel"$DaoMysql) toBo(gbo godal.IGenericBo) *$boname;format="Camel"$ {
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
func (dao *$boname;format="Camel"$DaoMysql) toGbo(bo *$boname;format="Camel"$) godal.IGenericBo {
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
func (dao *$boname;format="Camel"$DaoMysql) Delete(bo *$boname;format="Camel"$) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements $boname;format="Camel"$Dao.Create
func (dao *$boname;format="Camel"$DaoMysql) Create(bo *$boname;format="Camel"$) (bool, error) {
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements $boname;format="Camel"$Dao.Get
func (dao *$boname;format="Camel"$DaoMysql) Get(id string) (*$boname;format="Camel"$, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{colMysql$boname;format="Camel"$Id: id})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements $boname;format="Camel"$Dao.GetN
func (dao *$boname;format="Camel"$DaoMysql) GetN(fromOffset, maxNumRows int) ([]*$boname;format="Camel"$, error) {
	// order ascending by "id" column
	ordering := (&sql.GenericSorting{Flavor: dao.GetSqlFlavor()}).Add(colMysql$boname;format="Camel"$Id)
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, ordering, fromOffset, maxNumRows)
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
func (dao *$boname;format="Camel"$DaoMysql) GetAll() ([]*$boname;format="Camel"$, error) {
	return dao.GetN(0, 0)
}

// Update implements $boname;format="Camel"$Dao.Update
func (dao *$boname;format="Camel"$DaoMysql) Update(bo *$boname;format="Camel"$) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}
