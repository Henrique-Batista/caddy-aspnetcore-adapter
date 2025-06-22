package adapter

import (
	"encoding/xml"
	"fmt"
	"github.com/caddyserver/caddy/v2/caddyconfig"
)

type WebConfig struct {
	XMLName         xml.Name        `xml:"configuration"`
	SystemWebServer SystemWebServer `xml:"system.webServer"`
	AspNetCore      AspNetCore      `xml:"system.webServer/aspNetCore"`
}

type SystemWebServer struct {
	Rewrite         RewriteSection  `xml:"rewrite"`
	HttpRedirect    RedirectSection `xml:"httpRedirect"`
	DefaultDocument DefaultDocument `xml:"defaultDocument"`
}

type RewriteSection struct {
	Rules []RewriteRule `xml:"rules>rule"`
}

type RewriteRule struct {
	Name   string        `xml:"name,attr"`
	Match  RewriteMatch  `xml:"match"`
	Action RewriteAction `xml:"action"`
}

type RewriteMatch struct {
	URL string `xml:"url,attr"`
}

type RewriteAction struct {
	Type string `xml:"type,attr"`
	URL  string `xml:"url,attr"`
}

type RedirectSection struct {
	Enabled     string `xml:"enabled,attr"`
	Destination string `xml:"destination,attr"`
}

type DefaultDocument struct {
	Files []DefaultFile `xml:"files>add"`
}

type DefaultFile struct {
	Name string `xml:"value,attr"`
}

type AspNetCore struct {
	ProcessPath          string   `xml:"processPath,attr"`
	Arguments            string   `xml:"arguments,attr"`
	HostingModel         string   `xml:"hostingModel,attr"`
	EnvironmentVariables []EnvVar `xml:"environmentVariables>environmentVariable"`
}

type EnvVar struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func init() {
	caddyconfig.RegisterAdapter("webconfig", Adapter{})
}

type Adapter struct{}

func (Adapter) Adapt(body []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	var cfg WebConfig
	if err := xml.Unmarshal(body, &cfg); err != nil {
		return nil, nil, fmt.Errorf("erro ao parsear web.config: %v", err)
	}

	// Início do Caddyfile
	caddy := `{
	order webconfig before file_server
}
:80 {
	root * ./wwwroot
	file_server
}`
	// rewrite rules
	for _, rule := range cfg.SystemWebServer.Rewrite.Rules {
		caddy += fmt.Sprintf("    rewrite %s %s\n", rule.Match.URL, rule.Action.URL)
	}
	// redirect rules
	if cfg.SystemWebServer.HttpRedirect.Enabled == "true" {
		caddy += fmt.Sprintf("    redir %s\n", cfg.SystemWebServer.HttpRedirect.Destination)
	}
	// default documents
	if len(cfg.SystemWebServer.DefaultDocument.Files) > 0 {
		caddy += "    try_files "
		for _, f := range cfg.SystemWebServer.DefaultDocument.Files {
			caddy += f.Name + " "
		}
		caddy += "\n"
	}
	caddy += "    file_server\n}"

	return []byte(caddy), nil, nil
}
