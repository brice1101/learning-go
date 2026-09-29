# 01: Dynamic Array & Singly Linked List

**Phase:** 1 (Foundations)
**Status:** written up
**Code:** `datastructures/vector`, `datastructures/linkedlist`

---

## 0. Why these two first

The Phase 1 milestone is a hash map, a BST, and a graph traversal. The hash map is the wrong place to
start, because it's two unfamiliar problems stacked on top of each other: a hash function *and* a
collision strategy, debugged through pointer semantics you haven't used yet. When it breaks, you
won't know which layer broke.

Build the dependencies first:

```
dynamic array ──┬──> hash map (buckets are an array, chains are a list)
singly linked ──┘
    list        ───> BST        (same pointer + nil discipline, plus recursion)
dynamic array   ───> BFS / DFS  (a queue and a stack are both just this)
```

---

## 1. Test loop

| Command | What it does |
|---|---|
| `go test ./...` | Runs every test in the module |
| `go test ./datastructures/vector` | Just this package |
| `go test -v ./...` | Names each test as it runs |
| `go test -run TestPush ./...` | Only tests whose name matches the regex |
| `go test -timeout 5s ./...` | Kills a hung test (see §3) |

**Table-driven tests** are the Go convention: a slice of anonymous structs (name, inputs, expected
output), ranged over, with `t.Run` per case. Adding a case is one line instead of one function.

**Write the test first.** Writing the call site first is the fastest way to find out what API you
want.

---

## 2. Dynamic array

### The idea

A fixed-size block of memory, plus a counter tracking how much of it you're actually using:

- **length:** how many slots hold real values
- **capacity:** how many slots you've allocated

Length ≤ capacity, always. When full, allocate a bigger block and copy everything across. That copy
is the whole subject of this exercise. **No `append`**: it's the mechanism under study.

### API

| Method | Contract |
|---|---|
| `New` | Empty vector, some small starting capacity (or zero; decide, and note why) |
| `NewWithCapacity` | Empty vector, caller-chosen capacity. Length is still 0 |
| `Len` / `Cap` | Slots in use / slots allocated |
| `Get(i)` | Value at `i`. Out of range → error or panic (pick one, be consistent) |
| `Set(i, v)` | Overwrite `i`. Does **not** change length |
| `Push(v)` | Append to the end. Grows if full. Amortized O(1) |
| `Pop()` | Remove and return the last element. Empty → error |
| `Insert(i, v)` | Insert at `i`, shifting the tail right. Grows if full |
| `Remove(i)` | Remove at `i`, shifting the tail left |

### Growth policy: do this on paper

1. Growth by **+1 slot** when full. Push 8 elements from capacity 1. Count total element-copies.
   Then 16. Then generalise to N.
2. Growth by **doubling**. Same exercise.

You should land near N²/2 for the first and under 2N for the second. That second result is what
"amortized O(1) per push" means.

### Gotchas

- **Pointer receivers** for anything that mutates. A value receiver mutates a copy; with a slice
  backing store it half-works (elements change, `length` doesn't), which is worse than not working.
- **`copy` handles overlap.** If you hand-write the shift loop for `Insert`, iterate in the right
  direction or one value gets smeared across several slots.
- **Zero vacated slots** after `Pop`/`Remove`. Harmless for `int`, a real memory leak once the
  element type holds a pointer.

### Tests

1. New vector has length 0.
2. One push → length 1, `Get(0)` returns it.
3. Pushing past initial capacity preserves every earlier element in order.
4. Pushing past initial capacity doubles capacity (assert it explicitly).
5. `Pop` returns the last pushed value and decrements length.
6. `Pop` on empty errors, doesn't panic.
7. `Get` with negative index and index ≥ length both follow your policy.
8. `Set` overwrites in place, length unchanged.
9. `Insert` in the middle shifts right, increments length, preserves order.
10. `Remove` from the middle shifts left, decrements length, preserves order.
11. `Insert` at 0 and at `Len()` both work.
12. Push/pop 1000 times: capacity doesn't grow without bound.

### Stretch

- [ ] Make it generic. You'll need the zero value of a type parameter for zeroing slots.

### My notes: dynamic array

**Paper exercise results:**

Starting from capacity 1, a grow copies every element currently stored.

| Growth | Copies for N=8 | Copies for N=16 | General N |
|---|---|---|---|
| +1 slot | 1+2+…+7 = 28 | 1+2+…+15 = 120 | N(N−1)/2 → O(N²) total, O(N) per push |
| Doubling | 1+2+4 = 7 | 1+2+4+8 = 15 | < N → O(N) total, O(1) amortized per push |

The takeaway: with doubling, each grow is twice as expensive as the last but happens half as
often, so the total stays linear. With +1, *every* push past the start is a grow, so it's quadratic.
At N=16 that's already 120 vs 15, and the gap widens as N grows.

**Decisions I made and why:**

- **Starting capacity: 0.** `New()` returns an empty struct and the first `Push` grows 0 → 1.
  This means the zero value `Vector{}` also works, with no constructor required, which is the Go
  idiom (`bytes.Buffer`, `sync.Mutex` work the same way).
- **Out-of-range policy: return a `bool` (comma-ok), never panic.** `Get`/`Pop`/`Remove` return
  `(value, ok)` and `Set`/`Insert` return `ok`. It's the same shape as a Go map lookup
  (`v, ok := m[k]`), and it's consistent across every method. An `error` would carry more detail,
  but there's only one way any of these can fail, so a bool says everything.
- **`NewWithCapacity` with a negative capacity returns `nil`.** It's arguable. A caller who doesn't
  check will get a nil-pointer panic later, far from the cause. Something to revisit.
- **Capacity isn't a field.** My first version stored `capacity` alongside `data`. I refactored it
  to derive `Cap()` from `len(data)`, because two fields that must always agree will eventually
  disagree. One source of truth means that bug can't happen.
- **`copy` for `Insert`, hand-written loop for `Remove`.** The loop in `Remove` iterates forward
  (left to right), which is the safe direction for a left shift. `copy` handles overlap on its own.
- **`grow` takes a gap index.** When `Insert` triggers a reallocation, the gap is left open *during*
  the copy into the new array, instead of copying everything and then shifting again. One pass
  instead of two.

**Notes in my own words:**

- A Go slice is already this structure: a pointer to an array, a length, and a capacity. Building
  it by hand made it obvious why `append` sometimes returns a slice pointing at a *different* array,
  and why you must always write `s = append(s, x)`.
- Length and capacity are different questions: "how much is real" vs "how much is allocated". Every
  bounds check is against length. Capacity only matters to `grow`.
- "Amortized O(1)" doesn't mean every push is cheap. It means the occasional expensive push is paid
  for by all the cheap ones before it. A latency-sensitive system can still notice the spike.
- Zeroing a vacated slot is pointless for `int` but prevents a leak once the element holds a pointer,
  because the garbage collector can't free something the backing array still references.

---

## 3. Singly linked list

### The shape

A **node** holds a value and a pointer to the next node. The **list** holds a pointer to the head.
The last node's next is `nil`, and an empty list is just a `nil` head: the natural base case, not a
special one. "While current isn't nil, advance" handles empty, one-element and thousand-element
lists identically.

### API

| Method | Contract |
|---|---|
| `PushFront(v)` | New node becomes the head. O(1) |
| `PushBack(v)` | New node becomes the last. See the tail question |
| `PopFront()` | Remove and return the head's value. Empty → error |
| `Find(v)` | Report whether `v` is present (and at what position) |
| `Remove(v)` | Remove the first node holding `v`. Report whether anything was removed |
| `Len()` | Number of nodes |
| `Reverse()` | Reverse in place |

### Decisions to make deliberately

- **Tail pointer?** Makes `PushBack` O(1), but every remover must maintain it. Get it wrong and
  `PushBack` appends to a detached fragment, and the list silently stops growing.
- **`prev` variable or pointer-to-pointer for `Remove`?** `prev` special-cases the head;
  pointer-to-pointer makes the head link and every `next` link uniform.
- **Length: counter or walk?** Counter is O(1) but can drift if any mutator forgets it. Walking is
  O(n) but can't be wrong.

### Gotchas

- **A dropped link hangs, it doesn't crash.** A cycle means traversal never reaches `nil`, and
  `go test` just sits there. Use `-timeout 5s`, then check your assignment order.
- **`Reverse`:** save `next` *before* overwriting `current.next`.

### Tests

1. New list: length 0, `Find` reports not-found.
2. `PushFront` then `PopFront` round-trips.
3. `PushFront` ×3 gives reverse insertion order.
4. `PushBack` ×3 gives insertion order.
5. `PopFront` on empty errors.
6. `Remove` of head, middle, and last node (three separate code paths).
7. `Remove` of an absent value reports not-found, list unchanged.
8. `Remove` from a single-element list, then push again, and it works.
9. `Reverse` on empty and single-element lists doesn't crash.
10. `Reverse` twice restores the original order.
11. With a tail pointer: remove the last element, `PushBack`, confirm reachable from head.

Cases 6, 8 and 11 are where the bugs are.

### My notes: linked list

**Decisions I made and why:**

- **Tail pointer: yes.** `PushBack` is O(1). The cost: three methods have to maintain it.
  `PushFront` sets it on an empty list, `PopFront` clears it when the list empties, and `Remove`
  moves it back when the last node goes. That last case is the one that needed a test (#11).
- **`Remove`: started with `prev`, then switched to pointer-to-pointer.** The first version had
  three sections: a head special case (which delegated to `PopFront`), a loop for the rest, and tail
  fix-up. With pointer-to-pointer, `link` holds the *address of the pointer to change*, whether
  that's `l.head` or some node's `next`, so "unlink" is one line: `*link = (*link).next`, and the
  head case disappears.
- **...but I still kept a `prev` alongside it**, only for the tail pointer: when the removed node was
  last, the new tail is the node before it, and `link` points *into* that node but isn't the node
  itself. So the tail pointer brought back a little of the bookkeeping pointer-to-pointer removed.
- **`Len`: walk, don't count.** O(n), but it can't drift out of sync with the actual nodes. With no
  length field, there's no invariant for the mutators to break.
- **No constructor.** `Linkedlist{}` is a valid empty list (nil head, nil tail), so the zero value
  is ready to use.

**Notes in my own words:**

- In Python every variable is already a reference, so this is invisible. In Go I choose: `Node` is
  the value, `*Node` is its address. The list is just a chain of addresses, and `nil` means "no
  next node".
- A pointer to a pointer (`**Node`) is "the address of a box that holds an address." Writing through
  it (`*link = ...`) changes *which node* that box points to, which is exactly what unlinking needs.
- Loop on `n != nil`, not `n.next != nil`. My early `Find` and `Len` looped on `n.next` and then
  needed an extra check after the loop for the last node. Checking the current node instead handles
  empty, one-element and many-element lists with no special cases.
- `Reverse` is three pointers: `prev`, `curr`, `next`. Save `next` first, or you lose the rest of
  the list the moment you rewrite `curr.next`.

---

## 4. Comparison table

Fill this in from memory, not from a lookup.

| Operation | Dynamic array | Singly linked list |
|---|---|---|
| Access by index | O(1) | O(n): walk from head |
| Push to back | O(1) amortized | O(1) with tail pointer (O(n) without) |
| Push to front | O(n): shift everything right | O(1) |
| Insert in middle (position known) | O(n): shift the tail | O(1) if you hold the node before it |
| Remove from middle (position known) | O(n): shift the tail | O(1) if you hold the node before it; O(n) to find it in a singly linked list |
| Search for a value | O(n), fast in practice (contiguous, cache-friendly) | O(n), slower in practice (a pointer chase per node) |
| Memory overhead per element | ~0, plus up to 2× slack right after a grow | one `next` pointer (8 bytes) + a separate heap allocation per node |

The linked list's "O(1) insert in the middle" has a catch: you almost never *have* the position.
You get there by searching (O(n), the row below) or by indexing (O(n), the top row). So in practice
it's O(n) either way, and the array's O(n) is faster because it's one contiguous block of memory
the CPU can stream through.

> I would reach for a linked list over a dynamic array when **I already hold a reference to where
> the change happens and never need index access (e.g. a hash map bucket's collision chain, or an
> LRU cache's eviction order)**, because **then insert and remove really are O(1), and nothing has
> to be shifted or reallocated.** Otherwise, the dynamic array wins.

---

## 5. Retro

- **What broke, and how long until I saw why?**
  - A test named `PopFrontOnEmptyErrors` (missing the `Test` prefix) never ran. `go test` gave no
    warning. It just passed with one fewer test. Lesson: a passing suite only proves the tests
    that *ran* passed. `go test -v` lists every test by name, so a missing one is visible.
  - `Find` and `Len` originally looped on `n.next != nil`, which skipped the last node. I patched it
    with an extra check after the loop, and later fixed the loop condition itself so no patch was
    needed.
  - Hit a merge conflict in both linked list files when pulling from the remote, due to working on separate machines
- **Which gotchas did I hit anyway?**
  - The off-by-one on the last node (above), which is the "loop on the current node, not the next
    one" point from §3.
  - Pointer-to-pointer took real time to click.
- **Would I make the same decisions again?**
  - Tail pointer: yes. The O(1) `PushBack` matters when this becomes a queue for BFS.
  - Comma-ok bools instead of errors: yes for these, since each method only has one failure mode.
  - `NewWithCapacity` returning `nil` on bad input: probably not. Clamping to 0 or returning an
    error would fail closer to the cause.
- **What would I do differently?**
  - Run `go test -v` from the start, so a test that silently doesn't run is obvious.
  - Now that the loops start at `n := l.head` and check `n != nil`, the early
    `if l.head == nil` returns in `Find`, `Len` and `Remove` are redundant. The loop already
    handles the empty list. Worth deleting as a check that I trust that pattern.
  - Commit the `prev` version and the pointer-to-pointer version as separate commits (I did, so
    `git diff 00e2709 cdd4f68` shows the change), and do the same for future refactors.

---

## Next

**Hash map:** buckets are a dynamic array, and each bucket's collision chain is a linked list. The
new material is the hash function and the load-factor resize, which is the doubling logic from §2
applied one level up.