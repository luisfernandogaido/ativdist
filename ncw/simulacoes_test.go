package ncw

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestSimulacoesDoc(t *testing.T) {
	b, err := ioutil.ReadFile("./data/simulacoes-doc.html")
	if err != nil {
		t.Fatal(err)
	}
	simulacoes, err := simulacoesDoc(string(b), "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(simulacoes))
}

func TestSimulacaoDetalha(t *testing.T) {
	assembleia, err := time.Parse("2006-01-02 15:04:05 -0700", "2021-09-22 00:00:00 -0300")
	if err != nil {
		t.Fatal(err)
	}
	asseParticp, err := time.Parse("2006-01-02 15:04:05 -0700", "2021-09-22 00:00:00 -0300")
	if err != nil {
		t.Fatal(err)
	}
	sim := Simulacao{
		Consulta:    time.Now(),
		Conta:       "Javep",
		TipoPessoa:  "FÍSICA",
		TipoVenda:   "NACIONAL",
		RdbPrazo:    "1|383|84|A|1530|19001|000383|93% ONIX 1.0 - JOY|57279,00|110|216915|250|A|43",
		Grupo:       "019001 - A",
		Bem:         "93% ONIX 1.0 - JOY",
		Valor:       57279,
		Assembleia:  assembleia,
		AssParticp:  asseParticp,
		PrazoRateio: 43,
		PrazoTermo:  84,
		Particp:     250,
		Parcela:     1664.74,
		Seguro:      "Sim",
		Plano:       "TX ANT.2AUTOM.NACIONAIS 84/500",
		Vstate:      "",
		UrlForm:     "",
		Values:      nil,
	}
	doc := `
<!DOCTYPE html>

<html>
<head id="ctl00_Head1"><title>
    Intranet Newcon
</title><meta http-equiv="X-UA-Compatible" content="IE=EmulateIE11" /><meta http-equiv="Cache-Control" content="no-cache" />
    <script type="text/javascript">
        vallowNegative = false;
    </script>
    <link href="https://cncdealer.gmfinancial.com/NewconWeb/includes/Menu_Conat.css?v=31121600220000" type="text/css" rel="stylesheet" />
    <link href="https://cncdealer.gmfinancial.com/NewconWeb/js/css/overcast/jquery-ui.css?v=31121600220000" type="text/css" rel="stylesheet" />
    <link href="https://cncdealer.gmfinancial.com/NewconWeb/js/css/tab/tab.css?v=31121600220000" type="text/css" rel="stylesheet" />
    <link href="https://cncdealer.gmfinancial.com/NewconWeb/App_Themes/Default/Default.css?v=31121600220000" type="text/css" rel="stylesheet" />
    <style type='text/css'>                 div#fundoMenu {                         z-index: -2;                         position: absolute;                         margin-top: 110px;                         height: 6px;                         width: 100%;    	                 background: #989896;                         border-bottom-width: 3px;                         border-top-width: 0;                         border-right-width: 0;                         border-left-width: 0;                         border-style: solid;                         border-color: #B6862D;                 }                 div#fundoSubMenu{                     z-index: -1;                     position: absolute;                     margin-top: 120px;                     height: 38px;                     width: 100%;                     background: #F6F6F6;                     border-bottom-width: 1px;                     border-top-width: 0;                     border-right-width: 0;                     border-left-width: 0;                     border-style: solid;                     border-color: #c5c5c5;                 }                 div#tabs_principal {                     position: relative;                     z-index: 0;                 }                 @media screen and (max-width: 1200px) {                     div#fundoMenu , div#fundoSubMenu {                         width: 120%;                     }                 }                 @media screen and (max-width: 920px) {                     div#fundoMenu , div#fundoSubMenu {                         width: 140%;                     }                 }                 @media screen and (max-width: 720px) {                     div#fundoMenu , div#fundoSubMenu {                         width: 160%;                     }                 }                 @media screen and (max-width: 620px) {                     div#fundoMenu , div#fundoSubMenu {                         width: 200%;                     }                 }                 @media screen and (max-width: 480px) {                     div#fundoMenu , div#fundoSubMenu {                         width: 240%;                     }                 }                 div#logo {                     background-image: url('data:image/png;base64,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');                     background-repeat: no-repeat;                 }                 div.inner div.title, .titulo_pop_up {                      background: #989896;                      color: #FFFFFF;                 };             </style></head>

<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/jquery.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/jquery-migrate.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/Utils.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/w3cookies.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/menu.builder.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/jquery_maskedinput.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/PriceFormat.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/jquery.blockalpha.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/jquery.expose.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/jquery.loading.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/jquery.format.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/jquery-ui.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/MasterPageScript.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/cnpAjax.js?v=31121600220000'></script>
<script type="text/javascript" src='https://cncdealer.gmfinancial.com/NewconWeb/js/endereco.js?v=31121600220000'></script>



<body>
<form method="post" action="/Newconweb/CONPV/frmConPvCadProposta_Detalhes.aspx?applicationKey=Cyhprg/KUUmPjCBVS+oDpusbFelXH6zAQ0M0SdvGElsk1BPmsw1lDhmJcG+t4mQKiR+N/MVIZsQ08W+3qPs9fqnk4+tT6nz/kqDxbLdzmIfSmk0YDGuEpzmqlZGrZCdbU55RJvhYsX7LqFBNYHhhAvlRPB6lkiex5TryC+fwF3GNIlISmXXLF4iw6otboEjsHTJJwG1xhSJ/MzucSexxuNJmQzVEoLsd2/14I8xcquK985MInYtde0WnLfbglmRd6wVF5VzZ+4Ssxok/uvFWL2U3NB3G5yK9Vq0iwqPxUXE=" id="aspnetForm" novalidate="true">
    <div class="aspNetHidden">
        <input type="hidden" name="__VSTATE" id="__VSTATE" value="H4sIAAAAAAAEAKVVzZLaRhCOBsTCblz4p8BJOdayhV2urWjZGf3rCOw6ISExMXjjm2rMDIvKQiJIpDZvkHOSa3LJLW+Q58gpr5KL0yMhvF7AlapcYPTN1/1N93T3vJGqH0my4nndKEwWURA/598t/QUfRHHSoePXX/IfPK9elGvjJMD4kT+79NpBwhd0yMMplY8zuDejl7yzTJIobI95HNOhn/A2m/mhHycLyqIFlZvbqdGZH9MA1Pk4oe/IJDxk/oyHSVSVKyeOZdjY0nGV1dGkWq0XFTRmE+XDF8/73uc8mHvf8lcc9hACxv6kWkfK3oUf+68CPmUpWESSQCt+GPLFNJkF8mGXMgoHjBqMNwaLaA5B08ZJ44wnNJjymKEC2NWRXBWGxRG/SuSKabc021UxZgxVxMZhUW6ev+z2Xwx7F88aX7QvzgcNxzjVTAyu+u2vu+eNp72Xz4BeXdFLmLgYE0DurZAKdozGjMegydD9FbhP7FPsnmpYE9QHK7RsYhAn+DFgh2tzQ1+bN3NzTbtm/mSFFuBc8HkMgWlZQg68o17CZ91oGSaoUC9NsgTW9gsgU8BEKxHVNUxSJi1CDJXoJazarktKhmmrNimJ42BcFH9ibbmOrLuq6zBwX5NkU1NdTdxATUo5Iqk1CZxZlqHaBkPF1SchqqMJo7U80QvY0EqaatiaCxSDWKphgohmwymIaaoG/n/y5rvyOsvLR8gbQl5P5YmTyptwwkzeSeUtZ0PeMYQ8+Q/yumrha/IYvInd1S0TnHmGu5LzW05Lz1EdB8C9FbhHnJyXV2OF4JYO/jQBHuQ1p+e0W2uJFsSh6sa1yqxYbkvHtmpoAN5Zl2tumpdrFgpDtQ3G/RuMjzcYD24wHt74PtySgaNtGWhuZuDxtgw82TjC8Q3JTzdydLIlR6c5psEZTNW0ACO5oE1apqOJshPuxBwSN3lXUb6PL2COwUCcLQM6plGfhmPu9c68zxbLeYT+/EVSmjs57TjmM5hgPkV//XZLebSTOAhoGHnpJnqDdhGHI2/kzyNvIOYulaWnyic7PXb4DP1zoDzcTjgbeR0a89L8p29+/vX132VF3emox2CE+xN/LJ4B8cz4Ifrxd+c9UY/oFc0iQuH7YoYXK4/5D8Tq5XR0IPi9nU7/85BC6tiUsbqEJNExt2H5wds+kbL+EC+DnO5Jk6yvykq5G8fdgMaxXEtg6jfmfDGGMJY08AKhr1S+old9Hl4mU1ROXSlF72jYQWjdr2i6bsl0uZ+t6qVV695BjXV7Zuvq2/W/gUcUOZcHAAA=" />
        <input type="hidden" name="__VIEWSTATE" id="__VIEWSTATE" value="" />
    </div>


    <script src="../js/cnpAjax.js" type="text/javascript"></script>
    <script type="text/javascript">
        //<![CDATA[
        function PostAjax(arg, context, callback) {
            __theFormPostData = '';
            WebForm_InitCallback();
            WebForm_DoCallback('__Page',arg,RetornoAjax,context,RetornoErroAjax,true);
            if (callback != undefined){ callback(); }
        }//]]>
    </script>



    <div id="theme">
        <div id="fundoMenu"></div>
        <div id="fundoSubMenu"></div>
    </div>


    <div id="principal">
        <div id="logo">

        </div>

        <div id="menu">
            <div id="tabs_principal" class="tabs" style="position: relative; z-index: 0;">
                <ul>
                    <li id="Li1"><a href="../frmMain.aspx"><span>Página Inicial</span></a>

                    </li>
                    <li id="menu_trigger_0">
                        <a href="../frmMain.aspx">
                                                          <span id="VE">
                                                              Venda
                                                          </span>
                        </a>

                        <ul style="display:none">
                            <li><a href="../frmMainMenu.aspx?ID_Modulo=VE&ID_SubGrupo=RE">
                                Representante</a>
                                | </li>

                            <li><a href="../frmMainMenu.aspx?ID_Modulo=VE&ID_SubGrupo=VP">
                                Venda de Proposta</a>
                            </li>

                        </ul>

                    </li>
                    <li id="menu_trigger_2">
                        <a href="../frmMain.aspx">
                                                          <span id="CO">
                                                              Cobrança
                                                          </span>
                        </a>

                        <ul style="display:none">
                            <li><a href="../frmMainMenu.aspx?ID_Modulo=CO&ID_SubGrupo=RA">
                                Retenção / Recuperação</a>
                            </li>

                        </ul>

                    </li>
                    <li id="menu_trigger_3">
                        <a href="../frmMain.aspx">
                                                          <span id="CP">
                                                              Contemplação
                                                          </span>
                        </a>

                        <ul style="display:none">
                            <li><a href="../frmMainMenu.aspx?ID_Modulo=CP&ID_SubGrupo=PC">
                                Pré-Contemplação</a>
                                | </li>

                            <li><a href="../frmMainMenu.aspx?ID_Modulo=CP&ID_SubGrupo=CP">
                                Contemplação</a>
                                | </li>

                            <li><a href="../frmMainMenu.aspx?ID_Modulo=CP&ID_SubGrupo=RE">
                                Relatório</a>
                            </li>

                        </ul>

                    </li>
                    <li id="menu_trigger_4">
                        <a href="../frmMain.aspx">
                                                          <span id="CR">
                                                              Crédito
                                                          </span>
                        </a>

                        <ul style="display:none">
                            <li><a href="../frmMainMenu.aspx?ID_Modulo=CR&ID_SubGrupo=CR">
                                Concessão de Crédito</a>
                            </li>

                        </ul>

                    </li>
                    <li id="menu_trigger_6">
                        <a href="../frmMain.aspx">
                                                          <span id="SI">
                                                              Sistema
                                                          </span>
                        </a>

                        <ul style="display:none">
                            <li><a href="../frmMainMenu.aspx?ID_Modulo=SI&ID_SubGrupo=PE">
                                Pessoas</a>
                            </li>

                        </ul>

                    </li>
                    <li id="menu_trigger_7">
                        <a href="../frmMain.aspx">
                                                          <span id="AT">
                                                              Atendimento
                                                          </span>
                        </a>

                        <ul style="display:none">
                            <li><a href="../frmMainMenu.aspx?ID_Modulo=AT&ID_SubGrupo=AT">
                                Atendimento a Clientes</a>
                            </li>

                        </ul>


                        &nbsp;

                    </li>
                </ul>
                <input type="hidden" name="ctl00$hdnID_Modulo" id="ctl00_hdnID_Modulo" />
            </div>
            <div id="subs" class="sub"><ul></ul></div>
        </div>
        <div id="conteudo">

            <div id="formulario">
                <div id="divLoading" class="loading" style="display: none;">
                    <img alt="processando..." src='../App_Themes/Default/imagens/loading29.gif' /><br />
                    <strong>Processando...</strong><br />
                </div>

                <!--<formulario>-->
                <div class="theme_container" >
                    <div class="inner" style="z-index:0;position:relative;height: 892px;">
                        <div id="ctl00_Conteudo_lblTitle" class="title round_corner" style="width:935px; z-index: 100;">Cadastro de Proposta - Detalhes</div>
                        <br />

                        <div id="ctl00_Conteudo_divContent" class="content" style="position: relative; z-index: 1; top: 15px; width: 938px; height: 892px;">
                            <div style="position: relative; width: 938px; height: 77px; left: 0px; top: 0px;">
                                <span id="ctl00_Conteudo_lblMensagem" class="titulo_container" style="display:inline-block;width:936px;z-index: 105; position: absolute;">Detalhes do Bem</span>

                                <span id="ctl00_Conteudo_Label10" class="label" style="display:inline-block;width:100px;z-index: 101; position: absolute; top: 20px">Valor do Bem:</span>
                                <span id="ctl00_Conteudo_lblVL_Bem" class="label_destaque" style="display:inline-block;width:818px;z-index: 102; position: absolute; top: 22px; left: 112px;">57.279,00</span>

                                <span id="ctl00_Conteudo_Label8" class="label" style="display:inline-block;width:100px;z-index: 101; position: absolute; top: 48px; left: auto;">Característica:</span>

                                <textarea name="ctl00$Conteudo$lblNM_Caracteristica" rows="2" cols="20" readonly="readonly" id="ctl00_Conteudo_lblNM_Caracteristica" class="text" onblur="zeros_esquerda(this,6)" style="height:13px;width:816px;z-index: 101; left: 112px; position: absolute; height:50px; top: 50px; resize:none;  display:block;">
EXCLUSIVO JAVEP 84/250 - LANCE FIXO</textarea>
                            </div>
                            <div style="position: relative; width: 938px; height: 78px; left: 0px; top: 38px;">
                                <span id="ctl00_Conteudo_lblMensagem2" class="titulo_container" style="display:inline-block;width:936px;z-index: 105; position: absolute;">Detalhes do Grupo</span>

                                <span id="ctl00_Conteudo_Label11" class="label" style="display:inline-block;width:120px;z-index: 101; position: absolute; top: 20px">Grupo:</span>
                                <span id="ctl00_Conteudo_lblCD_Grupo" class="label_destaque" style="display:inline-block;width:120px;z-index: 102; position: absolute; top: 22px; left: 132px;">019001</span>

                                <span id="ctl00_Conteudo_Label13" class="label" style="display:inline-block;width:120px;z-index: 101; position: absolute; top: 20px; left: 342px;">PrazoRateio Termo:</span>
                                <span id="ctl00_Conteudo_lblPZ_Termo" class="label_destaque" style="display:inline-block;width:120px;z-index: 102; position: absolute; top: 22px; left: 474px;">084 meses</span>

                                <span id="ctl00_Conteudo_Label19" class="label" style="display:inline-block;width:120px;z-index: 101; position: absolute; top: 20px; left: 679px;">Vencimento:</span>
                                <span id="ctl00_Conteudo_lblDT_Vencimento" class="label_destaque" style="display:inline-block;width:120px;z-index: 102; position: absolute; top: 22px; left: 811px;">17/09/2021</span>

                                <span id="ctl00_Conteudo_Label22" class="label" style="display:inline-block;width:120px;z-index: 103; position: absolute; top: 47px">% Máximo do Lance:</span>
                                <span id="ctl00_Conteudo_lblVA_Lance_Maximo" class="label_destaque" style="display:inline-block;width:120px;z-index: 104; position: absolute; top: 49px; left: 132px;">50,0010%</span>

                                <span id="ctl00_Conteudo_Label17" class="label" style="display:inline-block;width:120px;z-index: 101; position: absolute; top: 47px; left: 342px;">PrazoRateio Rateio:</span>
                                <span id="ctl00_Conteudo_lblPZ_Rateio" class="label_destaque" style="display:inline-block;width:120px;z-index: 102; position: absolute; top: 49px; left: 474px;">043 meses</span>

                                <span id="ctl00_Conteudo_Label21" class="label" style="display:inline-block;width:120px;z-index: 101; position: absolute; top: 47px; left: 679px;">Assembleia:</span>
                                <span id="ctl00_Conteudo_lblDT_Assembleia" class="label_destaque" style="display:inline-block;width:120px;z-index: 102; position: absolute; top: 49px; left: 811px;">22/09/2021</span>

                                <span id="ctl00_Conteudo_Label15" class="label" style="display:inline-block;width:120px;z-index: 101; position: absolute; top: 74px">Participantes:</span>
                                <span id="ctl00_Conteudo_lblQT_Participante" class="label_destaque" style="display:inline-block;width:120px;z-index: 102; position: absolute; top: 76px; left: 132px;">250</span>
                            </div>
                            <div id="ctl00_Conteudo_divTabela" class="dashboard" style="position: relative; margin-top:67px; z-index: 104; width: 938px;">
                                <table width="100%"  cellpadding="1" cellspacing="1" border="0" style="background-color:#cccccc;">
                                    <tr align="center" class="fonte_8pt" style="background-color: #F7F7F7; font-weight: bold">
                                        <td colspan="2">Parcela</td>
                                        <td colspan="2">Fundo Comum</td>
                                        <td colspan="2">Taxa Administra&ccedil;&atilde;o</td>
                                        <td colspan="2">Taxa Ades&atilde;o</td>
                                        <td colspan="2">Fundo Reserva</td>
                                        <td>Seguro</td>
                                        <td>Total</td>
                                    </tr>
                                    <tr align="center" class="fonte_8pt" style="background-color: #F7F7F7; font-weight: bold">
                                        <td style="height: 13px">Inicial</td>
                                        <td style="height: 13px">Final</td>
                                        <td style="height: 13px">%</td>
                                        <td style="height: 13px">R$</td>
                                        <td style="height: 13px">%</td>
                                        <td style="height: 13px">R$</td>
                                        <td style="height: 13px">%</td>
                                        <td style="height: 13px">R$</td>
                                        <td style="height: 13px">%</td>
                                        <td style="height: 13px">R$</td>
                                        <td style="height: 13px">R$</td>
                                        <td style="height: 13px">R$</td>
                                    </tr>


                                    <tr style="background-color:#ffffff;" align="right" class="fonte_8pt">
                                        <td style="text-align:center;">001</td>
                                        <td style="text-align:center;">012</td>
                                        <td>1,9451</td>
                                        <td>1.114,13</td>
                                        <td>0,7991</td>
                                        <td>457,71</td>
                                        <td>0,0000</td>
                                        <td>0,00</td>
                                        <td>0,0698</td>
                                        <td>39,98</td>

                                        <td>
                                            52,92

                                        </td>

                                        <td>
                                            1.664,74

                                        </td>
                                    </tr>

                                    <tr style="background-color:#ffffff;" align="right" class="fonte_8pt">
                                        <td style="text-align:center;">013</td>
                                        <td style="text-align:center;">042</td>
                                        <td>2,4729</td>
                                        <td>1.416,45</td>
                                        <td>0,2713</td>
                                        <td>155,40</td>
                                        <td>0,0000</td>
                                        <td>0,00</td>
                                        <td>0,0698</td>
                                        <td>39,98</td>

                                        <td>
                                            52,92

                                        </td>

                                        <td>
                                            1.664,75

                                        </td>
                                    </tr>

                                    <tr style="background-color:#ffffff;" align="right" class="fonte_8pt">
                                        <td style="text-align:center;">043</td>
                                        <td style="text-align:center;">043</td>
                                        <td>2,4718</td>
                                        <td>1.415,82</td>
                                        <td>0,2718</td>
                                        <td>155,68</td>
                                        <td>0,0000</td>
                                        <td>0,00</td>
                                        <td>0,0684</td>
                                        <td>39,18</td>

                                        <td>
                                            52,92

                                        </td>

                                        <td>
                                            1.663,60

                                        </td>
                                    </tr>

                                    <tr class="fonte_8pt" style="background-color: #F7F7F7">
                                        <td align="center" colspan="2" style="background-color: buttonface; font-weight: bold">SubTotal</td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_FC_Sub">100,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_FC_Sub">57.278,88</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_TA_Sub">18,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_TA_Sub">10.310,20</span></td>
                                        <td align="right">
                                            0,0000</td>
                                        <td align="right">
                                            0,00</td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_FR_Sub">3,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_FR_Sub">1.718,34</span></td>
                                        <td align="right">
                                            0,00</td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_Total_Sub">69.307,42</span></td>
                                    </tr>
                                    <tr class="fonte_8pt" style="background-color: #F7F7F7">
                                        <td align="center" colspan="2" style="background-color: buttonface; font-weight: bold">% A Cobrar na Contemplação</td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_Contemplacao">0,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_Contemplacao">0,00</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_TA_Contemplacao">0,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_TA_Contemplacao">0,00</span></td>
                                        <td align="right">
                                            0,0000</td>
                                        <td align="right">
                                            0,00</td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_FR_Contemplacao">0,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_FR_Contemplacao">0,00</span></td>
                                        <td align="right">
                                            0,00</td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_Total_Contemplacao">0,00</span></td>
                                    </tr>
                                    <tr class="fonte_8pt" style="background-color: #F7F7F7">
                                        <td align="center" colspan="2" style="background-color: buttonface; font-weight: bold">Total</td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_FC_Total">100,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_FC_Total">57.278,88</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_TA_Total">18,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_TA_Total">10.310,20</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_AD_Total">0,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_AD_Total">0,00</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblPE_FR_Total">3,0000</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_FR_Total">1.718,34</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_SG_Total">2.275,56</span></td>
                                        <td align="right">
                                            <span id="ctl00_Conteudo_lblVL_Total">71.582,98</span></td>
                                    </tr>

                                </table>
                            </div>

                            <div class="dashboard" style="z-index: 103; position: relative; width: 938px; height: 35px;">
                                <input type="submit" name="ctl00$Conteudo$btnVoltar" value="« Voltar para a lista" id="ctl00_Conteudo_btnVoltar" class="button" style="width:140px;z-index: 888; position: absolute; top: 6px; right: 145px;" />
                                <input type="submit" name="ctl00$Conteudo$btnImprimir" value="Imprimir Detalhes" id="ctl00_Conteudo_btnImprimir" class="button" style="width:140px;z-index: 888; position: absolute; top: 6px; right: 0px;" />
                            </div>
                        </div>
                    </div>
                </div>
                <!--</formulario>-->

                <span id="ctl00_lblTipoBase" style="display:inline-block;font-weight:bold;width:330px;z-index: 101; font-size:14px; font-family:Arial; color: #8B0000; left: 552px; position: absolute; top: 13px; text-align:right; right: 78px;"></span>
                <input type="image" name="ctl00$img_AlteraSenha" id="ctl00_img_AlteraSenha" title="Editar a senha" src="../imagens/Tela_Inicial/img_AlteraSenha.png" style="z-index: 101; position:absolute; top:13px; right:14px;" />

                <input type="image" name="ctl00$img_Atendimento" id="ctl00_img_Atendimento" title="Atendimento" src="../imagens/Tela_Inicial/atendimento.png" style="z-index: 101; position:absolute; top:-37px; left:7px;" />

            </div>
        </div>
    </div>

    <div class="theme_container" >
        <div class="inner" style="z-index:0;position:relative;">

            <!-- POPUPS -->



        </div>
    </div>


    <div class="aspNetHidden">

        <input type="hidden" name="__EVENTVALIDATION" id="__EVENTVALIDATION" value="4Z4wWonQFjnyTvwv9NU2n9s3aGT1w/YCnSzagNhw+jkOyOGC/V+8LIADANP64oE6uUFKQO6ogkO7F+/Q3XXkAvT1B0ElJ9FKOTfx42ivFuZCtLdWfKOEoY5PkHGnCiL29ztDEtlQm7kNDg4r9J3Yv0glcPUPqKKhGf1V2eRC1PD4fqca0gqiCjRvdTBtzz2WfRWGEoDgGdEVpB7tR/efGLSYQWQ=" />
    </div></form>
</body>
<script type="text/javascript">
    if ($("#ctl00_img_Relatorio").length == 0) {
        $("#subs").css("left", "44px");
    }
    else
    {
        if ($("#ctl00_img_Relatorio").length == 0) {
            $("#subs").css("left", "80px");
        }
    }
</script>
</html>

`
	if err := simulacaoDetalha(&sim, doc); err != nil {
		t.Fatal(err)
	}
	fmt.Println(sim)

}
