package cellphone

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"context"
	"time"
	"log"
)

func Magalu(email string) string {
	maxTrys := 2
	url := "https://sacola.magazineluiza.com.br/n#/recuperar-senha/?"	
	//currentTime := time.Now()
	//formattedTime := currentTime.Format("2006-01-02 15:04:05")
	//fmt.Println("["+formattedTime+"]", "[URL] [TRY]", url)
	countBotsDetected := 0
	Leak_phoneNumber := ""
	var options []func(*chromedp.ExecAllocator)
	if countBotsDetected >= 1 {
		fmt.Println("[!!!] Required User Interaction")
		options = []chromedp.ExecAllocatorOption{
			chromedp.Flag("ignore-certificate-errors", "1"),
			chromedp.Flag("headless", false), // set headless to false
			chromedp.Flag("disable-gpu", true),
		}
	} else{
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
		errorUser := ""
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`#identificationReset`, chromedp.ByID), // substitua 'inputID' pelo ID do seu elemento de entrada
			chromedp.Sleep(1*time.Second),
			chromedp.SendKeys(`#identificationReset`, email, chromedp.ByID),
			chromedp.Sleep((15/10)*time.Second),
			chromedp.KeyEvent(kb.Enter),
			chromedp.WaitVisible(`.FormGroup-errorMessage, .SelectTruncatedPhoneOrEmail-PhoneNumber`, chromedp.ByQuery),
			chromedp.Evaluate(`document.getElementsByClassName("SelectTruncatedPhoneOrEmail-PhoneNumber")[0]?document.getElementsByClassName("SelectTruncatedPhoneOrEmail-PhoneNumber")[0].innerText:""`, &Leak_phoneNumber),
			chromedp.Evaluate(`document.getElementsByClassName("FormGroup-errorMessage")[0]?document.getElementsByClassName("FormGroup-errorMessage")[0].innerText:""`, &errorUser),
			)
		if err != nil {
			log.Println(err)
			if i != maxTrys {
				log.Println("[/] Tentando novamente:", email)
				continue
			}
		}
		defer cancel()
		break
	}
	return Leak_phoneNumber
}