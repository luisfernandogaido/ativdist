package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Escola struct {
	Categoria   string
	Nome        string
	Qs          string
	Professoras []Professora
}

type Professora struct {
	Escola     Escola
	Nome       string
	Qs         string
	Documentos []Documento
}

type Documento struct {
	Professora Professora
	Nome       string
	Href       string
}

func main() {
	if err := Coleta(); err != nil {
		log.Fatal(err)
	}
}

func Coleta() error {
	doc, newU, err := getRetry(urlBase + "/educacao/atividades_distancia.aspx")
	if err != nil {
		return fmt.Errorf("coleta: %w", err)
	}
	values, err := newValuesDoc(doc)
	if err != nil {
		return fmt.Errorf("coleta: %w", err)
	}
	values.Set("__EVENTTARGET", "ctl00$ctl00$ctl00$ExternoBody$content$content_educacao$lbtnMateriaisImpressao")
	doc, newU, err = postRetry(newU, values)
	if err != nil {
		return fmt.Errorf("coleta: %w", err)
	}
	escolas, err := Escolas(doc)
	if err != nil {
		return fmt.Errorf("coleta: %w", err)
	}
	for i, e := range escolas {
		fmt.Printf("escola %v/%v: %v\n", i+1, len(escolas), e.Nome)
		professoras, err := Professoras(e)
		if err != nil {
			return fmt.Errorf("coleta: %w", err)
		}
		for j, p := range professoras {
			fmt.Printf("professora %v/%v: %v\n", j+1, len(professoras), p.Nome)
			dirPro := filepath.Join(dir, fmt.Sprintf("%v/%v/%v", time.Now().Format("20060102"), e.Nome, p.Nome))
			if err := os.MkdirAll(dirPro, 0644); err != nil {
				return fmt.Errorf("coleta: %w", err)
			}
			documentos, err := Documentos(p)
			if err != nil {
				return fmt.Errorf("coleta: %w", err)
			}
			for k, d := range documentos {
				fmt.Printf("documento %v/%v: %v\n", k+1, len(documentos), d.Nome)
				if err := Download(d, dirPro); err != nil {
					return fmt.Errorf("coleta: %w", err)
				}
			}
			professoras[j].Documentos = documentos
		}
		escolas[i].Professoras = professoras
	}
	return nil
}

func Escolas(doc string) ([]Escola, error) {
	linhas := strings.Split(doc, "\n")
	var linha string
	for _, l := range linhas {
		if strings.Contains(l, `<span id="ctl00_ctl00_ctl00_ExternoBody_content_content_educacao_lblHTML">`) {
			linha = l
			break
		}
	}
	indexh3 := erH3.FindAllStringIndex(linha, -1)
	trechos := make([]string, 0)
	for i := range indexh3 {
		var trecho string
		if i != len(indexh3)-1 {
			trecho = linha[indexh3[i][0]:indexh3[i+1][0]]
		} else {
			trecho = linha[indexh3[i][0]:]
		}
		trechos = append(trechos, trecho)
	}
	escolas := make([]Escola, 0)
	for _, trecho := range trechos {
		matchesh3 := erH3.FindStringSubmatch(trecho)
		categoria := matchesh3[1]
		matches := erEscola.FindAllStringSubmatch(trecho, -1)
		for _, m := range matches {
			escola := Escola{
				Categoria: categoria,
				Nome:      m[2],
				Qs:        m[1],
			}
			escolas = append(escolas, escola)
		}

	}
	return escolas, nil
}

func Professoras(e Escola) ([]Professora, error) {
	doc, _, err := getRetry(urlBase + "/educacao/atividades_distancia.aspx?" + e.Qs)
	if err != nil {
		return nil, fmt.Errorf("professoras: %w", err)
	}
	matches := erProfessora.FindAllStringSubmatch(doc, -1)
	professoras := make([]Professora, len(matches))
	for i := range matches {
		professoras[i] = Professora{
			Escola: e,
			Nome:   matches[i][2],
			Qs:     matches[i][1],
		}
	}
	return professoras, nil
}

func Documentos(p Professora) ([]Documento, error) {
	doc, _, err := getRetry(urlBase + "/educacao/atividades_distancia.aspx?" + p.Qs)
	if err != nil {
		return nil, fmt.Errorf("documentos: %w", err)
	}
	matches := erDocumento.FindAllStringSubmatch(doc, -1)
	documentos := make([]Documento, 0)
	for _, m := range matches {
		d := Documento{
			Professora: p,
			Nome:       m[2],
			Href:       m[1],
		}
		documentos = append(documentos, d)
	}
	return documentos, nil
}

func Download(d Documento, dir string) error {
	filename := filepath.Join(dir, d.Nome)
	u := urlBase + d.Href
	res, err := client.Get(u)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	if err := ioutil.WriteFile(filename, b, 0644); err != nil {
		return fmt.Errorf("download: %w", err)
	}
	return nil
}
