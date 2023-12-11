package cellphone

import (
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"context"
	"time"
	"log"
)

func Pagbank(email string) string {
	url := "https://minhasenha.pagseguro.uol.com.br/recuperar-senha"	
	//currentTime := time.Now()
	//formattedTime := currentTime.Format("2006-01-02 15:04:05")
	//fmt.Println("["+formattedTime+"]", "[URL] [TRY]", url)
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
			chromedp.WaitVisible(`#credential-input`, chromedp.ByID),
			chromedp.Sleep(1*time.Second),
			chromedp.SendKeys(`#credential-input`, email, chromedp.ByID),
			chromedp.Sleep((15/10)*time.Second),
			chromedp.KeyEvent(kb.Enter),
			chromedp.WaitVisible(`[data-cy=credential-error-alert], #sms-factor`, chromedp.ByQuery),
			chromedp.Evaluate(`document.getElementById("sms-factor")?document.getElementById("sms-factor").innerText.split("\n")[1].replace(" ","").replace("(","").replace(")","").replace("-",""):""`, &PhoneNumber),
		)
	if err != nil {
		log.Fatal(err)
	}
	return PhoneNumber
}