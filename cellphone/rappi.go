package cellphone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Error struct {
		VerificationValue string `json:"verification_value"`
	} `json:"error"`
}

func Rappi(email string) string {
	url := "https://services.rappi.com.br/api/rocket/login/email/application_user"

	payload := map[string]string{
		"email": email,
		"scope": "all",
	}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Erro na requisição:", err)
		return ""
	}

	req.Header.Set("authority", "services.rappi.com.br")
	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-language", "pt-BR")
	req.Header.Set("access-control-allow-headers", "*")
	req.Header.Set("access-control-allow-origin", "*")
	req.Header.Set("app-version", "web_v1.40.2")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("needappsflyerid", "true")
	req.Header.Set("origin", "https://www.rappi.com.br")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://www.rappi.com.br/")
	req.Header.Set("sec-ch-ua", "\"Brave\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao enviar requisição:", err)
		return ""
	}
	defer resp.Body.Close()

	var responseObj Response
	err = json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		fmt.Println("Erro ao decodificar resposta:", err)
		return ""
	}

	if responseObj.Error.VerificationValue != "" {
		return responseObj.Error.VerificationValue
	} else {
		return ""
	}
}
