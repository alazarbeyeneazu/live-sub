package models

type ResponseMsg struct {
	ToolName string `csv:"tool_name"`
	FQDN     string `csv:"fqdn"`
}
