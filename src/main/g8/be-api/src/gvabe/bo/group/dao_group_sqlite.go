package group

import (
	"fmt"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
)

// NewGroupDaoSqlite is helper method to create SQLite-implementation of GroupDao
func NewGroupDaoSqlite(sqlc *prom.SqlConnect, tableName string) GroupDao {
	dao := &GroupDaoSqlite{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: mapSqliteFieldToColNameGroup},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: mapSqliteColNameToFieldGroup},
		ColumnsListMap:              map[string][]string{tableName: colsSqliteGroup},
	})
	dao.SetSqlFlavor(prom.FlavorDefault)
	return dao
}

func InitSqliteTableGroup(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s VARCHAR(255), PRIMARY KEY (%s))",
		tableName, colSqliteGroupId, colSqliteGroupName, colSqliteGroupId)
	if _, err := sqlc.GetDB().Exec(sql); err != nil {
		panic(err)
	}
}

const (
	colSqliteGroupId   = "gid"
	colSqliteGroupName = "gname"
)

var (
	colsSqliteGroup              = []string{colSqliteGroupId, colSqliteGroupName}
	mapSqliteFieldToColNameGroup = map[string]interface{}{
		fieldGroupId:   colSqliteGroupId,
		fieldGroupName: colSqliteGroupName,
	}
	mapSqliteColNameToFieldGroup = map[string]interface{}{
		colSqliteGroupId:   fieldGroupId,
		colSqliteGroupName: fieldGroupName,
	}
)

type GroupDaoSqlite struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *GroupDaoSqlite) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{colSqliteGroupId: bo.GboGetAttrUnsafe(fieldGroupId, reddo.TypeString)}
}

// ‭toBo transforms godal.IGenericBo to business object.
func (dao *GroupDaoSqlite) toBo(gbo godal.IGenericBo) *Group {
	if gbo == nil {
		return nil
	}
	bo := &Group{
		Id:   gbo.GboGetAttrUnsafe(fieldGroupId, reddo.TypeString).(string),
		Name: gbo.GboGetAttrUnsafe(fieldGroupName, reddo.TypeString).(string),
	}
	return bo
}

// toGbo transforms business object to godal.IGenericBo
func (dao *GroupDaoSqlite) toGbo(bo *Group) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldGroupId, bo.Id)
	gbo.GboSetAttr(fieldGroupName, bo.Name)
	return gbo
}

// Delete implements GroupDao.Delete
func (dao *GroupDaoSqlite) Delete(bo *Group) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements GroupDao.Create
func (dao *GroupDaoSqlite) Create(bo *Group) (bool, error) {
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements GroupDao.Get
func (dao *GroupDaoSqlite) Get(id string) (*Group, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{colSqliteGroupId: id})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements GroupDao.GetN
func (dao *GroupDaoSqlite) GetN(fromOffset, maxNumRows int) ([]*Group, error) {
	// order ascending by "id" column
	ordering := (&sql.GenericSorting{Flavor: dao.GetSqlFlavor()}).Add(colSqliteGroupId)
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, ordering, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*Group, 0)
	for _, gbo := range gboList {
		bo := dao.toBo(gbo)
		result = append(result, bo)
	}
	return result, nil
}

// GetAll implements GroupDao.GetAll
func (dao *GroupDaoSqlite) GetAll() ([]*Group, error) {
	return dao.GetN(0, 0)
}

// Update implements GroupDao.Update
func (dao *GroupDaoSqlite) Update(bo *Group) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}
