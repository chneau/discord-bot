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
type YT struct {
	playing chan interface{}
	stop    chan interface{}
}

// Name ...
func (y *YT) Name() string {
	return "yt"
}

// Description ...
func (y *YT) Description() string {
	return "__link-or-ID__ > play a yt video to first voice channel"
}

func (y *YT) linkFronChat(content string) string {
	strTok := strings.SplitN(content, " ", 2)
	ytlink := "JqRxmy1h5as"
	if len(strTok) == 2 {
		ytlink = strTok[1]
	}
	return ytlink
}

func (y *YT) dcaDefaultOpt() *dca.EncodeOptions {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = dca.AudioApplicationLowDelay
	return options
}

// MessageCreate ...
func (y *YT) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if len(y.playing) == 0 {
		y.stop <- nil
	}
	ytlink := y.linkFronChat(m.Content)
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return err
	}
	videoInfo, err := ytdl.GetVideoInfo(ytlink)
	if err != nil {
		return err
	}
	if videoInfo.Duration == 0 {
		return errors.New("video not found")
	}
	vc, err := s.ChannelVoiceJoin(g.ID, g.VoiceStates[0].ChannelID, false, true)
	if err != nil {
		return err
	}
	downloadURL, err := videoInfo.GetDownloadURL(videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0])
	if err != nil {
		return err
	}
	encodingSession, err := dca.EncodeFile(downloadURL.String(), y.dcaDefaultOpt())
	if err != nil {
		return err
	}
	defer encodingSession.Cleanup()
	message, _ := s.ChannelMessageSend(m.ChannelID, "Playing: "+videoInfo.Title)
	defer func() {
		s.ChannelMessageDelete(m.ChannelID, message.ID)
	}()
	err = y.play(encodingSession, vc.OpusSend)
	if err != nil {
		return err
	}
	return nil
}

func (y *YT) play(e *dca.EncodeSession, opusSend chan []byte) error {
	<-y.playing
	defer func() {
		y.playing <- nil
	}()
	for {
		opus, err := e.OpusFrame()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		select {
		case opusSend <- opus:
		case <-y.stop:
			return nil
		}
	}
}

// Init ...
func (y *YT) Init() {
	y.playing = make(chan interface{}, 1)
	y.stop = make(chan interface{})
	y.playing <- nil
}
