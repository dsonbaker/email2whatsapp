package automationWhatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func Run() {

	url := "https://web.whatsapp.com/"
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false), // set headless to false
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-cache", true),
	}
	ctx, cancel = chromedp.NewExecAllocator(ctx, options...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(url), chromedp.WaitVisible(`#pane-side`, chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}
	err := chromedp.Run(ctx,
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`document.querySelectorAll("[role=button]")[4]`, chromedp.ByJSPath),
		chromedp.Evaluate(`document.querySelectorAll("[role=button]")[4].click()`, nil),
		chromedp.Sleep(3*time.Second),
		//chromedp.SendKeys(`input:nth-of-type(1)`, str[len(str)-1:len(str)], chromedp.ByQuery),
		chromedp.WaitVisible(`document.querySelector('[role=textbox]')`, chromedp.ByJSPath),
		chromedp.Evaluate(`document.querySelectorAll('[role=textbox]')[0].click()`, nil),
		chromedp.SendKeys(`document.querySelectorAll('[role=textbox]')[0]`, `user`, chromedp.ByJSPath),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys(`document.querySelectorAll('[role=textbox]')[0]`, ` `, chromedp.ByJSPath),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`document.querySelectorAll('[role=textbox]')[0].click()`, nil),
		chromedp.SendKeys(`document.querySelectorAll('[role=textbox]')[0]`, "\b", chromedp.ByJSPath),
		chromedp.Sleep(2*time.Second),
		chromedp.SendKeys(`document.querySelectorAll('[role=textbox]')[0]`, ` `, chromedp.ByJSPath),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	var timeMax string
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`var positionMax = document.querySelector('[data-tab="4"]').parentElement.scrollHeight;
						   var timePerScroll = 2000; //timePerScroll in milliSecond
						   var positionPerScroll = 800;
						   var timeMax = (positionMax/positionPerScroll)*(timePerScroll/1000);
						   timeMax.toString()`, &timeMax),
	)
	if err != nil {
		log.Fatal(err)
	}
	flt, err := strconv.ParseFloat(timeMax, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	var sleepTime = ((time.Duration(flt) + 35) * time.Second)
	var TimeOut = ((time.Duration(flt) + 45) * time.Second)
	ctx, cancel = context.WithTimeout(ctx, TimeOut)
	defer cancel()
	fmt.Println("[-] Estimated time:", sleepTime)
	//document.querySelectorAll("div[role=listitem]").forEach((e)=>{if(e.innerText.includes("user")){
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`
		var users = [];
		var usersUnduplicate = [];
		var positionStart = 0;
		var positionMax = document.querySelector('[data-tab="4"]').parentElement.scrollHeight;
		
		function removeDuplicates(array) {
			let unique = {};
			let result = [];
			
			array.forEach(item => {
			  if (!unique[item.numberphone]) {
				unique[item.numberphone] = true;
				result.push(item);
			  }
			});
			
			return result;
		}

		function execScroll(params) {
		   positionMax = document.querySelector('[data-tab="4"]').parentElement.scrollHeight
		   document.querySelector('[data-tab="4"]').parentElement.scrollTop += positionPerScroll;
		   
		   if (positionStart < positionMax+positionPerScroll){
			   positionStart += positionPerScroll;
			   setTimeout(execScroll,timePerScroll)
			   console.log("scrolling")
			   document.querySelectorAll("div[role=listitem]").forEach((e)=>{if(e.innerText.includes("user")){
					srcImage=""
					if(e.getElementsByTagName("img").length >= 1) {
						if(!e.getElementsByTagName("img")[0].getAttribute("src").includes("data:")){
							srcImage = e.getElementsByTagName("img")[0].src
						}
					}
					users.push({numberphone:e.innerText.split('\n')[0],src:srcImage})
				}})
		   } else {
			   usersUnduplicate = removeDuplicates(users);
			   console.log(usersUnduplicate)
		   }
		}
		
		//execScroll();`, nil),
	)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(sleepTime)
	var usersINFO string
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`JSON.stringify(usersUnduplicate)`, &usersINFO),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Decodifique o JSON para uma estrutura Golang

	var listUsers []map[string]string
	if err := json.Unmarshal([]byte(usersINFO), &listUsers); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println(err)
	}
	quantityUsers := 0
	if _, err := os.Stat("numberphone.txt"); err == nil {
		os.Remove("numberphone.txt")
	}
	for _, user := range listUsers {
		// Acesso correto ao campo NumberPhone dentro de cada estrutura UserINFO
		if numberPhone, ok := user["numberphone"]; ok {
			numberPhone = strings.Split(numberPhone, " ")[1]
			WriteToFile("all-numbers.txt", numberPhone+"\n", "./numberphone/")
			if src, ok := user["src"]; ok {
				if src != "" {
					DownloadFile(src, numberPhone+".jpg", "./numberphone/profile/")
					WriteToFile("numbers-profile.txt", numberPhone+"\n", "./numberphone/")
				} else {
					WriteToFile("numbers-withoutProfile.txt", numberPhone+"\n", "./numberphone/")
				}
			}
			quantityUsers++
		} else {
			fmt.Println("Property 'numberphone' not found in JSON.")
		}
	}
	fmt.Println("[+] Number of users:", quantityUsers)
}

func WriteToFile(filename string, data string, folderName string) error {
	os.MkdirAll(folderName, os.ModePerm)
	filename = filepath.Join(folderName, filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		return err
	}
	return nil
}
