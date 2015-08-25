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

## Combat Tactics Sketch

Combat tactics is a huge problem. There are two main subproblems:

* When do we fight or run away?
* How do we fight or run away?

### Fight or Flight

My best guess for this question at the moment is as follows:

* Simulate the fight. If there is a high chance of winning without spending resources (most fights), fight.
* If we'll probably lose, or would spend large resources to win, look at the alternatives. Is there anything
else we could be doing right now elsewhere in the dungeon? If not, fight.
* Simulate running away. Is there a good chance we can run in a useful direction? If not, fight.
* Run.

### How to Fight/Run

We'll use minimax to determine how to fight. This approach is suitable for minimizing our
loss function during the fight. How exactly we'll define the loss function is unclear.
Obviously, `u.HP==0 -> -1`. We should also take away some points for using items -- the 
rarer the item, the more loss. Meanwhile, all enemies dead is obviously the 1 outcome, or
else us escaping the battle -- being able to leave the level -- if we're trying to run.

We can use our knowledge of enemy movement to prune a bunch of possible nodes. For example,
most enemies will not step away from us, and even if they would, it would be a worse move
for them, so we'll not have to consider it.

### Testing

One of the nice aspects of this approach to combat is that it will be very testable with
wizard mode. We can set up combat scenarios and run them repeatedly to see how well the algorithm
handles various scenarios. We can also test the algorithm's estimates about win chances.
It's not required that the win chance information be perfect, just that it be reasonable.