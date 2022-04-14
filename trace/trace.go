package trace

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var goroutineSpace = []byte("goroutine ")

func curGoroutingID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse gorouting ID out of %q: %v", b, err))
	}
	return n
}

func printTrace(gid uint64, name string, arrow string, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "    "
	}
	fmt.Printf("g[%05d]:%s%s%s\n", gid, indents, arrow, name)
}

var mu sync.Mutex
var depth = make(map[uint64]int)

func Trace() func() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}

	fn := runtime.FuncForPC(pc)
	name := fn.Name()

	gid := curGoroutingID()

	mu.Lock()
	indent := depth[gid]
	depth[gid] = indent + 1
	mu.Unlock()

	printTrace(gid, name, "->", indent+1)
	return func() {
		mu.Lock()
		indent := depth[gid]
		depth[gid] = indent - 1
		mu.Unlock()
		printTrace(gid, name, "<-", indent)
	}
}

func A1() {
	defer Trace()()
	B1()
}
func B1() {
	defer Trace()()
	C1()
}
func C1() {
	defer Trace()()
	D()
}

func D() {
	defer Trace()()
}

func A2() {
	defer Trace()()
	B2()
}
func B2() {
	defer Trace()()
	C2()
}
func C2() {
	defer Trace()()
	D()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		A2()
		wg.Done()
	}()
	A1()
	wg.Wait()
}
