package ping

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var (
	wg       sync.WaitGroup
	hostChan chan string
)

func Scan(tgt string, s int, e int) {
	hostChan = make(chan string, 20)
	for i := s; i <= e; i++ {
		host := fmt.Sprintf(tgt, strconv.Itoa(i))
		wg.Add(1)
		go ping(host, &wg)
	}
	wg.Wait()
	close(hostChan)
	save()

}

func save() {
	file, err := os.Create("../output.txt")
	if err != nil {
		fmt.Printf("failed to save data")
		return
	}
	for {
		if data, ok := <-hostChan; ok != false {
			_, err2 := file.Write([]byte(data))
			if err2 != nil {
				fmt.Println(err2)
				return
			}
		} else {
			break
		}
	}
	file.Close()
}

func ping(host string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("[*]testing " + host)
	res := exec.Command("ping", "-c 3", host)
	var stdout, stderr bytes.Buffer
	res.Stdout = &stdout
	res.Stderr = &stderr
	err := res.Run()
	if err != nil {
		//fmt.Println(err)
		return
	}
	if stderr.Len() != 0 {
		fmt.Printf("[-] %s is unknown\n", host)
	} else {
		fmt.Printf("[+] %s is alive\n", host)
		hostChan <- host + "\n"
		return
	}
	return
}
