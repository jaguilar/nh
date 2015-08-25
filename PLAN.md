# The Plan

The plan right now is to get this repository to the point where 
it can do basic state parsing and pass keyboard commands through
to Nethack.

After that, I will try to create some small tools that a human
player can use for assistance. For example, I might make
an auto-Elbereth tool that keeps trying to write Elbereth until
it works. Or I might make a price identification module
that automatically price-ids items from your inventory.
Little modules like this will expose flaws in the API that
will be easier to fix if there is not a large codebase depending
on them.

If all that is done, perhaps I'll start work on a real bot.

## Identification Sketch

A sketch of how an identification submodule would work.

You'll start out with a type that maps appearance to list of possible
actual class. It would look like this:

```Go
type ID struct {
	Wands   map[wand.Appearance][]*wand.Class
	Potions // ...
}
```

Identification would be eliminative. For example, suppose we're trying
to identify a platinum wand. We take it to a shop to price ID it. The shopkeeper
tells us it's worth 500 (this would be calculated automatically from
charisma, and we'd have the model get multiple prices if there were 
ambiguity). We look up in the wand price ID table:

```Go
var (
    WandPrices = [int][]*item.Class {
        // ...
        500: []*item.Class{wand.OfDeath, wand.OfWishing},
    }
)
```

We erase from `id.Wands[wand.Plat]` any elements that don't match the
price identification list. Later, we observe an oak wand being zapped and emitting
a death ray. This immediately identifies the oak wand as the wand of death.
So we set `id.Wands[wand.Oak] = []*wand.Class{wand.OfDeath}`. We also erase `wand.OfDeath`
from all of the other appearance lists. This leaves `id.Wands[wand.Platinum]` containing
only a single element: `wand.OfWishing`.

I believe that this eliminative approach will be sufficient for all identification purposes
in the game.