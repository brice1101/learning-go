# Backend & Systems Engineering Learning Path

For: self-taught programmer, no professional experience, backend/systems focus Cadence: 3 sessions × 1.5h/week (4.5h/week) — ~6 months, 24 weeks Format: each phase = weeks + goals + resources + a milestone project. Projects matter more than tutorials — you learn backend engineering by building things that have to survive real conditions (concurrency, failure, scale), not by reading about them.

## How to run each week
- Session 1 (1.5h): learn — read/watch the core material, take notes in your own words
- Session 2 (1.5h): build — apply it directly in your project
- Session 3 (1.5h): build/debug — finish the week's increment, write a short retro (what broke, what you'd do differently)

Skip anything you can already do fluently — this plan assumes real but informal programming experience, so don't re-learn syntax you already know. Move faster through familiar material and bank the extra time for the harder phases (3 and 4).

## Phase 1 — Foundations You Can't Skip (Weeks 1–4)

**Goal: lock in the CS fundamentals that backend work assumes you have, even if you learned to code without them.**

### Topics

- Data structures: arrays, hash maps, trees, graphs, heaps — when to use which
- Algorithmic complexity: Big-O, why it matters for backend (N+1 queries, scaling)
- Git beyond the basics: rebasing, bisecting, resolving conflicts, branching strategy
- Pick your primary language for this path if you haven't: Go is the most common recommendation for backend/systems learning (simple, built for concurrency, widely used in infra) — but if you're already strong in Python or Java, you can stay there and add Go later in Phase 3.

### Resources

- Grokking Algorithms (Aditya Bhargava) — fast, visual, no fluff
- NeetCode.io — data structure drills, use sparingly (10–15 problems total, not a grind)
- Official Git docs / Pro Git book, chapters 3 & 7

Milestone: implement a hash map, a binary search tree, and a basic graph traversal (BFS/DFS) from scratch, no libraries. This isn't busywork — it's the fastest way to internalize how these structures actually behave.

## Phase 2 — Core Backend Skills (Weeks 5–10)

**Goal: build a real API-backed service, end to end.**

### Topics

- HTTP deeply: methods, status codes, headers, statelessness, REST conventions
- Building a REST API in your chosen language (Go: net/http or chi; Python: FastAPI; Java: Spring Boot)
- Relational databases: schema design, normalization, indexes, transactions, EXPLAIN
- SQL: joins, aggregates, query optimization — not just CRUD
- Authentication basics: sessions vs JWT, password hashing (bcrypt/argon2)
- Testing: unit tests, integration tests, mocking external dependencies

### Resources

- Designing Data-Intensive Applications (Kleppmann) — start it now, read a chapter a week alongside building; it's the single best book for this whole path
- PostgreSQL official tutorial + "Use The Index, Luke" (use-the-index-luke.com)
- Your language's official testing docs

Milestone project: a task/notes API with user auth, PostgreSQL persistence, proper indexing, and a test suite. Deploy it somewhere (Railway, Fly.io, or a $5 VPS) so it's a real running service, not just code on your laptop.

## Phase 3 — Systems & Concurrency (Weeks 11–16)

**Goal: understand what happens under the API layer — this is what separates "can build a CRUD app" from "backend/systems engineer."**

### Topics

- Concurrency models: threads vs processes, async/await, goroutines/channels (if using Go), race conditions, locks
- Networking fundamentals: TCP/IP basics, DNS, load balancers, what actually happens on a request's journey
- Caching: in-memory (Redis), cache invalidation strategies, when caching helps vs hurts
- Message queues: why they exist, at-least-once vs exactly-once delivery (RabbitMQ or a lightweight option)
- Operating systems basics: processes, memory, file descriptors, why your service runs out of connections

### Resources

- Designing Data-Intensive Applications — continue (chapters on replication, partitioning)
- "Computer Networking: A Top-Down Approach" — skim relevant chapters, don't read cover to cover
- Redis docs + build something with it directly

Milestone project: extend your Phase 2 API — add Redis caching for hot reads, add a background job queue for a slow operation (e.g. sending emails, processing uploads), and load-test it (use hey or k6) to see where it breaks. Breaking it on purpose is the point.

## Phase 4 — Distributed Systems & System Design (Weeks 17–21)

**Goal: think at the level system design interviews and real infra decisions require.**

### Topics

- CAP theorem, consistency models, replication strategies
- Horizontal scaling: stateless services, sharding, database read replicas
- Observability: structured logging, metrics, basic tracing — you can't fix what you can't see
- Reading real system design case studies (how Twitter's timeline works, how Stripe handles idempotency, etc.)
- Practicing system design on paper: "design a URL shortener," "design a rate limiter," "design a chat system"

### Resources

- Designing Data-Intensive Applications — finish it
- "System Design Interview" (Alex Xu) Vol. 1 — good for structured practice
- ByteByteGo YouTube/newsletter for quick case-study exposure

Milestone: write up 3 system design docs from scratch (a design doc, not just a diagram) for classic problems, then containerize your Phase 3 service with Docker and add basic health checks + logging.

## Phase 5 — Capstone (Weeks 22–24)

Goal: one polished, portfolio-grade project that demonstrates the whole stack of skills — this is what you show in interviews or to collaborators.

Suggested capstone (pick one, or propose your own):

- A distributed rate limiter or job scheduler
- A simplified version of a real system (a Twitter clone with proper fan-out, a URL shortener with analytics, a chat backend with WebSockets)

Requirements to hit every phase's skills:

- REST or gRPC API, tested
- Postgres with a deliberate, indexed schema
- Redis caching or a queue, used for a real reason (not decoration)
- Dockerized, deployed somewhere real
- Basic logging/metrics
- A written design doc explaining trade-offs you made

Deliverable: README with architecture diagram, a live deployed link if possible, and the design doc. This becomes the centerpiece of a portfolio or GitHub profile.

## Ongoing habits (run throughout, not a separate phase)
- Read one well-known engineering blog post a week (Stripe, Cloudflare, Netflix, Discord engineering blogs are excellent for backend/systems specifically)
- Keep a running "things I learned" log — useful for interviews later, and for noticing your own progress
- Don't rush past Phase 1 if it's genuinely new — everything after depends on it being solid 

## If you want to go further after week 24
- Distributed consensus (Raft) and building a toy implementation
- Kubernetes fundamentals
- A systems language deep-dive (Rust, if you used Go, or vice versa)
- Contributing to an open-source backend project