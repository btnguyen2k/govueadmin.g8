package henge

import (
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"
)

func buildRowMapper(tableName string, extraColNameToFieldMappings map[string]string) *sql.GenericRowMapperSql {
	myCols := append([]string{}, columnNames...)
	myMapFieldToColName := cloneMap(mapFieldToColName)
	myMapColNameToField := cloneMap(mapColNameToField)
	for col, field := range extraColNameToFieldMappings {
		myCols = append(myCols, col)
		myMapColNameToField[col] = field
		myMapFieldToColName[field] = col
	}
	return &sql.GenericRowMapperSql{
		NameTransformation:          sql.NameTransfLowerCase,
		GboFieldToColNameTranslator: map[string]map[string]interface{}{tableName: myMapFieldToColName},
		ColNameToGboFieldTranslator: map[string]map[string]interface{}{tableName: myMapColNameToField},
		ColumnsListMap:              map[string][]string{tableName: myCols},
	}
}

// NewUniversalDaoSql is helper method to create SQL-implementation of UniversalDao
func NewUniversalDaoSql(sqlc *prom.SqlConnect, tableName string, extraColNameToFieldMappings map[string]string) UniversalDao {
	dao := &UniversalDaoSql{tableName: tableName}
	dao.GenericDaoSql = sql.NewGenericDaoSql(sqlc, godal.NewAbstractGenericDao(dao))
	dao.SetRowMapper(buildRowMapper(tableName, extraColNameToFieldMappings))
	dao.SetSqlFlavor(sqlc.GetDbFlavor())
	dao.SetTxModeOnWrite(true)
	return dao
}

const (
	ColId          = "zid"
	ColData        = "zdata"
	ColChecksum    = "zchecksum"
	ColTimeCreated = "ztcreated"
	ColTimeUpdated = "ztupdated"
	ColAppVersion  = "zaversion"
)

var (
	columnNames       = []string{ColId, ColData, ColAppVersion, ColChecksum, ColTimeCreated, ColTimeUpdated}
	mapFieldToColName = map[string]interface{}{
		FieldId:          ColId,
		FieldData:        ColData,
		FieldAppVersion:  ColAppVersion,
		FieldChecksum:    ColChecksum,
		FieldTimeCreated: ColTimeCreated,
		FieldTimeUpdated: ColTimeUpdated,
	}
	mapColNameToField = map[string]interface{}{
		ColId:          FieldId,
		ColData:        FieldData,
		ColAppVersion:  FieldAppVersion,
		ColChecksum:    FieldChecksum,
		ColTimeCreated: FieldTimeCreated,
		ColTimeUpdated: FieldTimeUpdated,
	}
)

type UniversalDaoSql struct {
	*sql.GenericDaoSql
	tableName string
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *UniversalDaoSql) GdaoCreateFilter(_ string, bo godal.IGenericBo) interface{} {
	return map[string]interface{}{ColId: bo.GboGetAttrUnsafe(FieldId, reddo.TypeString)}
}

// ToUniversalBo transforms godal.IGenericBo to business object.
func (dao *UniversalDaoSql) ToUniversalBo(gbo godal.IGenericBo) *UniversalBo {
	if gbo == nil {
		return nil
	}
	extraFields := make(map[string]interface{})
	gbo.GboTransferViaJson(&extraFields)
	for _, field := range topLevelFieldList {
		delete(extraFields, field)
	}
	return &UniversalBo{
		id:          gbo.GboGetAttrUnsafe(FieldId, reddo.TypeString).(string),
		dataJson:    gbo.GboGetAttrUnsafe(FieldData, reddo.TypeString).(string),
		checksum:    gbo.GboGetAttrUnsafe(FieldChecksum, reddo.TypeString).(string),
		timeCreated: gbo.GboGetAttrUnsafe(FieldTimeCreated, reddo.TypeTime).(time.Time),
		timeUpdated: gbo.GboGetAttrUnsafe(FieldTimeUpdated, reddo.TypeTime).(time.Time),
		appVersion:  gbo.GboGetAttrUnsafe(FieldAppVersion, reddo.TypeUint).(uint64),
		_extraAttrs: extraFields,
	}
}

// ToGenericBo transforms business object to godal.IGenericBo
func (dao *UniversalDaoSql) ToGenericBo(ubo *UniversalBo) godal.IGenericBo {
	if ubo == nil {
		return nil
	}
	gbo := godal.NewGenericBo()
	gbo.GboSetAttr(FieldId, ubo.id)
	gbo.GboSetAttr(FieldData, ubo.dataJson)
	gbo.GboSetAttr(FieldChecksum, ubo.checksum)
	gbo.GboSetAttr(FieldTimeCreated, ubo.timeCreated)
	gbo.GboSetAttr(FieldTimeUpdated, ubo.timeUpdated)
	gbo.GboSetAttr(FieldAppVersion, ubo.appVersion)
	for k, v := range ubo._extraAttrs {
		gbo.GboSetAttr(k, v)
	}
	return gbo
}

// Delete implements UniversalDao.Delete
func (dao *UniversalDaoSql) Delete(bo *UniversalBo) (bool, error) {
	numRows, err := dao.GdaoDelete(dao.tableName, dao.ToGenericBo(bo))
	return numRows > 0, err
}

// Create implements UniversalDao.Create
func (dao *UniversalDaoSql) Create(bo *UniversalBo) (bool, error) {
	numRows, err := dao.GdaoCreate(dao.tableName, dao.ToGenericBo(bo))
	return numRows > 0, err
}

// Get implements UniversalDao.Get
func (dao *UniversalDaoSql) Get(id string) (*UniversalBo, error) {
	gbo, err := dao.GdaoFetchOne(dao.tableName, map[string]interface{}{ColId: id})
	if err != nil {
		return nil, err
	}
	return dao.ToUniversalBo(gbo), nil
}

// GetN implements UniversalDao.GetN
func (dao *UniversalDaoSql) GetN(fromOffset, maxNumRows int) ([]*UniversalBo, error) {
	// order ascending by "id" column
	ordering := (&sql.GenericSorting{Flavor: dao.GetSqlFlavor()}).Add(ColId)
	gboList, err := dao.GdaoFetchMany(dao.tableName, nil, ordering, fromOffset, maxNumRows)
	if err != nil {
		return nil, err
	}
	result := make([]*UniversalBo, 0)
	for _, gbo := range gboList {
		bo := dao.ToUniversalBo(gbo)
		result = append(result, bo)
	}
	return result, nil
}

// GetAll implements UniversalDao.GetAll
func (dao *UniversalDaoSql) GetAll() ([]*UniversalBo, error) {
	return dao.GetN(0, 0)
}

// Update implements UniversalDao.Update
func (dao *UniversalDaoSql) Update(bo *UniversalBo) (bool, error) {
	numRows, err := dao.GdaoUpdate(dao.tableName, dao.ToGenericBo(bo))
	return numRows > 0, err
}

// Save implements UniversalDao.Save
func (dao *UniversalDaoSql) Save(bo *UniversalBo) (bool, error) {
	numRows, err := dao.GdaoSave(dao.tableName, dao.ToGenericBo(bo))
	return numRows > 0, err
}
