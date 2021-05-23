package user

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/btnguyen2k/consu/reddo"
)

const numSampleRows = 100

func initSampleRows(t *testing.T, testName string, dao UserDao) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numSampleRows; i++ {
		istr := fmt.Sprintf("%03d", i)
		_tagVersion := uint64(1337)
		_id := istr + "@local"
		_maskId := "admin" + istr
		_pwd := "mypassword"
		_displayName := "Administrator" + istr
		_isAdmin := i%7 == 0
		_email := istr + "@mydomain.com"
		_age := float64(18 + i)

		u := NewUser(_tagVersion, _id, _maskId)
		u.SetPassword(_pwd).SetDisplayName(_displayName).SetAdmin(_isAdmin)
		u.SetDataAttr("name.first", "Thanh"+istr)
		u.SetDataAttr("name.last", "Nguyen")
		u.SetDataAttr("email", _email)
		u.SetDataAttr("age", _age)
		if ok, err := dao.Create(u); err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", testName+"/Create", ok, err)
		}
	}
}

func doTestUserDaoCreateGet(t *testing.T, name string, dao UserDao) {
	_tagVersion := uint64(1337)
	_id := "admin@local"
	_maskId := "admin"
	_pwd := "mypassword"
	_displayName := "Administrator"
	_isAdmin := true
	_email := "myname@mydomain.com"
	_age := float64(35)

	user0 := NewUser(_tagVersion, _id, _maskId)
	user0.SetPassword(_pwd).SetDisplayName(_displayName).SetAdmin(_isAdmin)
	user0.SetDataAttr("name.first", "Thanh")
	user0.SetDataAttr("name.last", "Nguyen")
	user0.SetDataAttr("email", _email)
	user0.SetDataAttr("age", _age)
	if ok, err := dao.Create(user0); err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	if user1, err := dao.Get(_id); err != nil || user1 == nil {
		t.Fatalf("%s failed: nil or error %s", name+"/Get("+_id+")", err)
	} else {
		if v1, v0 := user1.GetDataAttrAsUnsafe("name.first", reddo.TypeString), "Thanh"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetDataAttrAsUnsafe("name.last", reddo.TypeString), "Nguyen"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetDataAttrAsUnsafe("email", reddo.TypeString), _email; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetDataAttrAsUnsafe("age", reddo.TypeInt), int64(_age); v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetTagVersion(), _tagVersion; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetId(), _id; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetDisplayName(), _displayName; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetMaskId(), _maskId; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetPassword(), _pwd; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.IsAdmin(), _isAdmin; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if t1, t0 := user1.GetTimeCreated(), user0.GetTimeCreated(); !t1.Equal(t0) {
			t.Fatalf("%s failed: expected %#v but received %#v", name, t0.Format(time.RFC3339), t1.Format(time.RFC3339))
		}
		if user1.GetChecksum() != user0.GetChecksum() {
			t.Fatalf("%s failed: expected %#v but received %#v", name, user0.GetChecksum(), user1.GetChecksum())
		}
	}
}

func doTestUserDaoCreateUpdateGet(t *testing.T, name string, dao UserDao) {
	_tagVersion := uint64(1337)
	_id := "admin@local"
	_maskId := "admin"
	_pwd := "mypassword"
	_displayName := "Administrator"
	_isAdmin := true
	_email := "myname@mydomain.com"
	_age := float64(35)

	user0 := NewUser(_tagVersion, _id, _maskId)
	user0.SetPassword(_pwd).SetDisplayName(_displayName).SetAdmin(_isAdmin)
	user0.SetDataAttr("name.first", "Thanh")
	user0.SetDataAttr("name.last", "Nguyen")
	user0.SetDataAttr("email", _email)
	user0.SetDataAttr("age", _age)
	if ok, err := dao.Create(user0); err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	user0.SetMaskId(_maskId + "-new").SetPassword(_pwd + "-new").SetDisplayName(_displayName + "-new").SetAdmin(!_isAdmin).SetTagVersion(_tagVersion + 3)
	user0.SetDataAttr("name.first", "Thanh2")
	user0.SetDataAttr("name.last", "Nguyen2")
	user0.SetDataAttr("email", _email+"-new")
	user0.SetDataAttr("age", _age+2)
	if ok, err := dao.Update(user0); err != nil {
		t.Fatalf("%s failed: %s", name+"/Update", err)
	} else if !ok {
		t.Fatalf("%s failed: cannot update record", name)
	}
	if user1, err := dao.Get(_id); err != nil || user1 == nil {
		t.Fatalf("%s failed: nil or error %s", name+"/Get("+_id+")", err)
	} else {
		if v1, v0 := user1.GetDataAttrAsUnsafe("name.first", reddo.TypeString), "Thanh2"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetDataAttrAsUnsafe("name.last", reddo.TypeString), "Nguyen2"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetDataAttrAsUnsafe("email", reddo.TypeString), _email+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetDataAttrAsUnsafe("age", reddo.TypeInt), int64(_age+2); v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetTagVersion(), _tagVersion+3; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetId(), _id; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetDisplayName(), _displayName+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetMaskId(), _maskId+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.GetPassword(), _pwd+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := user1.IsAdmin(), !_isAdmin; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if t1, t0 := user1.GetTimeCreated(), user0.GetTimeCreated(); !t1.Equal(t0) {
			t.Fatalf("%s failed: expected %#v but received %#v", name, t0.Format(time.RFC3339), t1.Format(time.RFC3339))
		}
		if user1.GetChecksum() != user0.GetChecksum() {
			t.Fatalf("%s failed: expected %#v but received %#v", name, user0.GetChecksum(), user1.GetChecksum())
		}
	}
}

func doTestUserDaoCreateDelete(t *testing.T, name string, dao UserDao) {
	_tagVersion := uint64(1337)
	_id := "admin@local"
	_maskId := "admin"
	_pwd := "mypassword"
	_displayName := "Administrator"
	_isAdmin := true
	_email := "myname@mydomain.com"
	_age := float64(35)

	user0 := NewUser(_tagVersion, _id, _maskId)
	user0.SetPassword(_pwd).SetDisplayName(_displayName).SetAdmin(_isAdmin)
	user0.SetDataAttr("name.first", "Thanh")
	user0.SetDataAttr("name.last", "Nguyen")
	user0.SetDataAttr("email", _email)
	user0.SetDataAttr("age", _age)
	if ok, err := dao.Create(user0); err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	if user1, err := dao.Get(_id); err != nil || user1 == nil {
		t.Fatalf("%s failed: nil or error %s", name+"/Get("+_id+")", err)
	} else if ok, err := dao.Delete(user1); !ok || err != nil {
		t.Fatalf("%s failed: not-ok or error %s", name+"/Delete("+_id+")", err)
	}

	if user1, err := dao.Get(_id); err != nil || user1 != nil {
		t.Fatalf("%s failed: not-nil or error %s", name+"/Get("+_id+")", err)
	}
}

func doTestUserDaoGetAll(t *testing.T, name string, dao UserDao) {
	initSampleRows(t, name, dao)
	userList, err := dao.GetAll(nil, nil)
	if err != nil || len(userList) != numSampleRows {
		t.Fatalf("%s failed: expected %#v but received %#v (error %s)", name+"/GetAll", numSampleRows, len(userList), err)
	}
}

func doTestUserDaoGetN(t *testing.T, name string, dao UserDao) {
	initSampleRows(t, name, dao)
	userList, err := dao.GetN(3, 5, nil, nil)
	if err != nil || len(userList) != 5 {
		t.Fatalf("%s failed: expected %#v but received %#v (error %s)", name+"/GetN", 5, len(userList), err)
	}
}
