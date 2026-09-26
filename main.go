package main

import (
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows internet API pointers
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
	fmt.Println("Opera Spoofer - by skipped.lol")
	fmt.Println("@jiggey on Discord")
	time.Sleep(1 * time.Second)

	fmt.Println("\n\nSending request #1")

	// build first req (should work for most people)
	ua := "Opera Installer/1.0\x00"
	host := "work.ink\x00"
	path := "/_api/v2/callback/operaGX\x00"
	reqmeth := "POST\x00"

	headers := "Content-Type: application/json\r\nCache-Control: no-cache\r\n\x00"
	body := `{"noteligible":true}`

	req1bod := doReq(ua, host, path, reqmeth, headers, body)

	fmt.Println("Request 1 Response:", req1bod)

	// build second req (force req)
	fmt.Println("\n\nSending request #2")

	host = "gig.workink.click\x00"
	path = "/v2/callback/operaGX?noteligible=1\x00"
	headers = "Accept: */*\r\nConnection: keep-alive\r\n\x00"
	reqmeth = "GET\x00"

	req2bod := doReq(ua, host, path, reqmeth, headers, "")
	fmt.Println("Request 2 Response:", req2bod)

	// build third req (force req)
	fmt.Println("\n\nSending request #3")

	host = "gig.workink.click\x00"
	path = "/v2/callback/operaGX\x00"
	headers = "Accept: */*\r\nConnection: keep-alive\r\n\x00"
	reqmeth = "GET\x00"

	req3bod := doReq(ua, host, path, reqmeth, headers, "")
	fmt.Println("Request 3 Response:", req3bod)

	fmt.Println("Finished! Refresh the work.ink page. This window will close in 5 seconds.")
	fmt.Println("If the work.ink page has not changed, please refresh and try again.")
	time.Sleep(5 * time.Second)
}

func doReq(ua string, host string, path string, reqmeth string, headers string, body string) string {
	uaBytes := []byte(ua)
	hostBytes := []byte(host)
	pathBytes := []byte(path)
	methodBytes := []byte(reqmeth)
	headerBytes := []byte(headers)
	bodyBytes := []byte(body)

	win, _, _ := intopen.Call(uintptr(unsafe.Pointer(&uaBytes[0])), 1, 0, 0, 0)
	if win == 0 {
		fmt.Println("failed.")
		return ""
	}
	defer closeint.Call(win)

	connect, _, _ := incon.Call(win, uintptr(unsafe.Pointer(&hostBytes[0])), 443, 0, 0, 3, 0, 0)
	if connect == 0 {
		fmt.Println("failed.")
		return ""
	}
	defer closeint.Call(connect)

	var dwFlags uintptr = 0x80800000

	reqwest, _, _ := openreq.Call(connect, uintptr(unsafe.Pointer(&methodBytes[0])), uintptr(unsafe.Pointer(&pathBytes[0])), 0, 0, 0, dwFlags, 0)
	if reqwest == 0 {
		fmt.Println("failed.")
		return ""
	}
	defer closeint.Call(reqwest)

	var bodyPtr uintptr
	if len(bodyBytes) > 0 {
		bodyPtr = uintptr(unsafe.Pointer(&bodyBytes[0]))
	}

	success, _, _ := sendreq.Call(reqwest, uintptr(unsafe.Pointer(&headerBytes[0])), uintptr(len(headerBytes)-1), bodyPtr, uintptr(len(bodyBytes)))

	if success == 0 {
		fmt.Println("req failed.")
		return ""
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

	return finalBody
}
