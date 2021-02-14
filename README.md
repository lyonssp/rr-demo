# rr-demo
This sample program randomly generates 10 numbers in parallel and prints them each to standard out.

See below show how to use rr to record a trace of a single program execution and replay it with delve.

### Build for debugging

```
$ go build -gcflags "all=-N -l"
```

### Generating a recording

```
$ rr record rr-demo
rr: Saving execution to trace directory `/home/user/.local/share/rr/rr-demo-0'.
21
81
85
70
61
56
30
71
82
42
```

### Replay the program

start the replay with delve and execute the replay
```
$ dlv replay ~/.local/share/rr/rr-demo-0
Type 'help' for list of commands.
(dlv) c
21
81
85
70
61
56
30
71
82
42
Process 4761 has exited with status 0
```

Notice how the non-deterministic logic is replicated exactly.

### Navigating the replay

Start the replay with delve and set a breakpoint for main
```
$ dlv replay ~/.local/share/rr/rr-demo-0
Type 'help' for list of commands.
(dlv) b main.main
Breakpoint 1 set at 0x4cf2db for main.main() ./main.go:11
(dlv) c
> main.main() ./main.go:11 (hits goroutine(1):1 total:1) (PC: 0x4cf2db)
Current event: 417
     6:		"sync"
     7:		"time"
     8:	)
     9:
    10:	// generate 10 random numbers in parallel and print them to stdout
=>  11:	func main() {
    12:		rand.Seed(time.Now().UnixNano())
    13:
    14:		n := 10
    15:		ch := produce(n)
    16:
```

If you stop the program on line 30 you can stop the program inside the goroutines responsible
for generating and sharing random numbers
```
(dlv) b main.go:30
Breakpoint 1 set at 0x4cf654 for main.produce.func1.1() ./main.go:30
(dlv) c
> main.produce.func1.1() ./main.go:30 (hits goroutine(16):1 total:1) (PC: 0x4cf654)
Current event: 421
    25:			var wg sync.WaitGroup
    26:			wg.Add(n)
    27:			for i := 0; i < n; i++ {
    28:				go func() {
    29:					x := rand.Intn(100)
=>  30:					ch <- x
    31:					wg.Done()
    32:				}()
    33:			}
    34:			wg.Wait()
    35:			close(ch)
(dlv) print x
21
```

Using the `rev` command we can step through the program execution backwards
```
(dlv) rev n
> main.produce.func1.1() ./main.go:29 (PC: 0x4cf63d)
Current event: 421
    24:		go func() {
    25:			var wg sync.WaitGroup
    26:			wg.Add(n)
    27:			for i := 0; i < n; i++ {
    28:				go func() {
=>  29:					x := rand.Intn(100)
    30:					ch <- x
    31:					wg.Done()
    32:				}()
    33:			}
    34:			wg.Wait()
(dlv) print x
Command failed: could not find symbol value for x
```

Imagine being able to stop a program's execution at a specific line where a program panics and then step back through the program to retrace how the program arrived at that panic statement.
