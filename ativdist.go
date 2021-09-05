package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	urlBase    = "https://www2.bauru.sp.gov.br"
	retries    = 40
	sleepRetry = time.Second * 5
	dir        = "C:\\GoPrograms\\ativdist"
)

var (
	client *http.Client
	cookie string
)

func init() {
	var err error
	cjar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Jar: cjar, Transport: tr}
}

type hashes struct {
	applicationKey     string
	lastFocus          string
	eventTarget        string
	eventArgument      string
	vState             string
	viewState          string
	viewStateGenerator string
	eventValidation    string
}

func getHashes(doc string) (hashes, error) {
	var h hashes
	linhas := strings.Split(doc, "\n")
	for _, l := range linhas {
		if strings.Contains(l, `?applicationKey=`) {
			matches := erAction.FindAllStringSubmatch(l, -1)
			if len(matches) != 1 {
				return hashes{}, fmt.Errorf("gethashes: len applicationKey != 1")
			}
			h.applicationKey = strings.Replace(
				matches[0][1],
				"/NewconWeb/frmCorCCCnsLogin.aspx?applicationKey=",
				"",
				1,
			)
			h.applicationKey = strings.Replace(
				h.applicationKey,
				"/Newconweb/frmCorCCCnsLogin.aspx?applicationKey=",
				"",
				1,
			)
		}
		if strings.Contains(l, `"__LASTFOCUS"`) {
			matches := erValue.FindAllStringSubmatch(l, -1)
			if len(matches) != 1 {
				return hashes{}, fmt.Errorf("gethashes: len last focus != 1")
			}
			h.lastFocus = matches[0][1]
		}
		if strings.Contains(l, `"__EVENTTARGET"`) {
			matches := erValue.FindAllStringSubmatch(l, -1)
			if len(matches) != 1 {
				return hashes{}, fmt.Errorf("gethashes: len event target != 1")
			}
			h.eventTarget = matches[0][1]
		}
		if strings.Contains(l, `"__EVENTARGUMENT"`) {
			matches := erValue.FindAllStringSubmatch(l, -1)
			if len(matches) != 1 {
				return hashes{}, fmt.Errorf("gethashes: len event argument != 1")
			}
			h.eventArgument = matches[0][1]
		}
		if strings.Contains(l, `"__VIEWSTATE"`) {
			matches := erValue.FindAllStringSubmatch(l, -1)
			if len(matches) != 1 {
				return hashes{}, fmt.Errorf("gethashes: len view state != 1")
			}
			h.viewState = matches[0][1]
		}
		if strings.Contains(l, `"__VSTATE"`) {
			matches := erValue.FindAllStringSubmatch(l, -1)
			if len(matches) != 1 {
				return hashes{}, fmt.Errorf("gethashes: len v state != 1")
			}
			h.vState = matches[0][1]
		}
		if strings.Contains(l, `"__VIEWSTATEGENERATOR"`) {
			matches := erValue.FindAllStringSubmatch(l, -1)
			if len(matches) != 1 {
				return hashes{}, fmt.Errorf("gethashes: len view state generator != 1")
			}
			h.viewStateGenerator = matches[0][1]
		}
		if strings.Contains(l, `"__EVENTVALIDATION"`) {
			matches := erValue.FindAllStringSubmatch(l, -1)
			if len(matches) != 1 {
				return hashes{}, fmt.Errorf("gethashes: len event validation != 1")
			}
			h.eventValidation = matches[0][1]
		}
		if strings.Contains(l, `"__VIEWSTATEENCRYPTED"`) {
			matches := erValue.FindAllStringSubmatch(l, -1)
			if len(matches) != 1 {
				return hashes{}, fmt.Errorf("gethashes: len view state encrypted != 1")
			}
			h.eventValidation = matches[0][1]
		}
	}
	return h, nil
}

func newValuesHashes(h hashes) url.Values {
	v := url.Values{}
	v.Set("__LASTFOCUS", h.lastFocus)
	v.Set("__EVENTTARGET", h.eventTarget)
	v.Set("__EVENTARGUMENT", h.eventArgument)
	v.Set("__VIEWSTATE", h.viewState)
	v.Set("__VSTATE", h.vState)
	v.Set("__VIEWSTATEGENERATOR", h.viewStateGenerator)
	v.Set("__EVENTVALIDATION", h.eventValidation)
	return v
}

func newValuesDoc(doc string) (url.Values, error) {
	h, err := getHashes(doc)
	if err != nil {
		return nil, fmt.Errorf("newvaluesdoc: %w", err)
	}
	return newValuesHashes(h), nil
}

func getRetry(u string) (doc string, newU string, finalErr error) {

	getReadAll := func(u string) ([]byte, string, error) {
		res, err := client.Get(u)
		if err != nil {
			return nil, "", fmt.Errorf("getreadall: %w", err)
		}
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, "", fmt.Errorf("getreadall: %w", err)
		}
		return b, res.Request.URL.String(), nil
	}

	for i := 0; i < retries; i++ {
		b, nu, err := getReadAll(u)
		if err != nil {
			log.Printf("getretry %v: %v\n", i, err)
			finalErr = err
			time.Sleep(sleepRetry)
			continue
		}
		newU = nu
		doc = string(b)
		return doc, newU, nil
	}
	return "", "", finalErr
}

func postRetry(u string, values url.Values) (doc string, newU string, finalErr error) {
	postReadAll := func(u string, data url.Values) ([]byte, string, error) {
		res, err := client.PostForm(u, data)
		if err != nil {
			return nil, "", fmt.Errorf("postreadall: %w", err)
		}
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, "", fmt.Errorf("postreadall: %w", err)
		}
		return b, res.Request.URL.String(), nil
	}

	for i := 0; i < retries; i++ {
		b, nu, err := postReadAll(u, values)
		if err != nil {
			log.Printf("postretry %v: %v\n", i, err)
			finalErr = err
			time.Sleep(sleepRetry)
			continue
		}
		newU = nu
		doc = string(b)
		return doc, newU, nil
	}
	return "", "", finalErr
}
