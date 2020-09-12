package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"os"
	"time"

	// "github.com/gorilla/websocket"
	govnc "github.com/mitchellh/go-vnc"
	"github.com/vmware/govmomi/find"
	"golang.org/x/net/websocket"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "    %s <path to VM>\n", os.Args[0])
		os.Exit(1)
	}
	vmPath := os.Args[1]

	// Get a connection context
	ctx := context.Background()

	// Connect to vCenter
	client, err := NewClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Find the virtual machine
	finder := find.NewFinder(client.Client)
	vm, err := finder.VirtualMachine(ctx, vmPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found machine: %s\n", vm.Name())

	// Get a webmks ticket from the server
	acquireTicket, err := vm.AcquireTicket(ctx, "webmks")
	if err != nil {
		fmt.Println(err)
		return
	}

	// create the websocket connection string from the ticket information
	url, err := url.Parse(fmt.Sprintf("wss://%s:%d/ticket/%s", acquireTicket.Host, acquireTicket.Port, acquireTicket.Ticket))
	if err != nil {
		fmt.Printf("parse url: %s\n", err)
		return
	}
	origin, err := url.Parse("http://localhost")
	if err != nil {
		fmt.Printf("parse origin: %s\n", err)
		return
	}

	// Create the websocket connection and set it to a BinaryFrame
	websocketConfig := &websocket.Config{
		Location:  url,
		Origin:    origin,
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
		Version:   websocket.ProtocolVersionHybi13,
		Protocol:  []string{"binary"},
	}
	nc, err := websocket.DialConfig(websocketConfig)
	if err != nil {
		fmt.Printf("Dial(): %s\n", err)
		return
	}
	nc.PayloadType = websocket.BinaryFrame

	// Setup the VNC connection over the websocket
	ccconfig := &govnc.ClientConfig{
		Auth:      []govnc.ClientAuth{new(govnc.ClientAuthNone)},
		Exclusive: false,
	}
	c, err := govnc.Client(nc, ccconfig)
	if err != nil {
		fmt.Printf("Client(): %s\n", err)
		return
	}
	defer c.Close()

	// Quick test to see if keystrokes go through...
	fmt.Printf("sending test command to console\n")
	lsString := "ls -l"
	for _, chr := range lsString {
		key := uint32(chr)
		SendKeySym(c, key)
	}
	SendKeySym(c, 0xFF0D)
	fmt.Printf("done\n")
}

func SendKeySym(c *govnc.ClientConn, key uint32) {
	delay := 100 * time.Millisecond
	err := c.KeyEvent(key, true)
	if err != nil {
		fmt.Printf("keydown failed: %s\n", err)
	}
	time.Sleep(delay)
	err = c.KeyEvent(key, false)
	if err != nil {
		fmt.Printf("keydown failed: %s\n", err)
	}
	time.Sleep(delay)
}
