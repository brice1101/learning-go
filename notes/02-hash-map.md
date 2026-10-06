# 02: Hash Map

**Phase:** 1 (Foundations)
**Status:** not started
**Code:** `datastructures/hashmap`

This is a spec, not a tutorial. There is no Go in this file on purpose.

**Suggested order:** one learn session on §1–§2 (theory + paper exercise, plus Grokking Algorithms
ch. 5), then build sessions on §3–§5.

---

## 1. Theory

### The problem

You want to store `key → value` pairs and look a value up by its key, fast. With what you've built
so far:

- **Linked list of pairs:** lookup walks the list: O(n).
- **Dynamic array of pairs:** lookup scans the array: O(n). (Sorted + binary search gets you
  O(log n), but then every insert has to shift to stay sorted.)

The one fast operation you have is **array access by index: O(1)**. A hash map is a trick for
turning a key into an array index.

### The hash function

A **hash function** takes a key (say, the string `"alice"`) and returns a big number. Then:

```
index = hash(key) mod number_of_slots
```

To look up `"alice"`, hash it, compute the index, and go straight to that slot. No scanning.

A hash function must be:

- **Deterministic:** the same key always gives the same number. Otherwise you'd store under one
  index and look under another.
- **Well spread:** similar keys (`"alice"`, `"alicf"`) should give very different numbers, so keys
  scatter evenly across the slots instead of piling into a few.

### Buckets

Each slot in the array is called a **bucket**. Why "bucket" and not just "slot"? Because a slot
sometimes has to hold more than one entry.

### Collisions

Two different keys can hash to the same index. That's a **collision**, and it's guaranteed to
happen. There are infinitely many possible strings and only so many buckets, so some keys must
share (the pigeonhole principle). Collisions also happen much sooner than intuition suggests: with
365 buckets, you'll probably get one after only 23 keys (the "birthday paradox").

So every hash map needs a collision strategy. There are two main families:

- **Separate chaining** (what you'll build): each bucket holds a **linked list** of the entries that
  landed there, called its **chain**. Insert: hash to the bucket, add to its chain. Lookup: hash to
  the bucket, walk its chain comparing keys. With a good hash function and enough buckets, chains
  stay very short (usually 0–2 entries), so the walk is effectively O(1).
- **Open addressing:** each bucket holds at most one entry. On a collision, probe for another free
  bucket (the next one, or some other rule). No lists, better for the CPU cache, but deletion gets
  tricky. Go's built-in `map` uses a variant of this. It's a stretch goal here.

```
buckets (a dynamic array)
 [0] → nil
 [1] → ("bob", 7) → nil
 [2] → nil
 [3] → ("alice", 3) → ("dave", 9) → nil      ← collision: both hashed to 3
 [4] → nil
 ...
```

This is why the dynamic array and linked list came first: **the bucket array is a dynamic array,
and each chain is a linked list.**

### Load factor and resizing

**Load factor** = number of entries ÷ number of buckets. At 0.5, there's an entry for every other
bucket, and chains are short. At 5, chains average 5 long and lookups slow down. If the bucket count
never changes, lookup degrades to O(n) as the map fills.

So when the load factor passes a threshold (0.75 is common), **double the bucket count**. That's
the same doubling as the vector, for the same amortized-O(1) reason.

One important difference from the vector: **you can't just copy entries across.** An entry's bucket
is `hash mod bucket_count`, and the bucket count just changed, so most entries now belong in a
*different* bucket. Resizing means re-inserting every entry into the new array. That's called
**rehashing**.

### Complexity

| Operation | Average | Worst case |
|---|---|---|
| Get / Put / Delete | O(1) | O(n): every key in one chain |
| Resize | O(n), but rare: amortized O(1) per Put | |

The worst case is real, not theoretical. If an attacker can choose your keys and knows your hash
function, they can send keys that all collide and turn every request into an O(n) scan (**hash
flooding**). This is why Go uses a random hash seed per map, and why iterating a Go `map` gives a
different order each time. Remember this for Phase 2: it's a backend security issue.

---

## 2. Paper exercise: do this before any code

Use a deliberately bad toy hash: **add up the letter positions (a=1, b=2, …, z=26), mod the bucket
count.**

1. With **4 buckets**, compute the bucket for: `cat`, `dog`, `act`, `bird`, `fish`. Draw the
   buckets and chains.
2. Which keys collided? Look at `cat` and `act`. Why is summing letters a bad hash function? What
   property from §1 does it violate?
3. What's the load factor after all 5 inserts? Assume a 0.75 threshold: at which insert should the
   map have resized?
4. **Resize to 8 buckets** and recompute every key's bucket. Which keys moved? This is why resizing
   is a rehash and not a copy.

Record your answers in **My notes** below.

---

## 3. The hash function you'll actually use: FNV-1a

A real, simple, well-spread hash function. It's a few lines, so write it yourself. 64-bit version:

1. Start with `hash = 14695981039346656037` (the "offset basis").
2. For each **byte** of the key:
   - `hash = hash XOR byte`
   - `hash = hash × 1099511628211` (the "FNV prime")
3. Return `hash`.

Use `uint64` throughout. The multiplication will overflow constantly, and that's intended: in Go,
unsigned integer overflow is defined behavior and simply wraps around. The wrapping is part of what
mixes the bits.

Sanity check: FNV-1a 64 of the empty string is the offset basis itself, and `"a"` should give
`0xaf63dc4c8601ec8c`. Make these your first two tests.

---

## 4. Build spec

### Scope

- Keys are `string`, values are `int`. Keep it concrete, as with the vector. Generics are a stretch
  goal.
- Don't use Go's `map`, obviously. Slices are fine for the bucket array (`make` at a fixed size, no
  `append`), since the bucket array only ever changes size when you rehash into a brand new one.

### The chain: don't reuse `linkedlist` (and notice why)

Your `linkedlist` package holds `int` values and finds by value. A chain entry needs a **key, a
value, and a next pointer**, and lookup compares **keys**. So you'll write a small unexported entry
type inside the `hashmap` package. That's a normal outcome: the *technique* carries over (the
`n != nil` loop, and pointer-to-pointer for deletion), even though the package doesn't. Write a
sentence in your notes about what would have to change for `linkedlist` to be reusable here.

### API

| Method | Contract |
|---|---|
| `New()` | Empty map with a small starting bucket count (8 is common). Decide if the zero value should work too |
| `Put(key, value)` | Insert, **or update** if the key already exists. Resizes when the load factor passes the threshold |
| `Get(key)` | `(value, ok)`, the same comma-ok shape as your vector |
| `Delete(key)` | Remove the entry. Report whether anything was removed |
| `Len()` | Number of entries, **O(1)** (see below) |

### Decisions to make deliberately

- **`Len`: counter this time, not a walk.** You need the entry count on every `Put` to compute the
  load factor, so walking all buckets every time would make `Put` O(n). This is the opposite
  choice from the linked list, and the reason differs. Every mutator must now keep the counter
  honest.
- **Where do new entries go in a chain: front or back?** One is O(1) with no tail pointer. But
  `Put` has to walk the chain anyway to check for an existing key...
- **Resize threshold and trigger:** check before inserting or after? Does an *update* (existing
  key) ever need to trigger a resize?

### Gotchas

- **Update vs insert.** `Put` on an existing key must overwrite, not add a second entry. Otherwise
  `Get` returns whichever it finds first, `Len` is wrong, and `Delete` leaves a zombie behind.
  The single most likely bug.
- **Negative indexes.** If a hash ever passes through a signed `int`, `%` in Go keeps the sign
  (`-7 % 4 == -3`), and indexing with that panics. Do the `mod` on the `uint64`, *then* convert.
- **Rehash, don't copy.** On resize, every entry gets its bucket recomputed against the new
  bucket count.
- **The counter drifts.** `Len` must go up only on a real insert and down only on a real delete.
- **Delete from a chain is linked-list `Remove` again.** Pointer-to-pointer makes deleting the
  first entry in a bucket the same as any other. You've already written this.

### Testing collisions: make the hash swappable

With a good hash, collisions are rare and random, so you can't reliably test chain handling. Make
the hash function **replaceable**: store it as a field on the map (a function value), which `New`
sets to FNV-1a and tests can override. A test that swaps in "always return 0" forces every key into
one bucket, and your map becomes a single linked list you can test deliberately.

This is your first taste of **dependency injection**: making a dependency swappable so tests can
control it. Phase 2 uses the same idea to mock databases and HTTP clients.

### Tests

1. FNV-1a: empty string and `"a"` match the known values (§3).
2. New map: `Len` is 0, `Get` of anything reports not-found.
3. `Put` then `Get` round-trips.
4. `Put` the same key twice: the second value wins, and `Len` is still 1.
5. `Delete` a present key: returns true, `Get` then reports not-found, `Len` decrements.
6. `Delete` an absent key: returns false, `Len` unchanged.
7. The empty string `""` works as a key.
8. Put 1,000 distinct keys: every one is retrievable with the right value (this forces several
   resizes).
9. After those 1,000, the load factor is at or below your threshold.
10. **Forced collisions** (hash always returns 0): put 3 keys, then delete the first, middle, and
    last of the chain, checking the others are still retrievable after each delete.
11. Forced collisions: update a key in the middle of a chain.
12. Put, delete, and re-put the same key: `Len` is correct at every step.

Tests 4, 10 and 12 are where the bugs are. Write them first.

### Stretch

- [ ] Power-of-two bucket counts let you replace `hash mod n` with `hash & (n-1)`. Work out why
  that's equivalent, and why it only works for powers of two.
- [ ] Benchmark your map against Go's built-in `map` with `go test -bench`. How far off are you?
- [ ] Generic version: `HashMap[K comparable, V any]`. The hard part: how do you hash an arbitrary
  `K`? Look at `hash/maphash`.
- [ ] Shrink when the load factor drops very low after many deletes.
- [ ] An open-addressing version (linear probing). Compare deletion complexity with chaining.

---

## 5. Resources

- **Grokking Algorithms, ch. 5 (Hash Tables):** read this in the learn session, before §2.
- After you've built it: "Faster Go maps with Swiss Tables" on the Go blog (go.dev/blog), which
  describes how Go's built-in map works today. It's much more readable once you've built the simple
  version.

---

## My notes

**Paper exercise (§2):**

| Key | Letter sum | Bucket (4) | Bucket (8) |
|---|---|---|---|
| cat | | | |
| dog | | | |
| act | | | |
| bird | | | |
| fish | | | |

- Collisions at 4 buckets:
- Why summing letters is a bad hash:
- Load factor after 5 inserts / when it should have resized:
- Keys that moved on resize:

**In my own words:**

- What a hash map is, in one sentence:
- Bucket:
- Chaining:
- Load factor:
- Why resize has to rehash:

**Decisions I made and why:**

- Starting bucket count / does the zero value work:
- Load factor threshold, and when it's checked:
- New entries at front or back of chain:
- What would need to change for `linkedlist` to be reusable as a chain:

**Complexity:**

| Operation | Average | Worst | Why |
|---|---|---|---|
| Get | | | |
| Put | | | |
| Delete | | | |
| Resize | | | |

---

## Retro

- **What broke, and how long until I saw why?**
  -
- **What surprised me?**
  -
- **What would I do differently?**
  -

---

## Next

BST: the same pointer and `nil` discipline as the linked list, plus recursion.