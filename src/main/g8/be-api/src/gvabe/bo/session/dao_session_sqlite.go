package session

import (
	"fmt"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
	"time"
)

// NewSessionDaoSqlite is helper method to create SQLite-implementation of SessionDao
func NewSessionDaoSqlite(sqlc *prom.SqlConnect, tableName string) SessionDao {
	dao := &SessionDaoSqlite{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: mapSqliteFieldToColNameSession},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: mapSqliteColNameToFieldSession},
		ColumnsListMap:              map[string][]string{tableName: colsSqliteSession},
	})
	dao.SetSqlFlavor(prom.FlavorDefault)
	return dao
}

func InitSqliteTableSession(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(255), %s TEXT, %s DATETIME, %s VARCHAR(255), PRIMARY KEY (%s))",
		tableName, colSqliteSessionId, colSqliteSessionData, colSqliteSessionExpiry, colSqliteSessionParentId, colSqliteSessionId)
	if _, err := sqlc.GetDB().Exec(sql); err != nil {
		panic(err)
	}
}

const (
	colSqliteSessionId       = "sid"
	colSqliteSessionData     = "sdata"
	colSqliteSessionExpiry   = "sexpiry"
	colSqliteSessionParentId = "spid"
)

var (
	colsSqliteSession              = []string{colSqliteSessionId, colSqliteSessionData, colSqliteSessionExpiry, colSqliteSessionParentId}
	mapSqliteFieldToColNameSession = map[string]interface{}{
		fieldSessionId:       colSqliteSessionId,
		fieldSessionData:     colSqliteSessionData,
		fieldSessionExpiry:   colSqliteSessionExpiry,
		fieldSessionParentId: colSqliteSessionParentId,
	}
	mapSqliteColNameToFieldSession = map[string]interface{}{
		colSqliteSessionId:       fieldSessionId,
		colSqliteSessionData:     fieldSessionData,
		colSqliteSessionExpiry:   fieldSessionExpiry,
		colSqliteSessionParentId: fieldSessionParentId,
	}
)

type SessionDaoSqlite struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *SessionDaoSqlite) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{colSqliteSessionId: bo.GboGetAttrUnsafe(fieldSessionId, reddo.TypeString)}
}

// ‭toBo transforms godal.IGenericBo to business object.
func (dao *SessionDaoSqlite) toBo(gbo godal.IGenericBo) *Session {
	if gbo == nil {
		return nil
	}
	bo := &Session{
		Id:       gbo.GboGetAttrUnsafe(fieldSessionId, reddo.TypeString).(string),
		Data:     gbo.GboGetAttrUnsafe(fieldSessionData, reddo.TypeString).(string),
		Expiry:   gbo.GboGetAttrUnsafe(fieldSessionExpiry, reddo.TypeTime).(time.Time),
		ParentId: gbo.GboGetAttrUnsafe(fieldSessionParentId, reddo.TypeString).(string),
	}
	return bo
}

// toGbo transforms business object to godal.IGenericBo
func (dao *SessionDaoSqlite) toGbo(bo *Session) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldSessionId, bo.Id)
	gbo.GboSetAttr(fieldSessionData, bo.Data)
	gbo.GboSetAttr(fieldSessionExpiry, bo.Expiry)
	gbo.GboSetAttr(fieldSessionParentId, bo.ParentId)
	return gbo
}

// Delete implements SessionDao.Delete
func (dao *SessionDaoSqlite) Delete(bo *Session) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements SessionDao.Create
func (dao *SessionDaoSqlite) Create(bo *Session) (bool, error) {
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements SessionDao.Get
func (dao *SessionDaoSqlite) Get(id string) (*Session, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{colSqliteSessionId: id})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements SessionDao.GetN
func (dao *SessionDaoSqlite) GetN(fromOffset, maxNumRows int) ([]*Session, error) {
	// order ascending by "id" column
	ordering := (&sql.GenericSorting{Flavor: dao.GetSqlFlavor()}).Add(colSqliteSessionId)
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, ordering, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*Session, 0)
	for _, gbo := range gboList {
		bo := dao.toBo(gbo)
		result = append(result, bo)
	}
	return result, nil
}

// GetAll implements SessionDao.GetAll
func (dao *SessionDaoSqlite) GetAll() ([]*Session, error) {
	return dao.GetN(0, 0)
}

// Update implements SessionDao.Update
func (dao *SessionDaoSqlite) Update(bo *Session) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}
