# Week 1, Session 2 — Dynamic Array & Linked List

**Format:** this is a spec, not a tutorial. There is no Go in this file on purpose. Every line of
code in `datastructures/` will be yours.

**Time:** ~1.5h. You will probably not finish both. That's expected — session 3 is for finishing.

---

## 0. Why these two, and why not the hash map

The Phase 1 milestone is a hash map, a BST, and a graph traversal. The hash map is listed first and
it is the wrong place to start, because it's two unfamiliar problems stacked on top of each other: a
hash function *and* a collision strategy, debugged through pointer semantics you haven't used yet.
When it breaks — and it will — you won't know which layer broke.

Build the dependencies first:

```
dynamic array ──┬──> hash map (buckets are an array, chains are a list)
singly linked ──┘
    list        ───> BST        (same pointer + nil discipline, plus recursion)
dynamic array   ───> BFS / DFS  (a queue and a stack are both just this)
```

Every later structure in Phase 1 is made of these two. Neither is hard. Both will expose exactly the
Go you're shaky on: pointers, `nil`, receivers, and the difference between a value and a reference
to a value.

---

## 1. Setup and the test loop

Create two package directories:

```
datastructures/vector/
datastructures/list/
```

Each directory is one package. Files inside `vector/` start with a package clause naming `vector`,
and the test file for them uses the same package name (you're testing unexported internals, so stay
in-package rather than using the `_test` suffix variant).

Commands you'll live in:

| Command | What it does |
|---|---|
| `go test ./...` | Runs every test in the module |
| `go test ./datastructures/vector` | Just this package |
| `go test -v ./...` | Names each test as it runs |
| `go test -run TestPush ./...` | Only tests whose name matches the regex |
| `go test -timeout 5s ./...` | Kills a hung test — you will need this in §3 |

**Table-driven tests** are the Go convention, and you should use them from the start. The shape,
described rather than written: declare a slice of anonymous structs where each element is one case
(a name, the inputs, the expected output). Range over the slice. Inside the loop, call `t.Run` with
the case's name and a closure holding the assertion. You get one named subtest per case, and adding
a case is one line instead of one function.

**Write the test first.** Not as a moral position — as a practical one. You don't yet know what API
you want, and writing the call site first is the fastest way to find out. Watch it fail to compile,
then make it compile, then make it pass. Phase 2 assumes this is already a habit.

---

## 2. Dynamic array

### The idea

A fixed-size block of memory, plus a counter tracking how much of it you're actually using. Two
numbers that people constantly conflate:

- **length** — how many slots hold real values
- **capacity** — how many slots you've allocated

Length ≤ capacity, always. When a push would make them equal and you need one more, you allocate a
bigger block and copy everything across. That copy is the entire subject of this exercise.

### The one rule

Go's slice already is a dynamic array, and `append` already does the growth. **You may not call
`append`.** It's the mechanism under study; calling it voids the exercise. Back your type with a
slice you only ever index into and `copy` from/to (or a genuinely fixed array if you prefer, though
that makes resizing awkward). Allocate with `make` at a chosen size, and treat that allocation as
your capacity.

### API to build

Decide the bounds behavior yourself, then be consistent. Both choices below are defensible; what
isn't defensible is `Get` panicking while `Set` returns an error.

| Method | Contract |
|---|---|
| `New` | Empty vector, some small starting capacity (or zero — decide, and note why) |
| `NewWithCapacity` | Empty vector, caller-chosen capacity. Length is still 0 |
| `Len` | Slots in use |
| `Cap` | Slots allocated |
| `Get(i)` | Value at index `i`. Out of range → error or panic (pick one) |
| `Set(i, v)` | Overwrite index `i`. Does **not** change length. Out of range → same policy as `Get` |
| `Push(v)` | Append to the end. Grows if full. Amortized O(1) |
| `Pop()` | Remove and return the last element. Empty → error |
| `Insert(i, v)` | Insert at `i`, shifting everything from `i` rightward. Grows if full |
| `Remove(i)` | Remove at `i`, shifting everything after `i` leftward |

Note that `Set` and `Insert` are different operations and a surprising number of people write one
when they meant the other. `Set` overwrites; `Insert` makes room.

### Growth policy

When full, allocate a new backing store **double** the current capacity and copy across. Before you
accept that doubling is right, do this on paper — it takes three minutes and it's the point of the
whole exercise:

1. Assume growth by **+1 slot** each time you're full. Push 8 elements starting from capacity 1.
   Count the total number of element-copies performed. Now do it for 16. Now generalize to N.
2. Assume **doubling**. Same exercise: 8 elements, then 16, then N.

You should end up with roughly N²/2 copies for the first and under 2N for the second. That second
result is what "amortized O(1) per push" means: individual pushes are occasionally expensive — one
push in the sequence copies N elements — but the total across N pushes is linear, so the *average*
is constant. Write both numbers down. This is the first time Big-O will feel like arithmetic rather
than vocabulary, and it's worth the three minutes.

### Go gotchas — read these before you start, not after

**Pointer receivers.** Any method that mutates the vector must take a pointer receiver. A value
receiver gets a *copy* of the struct; you will mutate the copy, the caller will see nothing change,
and nothing will error. This is the single most likely bug in this session. Symptom: `Push` appears
to work, `Len` always returns 0.

There's a nastier version of this: if your backing store is a slice, a value receiver will mutate
the *elements* successfully (both copies share the same underlying array) while failing to update
`length` or a reallocated backing store. So it half-works, which is worse than not working. Use
pointer receivers consistently for mutators and you sidestep the whole category.

**`copy`.** Go's builtin `copy` moves elements between slices and handles overlapping ranges
correctly, which matters for `Insert` and `Remove`. A hand-written loop is fine and arguably more
instructive the first time — but if you write one for `Insert`, think about which direction you have
to iterate. Going the wrong way overwrites elements before you've moved them, and you'll get a
vector where one value is smeared across several slots.

**Zeroing removed slots.** After `Pop` or `Remove`, the slot past the new length still holds the old
value. For an `int` vector this is harmless. The moment the element type is a pointer or a struct
holding one, that stale slot keeps an object alive that should have been collected — a real memory
leak, and a genuinely common one in production Go. Zero the vacated slot. It costs one assignment
and it's the right habit.

### Tests to write

Translate each into a case. Aim to have these before you have an implementation.

1. A new vector reports length 0.
2. After one push, length is 1 and `Get(0)` returns the pushed value.
3. Pushing past the initial capacity preserves every earlier element in order.
4. Pushing past the initial capacity increases capacity to the expected value (assert the doubling
   explicitly — this is the test that proves your growth policy).
5. `Pop` returns the most recently pushed value and decrements length.
6. `Pop` on an empty vector returns an error and doesn't panic.
7. `Get` with a negative index and with an index ≥ length both behave per your chosen policy.
8. `Set` overwrites in place and leaves length unchanged.
9. `Insert` in the middle shifts the tail rightward, increments length, and preserves order.
10. `Remove` from the middle shifts the tail leftward, decrements length, and preserves order.
11. `Insert` at index 0 and at index == length both work (the boundary cases where "shift the tail"
    means shifting everything or nothing).
12. Push/pop the same element 1000 times in a loop — capacity shouldn't grow without bound.

### Stretch

Make it generic: parameterise the element type. Two things change. First, comparisons in your tests
need care. Second, the zeroing advice in §2 gets more interesting — you can't assign `0` to an
arbitrary type, so you need the zero value of a type parameter. Work out how to express that.

---

## 3. Singly linked list

Do this one second, even if you don't reach it until session 3. It's where next week's hash map
comes from.

### The shape

A **node** holds a value and a pointer to the next node. The **list** holds a pointer to the first
node (the head). The last node's next pointer is `nil`, and `nil` is how you know you've reached the
end. An empty list is a list whose head is `nil` — not a special case, just the natural base case.

That `nil` is a usable, meaningful zero value rather than an error state is a genuinely Go idea, and
this structure is the cleanest place to internalise it. Your traversal loop is "while the current
pointer isn't nil, advance it" and that single pattern handles the empty list, the one-element list,
and the thousand-element list identically.

### API to build

| Method | Contract |
|---|---|
| `PushFront(v)` | New node becomes the head. O(1) |
| `PushBack(v)` | New node becomes the last. See the tail question below |
| `PopFront()` | Remove and return the head's value. Empty → error |
| `Find(v)` | Report whether `v` is present (and at what position, if you want) |
| `Remove(v)` | Remove the first node holding `v`. Report whether anything was removed |
| `Len()` | Number of nodes |
| `Reverse()` | Reverse the list in place |

### The tail-pointer question — decide deliberately

Without a tail pointer, `PushBack` walks the entire list to find the end: O(n). With a tail pointer,
it's O(1). Obviously keep a tail pointer, then — except now `Remove` and `PopFront` have to maintain
it, and getting that wrong gives you a tail pointing at a node that's no longer in the list. Every
subsequent `PushBack` then appends to a detached fragment that nothing can reach. Your list silently
stops growing.

Pick one. Then in your retro, write a sentence on what the choice cost you. If you take the tail
pointer, add a test that removes the last element and then pushes back — that's precisely the case
that breaks.

### Go gotchas

**Removal needs to see the previous node.** To unlink a node you have to change the *previous*
node's next pointer, but a singly linked list only lets you walk forward. Two solutions:

- Track a `prev` variable as you walk. Straightforward, but the head is a special case (it has no
  previous node) and you'll write an `if` for it.
- Walk a **pointer to a pointer** — hold the address of the link you'd have to modify, rather than
  the node before it. The head pointer and every next pointer are then uniform, and the special case
  disappears entirely.

Write the first version. Get it passing. Then try the second — it's the more idiomatic Go and the
moment it clicks is a genuine step up in how you read pointer code. Keeping both (one committed,
then replaced) also gives you something real to look at in `git diff`.

**A dropped link hangs, it doesn't crash.** If you reassign pointers in the wrong order during
`Reverse` or `Remove` and accidentally point a node at itself or at an earlier node, your traversal
loop never reaches `nil`. Go won't panic — `go test` just sits there. If a test hangs for more than
a couple of seconds, that's what happened; kill it with `-timeout 5s` and go look at the order of
your assignments. The rule for `Reverse` is: save the next pointer *before* you overwrite it.

**Length: counter or walk?** Storing a length field makes `Len` O(1) but creates an invariant every
mutator must maintain — miss one and the field silently drifts from reality. Walking the list makes
`Len` O(n) but impossible to get wrong. Both are legitimate. If you store the counter, make one of
your tests assert length after a long mixed sequence of pushes and removes, because that's the test
that catches drift.

### Tests to write

1. A new list has length 0 and `Find` on it reports not-found.
2. `PushFront` then `PopFront` round-trips the value.
3. `PushFront` three times produces reverse insertion order when traversed.
4. `PushBack` three times produces insertion order when traversed.
5. `PopFront` on an empty list errors rather than panicking or returning a garbage value.
6. `Remove` of the **head**, a **middle** node, and the **last** node each work (three cases — these
   are three different code paths in most implementations).
7. `Remove` of a value that isn't present reports not-found and leaves the list unchanged.
8. `Remove` from a single-element list leaves an empty, still-usable list — then push onto it again
   and confirm it works.
9. `Reverse` on empty and single-element lists is a no-op that doesn't crash.
10. `Reverse` twice returns the original order.
11. If you kept a tail pointer: remove the last element, then `PushBack`, then confirm the new
    element is reachable from the head.

Cases 6, 8 and 11 are where the bugs actually are. Write them first.

### Stretch

`Reverse` in place, with no allocation. Three pointers — previous, current, next — walked once
across the list. Don't look it up. Draw four nodes on paper, decide which pointer has to move first,
and work out why saving `next` before overwriting `current.next` is non-negotiable. It's a small
puzzle and solving it yourself is worth more than reading the eight lines.

---

## 4. The table you fill in

This is the actual deliverable of the week. The implementations are how you earn the right to fill
it in from memory rather than from a lookup.

| Operation | Dynamic array | Singly linked list |
|---|---|---|
| Access by index | | |
| Push to back | | |
| Push to front | | |
| Insert in middle (position known) | | |
| Remove from middle (position known) | | |
| Search for a value | | |
| Memory overhead per element | | |

Two of those cells should make you pause: "insert in middle, position known" is the one where the
linked list's theoretical advantage lives, and the reason it rarely materialises in practice is the
row above it plus the row below it. Once the table's filled, write two sentences below it:

> I would reach for a linked list over a dynamic array when ______, because ______.

If you can't complete that sentence honestly, that's a real finding and worth writing down as one —
the answer in modern practice is "almost never, and mostly as a component inside something else,"
which is exactly the role it plays in next week's hash map.

---

## 5. Session 3 and the retro

Finish what's unfinished and get `go test ./...` green in both packages.

**Commit per structure, not once at the end.** `git commit` after the vector passes, again after the
list passes, again after each stretch goal. Week 3's topic is git beyond the basics — `rebase`,
`bisect`, conflict resolution — and those are miserable to practice against a repo with one commit
in it. You're building the material you'll use later.

**Write the retro here**, at the bottom of this file, before you close the session:

- What broke, and how long until you saw why?
- Which of the gotchas in §2/§3 did you hit anyway?
- Tail pointer or not, `prev` or pointer-to-pointer — what did you pick and would you pick it again?
- What would you do differently?

Then add a line or two to `notes/learning-log.md`.

**Next week:** the hash map. Buckets are a dynamic array. Each bucket's collision chain is a linked
list. You'll have already written both halves — the new material is just the hash function and the
load-factor resize, which is the doubling logic from §2 applied one level up.

---

## Retro — Week 1, Session 2

<!-- write it here -->