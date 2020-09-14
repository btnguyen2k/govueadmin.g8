package henge

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/btnguyen2k/consu/reddo"
)

func TestNewUniversalBo(t *testing.T) {
	name := "TestNewUniversalBo"
	ubo := NewUniversalBo("id", 1357)
	if ubo == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if id := ubo.GetId(); id != "id" {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, "id", id)
	}
	if appVersion := ubo.GetAppVersion(); appVersion != 1357 {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, 1357, appVersion)
	}
}

func TestRowMapper(t *testing.T) {
	name := "TestRowMapper"
	tableName := "test_user"
	extraColNameToFieldMappings := map[string]string{"zuid": "owner_id"}
	rowMapper := buildRowMapper(tableName, extraColNameToFieldMappings)

	myColList := rowMapper.ColumnsList(tableName)
	expectedColList := append(columnNames, "zuid")
	if !reflect.DeepEqual(myColList, expectedColList) {
		t.Fatalf("%s failed: expected column list %#v but received %#v", name, expectedColList, myColList)
	}
}

func TestUniversalBo_datatypes(t *testing.T) {
	name := "TestUniversalBo_datatypes"
	ubo := NewUniversalBo("id", 1357)
	vInt := 123
	ubo.SetDataAttr("data.number[0]", vInt)
	vFloat := 45.6
	ubo.SetDataAttr("data.number[1]", vFloat)
	vBool := true
	ubo.SetDataAttr("data.bool", vBool)
	vString := "a string"
	ubo.SetDataAttr("data.string", vString)
	vTime := time.Now()
	ubo.SetDataAttr("data.time[0]", vTime)
	ubo.SetDataAttr("data.time[1]", vTime.Format(TimeLayout))

	if v, err := ubo.GetDataAttrAs("data.number[0]", reddo.TypeInt); err != nil {
		t.Fatalf("%s failed: %#e", name, err)
	} else if v != int64(vInt) {
		t.Fatalf("%s failed [int]: expected %#v but received %#v", name, vInt, v)
	}
	if v, err := ubo.GetDataAttrAs("data.number[0]", reddo.TypeUint); err != nil {
		t.Fatalf("%s failed: %#e", name, err)
	} else if v != uint64(vInt) {
		t.Fatalf("%s failed [uint]: expected %#v but received %#v", name, vInt, v)
	}
	if v, err := ubo.GetDataAttrAs("data.number[1]", reddo.TypeFloat); err != nil {
		t.Fatalf("%s failed: %#e", name, err)
	} else if v != float64(vFloat) {
		t.Fatalf("%s failed [float]: expected %#v but received %#v", name, vFloat, v)
	}
	if v, err := ubo.GetDataAttrAs("data.bool", reddo.TypeBool); err != nil {
		t.Fatalf("%s failed: %#e", name, err)
	} else if v != vBool {
		t.Fatalf("%s failed [bool]: expected %#v but received %#v", name, vBool, v)
	}
	if v, err := ubo.GetDataAttrAs("data.string", reddo.TypeString); err != nil {
		t.Fatalf("%s failed: %#e", name, err)
	} else if v != vString {
		t.Fatalf("%s failed [string]: expected %#v but received %#v", name, vString, v)
	}
	if v, err := ubo.GetDataAttrAsTimeWithLayout("data.time[0]", TimeLayout); err != nil {
		t.Fatalf("%s failed: %#e", name, err)
	} else if v.Format(TimeLayout) != vTime.Format(TimeLayout) {
		t.Fatalf("%s failed [time]: expected %#v but received %#v", name, vTime, v)
	}
	if v, err := ubo.GetDataAttrAsTimeWithLayout("data.time[1]", TimeLayout); err != nil {
		t.Fatalf("%s failed: %#e", name, err)
	} else if v.Format(TimeLayout) != vTime.Format(TimeLayout) {
		t.Fatalf("%s failed [time]: expected %#v but received %#v", name, vTime, v)
	}
}

func TestUniversalBo_json(t *testing.T) {
	name := "TestUniversalBo_json"
	ubo1 := NewUniversalBo("id", 1357)
	vInt := float64(123)
	ubo1.SetDataAttr("data.number[0]", vInt)
	vFloat := 45.6
	ubo1.SetDataAttr("data.number[1]", vFloat)
	vBool := true
	ubo1.SetDataAttr("data.bool", vBool)
	vString := "a string"
	ubo1.SetDataAttr("data.string", vString)
	vTime := time.Now()
	ubo1.SetDataAttr("data.time[0]", vTime.Format(TimeLayout))
	ubo1.SetDataAttr("data.time[1]", vTime.Format(TimeLayout))
	js1, _ := json.Marshal(ubo1)

	var ubo2 *UniversalBo
	err := json.Unmarshal(js1, &ubo2)
	if err != nil {
		t.Fatalf("%s failed: %e", name, err)
	}

	if ubo1.id != ubo2.id {
		t.Fatalf("%s failed [id]: expected %#v but received %#v", name, ubo1.id, ubo2.id)
	}
	if ubo1.appVersion != ubo2.appVersion {
		t.Fatalf("%s failed [appversion]: expected %#v but received %#v", name, ubo1.appVersion, ubo2.appVersion)
	}
	if !reflect.DeepEqual(ubo1._data, ubo2._data) {
		t.Fatalf("%s failed [data]: expected\n%#v\nbut received\n%#v", name, ubo1._data, ubo2._data)
	}
	if ubo1.checksum != ubo2.checksum {
		t.Fatalf("%s failed [checksum]: expected %#v but received %#v", name, ubo1.checksum, ubo2.checksum)
	}
}
