package KlaudiasGoTrace

import (
	"bytes"
	"fmt"
	"io"
	"os"
	_ "unsafe"
)

var ErrTimeOrder = fmt.Errorf("time stamps out of order")

var revents []rawEvent

/*func main() {
	f, _ := os.Open("C:/Users/klaud/go/src/examples/fanin1Trace.out")


	goVersion,rawEvents,uintStrings,_ := readTrace(f)
	//goVersion,_,_,_ := readTrace(f)
	//revents := rawEvents
	//fmt.Println("Goversion: ", goVersion)
	//events, stacks, err:= parseEvents(goVersion, rawEvents, uintStrings)
	events, _, _:= parseEvents(goVersion, rawEvents, uintStrings)
	nameEvents(events)

	//for _, rEvent := range revents {
	//	evName := EventDescriptions[rEvent.Typ].Name
	//	rEvent.Name = evName
	//	fmt.Printf("%+v\n", rEvent)
	//}

	for _, event := range events {
		fmt.Printf("%+v\n", event)
	}

}*/

func MyReadTrace(filename string) []*Event {
	f, _ := os.Open(filename)

	goVersion,rawEvents,uintStrings,_ := readTrace(f)

	events, _, _:= parseEvents(goVersion, rawEvents, uintStrings)
	nameEvents(events)

	//for _, event := range events {
	//	fmt.Printf("%+v\n", event)
	//}

	return events
}

type Event struct {
	Name       string
	Offset     int       // Offset in input file (for debugging and error reporting)
	Typ        byte      // one of Ev*
	seq        int64     // sequence number
	Timestmp   int64     // timestamp in nanoseconds
	ParentId   int       // ParentId on which the event happened (can be one of TimerP, NetpollP, SyscallP)
	GorutineId uint64    // GorutineId on which the event happened
	stackId    uint64    // unique stack ID
	stack      []*Frame  // stack trace (can be empty)
	args       [3]uint64 // event-type-specific arguments
	strArgs    []string  // event-type-specific string args
	// linked event (can be nil), depends on event type:
	// for GCStart: the GCStop
	// for GCSTWStart: the GCSTWDone
	// for GCSweepStart: the GCSweepDone
	// for GoCreate: first GoStart of the created goroutine
	// for GoStart/GoStartLabel: the associated GoEnd, GoBlock or other blocking event
	// for GoSched/GoPreempt: the next GoStart
	// for GoBlock and other blocking events: the unblock event
	// for GoUnblock: the associated GoStart
	// for blocking GoSysCall: the associated GoSysExit
	// for GoSysExit: the next GoStart
	// for GCMarkAssistStart: the associated GCMarkAssistDone
	// for UserTaskCreate: the UserTaskEnd
	// for UserRegion: if the start region, the corresponding UserRegion end event
	Link *Event
}
// Frame is a frame in stack traces.
type Frame struct {
	PC   uint64
	Fn   string
	File string
	Line int
}

// toto je Divanov kod
const (
	EvNone              = 0  // unused
	EvBatch             = 1  // start of per-ParentId batch of events [pid, timestamp]
	EvFrequency         = 2  // contains tracer timer frequency [frequency (ticks per second)]
	EvStack             = 3  // stack [stack id, number of PCs, array of {PC, func string ID, file string ID, line}]
	EvGomaxprocs        = 4  // current value of GOMAXPROCS [timestamp, GOMAXPROCS, stack id]
	EvProcStart         = 5  // start of ParentId [timestamp, thread id]
	EvProcStop          = 6  // stop of ParentId [timestamp]
	EvGCStart           = 7  // GC start [timestamp, seq, stack id]
	EvGCDone            = 8  // GC done [timestamp]
	EvGCSTWStart        = 9  // GC mark termination start [timestamp, kind]
	EvGCSTWDone         = 10 // GC mark termination done [timestamp]
	EvGCSweepStart      = 11 // GC sweep start [timestamp, stack id]
	EvGCSweepDone       = 12 // GC sweep done [timestamp, swept, reclaimed]
	EvGoCreate          = 13 // goroutine creation [timestamp, new goroutine id, new stack id, stack id]
	EvGoStart           = 14 // goroutine starts running [timestamp, goroutine id, seq]
	EvGoEnd             = 15 // goroutine ends [timestamp]
	EvGoStop            = 16 // goroutine stops (like in select{}) [timestamp, stack]
	EvGoSched           = 17 // goroutine calls Gosched [timestamp, stack]
	EvGoPreempt         = 18 // goroutine is preempted [timestamp, stack]
	EvGoSleep           = 19 // goroutine calls Sleep [timestamp, stack]
	EvGoBlock           = 20 // goroutine blocks [timestamp, stack]
	EvGoUnblock         = 21 // goroutine is unblocked [timestamp, goroutine id, seq, stack]
	EvGoBlockSend       = 22 // goroutine blocks on chan send [timestamp, stack]
	EvGoBlockRecv       = 23 // goroutine blocks on chan recv [timestamp, stack]
	EvGoBlockSelect     = 24 // goroutine blocks on select [timestamp, stack]
	EvGoBlockSync       = 25 // goroutine blocks on Mutex/RWMutex [timestamp, stack]
	EvGoBlockCond       = 26 // goroutine blocks on Cond [timestamp, stack]
	EvGoBlockNet        = 27 // goroutine blocks on network [timestamp, stack]
	EvGoSysCall         = 28 // syscall enter [timestamp, stack]
	EvGoSysExit         = 29 // syscall exit [timestamp, goroutine id, seq, real timestamp]
	EvGoSysBlock        = 30 // syscall blocks [timestamp]
	EvGoWaiting         = 31 // denotes that goroutine is blocked when tracing starts [timestamp, goroutine id]
	EvGoInSyscall       = 32 // denotes that goroutine is in syscall when tracing starts [timestamp, goroutine id]
	EvHeapAlloc         = 33 // memstats.heap_live change [timestamp, heap_alloc]
	EvNextGC            = 34 // memstats.next_gc change [timestamp, next_gc]
	EvTimerGoroutine    = 35 // denotes timer goroutine [timer goroutine id]
	EvFutileWakeup      = 36 // denotes that the previous wakeup of this goroutine was futile [timestamp]
	EvString            = 37 // string dictionary entry [ID, length, string]
	EvGoStartLocal      = 38 // goroutine starts running on the same ParentId as the last event [timestamp, goroutine id]
	EvGoUnblockLocal    = 39 // goroutine is unblocked on the same ParentId as the last event [timestamp, goroutine id, stack]
	EvGoSysExitLocal    = 40 // syscall exit on the same ParentId as the last event [timestamp, goroutine id, real timestamp]
	EvGoStartLabel      = 41 // goroutine starts running with label [timestamp, goroutine id, seq, label string id]
	EvGoBlockGC         = 42 // goroutine blocks on GC assist [timestamp, stack]
	EvGCMarkAssistStart = 43 // GC mark assist start [timestamp, stack]
	EvGCMarkAssistDone  = 44 // GC mark assist done [timestamp]
	EvUserTaskCreate    = 45 // trace.NewContext [timestamp, internal task id, internal parent id, stack, Name string]
	EvUserTaskEnd       = 46 // end of task [timestamp, internal task id, stack]
	EvUserRegion        = 47 // trace.WithRegion [timestamp, internal task id, mode(0:start, 1:end), stack, Name string]
	EvUserLog           = 48 // trace.Log [timestamp, internal id, key string id, stack, value string]
	EvGoSend            = 49 // goroutine chan send [timestamp, stack]
	EvGoRecv            = 50 // goroutine chan recv [timestamp, stack]
	EvCount             = 51
)

// toto je Divanov kod
var EventDescriptions = [EvCount]struct {
	Name       string
	minVersion int
	Stack      bool
	Args       []string
	SArgs      []string // string arguments
}{
	EvNone:              {"None", 1005, false, []string{}, nil},
	EvBatch:             {"Batch", 1005, false, []string{"p", "ticks"}, nil}, // in 1.5 format it was {"p", "seq", "ticks"}
	EvFrequency:         {"Frequency", 1005, false, []string{"freq"}, nil},   // in 1.5 format it was {"freq", "unused"}
	EvStack:             {"Stack", 1005, false, []string{"id", "siz"}, nil},
	EvGomaxprocs:        {"Gomaxprocs", 1005, true, []string{"procs"}, nil},
	EvProcStart:         {"ProcStart", 1005, false, []string{"thread"}, nil},
	EvProcStop:          {"ProcStop", 1005, false, []string{}, nil},
	EvGCStart:           {"GCStart", 1005, true, []string{"seq"}, nil}, // in 1.5 format it was {}
	EvGCDone:            {"GCDone", 1005, false, []string{}, nil},
	EvGCSTWStart:        {"GCSTWStart", 1005, false, []string{"kindid"}, []string{"kind"}}, // <= 1.9, args was {} (implicitly {0})
	EvGCSTWDone:         {"GCSTWDone", 1005, false, []string{}, nil},
	EvGCSweepStart:      {"GCSweepStart", 1005, true, []string{}, nil},
	EvGCSweepDone:       {"GCSweepDone", 1005, false, []string{"swept", "reclaimed"}, nil}, // before 1.9, format was {}
	EvGoCreate:          {"GoCreate", 1005, true, []string{"g", "stack"}, nil},
	EvGoStart:           {"GoStart", 1005, false, []string{"g", "seq"}, nil}, // in 1.5 format it was {"g"}
	EvGoEnd:             {"GoEnd", 1005, false, []string{}, nil},
	EvGoStop:            {"GoStop", 1005, true, []string{}, nil},
	EvGoSched:           {"GoSched", 1005, true, []string{}, nil},
	EvGoPreempt:         {"GoPreempt", 1005, true, []string{}, nil},
	EvGoSleep:           {"GoSleep", 1005, true, []string{}, nil},
	EvGoBlock:           {"GoBlock", 1005, true, []string{}, nil},
	EvGoUnblock:         {"GoUnblock", 1005, true, []string{"g", "seq"}, nil}, // in 1.5 format it was {"g"}
	EvGoBlockSend:       {"GoBlockSend", 1005, true, []string{}, nil},
	EvGoBlockRecv:       {"GoBlockRecv", 1005, true, []string{}, nil},
	EvGoBlockSelect:     {"GoBlockSelect", 1005, true, []string{}, nil},
	EvGoBlockSync:       {"GoBlockSync", 1005, true, []string{}, nil},
	EvGoBlockCond:       {"GoBlockCond", 1005, true, []string{}, nil},
	EvGoBlockNet:        {"GoBlockNet", 1005, true, []string{}, nil},
	EvGoSysCall:         {"GoSysCall", 1005, true, []string{}, nil},
	EvGoSysExit:         {"GoSysExit", 1005, false, []string{"g", "seq", "ts"}, nil},
	EvGoSysBlock:        {"GoSysBlock", 1005, false, []string{}, nil},
	EvGoWaiting:         {"GoWaiting", 1005, false, []string{"g"}, nil},
	EvGoInSyscall:       {"GoInSyscall", 1005, false, []string{"g"}, nil},
	EvHeapAlloc:         {"HeapAlloc", 1005, false, []string{"mem"}, nil},
	EvNextGC:            {"NextGC", 1005, false, []string{"mem"}, nil},
	EvTimerGoroutine:    {"TimerGoroutine", 1005, false, []string{"g"}, nil}, // in 1.5 format it was {"g", "unused"}
	EvFutileWakeup:      {"FutileWakeup", 1005, false, []string{}, nil},
	EvString:            {"String", 1007, false, []string{}, nil},
	EvGoStartLocal:      {"GoStartLocal", 1007, false, []string{"g"}, nil},
	EvGoUnblockLocal:    {"GoUnblockLocal", 1007, true, []string{"g"}, nil},
	EvGoSysExitLocal:    {"GoSysExitLocal", 1007, false, []string{"g", "ts"}, nil},
	EvGoStartLabel:      {"GoStartLabel", 1008, false, []string{"g", "seq", "labelid"}, []string{"label"}},
	EvGoBlockGC:         {"GoBlockGC", 1008, true, []string{}, nil},
	EvGCMarkAssistStart: {"GCMarkAssistStart", 1009, true, []string{}, nil},
	EvGCMarkAssistDone:  {"GCMarkAssistDone", 1009, false, []string{}, nil},
	EvUserTaskCreate:    {"UserTaskCreate", 1011, true, []string{"taskid", "pid", "typeid"}, []string{"Name"}},
	EvUserTaskEnd:       {"UserTaskEnd", 1011, true, []string{"taskid"}, nil},
	EvUserRegion:        {"UserRegion", 1011, true, []string{"taskid", "mode", "typeid"}, []string{"Name"}},
	EvUserLog:           {"UserLog", 1011, true, []string{"id", "keyid"}, []string{"category", "message"}},
	EvGoSend:            {"GoSend", 1006, false, []string{"eip", "cid", "val"}, nil}, // toto si mozno uz on dolnil ??
	EvGoRecv:            {"GoRecv", 1006, false, []string{"eid", "cid", "val"}, nil}, // toto si mozno uz on dolnil ??
}

// toto je Divanov kod
type rawEvent struct {
	name string
	off   int
	typ   byte
	args  []uint64
	sargs []string
}

// toto je Divanov kod
func parseHeader(buf []byte) (int, error) {
	if len(buf) != 16 {
		return 0, fmt.Errorf("bad header length")
	}
	if buf[0] != 'g' || buf[1] != 'o' || buf[2] != ' ' ||
		buf[3] < '1' || buf[3] > '9' ||
		buf[4] != '.' ||
		buf[5] < '1' || buf[5] > '9' {
		return 0, fmt.Errorf("not a trace file")
	}
	ver := int(buf[5] - '0')
	i := 0
	for ; buf[6+i] >= '0' && buf[6+i] <= '9' && i < 2; i++ {
		ver = ver*10 + int(buf[6+i]-'0')
	}
	ver += int(buf[3]-'0') * 1000
	if !bytes.Equal(buf[6+i:], []byte(" trace\x00\x00\x00\x00")[:10-i]) {
		return 0, fmt.Errorf("not a trace file")
	}
	return ver, nil
}

// toto je Divanov kod
func readTrace(r io.Reader) (ver int, events []rawEvent, strings map[uint64]string, err error) {
	// Read and validate trace header.
	var buf [16]byte
	off, err := io.ReadFull(r, buf[:])
	if err != nil {
		err = fmt.Errorf("failed to read header: read %v, err %v", off, err)
		return
	}
	ver, err = parseHeader(buf[:])
	if err != nil {
		fmt.Println("Je tam error")
		return
	}
	switch ver {
	case 1005, 1007, 1008, 1009, 1010, 1011:
		// Note: When adding a new version, add canned traces
		// from the old version to the test suite using mkcanned.bash.
		break
	default:
		err = fmt.Errorf("unsupported trace file version %v.%v (update Go toolchain) %v", ver/1000, ver%1000, ver)
		return
	}

	// Read events.
	strings = make(map[uint64]string)
	for {
		// Read event type and number of arguments (1 byte).
		off0 := off
		var n int
		n, err = r.Read(buf[:1])
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil || n != 1 {
			err = fmt.Errorf("failed to read trace at Offset 0x%x: n=%v err=%v", off0, n, err)
			return
		}
		off += n
		typ := buf[0] << 2 >> 2
		narg := buf[0]>>6 + 1
		inlineArgs := byte(4)
		if ver < 1007 {
			narg++
			inlineArgs++
		}
		if typ == EvNone || typ >= EvCount || EventDescriptions[typ].minVersion > ver {
			err = fmt.Errorf("unknown event type %v at Offset 0x%x", typ, off0)
			return
		}
		if typ == EvString {
			// String dictionary entry [ID, length, string].
			var id uint64
			id, off, err = readVal(r, off)
			if err != nil {
				return
			}
			if id == 0 {
				err = fmt.Errorf("string at Offset %d has invalid id 0", off)
				return
			}
			if strings[id] != "" {
				err = fmt.Errorf("string at Offset %d has duplicate id %v", off, id)
				return
			}
			var ln uint64
			ln, off, err = readVal(r, off)
			if err != nil {
				return
			}
			if ln == 0 {
				err = fmt.Errorf("string at Offset %d has invalid length 0", off)
				return
			}
			if ln > 1e6 {
				err = fmt.Errorf("string at Offset %d has too large length %v", off, ln)
				return
			}
			buf := make([]byte, ln)
			var n int
			n, err = io.ReadFull(r, buf)
			if err != nil {
				err = fmt.Errorf("failed to read trace at Offset %d: read %v, want %v, error %v", off, n, ln, err)
				return
			}
			off += n
			strings[id] = string(buf)
			continue
		}
		ev := rawEvent{typ: typ, off: off0}
		if narg < inlineArgs {
			for i := 0; i < int(narg); i++ {
				var v uint64
				v, off, err = readVal(r, off)
				if err != nil {
					err = fmt.Errorf("failed to read event %v argument at Offset %v (%v)", typ, off, err)
					return
				}
				ev.args = append(ev.args, v)
			}
		} else {
			// More than inlineArgs args, the first value is length of the event in bytes.
			var v uint64
			v, off, err = readVal(r, off)
			if err != nil {
				err = fmt.Errorf("failed to read event %v argument at Offset %v (%v)", typ, off, err)
				return
			}
			evLen := v
			off1 := off
			for evLen > uint64(off-off1) {
				v, off, err = readVal(r, off)
				if err != nil {
					err = fmt.Errorf("failed to read event %v argument at Offset %v (%v)", typ, off, err)
					return
				}
				ev.args = append(ev.args, v)
			}
			if evLen != uint64(off-off1) {
				err = fmt.Errorf("event has wrong length at Offset 0x%x: want %v, got %v", off0, evLen, off-off1)
				return
			}
		}
		switch ev.typ {
		case EvUserLog: // EvUserLog records are followed by a value string of length ev.args[len(ev.args)-1]
			var s string
			s, off, err = readStr(r, off)
			ev.sargs = append(ev.sargs, s)
		}
		events = append(events, ev)
	}

	return
}

// toto je Divanov kod
func readVal(r io.Reader, off0 int) (v uint64, off int, err error) {
	off = off0
	for i := 0; i < 10; i++ {
		var buf [1]byte
		var n int
		n, err = r.Read(buf[:])
		if err != nil || n != 1 {
			return 0, 0, fmt.Errorf("failed to read trace at Offset %d: read %v, error %v", off0, n, err)
		}
		off++
		v |= uint64(buf[0]&0x7f) << (uint(i) * 7)
		if buf[0]&0x80 == 0 {
			return
		}
	}
	return 0, 0, fmt.Errorf("bad value at Offset 0x%x", off0)
}

// toto je Divanov kod
func readStr(r io.Reader, off0 int) (s string, off int, err error) {
	var sz uint64
	sz, off, err = readVal(r, off0)
	if err != nil || sz == 0 {
		return "", off, err
	}
	if sz > 1e6 {
		return "", off, fmt.Errorf("string at Offset %d is too large (len=%d)", off, sz)
	}
	buf := make([]byte, sz)
	n, err := io.ReadFull(r, buf)
	if err != nil || sz != uint64(n) {
		return "", off + n, fmt.Errorf("failed to read trace at Offset %d: read %v, want %v, error %v", off, n, sz, err)
	}
	return string(buf), off + n, nil
}

const (
	// Special ParentId identifiers:
	FakeP    = 1000000 + iota
	TimerP   // depicts timer unblocks
	NetpollP // depicts network unblocks
	SyscallP // depicts returns from syscalls
	GCP      // depicts GC state
)

// Parse events transforms raw events into events.
// It does analyze and verify per-event-type arguments.
func parseEvents(ver int, rawEvents []rawEvent, strings map[uint64]string) (events []*Event, stacks map[uint64][]*Frame, err error) {
	var ticksPerSec, lastSeq, lastTs int64
	var lastG uint64
	var lastP int
	timerGoids := make(map[uint64]bool)
	lastGs := make(map[int]uint64) // last goroutine running on ParentId
	stacks = make(map[uint64][]*Frame)
	batches := make(map[int][]*Event) // events by ParentId
	for _, raw := range rawEvents {
		desc := EventDescriptions[raw.typ]
		if desc.Name == "" {
			err = fmt.Errorf("missing description for event type %v", raw.typ)
			return
		}
		narg := argNum(raw, ver)
		if len(raw.args) != narg {
			err = fmt.Errorf("%v has wrong number of arguments at Offset 0x%x: want %v, got %v",
				desc.Name, raw.off, narg, len(raw.args))
			return
		}
		switch raw.typ {
		case EvBatch:
			lastGs[lastP] = lastG
			lastP = int(raw.args[0])
			lastG = lastGs[lastP]
			if ver < 1007 {
				lastSeq = int64(raw.args[1])
				lastTs = int64(raw.args[2])
			} else {
				lastTs = int64(raw.args[1])
			}
		case EvFrequency:
			ticksPerSec = int64(raw.args[0])
			if ticksPerSec <= 0 {
				// The most likely cause for this is tick skew on different CPUs.
				// For example, solaris/amd64 seems to have wildly different
				// ticks on different CPUs.
				err = fmt.Errorf("time stamps out of order")
				return
			}
		case EvTimerGoroutine:
			timerGoids[raw.args[0]] = true
		case EvStack:
			if len(raw.args) < 2 {
				err = fmt.Errorf("EvStack has wrong number of arguments at Offset 0x%x: want at least 2, got %v",
					raw.off, len(raw.args))
				return
			}
			size := raw.args[1]
			if size > 1000 {
				err = fmt.Errorf("EvStack has bad number of frames at Offset 0x%x: %v",
					raw.off, size)
				return
			}
			want := 2 + 4*size
			if ver < 1007 {
				want = 2 + size
			}
			if uint64(len(raw.args)) != want {
				err = fmt.Errorf("EvStack has wrong number of arguments at Offset 0x%x: want %v, got %v",
					raw.off, want, len(raw.args))
				return
			}
			id := raw.args[0]
			if id != 0 && size > 0 {
				stk := make([]*Frame, size)
				for i := 0; i < int(size); i++ {
					if ver < 1007 {
						stk[i] = &Frame{PC: raw.args[2+i]}
					} else {
						pc := raw.args[2+i*4+0]
						fn := raw.args[2+i*4+1]
						file := raw.args[2+i*4+2]
						line := raw.args[2+i*4+3]
						stk[i] = &Frame{PC: pc, Fn: strings[fn], File: strings[file], Line: int(line)}
					}
				}
				stacks[id] = stk
			}
		default:
			e := &Event{Offset: raw.off, Typ: raw.typ, ParentId: lastP, GorutineId: lastG}
			var argOffset int
			if ver < 1007 {
				e.seq = lastSeq + int64(raw.args[0])
				e.Timestmp = lastTs + int64(raw.args[1])
				lastSeq = e.seq
				argOffset = 2
			} else {
				e.Timestmp = lastTs + int64(raw.args[0])
				argOffset = 1
			}
			lastTs = e.Timestmp
			for i := argOffset; i < narg; i++ {
				if i == narg-1 && desc.Stack {
					e.stackId = raw.args[i]
				} else {
					e.args[i-argOffset] = raw.args[i]
				}
			}
			switch raw.typ {
			case EvGoStart, EvGoStartLocal, EvGoStartLabel:
				lastG = e.args[0]
				e.GorutineId = lastG
				if raw.typ == EvGoStartLabel {
					e.strArgs = []string{strings[e.args[2]]}
				}
			case EvGCSTWStart:
				e.GorutineId = 0
				switch e.args[0] {
				case 0:
					e.strArgs = []string{"mark termination"}
				case 1:
					e.strArgs = []string{"sweep termination"}
				default:
					err = fmt.Errorf("unknown STW kind %d", e.args[0])
					return
				}
			case EvGCStart, EvGCDone, EvGCSTWDone:
				e.GorutineId = 0
			case EvGoEnd, EvGoStop, EvGoSched, EvGoPreempt,
				EvGoSleep, EvGoBlock, EvGoBlockSend, EvGoBlockRecv,
				EvGoBlockSelect, EvGoBlockSync, EvGoBlockCond, EvGoBlockNet,
				EvGoSysBlock, EvGoBlockGC:
				lastG = 0
			case EvGoSysExit, EvGoWaiting, EvGoInSyscall:
				e.GorutineId = e.args[0]
			case EvUserTaskCreate:
				// e.args 0: taskID, 1:parentID, 2:nameID
				e.strArgs = []string{strings[e.args[2]]}
			case EvUserRegion:
				// e.args 0: taskID, 1: mode, 2:nameID
				e.strArgs = []string{strings[e.args[2]]}
			case EvUserLog:
				// e.args 0: taskID, 1:keyID, 2: stackID
				e.strArgs = []string{strings[e.args[1]], raw.sargs[0]}
			}
			batches[lastP] = append(batches[lastP], e)
		}
	}
	if len(batches) == 0 {
		err = fmt.Errorf("trace is empty")
		return
	}
	if ticksPerSec == 0 {
		err = fmt.Errorf("no EvFrequency event")
		return
	}

	//if BreakTimestampsForTesting {
	//	var batchArr [][]*Event
	//	for _, batch := range batches {
	//		batchArr = append(batchArr, batch)
	//	}
	//	for i := 0; i < 5; i++ {
	//		batch := batchArr[rand.Intn(len(batchArr))]
	//		batch[rand.Intn(len(batch))].Timestmp += int64(rand.Intn(2000) - 1000)
	//	}
	//}

	if ver < 1007 {
		events, err = order1005(batches)
	} else {
		events, err = order1007(batches)
	}
	if err != nil {
		return
	}

	// Translate cpu ticks to real time.
	minTs := events[0].Timestmp
	// Use floating point to avoid integer overflows.
	freq := 1e9 / float64(ticksPerSec)
	for _, ev := range events {
		ev.Timestmp = int64(float64(ev.Timestmp-minTs) * freq)
		// Move timers and syscalls to separate fake Ps.
		if timerGoids[ev.GorutineId] && ev.Typ == EvGoUnblock {
			ev.ParentId = TimerP
		}
		if ev.Typ == EvGoSysExit {
			ev.ParentId = SyscallP
		}
	}

	return
}

// argNum returns total number of args for the event accounting for timestamps,
// sequence numbers and differences between trace format versions.
func argNum(raw rawEvent, ver int) int {
	desc := EventDescriptions[raw.typ]
	if raw.typ == EvStack {
		return len(raw.args)
	}
	narg := len(desc.Args)
	if desc.Stack {
		narg++
	}
	switch raw.typ {
	case EvBatch, EvFrequency, EvTimerGoroutine:
		if ver < 1007 {
			narg++ // there was an unused arg before 1.7
		}
		return narg
	}
	narg++ // timestamp
	if ver < 1007 {
		narg++ // sequence
	}
	switch raw.typ {
	case EvGCSweepDone:
		if ver < 1009 {
			narg -= 2 // 1.9 added two arguments
		}
	case EvGCStart, EvGoStart, EvGoUnblock:
		if ver < 1007 {
			narg-- // 1.7 added an additional seq arg
		}
	case EvGCSTWStart:
		if ver < 1010 {
			narg-- // 1.10 added an argument
		}
	}
	return narg
}


func nameEvents(evs []*Event){
	for _, s := range evs {
		eventName := EventDescriptions[s.Typ].Name
		s.Name = eventName
	}

}

