package ncw

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/luisfernandogaido/cnc/model"
)

const (
	cpfFake  = "965.894.880-45"
	cnpjFake = "46.236.911/0001-13"
)

type Simulacao struct {
	Consulta          time.Time          `json:"consulta"`
	Pagina            int                `json:"pagina"`
	Conta             string             `json:"conta"`
	TipoPessoa        string             `json:"tipo_pessoa"`
	TipoVenda         string             `json:"tipo_venda"`
	RdbPrazo          string             `json:"rdb_prazo"`
	Grupo             string             `json:"grupo"`
	Bem               string             `json:"bem"`
	Valor             float64            `json:"valor"`
	Assembleia        time.Time          `json:"assembleia"`
	AssParticp        time.Time          `json:"ass_particp"`
	PrazoRateio       int                `json:"prazo_rateio"`
	PrazoTermo        int                `json:"prazo_termo"`
	Particp           int                `json:"particp"`
	Parcela           float64            `json:"parcela"`
	Seguro            string             `json:"seguro"`
	Plano             string             `json:"plano"`
	Vstate            string             `json:"vstate"`
	UrlForm           string             `json:"url_form"`
	Values            url.Values         `json:"values"`
	Caracteristica    string             `json:"caracteristica"`
	Vencimento        time.Time          `json:"vencimento"`
	PerMaxLan         float64            `json:"per_max_lan"`
	IntervaloParcelas []IntevaloParcelas `json:"intervalo_parcelas"`
}

type IntevaloParcelas struct {
	Ini     int     `json:"ini"`
	Fim     int     `json:"fim"`
	FunComP float64 `json:"fun_com_p"`
	FunComR float64 `json:"fun_com_r"`
	TaxAdmP float64 `json:"tax_adm_p"`
	TaxAdmR float64 `json:"tax_adm_r"`
	TaxAdsP float64 `json:"tax_ads_p"`
	TaxAdsR float64 `json:"tax_ads_r"`
	FunResP float64 `json:"fun_res_p"`
	FunResR float64 `json:"fun_res_r"`
	Seguro  float64 `json:"seguro"`
	Total   float64 `json:"total"`
}

type formSimulacao struct {
	tipoPessoa string
	cpfCnpj    string
	tipoVenda  string
}

func Simulacoes() ([]Simulacao, error) {
	consulta := time.Now()
	vendas := map[string]string{
		"001": "NACIONAL",
		"002": "COLABORADOR/PJ",
		"003": "EXCLUSIVO",
		"004": "EXCLUSIVO/ESPECIAL",
		"005": "EXCLUSIVO R1",
		"006": "FUNCIONARIO GM",
	}
	pessoas := map[string]string{
		"rbtFisica":   "FÍSICA",
		"rbtJuridica": "JURÍDICA",
	}
	formsSimulacao := []struct {
		tipoPessoa string
		cpfCnpj    string
		tipoVenda  string
	}{
		{"rbtFisica", "965.894.880-45", "001"},
		{"rbtFisica", "965.894.880-45", "002"},
		{"rbtFisica", "965.894.880-45", "003"},
		{"rbtFisica", "965.894.880-45", "004"},
		{"rbtFisica", "965.894.880-45", "005"},
		{"rbtFisica", "965.894.880-45", "006"},
		{"rbtJuridica", "46.236.911/0001-13", "001"},
		{"rbtJuridica", "46.236.911/0001-13", "002"},
		{"rbtJuridica", "46.236.911/0001-13", "003"},
		{"rbtJuridica", "46.236.911/0001-13", "004"},
		{"rbtJuridica", "46.236.911/0001-13", "005"},
		{"rbtJuridica", "46.236.911/0001-13", "006"},
	}

	usuarios, err := model.Usuarios()
	if err != nil {
		return nil, fmt.Errorf("simulacoes: %w", err)
	}
	for _, u := range usuarios {
		if err := Login(u.Usuario, u.Senha); err != nil {
			if err := model.BadLoginInsert(u.Conta, u.Usuario, u.Senha); err != nil {
				log.Printf("simulacoes: conta: %v: %v\n", u.Conta, err)
				continue
			}
			log.Printf("simulacoes: conta: %v: %v\n", u.Conta, err)
			continue
		}
		doc, urlForm, err := getRetry(urlBase + "/CONPV/frmConPvCnsVendaSimulacao.aspx?cd=31821")
		if err != nil {
			log.Printf("simulacoes: conta: %v: %v\n", u.Conta, err)
			continue
		}
		values, err := newValuesDoc(doc)
		if err != nil {
			log.Printf("simulacoes: conta: %v: %v\n", u.Conta, err)
			continue
		}
		for _, fs := range formsSimulacao {
			values.Set("ctl00$hdnID_Modulo", "VE")
			values.Set("ctl00$Conteudo$Pessoa", fs.tipoPessoa)
			values.Set("ctl00$Conteudo$edtCD_Inscricao_Nacional", fs.cpfCnpj)
			values.Set("ctl00$Conteudo$cbxProdutos", "1")
			values.Set("ctl00$Conteudo$cbxSubProdutos", "1")
			values.Set("ctl00$Conteudo$edtBusca_TipoVenda", fs.tipoVenda)
			values.Set("ctl00$Conteudo$cbxTp_Venda", fs.tipoVenda[2:])
			values.Set("ctl00$Conteudo$cbxNegociacao", "RT")
			values.Set("ctl00$Conteudo$btnSimular2", "Simular")
			values.Set("__SCROLLPOSITIONX", "0")
			values.Set("__SCROLLPOSITIONY", "0")
			doc, urlForm, err = postRetry(urlForm, values)
			simPgs, err := simulacoesPaginas(urlForm, values)
			if err != nil {
				log.Printf("simulacoes: conta: %v: %v\n", u.Conta, err)
				continue
			}
			fmt.Println(time.Now().Format("15:04:05"), "len(simPgs):", len(simPgs))
			sims := make([]model.SimulacaoNcw, len(simPgs))
			for i := range simPgs {
				simPgs[i].Consulta = consulta
				simPgs[i].Conta = u.Conta
				simPgs[i].TipoPessoa = pessoas[fs.tipoPessoa]
				simPgs[i].TipoVenda = vendas[fs.tipoVenda]
				simPgs[i].Values = values

				if err := simulacoesDetalhes(&simPgs[i]); err != nil {
					log.Printf("simulacoes: conta: %v: %v\n", u.Conta, err)
					continue
				}

				sims[i].Consulta = simPgs[i].Consulta
				sims[i].Pagina = simPgs[i].Pagina
				sims[i].Conta = simPgs[i].Conta
				sims[i].TipoPessoa = simPgs[i].TipoPessoa
				sims[i].TipoVenda = simPgs[i].TipoVenda
				sims[i].RdbPrazo = simPgs[i].RdbPrazo
				sims[i].Grupo = simPgs[i].Grupo
				sims[i].Bem = simPgs[i].Bem
				sims[i].Valor = simPgs[i].Valor
				sims[i].Assembleia = simPgs[i].Assembleia
				sims[i].AssParticp = simPgs[i].AssParticp
				sims[i].PrazoRateio = simPgs[i].PrazoRateio
				sims[i].PrazoTermo = simPgs[i].PrazoTermo
				sims[i].Particp = simPgs[i].Particp
				sims[i].Parcela = simPgs[i].Parcela
				sims[i].Seguro = simPgs[i].Seguro
				sims[i].Plano = simPgs[i].Plano
				sims[i].Caracteristica = simPgs[i].Caracteristica
				sims[i].Vencimento = simPgs[i].Vencimento
				sims[i].PerMaxLan = simPgs[i].PerMaxLan
				for k := range simPgs[i].IntervaloParcelas {
					ip := model.IntevaloParcelas{
						Ini:     simPgs[i].IntervaloParcelas[k].Ini,
						Fim:     simPgs[i].IntervaloParcelas[k].Fim,
						FunComP: simPgs[i].IntervaloParcelas[k].FunComP,
						FunComR: simPgs[i].IntervaloParcelas[k].FunComR,
						TaxAdmP: simPgs[i].IntervaloParcelas[k].TaxAdmP,
						TaxAdmR: simPgs[i].IntervaloParcelas[k].TaxAdmR,
						TaxAdsP: simPgs[i].IntervaloParcelas[k].TaxAdsP,
						TaxAdsR: simPgs[i].IntervaloParcelas[k].TaxAdsR,
						FunResP: simPgs[i].IntervaloParcelas[k].FunResP,
						FunResR: simPgs[i].IntervaloParcelas[k].FunResR,
						Seguro:  simPgs[i].IntervaloParcelas[k].Seguro,
						Total:   simPgs[i].IntervaloParcelas[k].Total,
					}
					sims[i].IntervaloParcelas = append(sims[i].IntervaloParcelas, ip)
				}
				fmt.Println(time.Now().Format("15:04:05"), "sim: ", i, sims[i].RdbPrazo)
			}
			fmt.Println(time.Now().Format("15:04:05"), "insere sims:", len(sims))
			if err := model.SimulacoesNcwInsere(sims); err != nil {
				log.Printf("simulacoes: conta: %v: %v\n", u.Conta, err)
				continue
			}
		}
	}
	return nil, nil
}

func simulacoesPaginas(urlForm string, values url.Values) ([]Simulacao, error) {
	allSims := make([]Simulacao, 0)
	page := 1
	doc, urlForm, err := postRetry(urlForm, values)
	if err != nil {
		return nil, fmt.Errorf("simulacoespaginas: %w", err)
	}
	sims, err := simulacoesDoc(doc, urlForm)
	if err != nil {
		return nil, fmt.Errorf("simulacoespaginas: %w", err)
	}
	if len(sims) == 0 {
		return allSims, nil
	}
	for _, sim := range sims {
		sim.Pagina = 1
		allSims = append(allSims, sim)
	}
	for {
		page++
		sPage := fmt.Sprintf("Page$%v", page)
		if !strings.Contains(doc, sPage) {
			return allSims, nil
		}
		newValues, err := newValuesDoc(doc)
		if err != nil {
			return nil, fmt.Errorf("simulacoespaginas: %w", err)
		}
		newValues.Set("__EVENTTARGET", "ctl00$Conteudo$grdSimulacao")
		newValues.Set("__EVENTARGUMENT", sPage)
		newValues.Set("ctl00$hdnID_Modulo", values.Get("ctl00$hdnID_Modulo"))
		newValues.Set("ctl00$Conteudo$Pessoa", values.Get("ctl00$Conteudo$Pessoa"))
		newValues.Set("ctl00$Conteudo$edtCD_Inscricao_Nacional", values.Get("ctl00$Conteudo$edtCD_Inscricao_Nacional"))
		newValues.Set("ctl00$Conteudo$cbxProdutos", values.Get("ctl00$Conteudo$cbxProdutos"))
		newValues.Set("ctl00$Conteudo$cbxSubProdutos", values.Get("ctl00$Conteudo$cbxSubProdutos"))
		newValues.Set("ctl00$Conteudo$edtBusca_TipoVenda", values.Get("ctl00$Conteudo$edtBusca_TipoVenda"))
		newValues.Set("ctl00$Conteudo$cbxTp_Venda", values.Get("ctl00$Conteudo$cbxTp_Venda"))
		newValues.Set("ctl00$Conteudo$cbxNegociacao", values.Get("ctl00$Conteudo$cbxNegociacao"))
		newValues.Del("__VIEWSTATEGENERATOR")
		newValues.Del("__EVENTVALIDATION")
		doc, urlForm, err = postRetry(urlForm, newValues)
		if err != nil {
			return nil, fmt.Errorf("simulacoespaginas: %w", err)
		}
		sims, err := simulacoesDoc(doc, urlForm)
		if err != nil {
			return nil, fmt.Errorf("simulacoespaginas: %w", err)
		}
		if len(sims) == 0 {
			return allSims, nil
		}
		for _, sim := range sims {
			sim.Pagina = page
			allSims = append(allSims, sim)
		}
	}

}

func simulacoesDoc(doc string, urlForm string) ([]Simulacao, error) {
	values, err := newValuesDoc(doc)
	if err != nil {
		return nil, fmt.Errorf("simulacoesdoc: %w", err)
	}
	sims := make([]Simulacao, 0)
	linhas := strings.Split(doc, "\n")
	var (
		sim Simulacao
	)
	for _, linha := range linhas {
		linha = strings.TrimSpace(linha)
		if strings.Contains(linha, `<input name="rdbPrazo" `) {
			matches := erValue.FindStringSubmatch(linha)
			if len(matches) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: rdbprazo: len != 2")
			}
			sim = Simulacao{
				RdbPrazo: matches[1],
			}
		}
		if strings.Contains(linha, `<td align="center" style="height:30px;width:70px;">`) {
			matches := erSimGrupo.FindStringSubmatch(linha)
			if len(matches) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: grupo: len != 2")
			}
			sim.Grupo = matches[1]
			matches = erSimBem.FindStringSubmatch(linha)
			if len(matches) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: bem: len != 2")
			}
			sim.Bem = matches[1]
			matches = erSimValor.FindStringSubmatch(linha)
			if len(matches) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: valor: len != 2")
			}
			valor, err := strconv.ParseFloat(rplReal.Replace(matches[1]), 64)
			if err != nil {
				return nil, fmt.Errorf("simulacoesdoc: valor: %w", err)
			}
			sim.Valor = valor
		}
		if strings.Contains(linha, `style="z-index:999;position:relative; top: -4px;">`) {
			matches := erSimAssembleia.FindStringSubmatch(linha)
			if len(matches) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: assembleia: len != 2")
			}
			assembleia, err := time.ParseInLocation("02/01/2006", matches[1], time.Local)
			if err != nil {
				return nil, fmt.Errorf("simulacoesdoc: assembleia: %w", err)
			}
			sim.Assembleia = assembleia
		}
		if strings.Contains(linha, `<td align="center" style="width:85px;">`) {
			matches := erSimAssembleiaParticp.FindStringSubmatch(linha)
			if len(matches) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: assembleiaParticp: len != 2")
			}
			assembleiaParticp, err := time.ParseInLocation("02/01/2006", matches[1], time.Local)
			if err != nil {
				return nil, fmt.Errorf("simulacoesdoc: assembleiaParticp: %w", err)
			}
			sim.AssParticp = assembleiaParticp
		}
		if erSimPrazo.MatchString(linha) {
			matches := erSimPrazo.FindStringSubmatch(linha)
			if len(matches) != 3 {
				return nil, fmt.Errorf("simulacoesdoc: prazo: len != 3")
			}
			sim.PrazoRateio, _ = strconv.Atoi(matches[1])
			sim.PrazoTermo, _ = strconv.Atoi(matches[2])
		}
		if erSimParticp.MatchString(linha) {
			matches := erSimParticp.FindStringSubmatch(linha)
			if len(matches) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: particp: len != 2")
			}
			sim.Particp, _ = strconv.Atoi(matches[1])
		}
		if strings.Contains(linha, `<td align="right" style="width:55px;">`) {
			matches := erSimParcela.FindStringSubmatch(linha)
			if len(matches) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: parcela: len != 2")
			}
			parcela, err := strconv.ParseFloat(rplReal.Replace(matches[1]), 64)
			if err != nil {
				return nil, fmt.Errorf("simulacoesdoc: parcela: %w", err)
			}
			sim.Parcela = parcela
			matches = erSimSeguro.FindStringSubmatch(linha)
			if len(matches) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: seguro: len != 2")
			}
			sim.Seguro = matches[1]
			matches = erSimPlano.FindStringSubmatch(linha)
			matchesB := erSimPlanoB.FindStringSubmatch(linha)
			if len(matches) != 2 && len(matchesB) != 2 {
				return nil, fmt.Errorf("simulacoesdoc: plano: len != 2")
			}
			if len(matches) == 2 {
				sim.Plano = matches[1]
			} else {
				sim.Plano = matchesB[1]
			}
			sim.Vstate = values.Get("__VSTATE")
			sim.UrlForm = urlForm
			sims = append(sims, sim)
		}
	}
	return sims, nil
}

func simulacoesDetalhes(sim *Simulacao) error {
	values := url.Values{
		"__EVENTTARGET":         []string{""},
		"__EVENTARGUMENT":       []string{""},
		"__LASTFOCUS":           []string{""},
		"__VSTATE":              []string{sim.Vstate},
		"__VIEWSTATE":           []string{""},
		"ctl00$hdnID_Modulo":    []string{sim.Values.Get("ctl00$hdnID_Modulo")},
		"ctl00$Conteudo$Pessoa": []string{sim.Values.Get("ctl00$Conteudo$Pessoa")},
		"ctl00$Conteudo$edtCD_Inscricao_Nacional": []string{sim.Values.Get("ctl00$Conteudo$edtCD_Inscricao_Nacional")},
		"ctl00$Conteudo$cbxProdutos":              []string{sim.Values.Get("ctl00$Conteudo$cbxProdutos")},
		"ctl00$Conteudo$cbxSubProdutos":           []string{sim.Values.Get("ctl00$Conteudo$cbxSubProdutos")},
		"ctl00$Conteudo$edtBusca_TipoVenda":       []string{sim.Values.Get("ctl00$Conteudo$edtBusca_TipoVenda")},
		"ctl00$Conteudo$cbxTp_Venda":              []string{sim.Values.Get("ctl00$Conteudo$cbxTp_Venda")},
		"ctl00$Conteudo$cbxNegociacao":            []string{sim.Values.Get("ctl00$Conteudo$cbxNegociacao")},
		"rdbPrazo":                                []string{sim.RdbPrazo},
		"ctl00$Conteudo$btnVender":                []string{"Vender"},
	}
	doc, newU, err := postRetry(sim.UrlForm, values)
	if err != nil {
		return fmt.Errorf("simulacoesdetalhes: %w", err)
	}
	sim.UrlForm = newU
	newValues, err := newValuesDoc(doc)
	if err != nil {
		return fmt.Errorf("simulacoesdetalhes: %w", err)
	}
	linhas := strings.Split(doc, "\n")
	for i, l := range linhas {
		if strings.Contains(l, `ctl00$Conteudo$edtNM_Comissionado`) {
			matches := erValue.FindStringSubmatch(l)
			if len(matches) != 2 {
				return fmt.Errorf("simulacoesdetalhes: nm comissionado: len != 2")
			}
			newValues.Set("ctl00$Conteudo$edtNM_Comissionado", matches[1])
			continue
		}
		if strings.Contains(l, `ctl00$Conteudo$edtBusca_Ponto_Venda`) {
			matches := erValue.FindStringSubmatch(l)
			if len(matches) != 2 {
				return fmt.Errorf("simulacoesdetalhes: ponto venda: len != 2")
			}
			newValues.Set("ctl00$Conteudo$edtBusca_Ponto_Venda", matches[1])
			continue
		}
		if strings.Contains(l, `"ctl00$Conteudo$cbxEquipeVenda"`) {
			l = linhas[i+2]
			matches := erValue.FindStringSubmatch(l)
			if len(matches) != 2 {
				return fmt.Errorf("simulacoesdetalhes: equipe venda: len != 2")
			}
			newValues.Set("ctl00$Conteudo$cbxEquipeVenda", matches[1])
			continue
		}
		if strings.Contains(l, `"ctl00$Conteudo$edtComissionado"`) {
			matches := erValue.FindStringSubmatch(l)
			if len(matches) != 2 {
				return fmt.Errorf("simulacoesdetalhes: comissionado: len != 2")
			}
			newValues.Set("ctl00$Conteudo$edtBuscaEquipe_Venda", strings.TrimSpace(matches[1]))
			continue
		}
		if strings.Contains(l, `ctl00$Conteudo$edtBusca_TipoVenda`) {
			matches := erValue.FindStringSubmatch(l)
			if len(matches) != 2 {
				return fmt.Errorf("simulacoesdetalhes: comissionado: len != 2")
			}
			newValues.Set("ctl00$Conteudo$edtBusca_TipoVenda", matches[1])
			continue
		}
		if strings.Contains(l, `ctl00$Conteudo$edtBusca_Bem`) {
			matches := erValue.FindStringSubmatch(l)
			if len(matches) != 2 {
				return fmt.Errorf("simulacoesdetalhes: comissionado: len != 2")
			}
			newValues.Set("ctl00$Conteudo$edtBusca_Bem", matches[1])
			continue
		}
	}
	newValues.Set("__EVENTTARGET", "ctl00$Conteudo$btnAvancar")
	newValues.Set("ctl00$hdnID_Modulo", "VE")
	newValues.Set("ctl00$Conteudo$cbxSN_Sigilo", "N")
	newValues.Set("ctl00$Conteudo$edtCPFVendedor", "")
	newValues.Set("ctl00$Conteudo$rblTipoPesquisa", "B")
	newValues.Set("ctl00$Conteudo$cbxFormaRecebimento_1_Parc", "BB")
	newValues.Set("ctl00$Conteudo$hdnTokenRecaptcha", "")
	newValues.Set("__VIEWSTATEENCRYPTED", "")
	doc, newU, err = postRetry(newU, newValues)
	if err != nil {
		return fmt.Errorf("simulacoesdetalhes: %w", err)
	}
	newValues, err = newValuesDoc(doc)
	if err != nil {
		return fmt.Errorf("simulacoesdetalhes: %w", err)
	}
	newValues.Del("__LASTFOCUS")
	newValues.Del("__VIEWSTATEGENERATOR")
	newValues.Set("ctl00$hdnID_Modulo", "VE")
	newValues.Set("ctl00$Conteudo$grdProposta_Prazo$ctl02$btnDetalhesBem.x", "10")
	newValues.Set("ctl00$Conteudo$grdProposta_Prazo$ctl02$btnDetalhesBem.y", "13")
	doc, newU, err = postRetry(newU, newValues)
	if err != nil {
		return fmt.Errorf("simulacoesdetalhes: %w", err)
	}
	if err := simulacaoDetalha(sim, doc); err != nil {
		return fmt.Errorf("simulacoesdetalhes: %w", err)
	}
	return nil
}

func simulacaoDetalha(sim *Simulacao, doc string) error {
	linhas := strings.Split(doc, "\n")
	for i := range linhas {
		linhas[i] = strings.TrimSpace(linhas[i])
	}
	for i, l := range linhas {
		if strings.Contains(l, `</textarea>`) {
			sim.Caracteristica = strings.Replace(l, "</textarea>", "", 1)
			continue
		}
		if strings.Contains(l, `ctl00_Conteudo_lblDT_Vencimento`) {
			vencimento := strings.Replace(
				l,
				`<span id="ctl00_Conteudo_lblDT_Vencimento" class="label_destaque" style="display:inline-block;`+
					`width:120px;z-index: 102; position: absolute; top: 22px; left: 811px;">`,
				"",
				1,
			)
			vencimento = strings.Replace(vencimento, `</span>`, "", 1)
			vencimento = strings.TrimSpace(vencimento)
			v, err := time.ParseInLocation("02/01/2006", vencimento, time.Local)
			if err != nil {
				return fmt.Errorf("simulacao detalha: %w", err)
			}
			sim.Vencimento = v
			continue
		}
		if strings.Contains(l, `ctl00_Conteudo_lblVA_Lance_Maximo`) {
			perMaxLance := strings.Replace(
				l,
				`<span id="ctl00_Conteudo_lblVA_Lance_Maximo" class="label_destaque" style="display:inline-block;`+
					`width:120px;z-index: 104; position: absolute; top: 49px; left: 132px;">`,
				"",
				1,
			)
			perMaxLance = strings.Replace(perMaxLance, `</span>`, "", 1)
			perMaxLance = strings.Replace(perMaxLance, `%`, "", 1)
			perMaxLance = strings.TrimSpace(perMaxLance)
			perMaxLance = rplReal.Replace(perMaxLance)
			pml, err := strconv.ParseFloat(perMaxLance, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: %w", err)
			}
			sim.PerMaxLan = pml
			continue
		}
		if strings.Contains(l, `<tr style="background-color:#ffffff;" align="right" class="fonte_8pt">`) {
			ini := linhas[i+1]
			ini = strings.Replace(ini, `<td style="text-align:center;">`, "", 1)
			ini = strings.Replace(ini, `</td>`, "", 1)
			ini = strings.TrimSpace(ini)
			fim := linhas[i+2]
			fim = strings.Replace(fim, `<td style="text-align:center;">`, "", 1)
			fim = strings.Replace(fim, `</td>`, "", 1)
			fim = strings.TrimSpace(fim)
			funComP := linhas[i+3]
			funComP = strings.Replace(funComP, `<td>`, "", 1)
			funComP = strings.Replace(funComP, `</td>`, "", 1)
			funComP = rplReal.Replace(funComP)
			funComR := linhas[i+4]
			funComR = strings.Replace(funComR, `<td>`, "", 1)
			funComR = strings.Replace(funComR, `</td>`, "", 1)
			funComR = rplReal.Replace(funComR)
			taxAdmP := linhas[i+5]
			taxAdmP = strings.Replace(taxAdmP, `<td>`, "", 1)
			taxAdmP = strings.Replace(taxAdmP, `</td>`, "", 1)
			taxAdmP = rplReal.Replace(taxAdmP)
			taxAdmR := linhas[i+6]
			taxAdmR = strings.Replace(taxAdmR, `<td>`, "", 1)
			taxAdmR = strings.Replace(taxAdmR, `</td>`, "", 1)
			taxAdmR = rplReal.Replace(taxAdmR)
			taxAdsP := linhas[i+7]
			taxAdsP = strings.Replace(taxAdsP, `<td>`, "", 1)
			taxAdsP = strings.Replace(taxAdsP, `</td>`, "", 1)
			taxAdsP = rplReal.Replace(taxAdsP)
			taxAdsR := linhas[i+8]
			taxAdsR = strings.Replace(taxAdsR, `<td>`, "", 1)
			taxAdsR = strings.Replace(taxAdsR, `</td>`, "", 1)
			taxAdsR = rplReal.Replace(taxAdsR)
			funResP := linhas[i+9]
			funResP = strings.Replace(funResP, `<td>`, "", 1)
			funResP = strings.Replace(funResP, `</td>`, "", 1)
			funResP = rplReal.Replace(funResP)
			funResR := linhas[i+10]
			funResR = strings.Replace(funResR, `<td>`, "", 1)
			funResR = strings.Replace(funResR, `</td>`, "", 1)
			funResR = rplReal.Replace(funResR)
			seguro := linhas[i+13]
			seguro = rplReal.Replace(seguro)
			total := linhas[i+18]
			total = rplReal.Replace(total)
			var (
				err error
				ip  IntevaloParcelas
			)
			ip.Ini, err = strconv.Atoi(ini)
			if err != nil {
				return fmt.Errorf("simulacao detalha: ini: %w", err)
			}
			ip.Fim, err = strconv.Atoi(fim)
			if err != nil {
				return fmt.Errorf("simulacao detalha: fim: %w", err)
			}
			ip.FunComP, err = strconv.ParseFloat(funComP, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: fundo comum p: %w", err)
			}
			ip.FunComR, err = strconv.ParseFloat(funComR, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: fundo comum r: %w", err)
			}
			ip.TaxAdmP, err = strconv.ParseFloat(taxAdmP, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: taxa adm p: %w", err)
			}
			ip.TaxAdmR, err = strconv.ParseFloat(taxAdmR, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: taxa adm r: %w", err)
			}
			ip.TaxAdsP, err = strconv.ParseFloat(taxAdsP, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: taxa ads p: %w", err)
			}
			ip.TaxAdsR, err = strconv.ParseFloat(taxAdsR, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: taxa ads r: %w", err)
			}
			ip.FunResP, err = strconv.ParseFloat(funResP, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: fundo res p: %w", err)
			}
			ip.FunResR, err = strconv.ParseFloat(funResR, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: fundo res r: %w", err)
			}
			ip.Seguro, err = strconv.ParseFloat(seguro, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: seguro: %w", err)
			}
			ip.Total, err = strconv.ParseFloat(total, 64)
			if err != nil {
				return fmt.Errorf("simulacao detalha: total: %w", err)
			}
			sim.IntervaloParcelas = append(sim.IntervaloParcelas, ip)
		}
	}
	return nil
}
