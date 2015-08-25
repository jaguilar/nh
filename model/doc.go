/*
package model is a package designed to model the state of a
nethack game. It consists both of the model itself, and code
for parsing the model from the screen.

Operation is as follows. Begin by starting a Nethack subprocess.
You may have to lie to nethack about whether the input stream
is a terminal. Although how to do this is outside the scope
of this package, demo/ will eventually contain a main
function that will show you how.

After you have a running Nethack subprocess, you can start updating
the model with Game.Control(reader, writer). At this point, these
streams are owned by the Game and should not be used elsewhere.
Game.Control blocks until the game ends or an unrecoverable error
is encountered.

The general mode of operation for a bot is:

    if err := game.IdleLock(); err != nil {
      // An error here probably indicates the end of a play session.
      return err
    }
    var command model.Command = aistep(game)
    // Issue unlocks the game, even if there is an error.
    if err := game.Issue(command); err != nil {
	  return err
    }

Please read the documentation for IdleLock and Issue for more details.
You may not necessarily wish to give up every time Game.Issue returns
an error.
*/
package model
