package KlaudiasGoTrace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/trace"
	"strconv"
	"strings"
	"sync"
	"time"
)

var traceFileName string
var commands []*Command
var channels []*Chan
var mutex = &sync.Mutex{}

type Chan struct {
	Name string
	Ch   chan int
}

type Command struct {
	Time     int64       `json:"Time"`
	Command  string      `json:"Command"`
	Id       int64       `json:"Id"`
	ParentId int64       `json:"ParentId"`
	From     string      `json:"From"`
	To       string      `json:"To"`
	Channel  string      `json:"Channel"`
	Value    interface{} `json:"Value"`
	EventID  string      `json:"EventID"`
	Duration int64       `json:"Duration"`
}

func StartTrace() {
	traceFileName = ""
	_, file, _, ok := runtime.Caller(1)
	if ok {
		traceFileName = file
		arr := strings.Split(traceFileName, "/")
		filename := strings.Split(arr[len(arr)-1], ".")[0]
		traceFileName = filename + "Trace"

	} else {
		traceFileName = "trace"
	}

	f, err := os.Create(traceFileName + ".out")

	if err != nil {
		panic(err)
	}
	_ = trace.Start(f)
	StartGoroutine(0)
}

func Log(tag string, message string) {
	ctx := context.Background()
	trace.Log(ctx, tag, message)
}

func SendToChannel(value interface{}, channel chan int) { // prerobit aby to zvladlo aj <-chan (read only) alebo chan<- (write only)
	chanName := findChannel(channel)
	if chanName == "" {
		chanName = createChannel(channel)
	}
	//fmt.Println(channels)
	Log(fmt.Sprintf("%v", value)+"_"+chanName, "GoroutineSend")
}

func findChannel(ch chan int) string {
	mutex.Lock()
	for _, val := range channels {
		if val.Ch == ch {
			mutex.Unlock()
			return val.Name
		}
	}
	mutex.Unlock()
	return ""
}

func createChannel(ch chan int) string {
	mutex.Lock()
	chanName := randomString()
	for isName(chanName) {
		chanName = randomString()
	}
	channels = append(channels, &Chan{
		Name: chanName,
		Ch:   ch,
	})
	mutex.Unlock()
	return chanName
}

func isName(value string) bool {
	for _, val := range channels {
		if val.Name == value {
			return true
		}
	}
	return false
}

func ReceiveFromChannel(value interface{}, channel chan int) { // prerobit aby to zvladlo aj <-chan (read only) alebo chan<- (write only)
	chanName := findChannel(channel)
	Log(fmt.Sprintf("%v", value)+"_"+chanName, "GoroutineReceive")

}

func randomString() string {
	length := 10
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func StartGoroutine(parentId uint64) {
	parentIdStr := strconv.FormatUint(parentId, 10)
	Log(parentIdStr, "GoroutineStart")
}

func StopGoroutine() {
	Log("END", "GoroutineEnd")
}

func endGoroutinesManually(goroutines map[int64]bool, commands []*Command, mainEndCmd Command) []*Command {
	for gID, isEnded := range goroutines {
		if !isEnded {
			cmd := Command{
				Time:     mainEndCmd.Time,
				Command:  "GoroutineEnd",
				Id:       gID,
				ParentId: 0,
			}
			commands = append(commands, &cmd)
		}
	}
	return commands
}

func toJson(events []*Event) {
	goroutines := make(map[int64]bool)
	gParents := make(map[int64]int64)
	var mainEndCmd Command
	//fmt.Println(channels)
	for _, chann := range channels {
		fmt.Printf("Kanal: %p\n", chann.Name)
	}

	for _, event := range events {
		if event.Name == "UserLog" {
			fmt.Printf("%+v\n", event)
			parentId, _ := strconv.Atoi(event.strArgs[0])
			comm := event.strArgs[1]
			if strings.Contains(comm, "GoroutineStart") || strings.Contains(comm, "GoroutineEnd") {
				parentId64 := int64(parentId)
				cmd := Command{
					Time:     event.Timestmp,
					Command:  event.strArgs[1],
					Id:       int64(event.GorutineId),
					ParentId: parentId64,
				}
				commands = append(commands, &cmd)
				if event.strArgs[1] == "GoroutineStart" {
					goroutines[int64(event.GorutineId)] = false
					gParents[int64(event.GorutineId)] = int64(event.ParentId)
				} else if event.strArgs[1] == "GoroutineEnd" {
					goroutines[int64(event.GorutineId)] = true

					if event.GorutineId == uint64(1) {
						lastCmd := len(commands) - 1
						commands[lastCmd].ParentId = 0
						mainEndCmd = cmd
					}
				}
			} else if strings.Contains(comm, "GoroutineSend") || strings.Contains(comm, "GoroutineReceive") {
				args := strings.Split(event.strArgs[0], "_")
				value := args[0]
				//fmt.Println(value == "<nil>")
				if value == "<nil>" {
					value = "nil"
				}
				channel := args[1]
				cmd := Command{
					Time:    event.Timestmp,
					Command: event.strArgs[1],
					Id:      int64(event.GorutineId),
					Channel: channel,
					Value:   value,
				}
				commands = append(commands, &cmd)
			}

		}
	}

	commands = endGoroutinesManually(goroutines, commands, mainEndCmd)
	prettyJSON, err := json.MarshalIndent(commands, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	writeJson(string(prettyJSON))
}

func writeJson(json string) {
	f, err := os.Create(traceFileName + ".json")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(json)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func EndTrace() {
	StopGoroutine()
	trace.Stop()

	events := MyReadTrace(traceFileName + ".out")

	toJson(events)
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
