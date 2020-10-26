package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gdamore/tcell"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/rivo/tview"
)

// CLIUI is a Text User Interface (TUI) for Peers
type CLIUI struct {
	GeneralChannel   *Channel
	MiningChannel    *Channel
	FullNodesChannel *Channel
	app              *tview.Application
	peersList        *tview.TextView

	hostWindow *tview.TextView
	inputCh    chan string
	doneCh     chan struct{}
}

type Log struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../")
)

func NewCLIUI(generalChannel *Channel, miningChannel *Channel, fullNodesChannel *Channel) *CLIUI {
	app := tview.NewApplication()

	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle(fmt.Sprintf("HOST (%s)", strings.ToUpper(ShortID(generalChannel.self))))

	msgBox.SetChangedFunc(func() {
		app.Draw()
	})

	inputCh := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel(strings.ToUpper(ShortID(generalChannel.self)) + " > ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			// we don't want to do anything if they just tabbed away
			return
		}
		line := input.GetText()
		if len(line) == 0 {
			// ignore blank lines
			return
		}

		// bail if requested
		if line == "/quit" {
			app.Stop()
			return
		}

		inputCh <- line
		input.SetText("")
	})

	
	peersList := tview.NewTextView()
	peersList.SetBorder(true)
	peersList.SetTitle("Peers")


	chatPanel := tview.NewFlex().
		AddItem(msgBox, 0, 1, false).
		AddItem(peersList, 20, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false)

	app.SetRoot(flex, true)

	return &CLIUI{
		GeneralChannel:   generalChannel,
		MiningChannel:    miningChannel,
		FullNodesChannel: fullNodesChannel,
		app:              app,
		peersList:        peersList,
		hostWindow:       msgBox,
		inputCh:          inputCh,
		doneCh:           make(chan struct{}, 1),
	}
}

// Run starts the logs event loop in the background, then starts
// the event loop for the text UI.
func (ui *CLIUI) Run(net *Network) error {
	go ui.handleEvents(net)
	defer ui.end()

	return ui.app.Run()
}

// end signals the event loop to exit gracefully
func (ui *CLIUI) end() {
	ui.doneCh <- struct{}{}
}

// refreshPeers pulls the list of peers currently in the channel and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *CLIUI) refreshPeers() {
	peers := ui.GeneralChannel.ListPeers()
	minerPeers := ui.MiningChannel.ListPeers()
	idStrs := make([]string, len(peers))

	for i, p := range peers {
		peerId := strings.ToUpper(ShortID(p))
		if len(minerPeers) != 0 {
			isMiner := false
			for _, minerPeer := range minerPeers {
				if minerPeer == p {
					isMiner = true
					break
				}
			}
			if isMiner {
				idStrs[i] = "MINER: " + peerId
			} else {
				idStrs[i] = peerId
			}
		} else {
			idStrs[i] = peerId
		}
	}

	ui.peersList.SetText(strings.Join(idStrs, "\n"))
	ui.app.Draw()
}

func (ui *CLIUI) displaySelfMessage(msg string) {
	prompt := withColor("yellow", fmt.Sprintf("<%s>:", strings.ToUpper(ShortID(ui.GeneralChannel.self))))
	fmt.Fprintf(ui.hostWindow, "%s %s\n", prompt, msg)
}

func (ui *CLIUI) displayContent(content *ChannelContent) {
	prompt := withColor("green", fmt.Sprintf("<%s>:", strings.ToUpper(content.SendFrom)))
	fmt.Fprintf(ui.hostWindow, "%s %s\n", prompt, content.Message)
}

func (ui *CLIUI) HandleStream(net *Network, content *ChannelContent) {
	// ui.displayContent(content)
	if content.Payload != nil {
		command := BytesToCmd(content.Payload[:commandLength])
		log.Infof("Received  %s command \n", command)

		switch command {
		case "block":
			net.HandleBlocks(content)
		case "inv":
			net.HandleInv(content)
		case "getblocks":
			net.HandleGetBlocks(content)
		case "getdata":
			net.HandleGetData(content)
		case "tx":
			net.HandleTx(content)
		case "gettxfrompool":
			net.HandleGetTxFromPool(content)
		case "version":
			net.HandleVersion(content)
		default:
			log.Warn("Unknown Command")
		}
	}
}

func (ui *CLIUI) readFromLogs(instanceId string) {
	filename := "/logs/console.log"
	if instanceId != "" {
		filename = fmt.Sprintf("/logs/console_%s.log", instanceId)
	}

	logFile := path.Join(Root, filename)
	e := ioutil.WriteFile(logFile, []byte(""), 0644)
	if e != nil {
		panic(e)
	}
	log.SetOutput(ioutil.Discard)

	f, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	logLevels := map[string]string{
		"info":    "green",
		"warning": "brown",
		"error":   "red",
		"fatal":   "red",
	}
	oldSize := info.Size()
	for {
		for line, prefix, err := r.ReadLine(); err != io.EOF; line, prefix, err = r.ReadLine() {
			var data Log
			if err := json.Unmarshal(line, &data); err != nil {
				panic(err)
			}
			prompt := fmt.Sprintf("[%s]:", withColor(logLevels[data.Level], strings.ToUpper(data.Level)))
			if prefix {
				fmt.Fprintf(ui.hostWindow, "%s %s\n", prompt, data.Msg)
			} else {
				fmt.Fprintf(ui.hostWindow, "%s %s\n", prompt, data.Msg)
			}
			ui.hostWindow.ScrollToEnd()
		}
		pos, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			panic(err)
		}
		for {
			time.Sleep(time.Second)
			newinfo, err := f.Stat()
			if err != nil {
				panic(err)
			}
			newSize := newinfo.Size()
			if newSize != oldSize {
				if newSize < oldSize {
					f.Seek(0, 0)
				} else {
					f.Seek(pos, io.SeekStart)
				}
				r = bufio.NewReader(f)
				oldSize = newSize
				break
			}
		}
	}
}

// handleEvents runs an event loop that sends user input to the channel
// and displays messages received from the channel. It also periodically
// refreshes the list of peers in the UI.
func (ui *CLIUI) handleEvents(net *Network) {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()

	go ui.readFromLogs(net.Blockchain.InstanceId)
	log.Info("HOST ADDR: ", net.Host.Addrs())

	for {
		select {
		case input := <-ui.inputCh:

			err := ui.GeneralChannel.Publish(input, nil, "")
			if err != nil {
				log.Errorf("Publish error: %s", err)
			}
			ui.displaySelfMessage(input)

		case <-peerRefreshTicker.C:
			// refresh the list of peers in the chat room periodically
			ui.refreshPeers()

		case m := <-ui.GeneralChannel.Content:
			ui.HandleStream(net, m)

		case m := <-ui.MiningChannel.Content:
			ui.HandleStream(net, m)

		case m := <-ui.FullNodesChannel.Content:
			ui.HandleStream(net, m)

		case <-ui.GeneralChannel.ctx.Done():
			return

		case <-ui.doneCh:
			return
		}
	}
}

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}

// ShortID returns the last 8 chars of a base58-encoded peer id.
func ShortID(p peer.ID) string {
	pretty := p.Pretty()
	return pretty[len(pretty)-8:]
}
