package ncw

import (
	"fmt"
	"io"
)

func PosicaoConsorciado() error {
	res, err := client.Get(urlBase + "/CONAT/frmConAtSrcConsorciado.aspx?cd=15000")
	if err != nil {
		return fmt.Errorf("posicaoconsorciado: %w", err)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("posicaoconsorciado: %w", err)
	}
	h, err := getHashes(string(b))
	if err != nil {
		return fmt.Errorf("posicaoconsorciado: %w", err)
	}
	data := newValuesHashes(h)
	data.Set("ctl00$hdnID_Modulo", "AT")
	data.Set("ctl00$Conteudo$edtGrupo", "000000")
	data.Set("ctl00$Conteudo$edtCota:", "0000")
	data.Set("ctl00$Conteudo$edtVersao", "00")
	data.Set("ctl00$Conteudo$btnBuscaAvancada", "Busca avançada...")
	u := res.Request.URL.String()
	res2, err := client.PostForm(u, data)
	if err != nil {
		return fmt.Errorf("posicaoconsorciado: %w", err)
	}
	defer res2.Body.Close()
	b, err = io.ReadAll(res2.Body)
	if err != nil {
		return fmt.Errorf("posicaoconsorciado: %w", err)
	}
	h, err = getHashes(string(b))
	if err != nil {
		return fmt.Errorf("posicaoconsorciado: %w", err)
	}
	data = newValuesHashes(h)
	data.Set("ctl00$hdnID_Modulo", "AT")
	data.Set("ctl00$Conteudo$edtGrupo", "000000")
	data.Set("ctl00$Conteudo$edtCota:", "0000")
	data.Set("ctl00$Conteudo$edtVersao", "00")
	data.Set("ctl00$Conteudo$btnBuscaAvancada", "Busca avançada...")
	return nil
}
