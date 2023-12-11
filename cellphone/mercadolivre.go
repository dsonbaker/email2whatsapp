package cellphone

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

func Mercadolivre(email string) string {
	maxTrys := 2
	url := "https://www.mercadolivre.com.br/"
	countBotsDetected := 0
	PhoneNumber := ""
	var options []func(*chromedp.ExecAllocator)
	if countBotsDetected >= 1 {
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
		userNOTexist := ""
		cameraRequired := ""
		withoutCode := ""
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`body`, chromedp.ByQuery), // substitua 'inputID' pelo ID do seu elemento de entrada
			chromedp.Sleep(1*time.Second),
			chromedp.Evaluate(`document.body.querySelectorAll("a[data-link-id='login']")[0].click()`, nil),
		)
		if err != nil {
			log.Println(err)
			if i != maxTrys {
				log.Println("[/] Tentando novamente:", email)
				continue
			}
		}
		defer cancel()
		err = chromedp.Run(ctx,
			chromedp.WaitVisible(`#user_id`, chromedp.ByID),
			chromedp.Sleep(1*time.Second),
			chromedp.SendKeys(`#user_id`, email, chromedp.ByID),
			chromedp.Sleep(1*time.Second),
			chromedp.KeyEvent(kb.Enter),
			//chromedp.Evaluate(`document.body.querySelectorAll("button[type='submit']")[0].click()`, nil),
			chromedp.Sleep(1*time.Second),
			chromedp.WaitReady(`#rc-anchor-container, #code_validation, .input-error, .camera-not-found__image, [aria-labelledby=code-numeric]`, chromedp.ByQuery),
			chromedp.Evaluate(`document.querySelector("[aria-labelledby='code-numeric']")?"withoutCode":""`, &withoutCode),
			chromedp.Evaluate(`document.getElementsByClassName("recaptcha__error-icon")[0]?"botDetected":""`, &botDetected),
			chromedp.Evaluate(`document.getElementsByClassName("camera-not-found__image")[0]?"cameraRequired":""`, &cameraRequired),
			chromedp.Evaluate(`document.getElementsByClassName("input-error")[0].getElementsByClassName("ui-form__message")[0]?"notExist":""`, &userNOTexist),
		)
		if err != nil {
			log.Println(err)
			if i != maxTrys {
				log.Println("[-] Tentando novamente:", email)
				continue
			}
		}
		if withoutCode == "" {
			if cameraRequired == "" {
				if botDetected == "botDetected" {
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
								log.Println("[Error] Verfique o captcha: ", email)
								continue
							}
						}
						botDetected = ""
						err = chromedp.Run(ctx,
							chromedp.SendKeys(`#user_id`, "", chromedp.ByID),
							chromedp.Sleep(1*time.Second),
							chromedp.KeyEvent(kb.Enter),
							chromedp.WaitVisible(`#action-decline, #code_validation, .input-error, .camera-not-found__image, [aria-labelledby=code-numeric]`, chromedp.ByQuery),
							chromedp.Evaluate(`document.querySelector("[aria-labelledby='code-numeric']")?"withoutCode":""`, &withoutCode),
							chromedp.Evaluate(`document.getElementsByClassName("camera-not-found__image")[0]?(document.getElementsByClassName("camera-not-found__image")[0]?"cameraRequired":""):""`, &cameraRequired),
							chromedp.Evaluate(`document.getElementsByClassName("input-error")[0]?(document.getElementsByClassName("input-error")[0].getElementsByClassName("ui-form__message")[0]?"notExist":""):""`, &userNOTexist),
							chromedp.Evaluate(`document.getElementById("action-decline")?document.getElementById("action-decline").click():""`, nil),
							chromedp.WaitVisible(`#code_validation, .input-error, .camera-not-found__image, [aria-labelledby=code-numeric]`, chromedp.ByQuery),
							chromedp.Evaluate(`document.getElementById("whatsapp")?document.getElementById("whatsapp").innerText.split(" ")[document.getElementById("whatsapp").innerText.split(" ").length-1].replace(".",""):""`, &PhoneNumber),
						)
						if err != nil {
							log.Println(err)
							if i != maxTrys {
								log.Println("[Error] Verfique o captcha [1]: ", email)
								continue
							}
						}
					}
				} //recaptcha-checkbox-checkmark
				if botDetected != "botDetected" && userNOTexist == "" {
					countBotsDetected = 0
				}
				if botDetected == "" && userNOTexist == "notExist" {
					countBotsDetected = 0
				}
			} else {
				fmt.Println("[-] Camera Required.")
			}
		} else {
			fmt.Println("[-] Not Exist NumberPhone.")
		}
		defer cancel()
		break
	}
	return PhoneNumber
}
