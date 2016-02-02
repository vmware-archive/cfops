package cfopsplugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

//Server --
func (g BackupRestorePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &BackupRestoreRPCServer{Impl: g.P}, nil
}

//Client --
func (g BackupRestorePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &BackupRestoreRPC{client: c}, nil
}
