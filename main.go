package main

import (
	"context"
	"fmt"
	. "github.com/ezequielaguilera1993/tocadiscos"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var client *whatsmeow.Client

func sndMsg(v *events.Message, msg string) {
	responseMessage := msg
	pm := &proto.Message{
		Conversation: &responseMessage,
	}
	_, err := client.SendMessage(context.Background(), v.Info.Sender.ToNonAD(), pm)
	if err != nil {
		panic(err)
	}
}

func tocarNota(n Nota) {

	PlaySong(Song{
		NotesToPlay: []NoteToPlay{
			{Nota: n, FiguraMusical: N},
		},
		SampleRate: 44100,
		Tempo:      Tempo{BPM: 140},
	})
}
func tocarNotaCustom(n Nota, fm FiguraMusical) {

	PlaySong(Song{
		NotesToPlay: []NoteToPlay{
			{Nota: n, FiguraMusical: fm},
		},
		SampleRate: 44100,
		Tempo:      Tempo{BPM: 140},
	})
}

func eventHandler(evt interface{}) {

	switch v := evt.(type) {
	case *events.Message:
		fmt.Printf("%#v\n", v)
		var msg string
		if v.Message.ExtendedTextMessage != nil {
			msg = strings.ToLower(*v.Message.ExtendedTextMessage.Text)
		} else {
			msg = v.Message.GetConversation()
		}
		fmt.Println("Received a message!", msg)
		switch msg {
		case "f":
			sndMsg(v, "queselevahace")
			PlaySong(LooserSong)
		case "reto":
			tocarNotaCustom(6000, .5)
		default:
			noteWords := strings.Fields(msg)
			for _, noteWord := range noteWords {
				switch strings.ToLower(noteWord) {
				case ".do":
					tocarNota(C4)
				case ".do#":
					tocarNota(C4_SOSTENIDO)
				case ".re":
					tocarNota(D4)
				case ".re#":
					tocarNota(D4_SOSTENIDO)
				case ".mi":
					tocarNota(E4)
				case ".fa":
					tocarNota(F4)
				case ".fa#":
					tocarNota(F4_SOSTENIDO)
				case ".sol":
					tocarNota(G4)
				case ".sol#":
					tocarNota(G4_SOSTENIDO)
				case ".la":
					tocarNota(A4)
				case ".la#":
					tocarNota(A4_SOSTENIDO)
				case ".si":
					tocarNota(B4)
				}
			}

		}

	}

}

func main() {

	dbLog := waLog.Stdout("Database", "ERROR", true)
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
	clientLog := waLog.Stdout("Client", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
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

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
