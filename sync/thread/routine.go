package thread

import (
	"Go-Fish-Kit/core/builtin"
	"runtime"
	"strconv"
	"strings"
)

func Go(x func()) {
	go RunSafe(x)
}

// Only for debug, never use it in production
func RoutineId() uint64 {

	buf := make([]byte, 64)
	//通过 runtime.Stack 方法获取栈帧信息
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]

	id, _ := strconv.ParseUint(string(idField), 10, 64)
	//若是err了，那就return 0
	return id

}

func RunSafe(fn func()) {
	defer builtin.Recover()
	fn()
}
