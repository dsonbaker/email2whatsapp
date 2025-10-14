package cellphone

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
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
		errorUser := ""
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`#input-1`, chromedp.ByID), // substitua 'inputID' pelo ID do seu elemento de entrada
			chromedp.Sleep(1*time.Second),
			chromedp.SendKeys(`#input-1`, email, chromedp.ByID),
			chromedp.Sleep((15/10)*time.Second),
			chromedp.KeyEvent(kb.Enter),
			chromedp.WaitVisible(`.FormResetPasswordIdentification-Form .error, #SelectTruncatedPhoneOrEmail-phoneText`, chromedp.ByQuery),
			chromedp.Evaluate(`document.querySelector("#SelectTruncatedPhoneOrEmail-phoneText")?document.querySelector("#SelectTruncatedPhoneOrEmail-phoneText").innerText:""`, &Leak_phoneNumber),
			chromedp.Evaluate(`document.querySelector(".FormResetPasswordIdentification-Form .error")?document.querySelector(".FormResetPasswordIdentification-Form .error").innerText:""`, &errorUser),
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
