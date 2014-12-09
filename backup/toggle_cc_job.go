package backup

import (
	"bytes"
	"fmt"

	"github.com/pivotalservices/cfops/command"
)

var ToggleCCJobRunner = func(serverUrl, username, password string, exec command.CmdExecuter) (res string, err error) {
	var b bytes.Buffer
	formatString := `curl -v -XPUT -u "%s:%s" %s --insecure -H "Content-Type:text/yaml" -i -s | grep Location: | grep Location: | cut -d ' ' -f 2`
	cmd := fmt.Sprintf(formatString, username, password, serverUrl)
	err = exec.Execute(&b, cmd)
	res = b.String()
	return
}

//func toggleCCJobs(backupscript string, ip string, username string, password string, deploymentName string, ccjobs []string, state string) {
//serverURL := "https://" + ip + ":25555/"
//for i, ccjob := range ccjobs {
//connectionURL := serverURL + "deployments/" + deploymentName + "/jobs/" + ccjob + "/" + strconv.Itoa(i) + "?state=" + state

//output, cmderr := ToggleCCJobRunner(connectionURL, username, password, command.NewLocalExecuter())
//if cmderr != nil {
//fmt.Println(fmt.Sprintf("%s", cmderr))
//os.Exit(1)
//}

//updatedURL := strings.Replace(output, "https://"+ip+"/", serverURL, 1)
//fmt.Println(fmt.Sprintf("Fetching task status from %s", updatedURL))

//contents, taskerr := getTaskEvents("GET", updatedURL, username, password, false)
//if taskerr != nil {
//fmt.Printf("%s", taskerr)
//os.Exit(1)
//}

//var eventObject utils.EventObject
//utils.GetJSONObject(contents, &eventObject)
//if eventObject.State != "done" {
//fmt.Println(fmt.Sprintf("Attempting to change state to %s for cloud controller %s instance %v", state, ccjob, i))
//}

//for eventObject.State != "done" {
//fmt.Printf(".")
//time.Sleep(2 * time.Second)
//contents, taskerr := getTaskEvents("GET", updatedURL, username, password, false)
//if taskerr != nil {
//fmt.Printf("%s", taskerr)
//os.Exit(1)
//}

//utils.GetJSONObject(contents, &eventObject)
//}
//fmt.Println(fmt.Sprintf("Changed state to %s for cloud controller %s instance %v", state, ccjob, i))
//}
//}

//func getTaskEvents(method string, url string, username string, password string, isYaml bool) ([]byte, error) {
//resp, err := invoke(method, url, username, password, false)
//if err != nil {
//fmt.Printf("%s", err)
//os.Exit(1)
//}

//if resp.StatusCode != 200 {
//fmt.Println("Invalid Bosh Director Credentials")
//os.Exit(1)
//}

//defer resp.Body.Close()
//contents, err := ioutil.ReadAll(resp.Body)

//return contents, err
//}
