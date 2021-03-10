package KlaudiasGoTrace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
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
var traceBuffer bytes.Buffer

type Chan struct {
	Name string
	Ch   interface{}
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

	//f, err := os.Create(traceFileName + ".out")
	//if err != nil {
	//	panic(err)
	//}
	_ = trace.Start(&traceBuffer)

	StartGoroutine(0)
}

func EndTrace() {
	StopGoroutine()

	trace.Stop()

	//events := MyReadTrace(traceFileName + ".out")
	events := DoTrace(traceBuffer)

	//err := os.Remove(traceFileName + ".out")
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	toJson(events)
}

func Log(tag string, message string) {
	ctx := context.Background()
	trace.Log(ctx, tag, message)
}

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func isinstanceof(value interface{}, typ string) bool {
	return reflect.TypeOf(value).String() == typ
}

func isChannel(ch interface{}) bool {
	chann := reflect.TypeOf(ch).String()
	if strings.Contains(chann, "chan") {
		return true
	}
	return false
}

//func SendToChannel(value interface{}, channel chan int) {
func SendToChannel(value interface{}, channel interface{}) {

	if !isChannel(channel) {
		return
	}
	chanName := findChannel(channel)
	if chanName == "" {
		chanName = createChannel(channel)
	}

	Log(fmt.Sprintf("%v", value)+"_"+chanName, "GoroutineSend")
}

//func findChannel(ch chan int) string {
func findChannel(ch interface{}) string {
	mutex.Lock()
	chPtrStr := fmt.Sprintf("%v", ch)
	for _, val := range channels {
		ch2PtrStr := fmt.Sprintf("%v", val.Ch)
		if val.Ch == ch || chPtrStr == ch2PtrStr {
			mutex.Unlock()
			return val.Name
		}
	}
	mutex.Unlock()
	return ""
}

func createChannel(ch interface{}) string {
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

//func ReceiveFromChannel(value interface{}, channel chan int) {
func ReceiveFromChannel(value interface{}, channel interface{}) {
	if !isChannel(channel) {
		return
	}

	//fmt.Println("Prijate z: ", channel)
	//for i := range channels {
	//	fmt.Print(channels[i].Ch, ",")
	//}
	//fmt.Println("\n-------------------------------")
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

	time.Sleep(time.Duration(20) * time.Millisecond)
}

func StopGoroutine() {
	time.Sleep(time.Duration(20) * time.Millisecond)
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

	//for _, chann := range channels {
	//	fmt.Println("Kanal: ", chann.Name)
	//}

	//fmt.Println("Pole eventov je dlzky: ", len(events))
	for _, event := range events {
		//fmt.Printf("%+v\n", event)
		if event.Name == "UserLog" {
			//fmt.Printf("%+v\n", event)
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

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
