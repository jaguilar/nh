package model

import (
	"bufio"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/jaguilar/nethack/model/command"
	"github.com/jaguilar/nethack/model/level"
	"github.com/jaguilar/nethack/model/pc"
	"github.com/jaguilar/vt100"
)

var (
	// ErrUnexpectedResume is returned from Issue when the underlying game
	// resumed sending data unexpectedly. The correct thing to do when
	// this happens is to discard the previous decision you made and return
	// to LockWhenIdle.
	ErrUnexpectedResume = errors.New("game resumed unexpectedly")

	// ErrGameOver indicates that the game has ended. Any further commands will
	// result in the same error.
	ErrGameOver = errors.New("game has ended")
)

var (
	// The maximum number of turns for which we keep events. Turns are measured
	// not in nethack time units, but in number of opportunities for your bot
	// to issue commands.
	MaxEventLookback = 100
)

// WindowSize is the dimensions of the terminal that Nethack is running on.
// You must have set this on Nethack's output fd via termios or other similar
// functionality.
//
// Given default terminal window sizes (e.g 25x80), it is possible for nethack
// to truncate very long item names. We make no attempt to handle this case. Such
// items may either appear in your inventory with incorrect properties, or not
// appear at all. For example:
//
// * an uncursed greased thoroughly corroded thoroughly rusty rustproof +100 helm of opposite alignment
//
// This item name is already 100 columns wide. If Nethack were to truncate this
// item, it would possibly not appear in your inventory within the model.
//
// To prevent this problem from affecting you, you should use termios or similar
// to set a very wide Nethack window. 200 columns would likely be sufficient
// for the great majority of games. On the same token, a taller window will speed
// up the parsing of lists, because the model will not have to wait for as many
// frame strobes when advancing through each list.
type WindowSize struct {
	Y, X int
}

type Game struct {
	// The player. Tip: for a nethack-like programming experience, define
	// u in your package and have it point to this. This is assuming you are only
	// running one Game per process.
	pc.Player

	// Level contains all the levels we've seen.
	Level map[level.LevelID]*level.Level

	// TODO(jaguilar): adjust this to the Event interface once we actually parse
	// each event.
	Events [][]string

	// in and out are the input stream from and output stream to nethack.
	out io.Writer

	vt       *vt100.VT100
	mu, vtMu sync.Mutex

	inputCommands <-chan vt100.Command
	inputErrs     <-chan error
}

// NewGame makes a new Game. The game immediately starts parsing the Nethack input
// stream and building a model of the world. It also takes over control of the output
// stream to nethack. From this point, your only interaction with these IO objects
// should be through this instance of Game.
//
// It is safe to examine this between calls to Do, but not during any given Do
// call.
func NewGame(in io.Reader, out io.Writer, win WindowSize) (*Game, error) {
	cmds, errs := inputUntilClosed(in)
	g := &Game{
		Level:         make(map[level.LevelID]*level.Level),
		Events:        make([][]string, MaxEventLookback),
		in:            bufio.NewReader(in),
		out:           out,
		vt:            vt100.NewVT100(win.Y, win.X),
		err:           make(chan error),
		inputCommands: cmds,
		inputErrs:     errs,
	}

	return g, g.waitIdle()
}

func inputUntilClosed(in io.Reader) (<-chan vt100.Command, <-chan error) {
	buf := bufio.NewReaderSize(in, 512)
	cmds, errs := make(chan vt100.Command, 10), make(chan error)
	go func() {
		for {
			cmd, err := vt100.Decode(buf)
			if err == io.EOF {
				close(cmds)
				close(errs)
				return
			}
			if err != nil {
				errs <- err
				continue
			}
			cmds <- cmd
		}
	}()
	return cmds, errs
}

// waitIdle waits until a nethack frame has finished drawing.
func (g *Game) waitIdle() error {
	timeout := time.Millisecond * 20
	timer := time.NewTimer(timeout)
	for {
		select {
		case cmd, ok := <-g.inputCommands:
			if !ok {
				// Only way this can happen is io.EOF from the input channel.
				return ErrGameOver
			}

			g.vtMu.Lock()
			err := g.vt.Process(cmd)
			g.vtMu.Unlock()

			if err != nil {
				// We ignore unsupported errors from the VT100. If such an error results
				// in game state corruption, please file an issue with jaguilar/vt100.
				if _, ok := err.(vt100.UnsupportedError); ok {
					continue
				}
				return err
			}
		case err, ok := <-g.inputErrs:
			if !ok {
				return ErrGameOver
			}
			return err
		case <-timer.C:
			// First cut: we'll consider the frame idle if 20ms have passed since
			// the last update. This will likely not be sufficient for remote nethack.
			// Ideas for doing it better: given the last command, we know where the cursor
			// should end up (usually at the end of the first line, on the player, or at the
			// end of a list). We can use knowledge of where the cursor is to reduce
			// the timeout, and use a long timeout for when the cursor is anywhere else.
			return nil
		}
		timer.Reset(timeout)
	}
}

// Do a Command. errors are only returned if you issued an illegal command or
// something went wrong with nethack (e.g. it was killed out from under us).
// To see if your command worked or did what you intended, you'll need to check
// the event log.
//
// To exit the game (even if you died, or the game crashed), you need to send
// command.Quit.
func (g *Game) Do(c command.Command) error {

	if err := g.parse(); err != nil {
		return err
	}
	return g.waitIdle()
}
