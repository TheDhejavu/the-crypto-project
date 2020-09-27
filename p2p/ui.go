package p2p

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/rivo/tview"
)

// CLIUI is a Text User Interface (TUI) for the node.
// The Run method will draw the UI to the terminal in "fullsnodeeen"
// mode. You can quit with Ctrl-C, or by typing "/quit" into the
// chat prompt.
type CLIUI struct {
	GeneralChannel *Channel
	MiningChannel  *Channel
	app            *tview.Application
	peersList      *tview.TextView

	msgW    io.Writer
	inputCh chan string
	doneCh  chan struct{}
}

// NewCLIUI returns a new CLIUI struct that controls the text UI.
// It won't actually do anything until you call Run().
func NewCLIUI(generalChannel *Channel, miningChannel *Channel) *CLIUI {
	app := tview.NewApplication()

	// make a text view to contain our chat messages
	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle(fmt.Sprintf("Blockchain CLI UI"))

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgBox.SetChangedFunc(func() {
		app.Draw()
	})

	// an input field for typing messages into
	inputCh := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel(ShortID(generalChannel.self) + " > ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	// the done func is called when the user hits enter, or tabs out of the field
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

		// send the line onto the input chan and reset the field text
		inputCh <- line
		input.SetText("")
	})

	// make a text view to hold the list of peers in the room, updated by ui.refreshPeers()
	peersList := tview.NewTextView()
	peersList.SetBorder(true)
	peersList.SetTitle("Peers")

	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	chatPanel := tview.NewFlex().
		AddItem(msgBox, 0, 1, false).
		AddItem(peersList, 20, 1, false)

	// flex is a vertical box with the chatPanel on top and the input field at the bottom.

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false).
		AddItem(input, 1, 1, true)

	app.SetRoot(flex, true)

	return &CLIUI{
		GeneralChannel: generalChannel,
		MiningChannel:  miningChannel,
		app:            app,
		peersList:      peersList,
		msgW:           msgBox,
		inputCh:        inputCh,
		doneCh:         make(chan struct{}, 1),
	}
}

// Run starts the chat event loop in the background, then starts
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

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *CLIUI) refreshPeers() {
	peers := ui.GeneralChannel.ListPeers()
	minerPeers := ui.MiningChannel.ListPeers()
	idStrs := make([]string, len(peers))

	// fmt.Println(peers)
	for i, p := range peers {
		if len(minerPeers) != 0 {
			isMiner := false
			for _, minerPeer := range minerPeers {
				if minerPeer == p {
					isMiner = true
					break
				}
			}
			if isMiner {
				idStrs[i] = "MINER:" + ShortID(p)
			} else {
				idStrs[i] = ShortID(p)
			}
		} else {
			idStrs[i] = ShortID(p)
		}
	}

	ui.peersList.SetText(strings.Join(idStrs, "\n"))
	ui.app.Draw()
}

func (ui *CLIUI) displaySelfMessage(msg string) {
	prompt := withColor("yellow", fmt.Sprintf("<%s>:", ShortID(ui.GeneralChannel.self)))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

func (ui *CLIUI) displayContent(content *ChannelContent) {
	prompt := withColor("green", fmt.Sprintf("<%s>:", content.SenderNodeID))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, content.Message)
}

func (ui *CLIUI) HandleStream(net *Network, content *ChannelContent) {
	command := BytesToCmd(content.Payload[:commandLength])
	fmt.Printf("Received  %s command \n", command)

	ui.displayContent(content)

	switch command {
	case "block":
		net.HandleBlocks(content)
	case "inv":
		net.HandleInv(content)
	case "getblocks":
		net.HandleGetBlocks(content)
	case "getdata":
		net.HandleGetData(content)
	case "version":
		net.HandleVersion(content)
	default:
		fmt.Println("Unknown Command")
	}
}

// handleEvents runs an event loop that sends user input to the channel
// and displays messages received from the channel. It also periodically
// refreshes the list of peers in the UI.
func (ui *CLIUI) handleEvents(net *Network) {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()

	for {
		select {
		case input := <-ui.inputCh:

			if input == "miner" {
				err := ui.MiningChannel.Publish(input, nil, "")
				if err != nil {
					fmt.Sprintln("publish error: %s", err)
				}
			} else {
				err := ui.GeneralChannel.Publish(input, nil, "")
				if err != nil {
					fmt.Sprintln("publish error: %s", err)
				}
			}
			ui.displaySelfMessage(input)

		case <-peerRefreshTicker.C:
			// refresh the list of peers in the chat room periodically
			ui.refreshPeers()
		case m := <-ui.GeneralChannel.Content:
			ui.HandleStream(net, m)

		case m := <-ui.MiningChannel.Content:
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
