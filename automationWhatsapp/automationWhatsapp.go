package automationWhatsapp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	case *events.Connected:
		fmt.Println("[ v ] Conexão estabelecida com sucesso!")
	case *events.StreamReplaced:
		fmt.Println("[ x ] Stream foi substituído, reconectando...")
	}
}

func Run() {
	// Ler números do stdin
	listPhones := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		listPhones = append(listPhones, "+"+scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro de leitura:", err)
		return
	}

	if len(listPhones) == 0 {
		fmt.Println("Nenhum número foi fornecido")
		return
	}

	fmt.Printf("Números para processar: %d\n", len(listPhones))

	// Configurar banco de dados
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New(
		context.Background(),
		"sqlite3",
		"file:examplestore.db?_foreign_keys=on",
		dbLog,
	)
	if err != nil {
		panic(err)
	}

	// Obter dispositivo
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		panic(err)
	}

	// Criar cliente
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	// Verificar se precisa fazer login
	if client.Store.ID == nil {
		fmt.Println("Nenhuma sessão encontrada. Faça login escaneando o QR Code:")

		qrChan, _ := client.GetQRChannel(context.TODO())
		err = client.Connect()
		if err != nil {
			panic(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println("\nEscaneie o QR code acima com seu WhatsApp")
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}

		fmt.Println("[ v ] Login realizado com sucesso!")
		fmt.Println("Aguardando sincronização completa do WhatsApp...")
		time.Sleep(15 * time.Second)

	} else {
		fmt.Println("Sessão encontrada. Conectando...")
		err = client.Connect()
		if err != nil {
			panic(err)
		}

		fmt.Println("Aguardando estabilização da conexão...")
		time.Sleep(5 * time.Second)
	}

	// Verificar se está conectado
	if !client.IsConnected() {
		fmt.Println("[ x ] Erro: Cliente não está conectado")
		return
	}

	fmt.Println("[ v ] Cliente conectado e pronto!")
	fmt.Println("Iniciando processamento dos números...")

	// Limpar arquivos anteriores
	quantityUsers := 0
	RemoveFile("./numberphone/all-numbers.txt")
	RemoveFile("./numberphone/numbers-profile.txt")
	RemoveFile("./numberphone/numbers-withoutProfile.txt")

	// Processar cada número
	for i, numberphone := range listPhones {
		fmt.Printf("[%d/%d] Verificando %s... ", i+1, len(listPhones), numberphone)

		// Verificar se o número está no WhatsApp (como no código original)
		IsOnWhatsAppResponse, errIsOnWhatsApp := client.IsOnWhatsApp([]string{numberphone})
		if errIsOnWhatsApp != nil {
			fmt.Printf("[ x ] Erro: %v\n", errIsOnWhatsApp)
			continue
		}

		if !IsOnWhatsAppResponse[0].IsIn {
			fmt.Println("[ x ] Não está no WhatsApp")
			continue
		}

		// Número está no WhatsApp
		quantityUsers++
		fmt.Print("[ v ] Tem WhatsApp:", numberphone)

		WriteToFile("all-numbers.txt", numberphone+"\n", "./numberphone/")

		// Pequeno delay antes de pegar foto
		time.Sleep(500 * time.Millisecond)

		// Tentar obter foto de perfil
		errorProfileHidden := false
		GetProfilePictureInfoResponse, errGetProfile := client.GetProfilePictureInfo(IsOnWhatsAppResponse[0].JID, nil)

		if errGetProfile != nil {
			if strings.Contains(errGetProfile.Error(), "hidden their profile") ||
				strings.Contains(errGetProfile.Error(), "group does not have a profile") {
				errorProfileHidden = true
			} else {
				fmt.Printf(" | [ x ] Erro ao obter foto: %v", errGetProfile)
				errorProfileHidden = true
			}
		}

		if !errorProfileHidden && GetProfilePictureInfoResponse != nil && GetProfilePictureInfoResponse.URL != "" {
			err := DownloadFile(GetProfilePictureInfoResponse.URL, numberphone+".jpg", "./numberphone/profile/")
			if err == nil {
				WriteToFile("numbers-profile.txt", numberphone+"\n", "./numberphone/")
				fmt.Println(" | [ v ] Foto baixada")
			} else {
				WriteToFile("numbers-withoutProfile.txt", numberphone+"\n", "./numberphone/")
				fmt.Printf(" | [ x ] Erro ao baixar foto: %v\n", err)
			}
		} else {
			WriteToFile("numbers-withoutProfile.txt", numberphone+"\n", "./numberphone/")
			fmt.Println(" | [ # ] Sem foto/perfil oculto")
		}

		// Delay entre requisições
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("\033[32m[ v ] Processamento concluído!\033[0m\n")
	fmt.Printf("\033[32m  Total de números no WhatsApp: %d de %d\033[0m\n", quantityUsers, len(listPhones))
	fmt.Println(strings.Repeat("=", 50))

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

func DownloadFile(targetURL string, nameFile string, folderName string) error {
	response, err := http.Get(targetURL)
	if err != nil {
		return fmt.Errorf("erro ao fazer request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("status code: %d", response.StatusCode)
	}

	os.MkdirAll(folderName, os.ModePerm)
	arquivo, err := os.Create(filepath.Join(folderName, nameFile))
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer arquivo.Close()

	_, err = io.Copy(arquivo, response.Body)
	if err != nil {
		return fmt.Errorf("erro ao copiar dados: %w", err)
	}

	return nil
}

func RemoveFile(filename string) {
	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}
}
