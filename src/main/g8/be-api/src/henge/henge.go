// package henge offers universal data access layer implementation
package henge

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/btnguyen2k/consu/checksum"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/consu/semita"
	"github.com/btnguyen2k/godal"
)

func cloneMap(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range src {
		result[k] = v
	}
	return result
}

// NewUniversalBo is helper function to create a new UniversalBo instance
func NewUniversalBo(id string, appVersion uint64) *UniversalBo {
	now := time.Now()
	bo := &UniversalBo{
		id:          strings.ToLower(strings.TrimSpace(id)),
		timeCreated: now,
		timeUpdated: now,
		appVersion:  appVersion,
		_dirty:      true,
		_extraAttrs: make(map[string]interface{}),
	}
	return bo.Sync()
}

const (
	FieldId          = "id"
	FieldData        = "data"
	FieldAppVersion  = "aver"
	FieldChecksum    = "csum"
	FieldTimeCreated = "tcre"
	FieldTimeUpdated = "tupd"
	FieldExtras      = "_ext"

	TimeLayout = "2006-01-02T15:04:05.999999-07:00"
)

var (
	topLevelFieldList = []string{FieldId, FieldData, FieldChecksum, FieldAppVersion, FieldTimeCreated, FieldTimeUpdated}
)

// UniversalBo is the "universal" business object
type UniversalBo struct {
	/* top level attributes */
	id          string    `json:"id"`   // bo's unique identifier
	dataJson    string    `json:"data"` // bo's attributes encoded as JSON string
	appVersion  uint64    `json:"aver"` // for application internal use
	checksum    string    `json:"csum"` // bo's checksum (should not take update-time into account)
	timeCreated time.Time `json:"tcre"` // bo's creation timestamp
	timeUpdated time.Time `json:"tupd"` // bo's last update timestamp

	/* computed attributes */
	_data       interface{}            `json:"-"`    // deserialized form of data-json
	_sdata      *semita.Semita         `json:"-"`    // used to access data in hierarchy manner
	_extraAttrs map[string]interface{} `json:"_ext"` // other top-level arbitrary attributes
	_lock       sync.RWMutex
	_dirty      bool
}

// FuncPreUboToMap is used by UniversalBo.ToMap to export a UniversalBo to a map[string]interface{}
type FuncPreUboToMap func(*UniversalBo) map[string]interface{}

// FuncPostUboToMap is used by UniversalBo.ToMap to transform the result map
type FuncPostUboToMap func(map[string]interface{}) map[string]interface{}

// DefaultFuncPreUboToMap is default implementation of FuncPreUboToMap
//
// This function exports the input UniversalBo to a map with following fields:
// { FieldId (string), FieldData (string), FieldAppVersion (uint64), FieldChecksum (string),
// FieldTimeCreated (time.Time), FieldTimeUpdated (time.Time), FieldExtras (map[string]interface{}) }
var DefaultFuncPreUboToMap FuncPreUboToMap = func(_ubo *UniversalBo) map[string]interface{} {
	ubo := _ubo.Clone()
	return map[string]interface{}{
		FieldId:          ubo.id,
		FieldData:        ubo.dataJson,
		FieldAppVersion:  ubo.appVersion,
		FieldChecksum:    ubo.checksum,
		FieldTimeCreated: ubo.timeCreated,
		FieldTimeUpdated: ubo.timeUpdated,
		FieldExtras:      cloneMap(ubo._extraAttrs),
	}
}

// ToMap exports the UniversalBo to a map[string]interface{}
//  - preFunc is used to export UniversalBo to a map. If not supplied, DefaultFuncPreUboToMap is used.
//  - postFunc is used to transform the result map (output from preFunc) to the final result. If not supplied, the result from preFunc is returned.
func (ubo *UniversalBo) ToMap(preFunc FuncPreUboToMap, postFunc FuncPostUboToMap) map[string]interface{} {
	if preFunc == nil {
		preFunc = DefaultFuncPreUboToMap
	}
	result := preFunc(ubo.Clone())
	if postFunc != nil {
		result = postFunc(result)
	}
	return result
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON
func (ubo *UniversalBo) MarshalJSON() ([]byte, error) {
	ubo.Sync()
	ubo._lock.RLock()
	defer ubo._lock.RUnlock()
	m := map[string]interface{}{
		FieldId:          ubo.id,
		FieldData:        ubo.dataJson,
		FieldAppVersion:  ubo.appVersion,
		FieldChecksum:    ubo.checksum,
		FieldTimeCreated: ubo.timeCreated,
		FieldTimeUpdated: ubo.timeUpdated,
		FieldExtras:      cloneMap(ubo._extraAttrs),
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON
func (ubo *UniversalBo) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var err error
	if m[FieldId], err = reddo.ToString(m[FieldId]); err != nil {
		return err
	}
	if m[FieldData], err = reddo.ToString(m[FieldData]); err != nil {
		return err
	}
	if m[FieldAppVersion], err = reddo.ToUint(m[FieldAppVersion]); err != nil {
		return err
	}
	if m[FieldChecksum], err = reddo.ToString(m[FieldChecksum]); err != nil {
		return err
	}
	if m[FieldTimeCreated], err = reddo.ToTimeWithLayout(m[FieldTimeCreated], TimeLayout); err != nil {
		return err
	}
	if m[FieldTimeUpdated], err = reddo.ToTimeWithLayout(m[FieldTimeUpdated], TimeLayout); err != nil {
		return err
	}
	if m[FieldExtras], err = reddo.ToMap(m[FieldExtras], reflect.TypeOf(map[string]interface{}{})); err != nil {
		return err
	}

	ubo._lock.Lock()
	defer ubo._lock.Unlock()
	ubo.id = m[FieldId].(string)
	ubo.appVersion = m[FieldAppVersion].(uint64)
	ubo.checksum = m[FieldChecksum].(string)
	ubo.timeCreated = m[FieldTimeCreated].(time.Time)
	ubo.timeUpdated = m[FieldTimeUpdated].(time.Time)
	if m[FieldExtras] != nil {
		ubo._extraAttrs = m[FieldExtras].(map[string]interface{})
	} else {
		ubo._extraAttrs = make(map[string]interface{})
	}
	ubo._setDataJson(m[FieldData].(string))
	ubo._sync()
	return nil
}

// GetId returns value of bo's 'id' field
func (ubo *UniversalBo) GetId() string {
	return ubo.id
}

// SetId sets value of bo's 'id' field
func (ubo *UniversalBo) SetId(value string) *UniversalBo {
	ubo._lock.Lock()
	defer ubo._lock.Unlock()
	ubo.id = strings.TrimSpace(value)
	ubo._dirty = true
	return ubo
}

// GetDataJson returns value of bo's 'data-json' field
func (ubo *UniversalBo) GetDataJson() string {
	return ubo.dataJson
}

type dataInitType int

const (
	dataInitNone dataInitType = iota
	dataInitMap
	dataInitSlice
)

var (
	ErrorDataInitedAsMap   = errors.New("data is initialized as empty map")
	ErrorDataInitedAsSlice = errors.New("data is initialized as empty slice")
)

func (ubo *UniversalBo) _parseDataJson(dataInit dataInitType) error {
	err := json.Unmarshal([]byte(ubo.dataJson), &ubo._data)
	if err != nil || ubo._data == nil {
		if dataInit == dataInitMap {
			ubo._data = make(map[string]interface{})
			err = ErrorDataInitedAsMap
		} else if dataInit == dataInitSlice {
			ubo._data = make([]interface{}, 0)
			err = ErrorDataInitedAsMap
		} else {
			ubo._data = nil
		}
	}
	if ubo._data != nil {
		ubo._sdata = semita.NewSemita(&ubo._data)
	} else {
		ubo._sdata = nil
	}
	return err
}

func (ubo *UniversalBo) _setDataJson(value string) *UniversalBo {
	ubo.dataJson = strings.TrimSpace(value)
	ubo._parseDataJson(dataInitNone)
	ubo._dirty = true
	return ubo
}

// SetDataJson sets value of bo's 'data-json' field
func (ubo *UniversalBo) SetDataJson(value string) *UniversalBo {
	ubo._lock.Lock()
	defer ubo._lock.Unlock()
	return ubo._setDataJson(value)
}

// GetAppVersion returns value of bo's 'app-version' field
func (ubo *UniversalBo) GetAppVersion() uint64 {
	return ubo.appVersion
}

// SetAppVersion sets value of bo's 'id' field
func (ubo *UniversalBo) SetAppVersion(value uint64) *UniversalBo {
	ubo._lock.Lock()
	defer ubo._lock.Unlock()
	ubo.appVersion = value
	ubo._dirty = true
	return ubo
}

// GetChecksum returns value of bo's 'checksum' field
func (ubo *UniversalBo) GetChecksum() string {
	return ubo.checksum
}

// GetChecksum returns value of bo's 'timestamp-created' field
func (ubo *UniversalBo) GetTimeCreated() time.Time {
	return ubo.timeCreated
}

// GetChecksum returns value of bo's 'timestamp-updated' field
func (ubo *UniversalBo) GetTimeUpdated() time.Time {
	return ubo.timeUpdated
}

// GetChecksum returns value of bo's 'timestamp-updated' field
func (ubo *UniversalBo) SetTimeUpdated(value time.Time) *UniversalBo {
	ubo.timeCreated = value
	return ubo
}

// IsDirty returns 'true' if bo's data has been modified
func (ubo *UniversalBo) IsDirty() bool {
	return ubo._dirty
}

// GetDataAttr is alias of GetDataAttrAs(path, nil)
func (ubo *UniversalBo) GetDataAttr(path string) (interface{}, error) {
	return ubo.GetDataAttrAs(path, nil)
}

// GetDataAttrUnsafe is similar to GetDataAttr but ignoring error
func (ubo *UniversalBo) GetDataAttrUnsafe(path string) interface{} {
	return ubo.GetDataAttrAsUnsafe(path, nil)
}

// GetDataAttrAsUnsafe is similar to GetDataAttrAs but ignoring error
func (ubo *UniversalBo) GetDataAttrAsUnsafe(path string, typ reflect.Type) interface{} {
	v, _ := ubo.GetDataAttrAs(path, typ)
	return v
}

// GetDataAttrAsTimeWithLayout returns value, converted to time, of a data attribute located at 'path'
func (ubo *UniversalBo) GetDataAttrAsTimeWithLayout(path, layout string) (time.Time, error) {
	v, _ := ubo.GetDataAttr(path)
	return reddo.ToTimeWithLayout(v, layout)
}

// GetDataAttrAsTimeWithLayoutUnsafe is similar to GetDataAttrAsTimeWithLayout but ignoring error
func (ubo *UniversalBo) GetDataAttrAsTimeWithLayoutUnsafe(path, layout string) time.Time {
	t, _ := ubo.GetDataAttrAsTimeWithLayout(path, layout)
	return t
}

func (ubo *UniversalBo) _initSdata(path string) {
	if ubo._sdata == nil {
		dataInit := dataInitMap
		if strings.HasSuffix(path, "[") {
			dataInit = dataInitSlice
		}
		ubo._parseDataJson(dataInit)
	}
}

// GetDataAttrAs returns value, converted to the specified type, of a data attribute located at 'path'
func (ubo *UniversalBo) GetDataAttrAs(path string, typ reflect.Type) (interface{}, error) {
	ubo._lock.RLock()
	defer ubo._lock.RUnlock()
	ubo._initSdata(path)
	if ubo._sdata == nil {
		return nil, errors.New("cannot get data at path [" + path + "]")
	}
	return ubo._sdata.GetValueOfType(path, typ)
}

// SetDataAttr sets value of a data attribute located at 'path'
func (ubo *UniversalBo) SetDataAttr(path string, value interface{}) error {
	ubo._lock.Lock()
	defer ubo._lock.Unlock()
	ubo._dirty = true
	ubo._initSdata(path)
	if ubo._sdata == nil {
		return errors.New("cannot set data at path [" + path + "]")
	}
	switch value.(type) {
	case time.Time:
		value, _ = time.Parse(TimeLayout, value.(time.Time).Format(TimeLayout))
	case *time.Time:
		value, _ = time.Parse(TimeLayout, value.(*time.Time).Format(TimeLayout))
	}
	return ubo._sdata.SetValue(path, value)
}

// GetExtraAttrs returns the 'extra-attrs' map
func (ubo *UniversalBo) GetExtraAttrs() map[string]interface{} {
	ubo._lock.RLock()
	defer ubo._lock.RUnlock()
	return cloneMap(ubo._extraAttrs)
}

// GetExtraAttr returns value of an 'extra' attribute specified by 'key'
func (ubo *UniversalBo) GetExtraAttr(key string) interface{} {
	v := ubo._extraAttrs[key]
	return v
}

// GetExtraAttrAs returns value, converted to the specified type, of an 'extra' attribute specified by 'key'
func (ubo *UniversalBo) GetExtraAttrAs(key string, typ reflect.Type) (interface{}, error) {
	v := ubo.GetExtraAttr(key)
	return reddo.Convert(v, typ)
}

// GetExtraAttrAsTimeWithLayout returns value, converted to time, of an 'extra' attribute specified by 'key'
func (ubo *UniversalBo) GetExtraAttrAsTimeWithLayout(key, layout string) (time.Time, error) {
	v := ubo.GetExtraAttr(key)
	return reddo.ToTimeWithLayout(v, layout)
}

// GetExtraAttrAsUnsafe is similar to GetExtraAttrAs but no error is returned
func (ubo *UniversalBo) GetExtraAttrAsUnsafe(key string, typ reflect.Type) interface{} {
	v, _ := ubo.GetExtraAttrAs(key, typ)
	return v
}

// GetExtraAttrAsTimeWithLayoutUnsafe is similar to GetExtraAttrAsTimeWithLayout but no error is returned
func (ubo *UniversalBo) GetExtraAttrAsTimeWithLayoutUnsafe(key, layout string) time.Time {
	t, _ := ubo.GetExtraAttrAsTimeWithLayout(key, layout)
	return t
}

// SetExtraAttr sets value of an 'extra' attribute specified by 'key'
func (ubo *UniversalBo) SetExtraAttr(key string, value interface{}) *UniversalBo {
	ubo._lock.Lock()
	defer ubo._lock.Unlock()
	if ubo._extraAttrs == nil {
		ubo._extraAttrs = make(map[string]interface{})
	}
	ubo._dirty = true
	switch value.(type) {
	case time.Time:
		value, _ = time.Parse(TimeLayout, value.(time.Time).Format(TimeLayout))
	case *time.Time:
		value, _ = time.Parse(TimeLayout, value.(*time.Time).Format(TimeLayout))
	}
	ubo._extraAttrs[key] = value
	return ubo
}

func (ubo *UniversalBo) _sync() *UniversalBo {
	if ubo._dirty {
		ubo.timeUpdated = time.Now()
		csumMap := map[string]interface{}{
			"id":          ubo.id,
			"app_version": ubo.appVersion,
			"t_created":   ubo.timeCreated.Format(TimeLayout),
			"data":        ubo._data,
			"extra":       ubo._extraAttrs,
		}
		ubo.checksum = fmt.Sprintf("%x", checksum.Md5Checksum(csumMap))
		js, _ := json.Marshal(ubo._data)
		ubo.dataJson = string(js)
		ubo._dirty = false
	}
	return ubo
}

// Sync syncs bo's '_data' object and 'data-json' attribute and returns itself.
func (ubo *UniversalBo) Sync() *UniversalBo {
	ubo._lock.Lock()
	defer ubo._lock.Unlock()
	return ubo._sync()
}

// Clone creates a cloned copy of the business object
func (ubo *UniversalBo) Clone() *UniversalBo {
	ubo.Sync()
	ubo._lock.RLock()
	defer ubo._lock.RUnlock()
	return &UniversalBo{
		id:          ubo.id,
		dataJson:    ubo.dataJson,
		appVersion:  ubo.appVersion,
		checksum:    ubo.checksum,
		timeCreated: ubo.timeCreated,
		timeUpdated: ubo.timeUpdated,
		_data:       nil,
		_sdata:      nil,
		_extraAttrs: cloneMap(ubo._extraAttrs),
		_dirty:      false,
	}
}

// UniversalDao defines API to access UniversalBo storage
type UniversalDao interface {
	// ToUniversalBo transforms godal.IGenericBo to business object.
	ToUniversalBo(gbo godal.IGenericBo) *UniversalBo

	// ToGenericBo transforms business object to godal.IGenericBo.
	ToGenericBo(ubo *UniversalBo) godal.IGenericBo

	// Delete removes the specified business object from storage.
	// This function returns true if number of deleted record is non-zero.
	Delete(bo *UniversalBo) (bool, error)

	// Create persists a new business object to storage.
	// This function returns true if number of inserted record is non-zero.
	Create(bo *UniversalBo) (bool, error)

	// Get retrieves a business object from storage.
	Get(id string) (*UniversalBo, error)

	// GetN retrieves N business objects from storage.
	GetN(fromOffset, maxNumRows int, filter interface{}, sorting interface{}) ([]*UniversalBo, error)

	// GetAll retrieves all available business objects from storage.
	GetAll(filter interface{}, sorting interface{}) ([]*UniversalBo, error)

	// Update modifies an existing business object.
	// This function returns true if number of updated record is non-zero.
	Update(bo *UniversalBo) (bool, error)

	// Save creates new business object or updates an existing one.
	// This function returns the existing record along with value true if number of inserted/updated record is non-zero.
	Save(bo *UniversalBo) (bool, *UniversalBo, error)
}
