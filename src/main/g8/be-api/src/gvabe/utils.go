package gvabe

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/consu/semita"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	systemGroupId = "administrator"

	systemAdminUsername = "admin"
	systemAdminName     = "Adam Local"
)

func encryptPassword(username, rawPassword string) string {
	saltAndPwd := username + "." + rawPassword
	out := sha1.Sum([]byte(saltAndPwd))
	return strings.ToLower(hex.EncodeToString(out[:]))
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

/*
randomString generates a random string with specified length.
*/
func randomString(l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

/*----------------------------------------------------------------------*/

var muxSytemInfo sync.Mutex
var systemInfoArr = make([]map[string]interface{}, 0)

func startUpdateSystemInfo() {
	for {
		go doUpdateSystemInfo()
		<-time.After(10 * time.Second)
	}
}

func lastSystemInfo() map[string]interface{} {
	muxSytemInfo.Lock()
	defer muxSytemInfo.Unlock()
	return systemInfoArr[len(systemInfoArr)-1]
}

func doUpdateSystemInfo() {
	muxSytemInfo.Lock()
	defer muxSytemInfo.Unlock()

	data := make(map[string]interface{})
	{
		load, err := load.Avg()
		cpuLoad := -1.0
		if err == nil && load != nil {
			cpuLoad = math.Floor(load.Load1*100) / 100
		}
		historyLoad := make([]float64, 0)
		for _, data := range systemInfoArr {
			s := semita.NewSemita(data)
			load, _ := s.GetValueOfType("cpu.load", reddo.TypeFloat)
			historyLoad = append(historyLoad, load.(float64))
		}
		data["cpu"] = map[string]interface{}{
			"cores":        runtime.NumCPU(),
			"load":         cpuLoad,
			"history_load": historyLoad,
		}
	}

	{
		mem, err := mem.VirtualMemory()
		memFree := uint64(0)
		if err == nil && mem != nil {
			memFree = mem.Free
		}
		historyFree := make([]uint64, 0)
		historyFreeKb := make([]float64, 0)
		historyFreeMb := make([]float64, 0)
		historyFreeGb := make([]float64, 0)
		for _, data := range systemInfoArr {
			s := semita.NewSemita(data)
			free, _ := s.GetValueOfType("memory.free", reddo.TypeUint)
			historyFree = append(historyFree, free.(uint64))
			freeKb, _ := s.GetValueOfType("memory.freeKb", reddo.TypeFloat)
			historyFreeKb = append(historyFreeKb, freeKb.(float64))
			freeMb, _ := s.GetValueOfType("memory.freeMb", reddo.TypeFloat)
			historyFreeMb = append(historyFreeMb, freeMb.(float64))
			freeGb, _ := s.GetValueOfType("memory.freeGb", reddo.TypeFloat)
			historyFreeGb = append(historyFreeGb, freeGb.(float64))
		}
		data["memory"] = map[string]interface{}{
			"free":           memFree,
			"freeKb":         math.Floor(100.0*(float64(memFree)/1024)) / 100.0,
			"freeMb":         math.Floor(100.0*(float64(memFree)/1024/1024)) / 100.0,
			"freeGb":         math.Floor(100.0*(float64(memFree)/1024/1024/1024)) / 100.0,
			"history_free":   historyFree,
			"history_freeKb": historyFreeKb,
			"history_freeMb": historyFreeMb,
			"history_freeGb": historyFreeGb,
		}
	}

	{
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		historyUsed := make([]uint64, 0)
		historyUsedKb := make([]float64, 0)
		historyUsedMb := make([]float64, 0)
		historyUsedGb := make([]float64, 0)
		for _, data := range systemInfoArr {
			s := semita.NewSemita(data)
			used, _ := s.GetValueOfType("app_memory.used", reddo.TypeUint)
			historyUsed = append(historyUsed, used.(uint64))
			usedKb, _ := s.GetValueOfType("app_memory.usedKb", reddo.TypeFloat)
			historyUsedKb = append(historyUsedKb, usedKb.(float64))
			usedMb, _ := s.GetValueOfType("app_memory.usedMb", reddo.TypeFloat)
			historyUsedMb = append(historyUsedMb, usedMb.(float64))
			usedGb, _ := s.GetValueOfType("app_memory.usedGb", reddo.TypeFloat)
			historyUsedGb = append(historyUsedGb, usedGb.(float64))
		}
		data["app_memory"] = map[string]interface{}{
			"used":           m.Alloc,
			"usedKb":         math.Floor(100.0*(float64(m.Alloc)/1024)) / 100.0,
			"usedMb":         math.Floor(100.0*(float64(m.Alloc)/1024/1024)) / 100.0,
			"usedGb":         math.Floor(100.0*(float64(m.Alloc)/1024/1024/1024)) / 100.0,
			"history_used":   historyUsed,
			"history_usedKb": historyUsedKb,
			"history_usedMb": historyUsedMb,
			"history_usedGb": historyUsedGb,
		}
	}

	{
		history := make([]int, 0)
		for _, data := range systemInfoArr {
			s := semita.NewSemita(data)
			n, _ := s.GetValueOfType("go_routines.num", reddo.TypeInt)
			history = append(history, int(n.(int64)))
		}
		data["go_routines"] = map[string]interface{}{
			"num":     runtime.NumGoroutine(),
			"history": history,
		}
	}

	systemInfoArr = append(systemInfoArr, data)
	if len(systemInfoArr) > 10 {
		systemInfoArr[0] = nil
		systemInfoArr = systemInfoArr[1:]
	}
}
