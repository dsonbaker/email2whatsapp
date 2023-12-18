package automationWhatsapp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"

	waLog "go.mau.fi/whatsmeow/util/log"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func Run() {
	listPhones := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		listPhones = append(listPhones, "+"+scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro de leitura:", err)
	}
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
	quantityUsers := 0
	RemoveFile("all-numbers.txt")
	RemoveFile("numbers-profile.txt")
	RemoveFile("numbers-withoutProfile.txt")
	for _, numberphone := range listPhones {
		IsOnWhatsAppResponse, errIsOnWhatsApp := client.IsOnWhatsApp([]string{numberphone})
		if errIsOnWhatsApp != nil {
			panic(errIsOnWhatsApp)
		}
		if IsOnWhatsAppResponse[0].IsIn {
			GetProfilePictureInfoResponse, errGetProfile := client.GetProfilePictureInfo(IsOnWhatsAppResponse[0].JID, nil)
			if errGetProfile != nil {
				fmt.Println("error:", errGetProfile.Error())
				panic(errGetProfile)
			}
			WriteToFile("all-numbers.txt", numberphone+"\n", "./numberphone/")
			if GetProfilePictureInfoResponse.URL != "" {
				DownloadFile(GetProfilePictureInfoResponse.URL, numberphone+".jpg", "./numberphone/profile/")
				WriteToFile("numbers-profile.txt", numberphone+"\n", "./numberphone/")
				fmt.Println(GetProfilePictureInfoResponse.URL)
			} else {
				WriteToFile("numbers-withoutProfile.txt", numberphone+"\n", "./numberphone/")
			}
			quantityUsers++
		}
	}
	fmt.Println("\033[32m[+] Number of users:", quantityUsers, "\033[0m")
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
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
func DownloadFile(targetURL string, nameFile string, folderName string) {

	response, err := http.Get(targetURL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	os.MkdirAll(folderName, os.ModePerm)
	arquivo, err := os.Create(filepath.Join(folderName, nameFile))
	if err != nil {
		panic(err)
	}
	defer arquivo.Close()

	_, err = io.Copy(arquivo, response.Body)
	if err != nil {
		panic(err)
	}
}

func RemoveFile(filename string) {
	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}
}
