package main

const (
	//CfopsHelpTemplate holds the help structure of CLI
	CfopsHelpTemplate = `
NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if .Version}}
VERSION:
   {{.Version}}
   {{end}}{{if len .Authors}}
AUTHOR(S):
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .Flags}}
   {{end}}
`
)
