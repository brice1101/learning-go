# Backend & Systems Engineering Learning Path

For: self-taught programmer, no professional experience, backend/systems focus. Primary language: Go.

This is a **progression, not a schedule**. Work through the phases in order, tick things off as they
become solid, and move on when the milestone is done — not when a calendar says so. Projects matter
more than tutorials: you learn backend engineering by building things that have to survive real
conditions (concurrency, failure, scale), not by reading about them.

Skip anything you can already do fluently. Move faster through familiar material and spend the time
you save on the harder phases (3 and 4).

---

## Working with Claude in this repo

- **Specs, not solutions.** When asked for guidance on a new topic, write a spec in `notes/` (what to
  build, API contract, gotchas, tests to write), with no Go in it. Every line of code under
  `datastructures/` (and later projects) is written by me.
- **Reviews are findings only.** When asked to review code, report issues; don't apply fixes.
- **Keep this file's "Current focus" and checklists current** when something is finished.

---

## How a session runs

Sessions are ~1.5h, roughly three a week, but the unit of progress is the topic, not the week.
Each session is one of:

- **Learn:** read or watch the core material, take notes in your own words (in the topic's note file)
- **Build:** apply it directly in code
- **Debug / finish:** get it green, commit, write the retro in the topic's note file

At the end of every session, add a dated entry to `notes/learning-log.md`.

### Notes layout

| File | Purpose |
|---|---|
| `notes/learning-log.md` | Running dated log: one short entry per session |
| `notes/NN-topic.md` | One file per topic: spec, my notes, decisions, retro |
| `notes/_template.md` | Copy this to start a new topic file |

---

## Current focus

> **Now:** Phase 1: hash map (next milestone structure)
> **Last finished:** dynamic array + singly linked list (`notes/01-vector-and-linked-list.md`)
> **Blocked on / questions:** _

---

## Phase 1: Foundations You Can't Skip

**Goal: lock in the CS fundamentals that backend work assumes you have, even if you learned to code
without them.**

### Progress

- [x] Dynamic array (`datastructures/vector`)
- [x] Singly linked list (`datastructures/linkedlist`)
- [ ] Hash map: buckets + chaining, load-factor resize
- [ ] Binary search tree: insert, search, delete, in-order traversal
- [ ] Graph + BFS / DFS
- [ ] Heap / priority queue
- [ ] Big-O: can explain amortized cost, and N+1 queries as a Big-O problem
- [ ] Git: rebase, bisect, resolving a conflict, a branching strategy I'd defend
- [ ] NeetCode drills: 0 / 10–15

### Topics

- Data structures: arrays, hash maps, trees, graphs, heaps; when to use which
- Algorithmic complexity: Big-O, why it matters for backend (N+1 queries, scaling)
- Git beyond the basics: rebasing, bisecting, resolving conflicts, branching strategy

### Resources

- Grokking Algorithms (Aditya Bhargava): fast, visual, no fluff
- NeetCode.io: data structure drills, use sparingly (10–15 problems total, not a grind)
- Official Git docs / Pro Git book, chapters 3 & 7

**Milestone:** implement a hash map, a binary search tree, and a basic graph traversal (BFS/DFS)
from scratch, no libraries.

---

## Phase 2: Core Backend Skills

**Goal: build a real API-backed service, end to end.**

### Progress

- [ ] HTTP: methods, status codes, headers, statelessness, REST conventions
- [ ] REST API in Go (net/http or chi)
- [ ] Postgres schema design: normalization, indexes, transactions, EXPLAIN
- [ ] SQL beyond CRUD: joins, aggregates, query optimization
- [ ] Auth: sessions vs JWT, password hashing (bcrypt/argon2)
- [ ] Testing: unit, integration, mocking external dependencies
- [ ] **Milestone:** task/notes API, deployed

### Resources

- Designing Data-Intensive Applications (Kleppmann): start here, read a chapter alongside each
  chunk of building; it's the single best book for this whole path
- PostgreSQL official tutorial + "Use The Index, Luke" (use-the-index-luke.com)
- Go's `testing` package docs

**Milestone:** a task/notes API with user auth, PostgreSQL persistence, proper indexing, and a test
suite. Deploy it somewhere (Railway, Fly.io, or a $5 VPS) so it's a real running service.

---

## Phase 3: Systems & Concurrency

**Goal: understand what happens under the API layer. This is what separates "can build a CRUD app"
from "backend/systems engineer."**

### Progress

- [ ] Concurrency: goroutines/channels, threads vs processes, race conditions, locks
- [ ] Networking: TCP/IP basics, DNS, load balancers, a request's full journey
- [ ] Caching: Redis, invalidation strategies, when caching helps vs hurts
- [ ] Message queues: why they exist, at-least-once vs exactly-once delivery
- [ ] OS basics: processes, memory, file descriptors, running out of connections
- [ ] **Milestone:** Phase 2 API + Redis cache + job queue, load-tested until it breaks

### Resources

- Designing Data-Intensive Applications: continue (replication, partitioning)
- "Computer Networking: A Top-Down Approach": skim relevant chapters
- Redis docs + build something with it directly

**Milestone:** extend the Phase 2 API: add Redis caching for hot reads, add a background job queue
for a slow operation (e.g. sending emails, processing uploads), and load-test it (hey or k6) to see
where it breaks. Breaking it on purpose is the point.

---

## Phase 4: Distributed Systems & System Design

**Goal: think at the level system design interviews and real infra decisions require.**

### Progress

- [ ] CAP theorem, consistency models, replication strategies
- [ ] Horizontal scaling: stateless services, sharding, read replicas
- [ ] Observability: structured logging, metrics, basic tracing
- [ ] Case studies read (Twitter timeline, Stripe idempotency, ...): 0
- [ ] Design doc 1: _
- [ ] Design doc 2: _
- [ ] Design doc 3: _
- [ ] **Milestone:** Phase 3 service containerized with health checks + logging

### Resources

- Designing Data-Intensive Applications: finish it
- "System Design Interview" (Alex Xu) Vol. 1: good for structured practice
- ByteByteGo YouTube/newsletter for quick case-study exposure

**Milestone:** write 3 system design docs from scratch (a design doc, not just a diagram) for classic
problems (URL shortener, rate limiter, chat system), then containerize the Phase 3 service with
Docker, and add basic health checks + logging.

---

## Phase 5: Capstone

**Goal: one polished, portfolio-grade project that demonstrates the whole stack of skills.**

Options (pick one, or propose your own):

- A distributed rate limiter or job scheduler
- A simplified real system: Twitter clone with proper fan-out, URL shortener with analytics, chat
  backend with WebSockets

Chosen: _

### Requirements

- [ ] REST or gRPC API, tested
- [ ] Postgres with a deliberate, indexed schema
- [ ] Redis caching or a queue, used for a real reason (not decoration)
- [ ] Dockerized, deployed somewhere real
- [ ] Basic logging/metrics
- [ ] Written design doc explaining the trade-offs
- [ ] README with architecture diagram + live link

---

## Ongoing habits

- Read one well-known engineering blog post a week (Stripe, Cloudflare, Netflix, Discord). Note it
  in the learning log.
- Keep `notes/learning-log.md` going. It's useful for interviews and for noticing your own progress.
- Commit per working increment, not once at the end. The history is practice material for git.
- Don't rush past Phase 1 if it's genuinely new. Everything after depends on it being solid.

## After the capstone

- Distributed consensus (Raft) and a toy implementation
- Kubernetes fundamentals
- A systems language deep-dive (Rust)
- Contributing to an open-source backend project