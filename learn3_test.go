package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Timers 可以让程序等一段时间再执行下去，如果只是想等待，建议使用 time.Sleep;
// 			Timer 有个不同的地方是，它可以取消 time.Stop()
func TestTimers(t *testing.T) {
	timer1 := time.NewTimer(2 * time.Second)

	fmt.Println("Timer 1 expired")
	<-timer1.C

	time2 := time.NewTimer(time.Second)
	go func() {
		<-time2.C
		fmt.Println("Timer 2 expired")
	}()
	stop2 := time2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}
}

// Tickers 定时器，可以隔一段时间执行一次
func TestTickers(t *testing.T) {
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

// 协程池
func Worker2(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func TestWorkerPool(t *testing.T) {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	for w := 1; w <= 3; w++ {
		go Worker2(w, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= 5; a++ {
		<-results
	}
}

// 限制速率 通过 tick 限制速率
func TestRateLimiting(*testing.T) {
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.Tick(200 * time.Millisecond)

	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}
	// 上面这个是单纯的按 200 毫秒进行限制
	// 下面这个预先有 3 个缓存，当 3 个缓存瞬间消耗完的时候，再按照 200 毫秒的速率进行限制
	burstryLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstryLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstryLimiter <- t
		}
	}()

	burstryRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstryRequests <- i
	}
	close(burstryRequests)
	for req := range burstryRequests {
		<-burstryLimiter
		fmt.Println("request", req, time.Now())
	}
}

// 原子计数
func TestAtomicCounters(t *testing.T) {
	var ops uint64

	for i := 0; i < 50; i++ {
		go func() {
			for {
				atomic.AddUint64(&ops, 1)

				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops:", opsFinal)
}

// 互斥锁
func TestMutexes(t *testing.T) {
	var state = make(map[int]int)

	var mutex = &sync.Mutex{}

	var readOps uint64
	var writeOps uint64

	for r := 0; r < 100; r++ {
		go func() {
			total := 0
			for {
				key := rand.Intn(5)
				mutex.Lock()
				total += state[key]
				mutex.Unlock()
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)
				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)

	mutex.Lock()
	fmt.Println("state:", state)
	mutex.Unlock()
}

// 有状态的协程 利用协程阻塞，实现跨多个协程同步对共享状态访问 todo
type readOp struct {
	key  int
	resp chan int
}

type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func TestStatefulGoroutines(t *testing.T) {
	var readOps uint64
	var writeOps uint64

	reads := make(chan *readOp)
	writes := make(chan *writeOp)

	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := &readOp{
					key:  rand.Intn(5),
					resp: make(chan int),
				}
				reads <- read
				<-read.resp
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := &writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool),
				}
				writes <- write
				<-write.resp
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
}

// 排序
func TestSorting(t *testing.T) {
	strs := []string{"c", "a", "b"}
	sort.Strings(strs)
	fmt.Println("strings:", strs)

	ints := []int{7, 2, 4}
	sort.Ints(ints)
	fmt.Println("ints:	", ints)

	s := sort.IntsAreSorted(ints)
	fmt.Println("Sorted:	", s)
}

// 排序函数
type byLength []string

func (s byLength) Len() int {
	return len(s)
}

func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func TestSortingByFunction(t *testing.T) {
	fruits := []string{"peach", "banana", "kiwi"}
	sort.Sort(byLength(fruits))
	fmt.Println(fruits)
}

// panic
func TestPanic(t *testing.T) {
	panic("a problem")

	_, err := os.Create("/tmp/file")
	if err != nil {
		panic(err)
	}
}

// defer
func createFile(p string) *os.File {
	fmt.Println("creating")
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}

func writeFile(f *os.File) {
	fmt.Println("writing")
	fmt.Fprintln(f, "data")
}

func closeFile(f *os.File) {
	fmt.Println("closing")
	f.Close()
}

func TestDefer(t *testing.T) {
	f := createFile("/tmp/defer.txt")
	defer closeFile(f)
	writeFile(f)
}

// 集合函数
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func TestCollectionFunctions(t *testing.T) {
	var strs = []string{"peach", "apple", "pear", "plum"}

	fmt.Println(Index(strs, "pear"))

	fmt.Println(Include(strs, "grape"))

	fmt.Println(Any(strs, func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))

	fmt.Println(All(strs, func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))

	fmt.Println(Filter(strs, func(v string) bool {
		return strings.Contains(v, "e")
	}))

	fmt.Println(Map(strs, strings.ToUpper))
}

// 字符串函数
func TestStringFunctions(t *testing.T) {
	var p = fmt.Println

	p("Contains:	", strings.Contains("test", "es"))
	p("Count:		", strings.Count("test", "t"))
	p("HasPrefix:	", strings.HasPrefix("test", "te"))
	p("HasSuffix:	", strings.HasSuffix("test", "st"))
	p("Index:		", strings.Index("test", "e"))
	p("Join:		", strings.Join([]string{"a", "b"}, "-"))
	p("Repeat:		", strings.Repeat("a", 5))
	p("Replace:		", strings.Replace("fool", "o", "0", -1))
	p("Replace:		", strings.Replace("fool", "o", "0", 1))
	p("Split:		", strings.Split("a-b-c-d-e", "-"))
	p("ToLower:		", strings.ToLower("TEST"))
	p("ToUpper:		", strings.ToUpper("test"))
	p()

	p("Len:	", len("hello"))
	p("Char:", "hello"[1])
}

// 字符串格式化
type point struct {
	x, y int
}

func TestStringFormatting(t *testing.T) {
	p := point{1, 2}
	fmt.Printf("%v\n", p)

	fmt.Printf("%+v\n", p)

	fmt.Printf("%#v\n", p)

	fmt.Printf("%T\n", p)

	fmt.Printf("%t\n", true)

	fmt.Printf("%d\n", 123)

	fmt.Printf("%b\n", 14)

	fmt.Printf("%c\n", 33)

	fmt.Printf("%x\n", 456)

	fmt.Printf("%f\n", 78.9)

	fmt.Printf("%e\n", 123400000.0)
	fmt.Printf("%E\n", 123400000.0)

	fmt.Printf("%s\n", "\"string\"")

	fmt.Printf("%q\n", "\"string\"")

	fmt.Printf("%x\n", "hex this")

	fmt.Printf("%p\n", &p)

	fmt.Printf("|%6d|%6d|\n", 12, 345)

	fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)

	fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45)

	fmt.Printf("|%6s|%6s|\n", "foo", "b")

	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

	s := fmt.Sprintf("a %s", "string")
	fmt.Println(s)

	fmt.Fprintf(os.Stderr, "an %s\n", "error")
}
