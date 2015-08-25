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