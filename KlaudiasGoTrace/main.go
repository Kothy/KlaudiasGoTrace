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
var outputDir string
var commands []*Command
var channels []*Chan
var sleeps = make(map[int64]int64)
var blocks = make(map[int64]int64)
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
	Channel  string      `json:"Channel"`
	Value    interface{} `json:"Value"`
	Duration int64       `json:"Duration"`
	//From     string      `json:"From"`
	//To       string      `json:"To"`
	//EventID  string      `json:"EventID"`
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

	// ak nebude fungovat buffer odkomentuj
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

	// ak nebude fungovat buffer odkomentuj
	//events := MyReadTrace(traceFileName + ".out")
	events := DoTrace(traceBuffer)

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

func SendToChannel(value interface{}, channel interface{}) {
	if !isChannel(channel) {
		return
	}
	chanName := findChannel(channel)
	if chanName == "" {
		chanName = createChannel(channel)
	}

	Log(chanName+"_"+fmt.Sprintf("%v", value), "GoroutineSend")
}

func ReceiveFromChannel(value interface{}, channel interface{}) interface{} {
	if !isChannel(channel) {
		return nil
	}
	chanName := findChannel(channel)

	Log(chanName+"_"+fmt.Sprintf("%v", value), "GoroutineReceive")

	return value
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

func SetOutputDirectory(directory string) {
	outputDir = directory
}

func toJson(events []*Event) {
	goroutines := make(map[int64]bool)
	gParents := make(map[int64]int64)
	var mainEndCmd Command

	for _, event := range events {
		//fmt.Printf("%+v\n", event)
		gId := int64(event.GorutineId)
		_, ok := goroutines[gId]
		if event.Name == "UserLog" {
			//fmt.Printf("%+v\n", event)
			parentId, _ := strconv.Atoi(event.strArgs[0])
			command := event.strArgs[1]
			if strings.Contains(command, "GoroutineStart") || strings.Contains(command, "GoroutineEnd") {
				parentId64 := int64(parentId)
				cmd := Command{
					Time:     event.Timestmp,
					Command:  event.strArgs[1],
					Id:       gId,
					ParentId: parentId64,
				}
				commands = append(commands, &cmd)
				if event.strArgs[1] == "GoroutineStart" {
					goroutines[gId] = false
					gParents[gId] = int64(event.ParentId)
					sleeps[gId] = 0
				} else if event.strArgs[1] == "GoroutineEnd" {
					goroutines[gId] = true

					if event.GorutineId == uint64(1) {
						lastCmd := len(commands) - 1
						commands[lastCmd].ParentId = 0
						mainEndCmd = cmd
					}
				}
			} else if strings.Contains(command, "GoroutineSend") || strings.Contains(command, "GoroutineReceive") {
				channel := event.strArgs[0][:10]
				value := event.strArgs[0][11:]
				if value == "<nil>" {
					value = "nil"
				}

				cmd := Command{
					Time:    event.Timestmp,
					Command: event.strArgs[1],
					Id:      gId,
					Channel: channel,
					Value:   value,
				}
				commands = append(commands, &cmd)
			}
		} else if event.Name == "GoSleep" {
			//fmt.Printf("%+v\n", event)
			sleeps[gId] = event.Timestmp

		} else if event.Name == "GoStart" {
			//fmt.Printf("%+v\n", event)
			if sleeps[gId] > 0 && ok {
				cmd := Command{
					Time:     sleeps[gId],
					Command:  "GoroutineSleep",
					Id:       gId,
					Duration: event.Timestmp - sleeps[gId],
				}
				commands = append(commands, &cmd)
				sleeps[gId] = 0

			}
		} else if event.Name == "GoBlockRecv" || event.Name == "GoBlockSend" {
			//fmt.Printf("%+v\n", event)
			if ok {
				blocks[gId] = event.Timestmp
			}

		} else if event.Name == "GoUnblock" {
			//fmt.Printf("%+v\n", event)
			_, ok2 := blocks[gId]
			if ok && ok2 {

				cmd := Command{
					Time:     blocks[gId],
					Command:  "GoroutineBlock",
					Id:       gId,
					Duration: event.Timestmp - blocks[gId],
				}
				commands = append(commands, &cmd)
				blocks[gId] = 0
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
	_, err3 := os.Stat("jsons")

	if os.IsNotExist(err3) {
		errDir := os.MkdirAll("jsons", 0755)
		if errDir != nil {
			log.Fatal(err3)
		}
	}

	f, err := os.Create("jsons/" + traceFileName + ".json")
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
