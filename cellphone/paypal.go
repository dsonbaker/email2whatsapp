package cellphone

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

func Paypal(email string) string {
	url := "https://www.paypal.com/authflow/password-recovery/?country.x=BR&locale.x=pt_BR&redirectUri=%252Fsignin"
	var options []func(*chromedp.ExecAllocator)
	options = []chromedp.ExecAllocatorOption{
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
	ctx, cancel = context.WithTimeout(ctx, 80*time.Second)
	defer cancel()
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	)
	if err != nil {
		log.Fatal(err)
	}
	PhoneNumber := ""
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`#pwrStartPageEmail`, chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys(`#pwrStartPageEmail`, email, chromedp.ByID),
		chromedp.Sleep((15/10)*time.Second),
		chromedp.KeyEvent(kb.Enter),
		chromedp.WaitReady(`#message_pwrStartPageEmail, .verification-method`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
		chromedp.Evaluate(`document.getElementsByClassName("verification-method")[0]?document.getElementsByClassName("verification-method")[0].innerText.split(" ").slice(document.getElementsByClassName("verification-method")[0].innerText.split(" ").length-2,document.getElementsByClassName("verification-method")[0].innerText.split(" ").length).join("").replaceAll("•","*").replaceAll("-","").replace(/.{1}$/,""):""`, &PhoneNumber),
	)
	if err != nil {
		log.Fatal(err)
	}
	return PhoneNumber
}
