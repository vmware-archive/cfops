package plugin

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync/atomic"

	"github.com/xchapter7x/lo"
)

//IsPluginMetaCall - determines if the call to the plugin is for meta data or execution
var IsPluginMetaCall atomic.Value

func init() {
	IsPluginMetaCall.Store(func() bool {
		return len(PluginArgs) == 2 && PluginArgs[1] == PluginMeta
	})
}

//Start - call on a plugin within a main to make your app a cfops plugin
func Start(plgn Plugin) {

	if IsPluginMetaCall.Load().(func() bool)() {
		b, _ := json.Marshal(plgn.GetMeta())
		UIOutput(string(b))

	} else {
		err := rpc.RegisterName(plgn.GetMeta().Name, wrapPlugin(plgn))

		if err != nil {
			lo.G.Fatal("register err: ", err)
		}
		rpc.HandleHTTP()
		l, e := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", PluginPort))

		if e != nil {
			lo.G.Fatal("listen error:", e)
		}
		http.Serve(l, nil)
	}
}

//NewPivotalCF - creates a pivotalcf from a given object and registers the gob
func NewPivotalCF(pcf PivotalCF) *PivotalCF {
	gob.Register(pcf)
	return &pcf
}

func wrapPlugin(p Plugin) *wrappedPlugin {
	return &wrappedPlugin{plugin: p}
}

func (s *wrappedPlugin) Run(pcf PivotalCF, args *[]string) (err error) {
	return s.plugin.Run(pcf, &os.Args)
}
