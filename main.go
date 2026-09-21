package main

import (
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	inetdll  = windows.NewLazySystemDLL("wininet.dll")
	intopen  = inetdll.NewProc("InternetOpenA")
	incon    = inetdll.NewProc("InternetConnectA")
	openreq  = inetdll.NewProc("HttpOpenRequestA")
	sendreq  = inetdll.NewProc("HttpSendRequestA")
	closeint = inetdll.NewProc("InternetCloseHandle")
	intread  = inetdll.NewProc("InternetReadFile")
)

func main() {
	fmt.Println("Spoofing Opera task...")

	ua := "Opera Installer/1.0\x00"
	host := "work.ink\x00"
	path := "/_api/v2/callback/operaGX\x00"
	verb := "POST\x00"

	win, _, _ := intopen.Call(uintptr(unsafe.Pointer(&[]byte(ua)[0])), 1, 0, 0, 0)
	if win == 0 {
		fmt.Println("failed.")
		return
	}
	defer closeint.Call(win)

	connect, _, _ := incon.Call(win, uintptr(unsafe.Pointer(&[]byte(host)[0])), 443, 0, 0, 3, 0, 0)
	if connect == 0 {
		fmt.Println("failed.")
		return
	}
	defer closeint.Call(connect)

	var dwFlags uintptr = 0x80800000

	reqwest, _, _ := openreq.Call(connect, uintptr(unsafe.Pointer(&[]byte(verb)[0])), uintptr(unsafe.Pointer(&[]byte(path)[0])), 0, 0, 0, dwFlags, 0)
	if reqwest == 0 {
		fmt.Println("Failed to open HTTP request")
		return
	}
	defer closeint.Call(reqwest)

	headers := "Content-Type: application/json\r\nCache-Control: no-cache\r\n\x00"
	body := `{"noteligible":false}`

	success, _, _ := sendreq.Call(reqwest, uintptr(unsafe.Pointer(&[]byte(headers)[0])), uintptr(len(headers)-1), uintptr(unsafe.Pointer(&[]byte(body)[0])), uintptr(len(body)))

	if success == 0 {
		fmt.Println("req failed.")
		return
	}

	var buf [4096]byte
	var bytesRead uint32
	var finalBody string

	for {
		res, _, _ := intread.Call(reqwest, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), uintptr(unsafe.Pointer(&bytesRead)))
		if res == 0 || bytesRead == 0 {
			break
		}
		finalBody += string(buf[:bytesRead])
	}

	fmt.Println("Response:", finalBody)

	fmt.Println("Finished! Refresh the work.ink page. This window will close in 5 seconds.")
	time.Sleep(5 * time.Second)
}
