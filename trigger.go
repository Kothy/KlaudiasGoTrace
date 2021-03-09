package KlaudiasGoTrace

import (
	"fmt"
	"github.com/sqweek/dialog"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	getPathAndTraceFile()
}

func getPathAndTraceFile() {
	filepath, err := dialog.File().Title("Select Go file").Filter("Go Files", "go").Load()

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		parsedProgram := Parse(filepath)

		cmd := exec.Command("go", "run", parsedProgram)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		} else {
			e := os.Remove(parsedProgram)
			if e != nil {
				log.Fatal(e)
			}

			e = os.Remove(strings.ReplaceAll(parsedProgram, ".go", "Trace.out"))
			if e != nil {
				log.Fatal(e)
			}
		}
	}
}
