package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/bwmarrin/discordgo"
)

// loadSound attempts to load an encoded sound file from disk.
func loadSound() error {

	file, err := os.Open("airhorn.dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		airHornBuffer = append(airHornBuffer, InBuf)
	}
}

func loadYoutube(url string) error {
	// Open a youtube stream
	yt, err := youtubePy(url)
	if err != nil {
		return err
	}

	// Create opus stream
	stream, err := convertToOpus(yt)
	if err != nil {
		return err
	}

	youtubeTempBuffer := new(bytes.Buffer)
	youtubeTempBuffer.ReadFrom(stream)

	// append youtube vid to buffer
	youtubeBuffer = append(youtubeBuffer, youtubeTempBuffer.Bytes())

	return nil
}

// convertToOpus converts the given io.Reader stream to an Opus stream
// Using ffmpeg and dca-rs
func convertToOpus(rd io.Reader) (io.Reader, error) {

	// Convert to a format that can be passed to dca-rs
	fmt.Printf("runtime is: %s\n", runtime.GOOS)
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpeg.Stdin = rd
	ffmpegout, err := ffmpeg.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// get the proper dca-rs binary
	dcaBinary := ""
	if runtime.GOOS == "darwin" {
		dcaBinary = "dependencies/osx-bin/dca-rs"
	} else if runtime.GOOS == "linux" {
		dcaBinary = "dependencies/linux-bin/dca-rs"
	} else {
		fmt.Printf("dca-rs not compatiable with current OS (osx and linux currently supported)\n")
	}

	// Convert to opus
	dca := exec.Command(dcaBinary, "--raw", "-i", "pipe:0")
	dca.Stdin = ffmpegout
	dcaout, err := dca.StdoutPipe()
	dcabuf := bufio.NewReaderSize(dcaout, 1024)
	if err != nil {
		return nil, err
	}

	// Start ffmpeg
	err = ffmpeg.Start()
	if err != nil {
		return nil, err
	}

	// Start dca-rs
	err = dca.Start()
	if err != nil {
		return nil, err
	}

	// Returns a stream of opus data
	fmt.Printf("Im in covertToOpus\n")
	return dcabuf, nil
}

// youtubePy downloads a URL using the python ytdl
func youtubePy(url string) (io.Reader, error) {
	fmt.Printf("Url is: %s\n", url)

	ytdl := exec.Command("youtube-dl", "-f", "bestaudio", "-o", "-", url)
	ytdlout, err := ytdl.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = ytdl.Start()
	if err != nil {
		return nil, err
	}
	return ytdlout, nil
}

// Reads an opus packet to send over the vc.OpusSend channel
func readOpus(source io.Reader) ([]byte, error) {
	var opuslen int16
	err := binary.Read(source, binary.LittleEndian, &opuslen)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, err
		}
		return nil, errors.New("ERR reading opus header")
	}

	var opusframe = make([]byte, opuslen)
	err = binary.Read(source, binary.LittleEndian, &opusframe)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, err
		}
		return nil, errors.New("ERR reading opus frame")
	}

	return opusframe, nil
}

func findUserVoiceState(session *discordgo.Session, userid string) (*discordgo.VoiceState, error) {
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				return vs, nil
			}
		}
	}
	return nil, errors.New("Could not find user's voice state")
}

func joinUserVoiceChannel(session *discordgo.Session, userID string) (*discordgo.VoiceConnection, error) {
	// Find a user's current voice channel
	vs, err := findUserVoiceState(session, userID)
	if err != nil {
		return nil, err
	}

	// Join the user's channel and start unmuted and deafened.
	return session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
}
