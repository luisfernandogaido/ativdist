package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Arquivo struct {
	Nome string
	Href string
}

type Escola struct {
	Nome     string
	A        string
	Arquivos []Arquivo
}

const (
	urlBase = "https://www2.bauru.sp.gov.br"
	dirBase = "./infantil"
	dirArqs = "/arquivos/arquivos_site/sec_educacao/atividades_pedagogica_distancia/"
)

var (
	reEscInf = regexp.MustCompile(`<a href="atividades_distancia\.aspx\?t=1&a=([\d]+)#1">([^<]+)</a>`)
	reArqInf = regexp.MustCompile(
		`<a href="/arquivos/arquivos_site/sec_educacao/atividades_pedagogica_distancia/([^"]+)"><b>([^<]+)</b></a>`,
	)
	reHidden = regexp.MustCompile(
		`<input type="hidden" name="([^"]+)" id="[^"]+" value="([^"]+)" />`,
	)
)

func main() {
	t0 := time.Now()
	if err := os.Mkdir(dirBase, 0744); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	escolas, err := escolasInfantil()
	if err != nil {
		log.Fatal(err)
	}
	for i, escola := range escolas {
		dir := filepath.Join(dirBase, escola.Nome)
		_, err := os.Stat(dir)
		if err != nil && !os.IsExist(err) {
			if err := os.Mkdir(dir, 0744); err != nil && err != os.ErrExist {
				log.Fatal(err)
			}
		}
		if err := arquivos(&escola); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v) %v (%v arquivos):\n", i+1, escola.Nome, len(escola.Arquivos))
		for _, a := range escola.Arquivos {
			file := filepath.Join(dir, a.Nome)
			fmt.Println(a.Nome)
			if err := download(a.Href, file); err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("---")
	}
	fmt.Println(time.Since(t0))
}

func escolasInfantil() ([]Escola, error) {
	res, err := http.Get(urlBase + "/educacao/atividades_distancia.aspx")
	if err != nil {
		return nil, fmt.Errorf("escolas infantil: %w", err)
	}
	s := bufio.NewScanner(res.Body)
	values := url.Values{
		"__EVENTTARGET":                        []string{"ctl00$ctl00$ctl00$ExternoBody$content$content_educacao$lbtnMateriaisImpressao"},
		"__EVENTARGUMENT":                      []string{""},
		"_SCROLLPOSITIONX":                     []string{"0"},
		"__SCROLLPOSITIONY":                    []string{"0"},
		"ctl00$ctl00$ctl00$cabecalho$txtBusca": []string{""},
	}
	for s.Scan() {
		matches := reHidden.FindStringSubmatch(s.Text())
		if len(matches) > 0 {
			values.Set(matches[1], matches[2])
		}
	}
	res.Body.Close()
	res, err = http.PostForm(urlBase+"/educacao/atividades_distancia.aspx", values)
	if err != nil {
		return nil, fmt.Errorf("escolas infantil: %w", err)
	}
	s = bufio.NewScanner(res.Body)
	var linha string
	for s.Scan() {
		if strings.Contains(s.Text(), "<span id=\"ctl00_ctl00_ctl00_ExternoBody_content_content_educacao_lblHTML\">") {
			linha = s.Text()
			break
		}
	}
	q := strings.Index(linha, "<h3>Fundamental</h3>")
	linha = linha[:q]
	matches := reEscInf.FindAllStringSubmatch(linha, -1)
	escolas := make([]Escola, 0)
	for _, m := range matches {
		if len(m) != 3 {
			return nil, fmt.Errorf("escolas infantil: %w", err)
		}
		escola := Escola{
			Nome: m[2],
			A:    m[1],
		}
		escolas = append(escolas, escola)
	}
	return escolas, nil
}

func arquivos(e *Escola) error {
	res, err := http.Get(fmt.Sprintf("%v/educacao/atividades_distancia.aspx?t=1&a=%v", urlBase, e.A))
	if err != nil {
		return fmt.Errorf("arquivos: %w", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("arquivos: %w", err)
	}
	doc := string(b)
	matches := reArqInf.FindAllStringSubmatch(doc, -1)
	e.Arquivos = make([]Arquivo, 0)
	for _, m := range matches {
		if len(m) != 3 {
			return fmt.Errorf("arquivos: %w", err)
		}
		arq := Arquivo{
			Nome: m[2],
			Href: urlBase + dirArqs + url.PathEscape(m[1]),
		}
		e.Arquivos = append(e.Arquivos, arq)
	}
	return nil
}

func download(path, file string) error {
	res, err := http.Get(path)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer res.Body.Close()
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, res.Body); err != nil {
		return fmt.Errorf("download: %w", err)
	}
	return nil
}
