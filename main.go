package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jonas747/dca"
	"github.com/namsral/flag"
	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/sardap/discom"
	"github.com/sardap/dr-michael-morseies/morse"
)

var (
	userLocation          cmap.ConcurrentMap[string] = cmap.New[string]()
	activeSessions        map[string]bool            = make(map[string]bool)
	activeLock            *sync.Mutex                = &sync.Mutex{}
	discord               *discordgo.Session
	dashSoundFile         string
	dotSoundFile          string
	dotDashBreakSoundFile string
	letterBreakSoundFile  string
	wordBreakSoundFile    string
)

func toMorseCmd(s *discordgo.Session, i discom.Interaction) error {
	text := i.Option("text").StringValue()
	i.Respond(s, discom.Response{Content: morse.ToMorseCode(text)})
	return nil
}

func formMorseCmd(s *discordgo.Session, i discom.Interaction) error {
	text := i.Option("text").StringValue()
	lang := i.Option("lang").StringValue()
	i.Respond(s, discom.Response{Content: morse.FromMorseCode(text, lang)})
	return nil
}

func playVideo(conn *discordgo.VoiceConnection, errCh chan error, encodingSession *dca.EncodeSession) {
	conn.Speaking(true)
	defer conn.Speaking(false)
	done := make(chan error)
	dca.NewStream(encodingSession, conn, done)
	err := <-done
	if err != nil && err != io.EOF {
		log.Printf("Play video error: %s", err.Error())
		errCh <- err
		return
	}

	errCh <- nil
}

func playFile(conn *discordgo.VoiceConnection, errCh chan error, path string, volume int) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	options.Volume = volume

	encodingSession, err := dca.EncodeFile(path, options)
	if err != nil {
		errCh <- errors.Wrapf(err, "encdoing file")
		return
	}

	playVideo(conn, errCh, encodingSession)
}

func playMorseCode(conn *discordgo.VoiceConnection, i discom.Interaction, morseCode string) {
	defer func() {
		activeLock.Lock()
		delete(activeSessions, i.GetPayload().GuildId)
		activeLock.Unlock()
	}()

	morseCode = strings.ReplaceAll(morseCode, " / ", "/")

	seq := make([]string, 0)
	for _, c := range morseCode {
		switch c {
		case '.':
			seq = append(seq, dotSoundFile)
			seq = append(seq, dotDashBreakSoundFile)
		case '-':
			seq = append(seq, dashSoundFile)
			seq = append(seq, dotDashBreakSoundFile)
		case ' ':
			seq = append(seq, letterBreakSoundFile)
		case '/':
			seq = append(seq, wordBreakSoundFile)
		}

	}

	playFilename := fmt.Sprintf("%d.wav", rand.Int())
	defer os.Remove(playFilename)
	seq = append(seq, playFilename)

	cmd := exec.Command("sox", seq...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		i.Respond(discord, discom.Response{Content: "Error playing please contact admin"})
		log.Printf("playMorseCode error while running sox %v", err)
		return
	}

	i.Respond(discord, discom.Response{Content: "Playing"})
	ch := make(chan error)
	go playFile(conn, ch, playFilename, 80)
	err = <-ch
	close(ch)
	if err != nil {
		log.Printf("Error playing video %s", err)
		return
	}
	i.Respond(discord, discom.Response{Content: "Done playing"})

	conn.Disconnect()
}

func playCmd(s *discordgo.Session, i discom.Interaction) error {
	i.Respond(s, discom.Response{Content: "processing"})

	activeLock.Lock()
	defer activeLock.Unlock()

	if _, ok := activeSessions[i.GetPayload().GuildId]; ok {
		i.Respond(s, discom.Response{Content: "Already playing on this server"})
		return nil
	}

	text := i.Option("text").StringValue()
	morseCode := morse.ToMorseCode(text)

	chId, ok := userLocation.Get(i.GetPayload().AuthorId)
	if !ok {
		i.Respond(s, discom.Response{Content: "Can't find you in a voice channel"})
		return nil
	}

	i.Respond(s, discom.Response{Content: "Joining"})
	conn, _ := s.ChannelVoiceJoin(i.GetPayload().GuildId, chId, false, true)

	activeSessions[i.GetPayload().GuildId] = true
	go playMorseCode(conn, i, morseCode)

	return nil
}

func errorHandler(s *discordgo.Session, i discom.Interaction, cmdErr error) {
	i.Respond(s, discom.Response{Content: cmdErr.Error()})
}

func createCommandSet(prefix string) *discom.CommandSet {
	commandSet, _ := discom.CreateCommandSet(prefix, errorHandler)

	commandSet.AddCommand(discom.Command{
		Name: "to_morse", Handler: toMorseCmd,
		Description: "convertes text to morse code",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "text to turn into morse code",
				Required:    true,
			},
		},
	})

	commandSet.AddCommand(discom.Command{
		Name: "from_morse", Handler: formMorseCmd,
		Description: "Morse code to english",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "text which is morse code",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "lang",
				Description: "language to decode from",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "English",
						Value: morse.LangEnglish,
					},
					{
						Name:  "한글",
						Value: morse.LangKorean,
					},
				},
			},
		},
	})

	commandSet.AddCommand(discom.Command{
		Name: "play", Handler: playCmd,
		Description: "Plays morse code in voice channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "english text",
				Required:    true,
			},
		},
	})

	return commandSet
}

func voiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.ChannelID != "" {
		userLocation.Set(v.UserID, v.ChannelID)
	} else {
		userLocation.Remove(v.UserID)
	}
}

func main() {
	rand.Seed(time.Now().UnixMicro())

	// Flags
	var discordAuth string
	flag.StringVar(&discordAuth, "discord_auth", "", "the discord auth token")
	flag.StringVar(&dashSoundFile, "dash_sound_file", "", "dash sound file")
	flag.StringVar(&dotSoundFile, "dot_sound_file", "", "dot sound file")
	flag.StringVar(&dotDashBreakSoundFile, "dot_dash_break_sound_file", "", "dot dash break sound file")
	flag.StringVar(&letterBreakSoundFile, "letter_break_sound_file", "", "letter break sound file")
	flag.StringVar(&wordBreakSoundFile, "word_break_sound_file", "", "word break sound file")
	flag.Parse()

	prefix := "-dm"
	//Create command set
	cs := createCommandSet(prefix)

	var err error
	discord, err = discordgo.New("Bot " + discordAuth)
	if err != nil {
		log.Printf("unable to create new discord instance")
		log.Fatal(err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(voiceStateUpdate)
	discord.AddHandler(cs.Handler)
	discord.AddHandler(cs.IntreactionHandler)

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		discord.UpdateListeningStatus(fmt.Sprintf("%s help", prefix))
		log.Println("Bot is up!")
	})

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	defer discord.Close()

	cs.SyncAppCommands(discord)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()

}
