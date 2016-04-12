package load

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pivotalservices/cfbackup/tileregistry"
	"github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/xchapter7x/lo"
)

func init() {

	if dir := os.Getenv(PluginDirEnv); dir != "" {
		PluginDir = dir
	}
	Plugins(PluginDir)
}

//Plugins - function to register plugins residing in a given directory with cfops
func Plugins(dir string) (err error) {
	var fileInfoArray []os.FileInfo
	var fileInfo os.FileInfo

	if fileInfoArray, err = ioutil.ReadDir(dir); err == nil && len(fileInfoArray) > 0 {
		for _, fileInfo = range fileInfoArray {
			if err = loadPlugin(dir, fileInfo); err != nil {
				lo.G.Debug("error loading plugin: ", err, fileInfo)
				break
			}
			lo.G.Debug("loading plugin from: ", fileInfo)
		}

	} else if err != nil {
		lo.G.Debug("not loading plugins: ", err, fileInfoArray)
		err = fmt.Errorf("error loading plugins: %v %v", err, fileInfoArray)
	}
	return
}

func loadPlugin(dir string, fileInfo os.FileInfo) (err error) {
	var bytes []byte
	filePath := dir + "/" + fileInfo.Name()

	if bytes, err = exec.Command(filePath, cfopsplugin.PluginMeta).Output(); err == nil {
		meta := cfopsplugin.Meta{}
		if err = json.Unmarshal(bytes, &meta); err == nil {

			if meta.Name == "" {
				lo.G.Debug("plugin meta busted", meta)
				err = ErrInvalidPluginMeta

			} else {
				ptb := &cfopsplugin.PluginTileBuilder{
					FilePath:   filePath,
					Meta:       meta,
					CmdBuilder: cfopsplugin.DefaultCmdBuilder,
				}
				lo.G.Debug("registering plugin: ", ptb)
				tileregistry.Register(meta.Name, ptb)
			}

		} else {
			lo.G.Error("plugin load error: ", filePath, err)
			err = fmt.Errorf("plugin load error: %v %v ", filePath, err)
		}
	}
	return
}
