package main

import (
	"KlaudiasGoTrace/KlaudiasGoTrace"
	"fmt"
	"github.com/sqweek/dialog"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func getPathAndTraceFile() {
	_, filename, _, _ := runtime.Caller(1)
	filepath2 := path.Join(path.Dir(filename))

	filepath, err := dialog.File().Title("Select Go file").Filter("Go Files", "go").Load()

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Running program")
		arrDir := strings.Split(strings.ReplaceAll(filepath, "\\", "/"), "/")
		arrDir = arrDir[:len(arrDir)-1]
		dir := strings.Join(arrDir[:], "/")
		KlaudiasGoTrace.SetOutputDirectory(dir)
		parsedProgram := KlaudiasGoTrace.Parse(filepath)
		cmd := exec.Command("go", "run", parsedProgram)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr

		//rescueStdout := os.Stdout
		//r, w, _ := os.Pipe()
		//cmd.Stdout = w

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		} else {
			html := "<html>" +
				"<head>" +
				"<title>Diploma Thesis</title>" +
				"<link rel='stylesheet' media='screen' href='mystyle.css'>" +
				"</head>" +
				"<body>" +
				"<script id=\"myJson\" type=\"application/json\">_</script>" +
				"<script src='js/Notiflix/dist/notiflix-aio-2.7.0.min.js'></script>" +
				"<script src='js/three.js'></script>" +
				"<script src='js/OrbitControls.js'></script>" +
				"<script src='js/font.js'></script>" +
				"<script src='js/tweakpane-2.1.1.js'></script>" +
				"<script src='js/app.js'></script>" +
				"<input type='file' accept='json' onchange='openFile(event)'></input><br>" +
				"<script> jsonFromHtml(); </script>" +
				"</body>" +
				"</html>"

			outSlice := strings.Split(string(out), "\n")
			numberLines, _ := strconv.Atoi(outSlice[len(outSlice)-1])
			outSlice = outSlice[(len(outSlice) - 2 - numberLines) : len(outSlice)-1]

			html = strings.ReplaceAll(html, "_", strings.Join(outSlice, "\n"))
			writeFile(filepath2+"/webpage/index2.html", html)
			fmt.Println("Open url: http://localhost:8080/index2.html")
			http.ListenAndServe(":8080", http.FileServer(http.Dir(filepath2+"/webpage/")))

		}
	}
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func writeFile(path string, dataStr string) {
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(dataStr)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func main() {
	getPathAndTraceFile()
}
