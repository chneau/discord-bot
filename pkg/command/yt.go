package command

import (
	"errors"
	"io"
	"strings"

	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"

	"github.com/bwmarrin/discordgo"
)

// YT ...
type YT struct{}

// Name ...
func (y *YT) Name() string {
	return "yt"
}

// Description ...
func (y *YT) Description() string {
	return "> echo a yt video to first voice channel"
}

// MessageCreate ...
func (y *YT) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	strTok := strings.SplitN(m.Content, " ", 2)
	if len(strTok) < 2 {
		return errors.New("Not enough arguments to command")
	}
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return err
	}
	vs := g.VoiceStates[0]
	vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
	// Change these accordingly
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	videoInfo, err := ytdl.GetVideoInfo(strTok[1])
	if err != nil {
		return err
	}

	format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	downloadURL, err := videoInfo.GetDownloadURL(format)
	if err != nil {
		return err
	}

	encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
	if err != nil {
		return err
	}
	defer encodingSession.Cleanup()
	vc.Speaking(true)
	for {
		opus, err := encodingSession.OpusFrame()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		vc.OpusSend <- opus
	}
	vc.Speaking(false)
	vc.Disconnect()

	// done := make(chan error)
	// y.ss = dca.NewStream(encodingSession, vc, done)
	// y.ss.Unlock
	// err = <-done
	// if err != nil && err != io.EOF {
	// 	return err
	// }
	return nil
}
