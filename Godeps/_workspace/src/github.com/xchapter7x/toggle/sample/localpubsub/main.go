package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/xchapter7x/goutil/unpack"
	"github.com/xchapter7x/toggle"
	"github.com/xchapter7x/toggle/engines/localpubsub"
)

func TestA(s string) (r string) {
	r = fmt.Sprintln("testa", s)
	fmt.Println(r)
	return
}

func TestB(s string) (r string) {
	r = fmt.Sprintln("testb", s)
	fmt.Println(r)
	return
}

func main() {
	c, _ := redis.Dial("tcp", "localhost:6379")
	defer c.Close()
	c2, _ := redis.Dial("tcp", "localhost:6379")
	defer c2.Close()

	psc := redis.PubSubConn{Conn: c}
	lps := localpubsub.NewLocalPubSubEngine(psc, toggle.ShowFeatures())
	defer lps.Close()

	toggle.Init("NS", lps)
	toggle.RegisterFeature("test")
	f := toggle.Flip("test", TestA, TestB, "argstring")
	var output string
	unpack.Unpack(f, &output)
	fmt.Println(output)
	fmt.Println("publish state change")
	c2.Do("PUBLISH", "NS_test", "true")
	time.Sleep(1000 * time.Millisecond)
	f = toggle.Flip("test", TestA, TestB, "argstring")
	unpack.Unpack(f, &output)
	fmt.Println(output)
}
