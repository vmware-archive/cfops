package system

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func FindFreePort() (int, error) {
	ln, _ := net.Listen("tcp", ":0")
	defer ln.Close()

	fmt.Fprintf(os.Stdout, "Listening for tcp traffic on %q", ln.Addr().String())

	port, err := strconv.ParseInt(ln.Addr().String()[5:], 10, 32)
	if err != nil {
		return -1, err
	}

	return int(port), nil
}
