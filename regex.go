package main

import "regexp"

var (
	erValue  = regexp.MustCompile(`value="([^"]*)"`)
	erAction = regexp.MustCompile(`action="([^"]*)"`)
	erH3     = regexp.MustCompile(`<h3>([^<]+)</h3>`)
	erEscola = regexp.MustCompile(
		`<div style="margin-left:20px;"><a href="atividades_distancia.aspx\?([^"]+)"><h4>([^<]+)</h4></a></div>`,
	)
	erProfessora = regexp.MustCompile(`<a href="atividades_distancia.aspx\?(t=\d+&a=\d+&b=\d+#\d+")>([^<]+)</a>`)
	erDocumento  = regexp.MustCompile(`<li><a href="([^"]+)"><b>([^<]+)</b><\/a>[^<]+</li>`)
)
