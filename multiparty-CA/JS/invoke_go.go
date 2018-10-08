package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)
const BUFNUM=3
func invokejs(line [BUFNUM]string) {
	cmd := "node"
	args := []string{"invoke.js"}
	for _, element := range line{
	args = append(args, element)
	}
	process := exec.Command(cmd, args...)
	stdin, err := process.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()
	buf := new(bytes.Buffer) // THIS STORES THE NODEJS OUTPUT
	process.Stdout = buf
	process.Stderr = os.Stderr

	if err = process.Start(); err != nil {
		fmt.Println("An error occured: ", err)
	}

	process.Wait()
	fmt.Println("Generated string:", buf)
}

func main() {
	var sendbuf [BUFNUM]string
	bufcount := 0
	out, err := exec.Command("wc", "-l", "./url.dat").Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	num, err := strconv.Atoi(strings.Split(string(out), " ")[0])
	fmt.Println(num)
	bar := pb.StartNew(num)
	f, err := os.Open("url.dat")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	line, err := r.ReadString('\n') // line defined once
	sendbuf[bufcount] = strings.TrimSuffix(line, "\n")
	bufcount ++
	for err != io.EOF {
		line, err = r.ReadString('\n') // line defined once
		sendbuf[bufcount] = strings.TrimSuffix(line, "\n")

		bufcount ++
		bar.Increment()
		// fmt.Print(line)              // or any stuff
		if bufcount%3 == 0 {
			bufcount = 0
			invokejs(sendbuf)
		}

	}
	invokejs(sendbuf)
	bar.FinishPrint("The End!")
	return

}
