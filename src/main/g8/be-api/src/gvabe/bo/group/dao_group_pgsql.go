package group

import (
	"fmt"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
)

// NewGroupDaoPgsql is helper method to create PostgreSQL-implementation of GroupDao
func NewGroupDaoPgsql(sqlc *prom.SqlConnect, tableName string) GroupDao {
	dao := &GroupDaoPgsql{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(&sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: mapPgsqlFieldToColNameGroup},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: mapPgsqlColNameToFieldGroup},
		ColumnsListMap:              map[string][]string{tableName: colsPgsqlGroup},
	})
	dao.SetSqlFlavor(prom.FlavorPgSql)
	return dao
}

func InitPgsqlTableGroup(sqlc *prom.SqlConnect, tableName string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(64), %s VARCHAR(255), PRIMARY KEY (%s))",
		tableName, colPgsqlGroupId, colPgsqlGroupName, colPgsqlGroupId)
	if _, err := sqlc.GetDB().Exec(sql); err != nil {
		panic(err)
	}
}

const (
	colPgsqlGroupId   = "gid"
	colPgsqlGroupName = "gname"
)

var (
	colsPgsqlGroup              = []string{colPgsqlGroupId, colPgsqlGroupName}
	mapPgsqlFieldToColNameGroup = map[string]interface{}{
		fieldGroupId:   colPgsqlGroupId,
		fieldGroupName: colPgsqlGroupName,
	}
	mapPgsqlColNameToFieldGroup = map[string]interface{}{
		colPgsqlGroupId:   fieldGroupId,
		colPgsqlGroupName: fieldGroupName,
	}
)

type GroupDaoPgsql struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *GroupDaoPgsql) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{colPgsqlGroupId: bo.GboGetAttrUnsafe(fieldGroupId, reddo.TypeString)}
}

// ‭toBo transforms godal.IGenericBo to business object.
func (dao *GroupDaoPgsql) toBo(gbo godal.IGenericBo) *Group {
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
func (dao *GroupDaoPgsql) toGbo(bo *Group) godal.IGenericBo {
	if bo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(fieldGroupId, bo.Id)
	gbo.GboSetAttr(fieldGroupName, bo.Name)
	return gbo
}

// Delete implements GroupDao.Delete
func (dao *GroupDaoPgsql) Delete(bo *Group) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Create implements GroupDao.Create
func (dao *GroupDaoPgsql) Create(bo *Group) (bool, error) {
	numRows, err := dao.GdaoCreate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}

// Get implements GroupDao.Get
func (dao *GroupDaoPgsql) Get(id string) (*Group, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{colPgsqlGroupId: id})
	if err != nil {
		return nil, err
	}
	return dao.toBo(gbo), nil
}

// GetN implements GroupDao.GetN
func (dao *GroupDaoPgsql) GetN(fromOffset, maxNumRows int) ([]*Group, error) {
	// order ascending by "id" column
	ordering := (&sql.GenericSorting{Flavor: dao.GetSqlFlavor()}).Add(colPgsqlGroupId)
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
func (dao *GroupDaoPgsql) GetAll() ([]*Group, error) {
	return dao.GetN(0, 0)
}

// Update implements GroupDao.Update
func (dao *GroupDaoPgsql) Update(bo *Group) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.toGbo(bo))
	return numRows > 0, err
}
