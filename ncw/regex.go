package ncw

import "regexp"

var (
	erValue                = regexp.MustCompile(`value="([^"]*)"`)
	erAction               = regexp.MustCompile(`action="([^"]*)"`)
	erSimGrupo             = regexp.MustCompile(`<td align="center" style="height:30px;width:70px;">([^<]*)</td>`)
	erSimBem               = regexp.MustCompile(`<td title="([^"]*)" align="left" style="height:30px;width:175px;">`)
	erSimValor             = regexp.MustCompile(`<td align="right" style="width:70px;">([^<]*)</td>`)
	erSimAssembleia        = regexp.MustCompile(`style="z-index:999;position:relative; top: -4px;">([^<]*)</span>`)
	erSimAssembleiaParticp = regexp.MustCompile(`<td align="center" style="width:85px;">([^<]*)</td>`)
	erSimPrazo             = regexp.MustCompile(
		`<span id="ctl00_Conteudo_grdSimulacao_ctl\d{2}_Label1">(\d+)</span>` +
			`/` +
			`<span id="ctl00_Conteudo_grdSimulacao_ctl\d{2}_Label2">(\d+)</span>`,
	)
	erSimParticp = regexp.MustCompile(`<td align="center" style="width:50px;">(\d+)</td>`)
	erSimParcela = regexp.MustCompile(`<td align="right" style="width:55px;">([^<]*)</td>`)
	erSimSeguro  = regexp.MustCompile(`<td align="center" style="width:60px;">([^<]*)</td>`)
	erSimPlano   = regexp.MustCompile(`<td title="([^"]*)" align="left" style="width:190px;">`)
	erSimPlanoB  = regexp.MustCompile(`<td align="left" style="width:190px;">([^<]*)</td>`)
)
