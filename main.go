package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dsonbaker/email2whatsapp/automationWhatsapp"
	"github.com/dsonbaker/email2whatsapp/bruteforceSite"
	"github.com/dsonbaker/email2whatsapp/cellphone"
)

func main() {
	verde := "\033[32m"
	email := flag.String("email", "", "Target email")
	whatsapp := flag.Bool("whatsapp", false, "Whatsapp Automation Mode")
	bruteforce := flag.String("bruteforce", "", "Select one of the sites for bruteforce: [paypal, meli, twitter, google]")

	flag.Parse()
	if *email == "" && !*whatsapp && *bruteforce == "" {
		fmt.Println("[-] You must provide the --email flag or the --whatsapp flag.")
		os.Exit(1)
	}
	if *email != "" {
		PrintInfo(verde, "[+] Looking for Email: "+*email)
		searchLeakedNumbers(*email)
	}

	if *whatsapp {
		fmt.Println("[+] Automate Whatsapp.")
		automationWhatsapp.Run()
	}
	if *bruteforce != "" {
		PrintInfo(verde, "[+] Looking for Email: "+*bruteforce)
		if *bruteforce != "paypal" && *bruteforce != "meli" && *bruteforce != "twitter" && *bruteforce != "google" {
			fmt.Println("[-] Insert paypal, meli, twitter or google")
			os.Exit(1)
		}
		if *bruteforce == "paypal" {
			bruteforceSite.BrutePaypal()
		}
		if *bruteforce == "meli" {
			bruteforceSite.BruteMercadoLivre()
		}
		if *bruteforce == "twitter" {
			bruteforceSite.BruteTwitter()
		}
		if *bruteforce == "google" {
			bruteforceSite.BruteGoogle()
		}
	}
}

func searchLeakedNumbers(email string) {
	numberphoneBR := [][]string{{"*", "*"}, {"9", "*", "*", "*", "*", "*", "*", "*", "*"}}
	possibleNumbers := []string{}
	vermelho := "\033[31m"
	verde := "\033[32m"
	numberShow := ""
	// Magazine Luiza
	PrintInfo(verde, "[+] Searching on MagazineLuiza.")
	magaluPhone := cellphone.Magalu(email)
	if magaluPhone != "" {
		PrintInfo(vermelho, "[!] Found Number: "+magaluPhone)
	}
	// Paypal
	PrintInfo(verde, "[+] Searching on Paypal.")
	paypalPhone := cellphone.Paypal(email)
	if paypalPhone != "" {
		PrintInfo(vermelho, "[!] Found Number: "+paypalPhone)
	}
	// PagBank
	PrintInfo(verde, "[+] Searching on PagBank.")
	pagbankPhone := cellphone.Pagbank(email)
	if pagbankPhone != "" {
		PrintInfo(vermelho, "[!] Found Number: "+pagbankPhone)
	}
	// Mercado Livre
	PrintInfo(verde, "[+] Searching on MercadoLivre.")
	mercadolivrePhone := cellphone.Mercadolivre(email)
	if mercadolivrePhone != "" {
		PrintInfo(vermelho, "[!] Found Number: "+mercadolivrePhone)
	}
	// Rappi
	PrintInfo(verde, "[+] Searching on Rappi.")
	rappiPhone := cellphone.Rappi(email)
	if rappiPhone != "" {
		PrintInfo(vermelho, "[!] Found Number: "+rappiPhone)
	}

	if len(magaluPhone) > 1 {
		numberphoneBR[0][0] = string(magaluPhone[0])
		numberphoneBR[0][1] = string(magaluPhone[1])
		numberphoneBR[1][1] = string(magaluPhone[3])
		numberphoneBR[1][2] = string(magaluPhone[4])
		numberphoneBR[1][3] = string(magaluPhone[5])
		numberShow = showNumberPhoneBR(numberphoneBR)
		PrintInfo(verde, "[+] Magalu, Possible Combination: "+numberShow)
		//possibleNumbers = append(possibleNumbers, numberShow)
		numberShow = ""
	}
	if len(paypalPhone) > 1 {
		diffNumbers := true
		numberphoneBR[0][0] = string(paypalPhone[0])
		numberphoneBR[1][4] = string(paypalPhone[len(paypalPhone)-5])
		numberphoneBR[1][5] = string(paypalPhone[len(paypalPhone)-4])
		numberphoneBR[1][6] = string(paypalPhone[len(paypalPhone)-3])
		numberphoneBR[1][7] = string(paypalPhone[len(paypalPhone)-2])
		numberphoneBR[1][8] = string(paypalPhone[len(paypalPhone)-1])
		if len(magaluPhone) > 1 {
			if string(paypalPhone[0]) == string(magaluPhone[0]) {
				diffNumbers = false
				numberphoneBR[0][1] = string(magaluPhone[1]) //magalu
			}
		}
		if len(pagbankPhone) > 1 {
			if string(paypalPhone[len(paypalPhone)-4:]) == string(pagbankPhone[len(pagbankPhone)-4:]) {
				diffNumbers = false
				numberphoneBR[0][1] = string(pagbankPhone[1])
			}
		}
		if diffNumbers {
			numberphoneBR[0][1] = "*"
		}
		numberShow = showNumberPhoneBR(numberphoneBR)
		PrintInfo(verde, "[+] Paypal, Possible Combination: "+numberShow)
		possibleNumbers = append(possibleNumbers, numberShow)
		numberShow = ""
	}
	if len(pagbankPhone) > 1 {
		newNumber := false
		if len(paypalPhone) > 1 {
			if string(pagbankPhone[len(pagbankPhone)-4:]) != string(paypalPhone[len(paypalPhone)-4:]) {
				newNumber = true
			}
		}
		if len(paypalPhone) < 1 {
			newNumber = true
		}
		if newNumber {
			numberphoneBR[0][0] = string(pagbankPhone[0])
			numberphoneBR[0][1] = string(pagbankPhone[1])
			numberphoneBR[1][4] = "*"
			numberphoneBR[1][5] = string(pagbankPhone[len(pagbankPhone)-4])
			numberphoneBR[1][6] = string(pagbankPhone[len(pagbankPhone)-3])
			numberphoneBR[1][7] = string(pagbankPhone[len(pagbankPhone)-2])
			numberphoneBR[1][8] = string(pagbankPhone[len(pagbankPhone)-1])
			numberShow = showNumberPhoneBR(numberphoneBR)
			PrintInfo(verde, "[+] Pagbank, Possible Combination: "+numberShow)
			possibleNumbers = append(possibleNumbers, numberShow)
			numberShow = ""
		}
	}
	if len(mercadolivrePhone) > 1 {
		newNumber := false
		if len(paypalPhone) > 1 {
			if string(mercadolivrePhone[len(mercadolivrePhone)-4:]) != string(paypalPhone[len(paypalPhone)-4:]) {
				newNumber = true
			}
		}
		if len(pagbankPhone) > 1 {
			if string(mercadolivrePhone[len(mercadolivrePhone)-4:]) != string(pagbankPhone[len(pagbankPhone)-4:]) {
				newNumber = true
			}
		}
		if len(paypalPhone) < 1 && len(pagbankPhone) < 1 {
			newNumber = true
		}
		if newNumber {
			numberphoneBR[1][5] = string(mercadolivrePhone[len(mercadolivrePhone)-4])
			numberphoneBR[1][6] = string(mercadolivrePhone[len(mercadolivrePhone)-3])
			numberphoneBR[1][7] = string(mercadolivrePhone[len(mercadolivrePhone)-2])
			numberphoneBR[1][8] = string(mercadolivrePhone[len(mercadolivrePhone)-1])
			numberShow = showNumberPhoneBR(numberphoneBR)
			PrintInfo(verde, "[+] Meli, Possible Combination: "+numberShow)
			possibleNumbers = append(possibleNumbers, numberShow)
			numberShow = ""
		}
	}
	if len(rappiPhone) > 1 {
		newNumber := false
		if len(paypalPhone) > 1 {
			if string(rappiPhone[len(rappiPhone)-4:]) != string(paypalPhone[len(paypalPhone)-4:]) {
				newNumber = true
			}
		}
		if len(pagbankPhone) > 1 {
			if string(rappiPhone[len(rappiPhone)-4:]) != string(pagbankPhone[len(pagbankPhone)-4:]) {
				newNumber = true
			}
		}
		if len(mercadolivrePhone) > 1 {
			if string(rappiPhone[len(rappiPhone)-4:]) != string(mercadolivrePhone[len(mercadolivrePhone)-4:]) {
				newNumber = true
			}
		}
		if newNumber {
			numberphoneBR[1][5] = string(mercadolivrePhone[len(mercadolivrePhone)-4])
			numberphoneBR[1][6] = string(mercadolivrePhone[len(mercadolivrePhone)-3])
			numberphoneBR[1][7] = string(mercadolivrePhone[len(mercadolivrePhone)-2])
			numberphoneBR[1][8] = string(mercadolivrePhone[len(mercadolivrePhone)-1])
			numberShow = showNumberPhoneBR(numberphoneBR)
			PrintInfo(verde, "[+] Rappi, Possible Combination: "+numberShow)
			possibleNumbers = append(possibleNumbers, numberShow)
			numberShow = ""
		}
	}

	if len(possibleNumbers) > 0 {
		numberUsers := exportContactsBR(possibleNumbers)
		PrintInfo(verde, "[+] The contact list has \""+strconv.Itoa(numberUsers)+"\" cellphone numbers.")
	} else {
		PrintInfo(vermelho, "[+] Unable to find result for email: "+email)
	}
}

func PrintInfo(color string, text string) {
	fmt.Println(color + text + "\033[0m")
}

func showNumberPhoneBR(numberphoneBR [][]string) string {
	numberShow := ""
	for _, ddd := range numberphoneBR[0] {
		numberShow += ddd
	}
	for _, number := range numberphoneBR[1] {
		numberShow += number
	}
	return numberShow
}

func generateDDD_BR(ddd string, wildcardNumber string) []string {
	listDDD := []string{"11", "12", "13", "14", "15", "16", "17", "18", "19", "21", "22", "24", "27", "28", "31", "32", "33", "34", "35", "37", "38", "41", "42", "43", "44", "45", "46", "47", "48", "49", "51", "53", "54", "55", "61", "62", "63", "64", "65", "66", "67", "68", "69", "71", "73", "74", "75", "77", "79", "81", "82", "83", "84", "85", "86", "87", "88", "89", "91", "92", "93", "94", "95", "96", "97", "98", "99"}
	possibleDDD := []string{}
	vermelho := "\033[31m"

	if ddd == "**" {
		var num int
		fmt.Print(vermelho, "[!] No DDD digit was found for the number, try to find the possible state of the person, using other OSINT techniques:", "\033[0m")
		_, err := fmt.Scan(&num)
		if err != nil {
			log.Fatal(err)
		}
		if num >= 10 || num <= 99 {
			ddd = strconv.Itoa(num)
		}
		fmt.Println()
	}

	if string(ddd[0]) != "*" && string(ddd[1]) == "*" {
		for _, selectDDD := range listDDD {
			if ddd[0] == selectDDD[0] {
				possibleDDD = append(possibleDDD, selectDDD+wildcardNumber)
				//fmt.Println("[+] Possibilidade DDD: " + selectDDD)
			}
		}
	}
	if string(ddd[0]) != "*" && string(ddd[1]) != "*" {
		for _, selectDDD := range listDDD {
			if ddd == selectDDD {
				possibleDDD = append(possibleDDD, selectDDD+wildcardNumber)
				//fmt.Println("[+] DDD Encontrado: " + selectDDD)
			}
		}
	}
	return possibleDDD
}

func generateCombinationsNumber_BR(numberUnknown string) []string {
	var combinations []string

	index := strings.Index(numberUnknown, "*")
	if index == -1 {
		combinations = append(combinations, numberUnknown)
		return combinations
	}

	for i := 0; i <= 9; i++ {
		newInput := strings.Replace(numberUnknown, "*", strconv.Itoa(i), 1)
		combinations = append(combinations, generateCombinationsNumber_BR(newInput)...)
	}

	return combinations
}

func exportContactsBR(possibleNumbers []string) int {
	numberUsers := 0
	if _, err := os.Stat("possible_numbers.txt"); err == nil {
		os.Remove("possible_numbers.txt")
	}
	for _, number := range possibleNumbers {
		numbersWithDDD := generateDDD_BR(string(number[0])+string(number[1]), string(number[2:]))
		for _, numberWithDDD := range numbersWithDDD {
			combinationNumbers := generateCombinationsNumber_BR(numberWithDDD)
			for _, combo := range combinationNumbers {
				combo = "55" + combo
				err := WriteToFile("possible_numbers.txt", combo+"\n")
				if err != nil {
					log.Fatal(err)
				}
				numberUsers++
			}
		}
	}
	return numberUsers
}

func WriteToFile(filename string, data string) error {
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
