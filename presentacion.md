# POS AI-First MVP

## Bootcamp Kiro × Código Facilito, Hackathon 2026

**Presenter:** Gabriel Magallon Sanchez

The app defaults to English and supports Spanish through its language switch. Render deployment remains pending. Keep the existing 22-slide design and use the English PowerPoint alongside these notes.

## Slide 1: POS AI-First MVP

**On-slide copy:**

- HACKATHON 2026
- POS AI-First MVP
- Your business answers questions.
- A point of sale where asking your data feels as natural as sending a message.
- código facilito
- ×
- kiro
- Bootcamp Kiro × Código Facilito
- Gabriel Magallon Sanchez

**Speaker notes:** Hello, I am Gabriel Magallon Sanchez. POS AI-First helps a small business owner record sales and ask questions about their data. The interface defaults to English and also supports Spanish.

---

## Slide 2: The business owner needs answers

**On-slide copy:**

- 01 · PROBLEM
- The owner needs answers, not more reports
- “What did I sell today?”
- Without spreadsheet filters or manual totals.
- A direct answer
- Business data in natural language.

**Speaker notes:** Imagine a taquería owner closing the day. They want to understand sales and inventory without searching through spreadsheet rows. This scenario motivates the MVP, rather than a measured market statistic.

---

## Slide 3: Questions in English or Spanish

**On-slide copy:**

- 02 · SOLUTION
- Ask in English or Spanish. Explore your data.
- NL-to-SQL chat for small businesses, with checks before execution.
- 01
- QUESTION
- Which product sold the most?
- 02
- AI
- Generates a SQL query
- 03
- DATA
- SQLite read-only
- 04
- ANSWER
- Selected-language prompt
- “What did I sell today?”
- →
- Example: $12,480 across 37 sales.

**Speaker notes:** The application sends a selected-language instruction to the AI provider, validates the generated query and executes accepted queries against business data. English and Spanish UI rendering and prompt selection pass local tests. A live external-provider response still needs verification. The figures on this slide are illustrative, not measured demo results.

---

## Slide 4: Hackathon MVP scope

**On-slide copy:**

- 03 · CHALLENGE
- A working MVP for the hackathon
- Business value, working software and a documented Kiro workflow.
- 01
- Define
- Specs
- 02
- Build
- Core POS
- 03
- Extend
- AI chat
- 04
- Check
- Tests
- 05
- Prepare
- Deploy
- POS
- Sales
- AI
- Questions
- UI
- EN / ES
- GO
- Backend

**Speaker notes:** The original project plan targeted a five-day hackathon build. This slide presents the scope rather than asserting a measured development duration or unverified judging weights. The deliverable combines point-of-sale operations, AI questions and a bilingual interface. Deployment is still pending.

---

## Slide 5: Project selection criteria

**On-slide copy:**

- 04 · DECISION
- Why this MVP?
- A concrete business workflow with a focused AI use case.
- CRITERION
- Feasibility
- Scope
- Focused
- Stack
- Go + HTMX
- CRITERION
- Business value
- Users
- Small shops
- Workflow
- Sales
- SELECTED
- POS AI-First
- Data
- Sales
- Feature
- AI chat
- Small business + NL-to-SQL
- Project analysis defines the scope and the reason for choosing this MVP.

**Speaker notes:** The previous material describes evaluating several ideas. This slide preserves the decision story without presenting fictional Idea A and Idea B ratings. A manageable scope and a useful question over sales data explain the POS choice.

---

## Slide 6: Research before implementation

**On-slide copy:**

- 05 · DISCOVERY
- Research before implementation
- Understanding the problem helps define a manageable MVP.
- 01
- Review
- Existing POS workflows
- 02
- Analysis
- Potential differentiators
- 03
- Constitution
- Scope and exclusions
- 04
- Questions
- Data model and security
- FOUNDATIONS BEFORE IMPLEMENTATION

**Speaker notes:** The project began with analysis and scope decisions. The useful questions concerned the sales data model, the dashboard and the risk of executing generated SQL. Keep the walkthrough grounded in the analysis files you can actually show.

---

## Slide 7: Kiro development workflow

**On-slide copy:**

- 06 · KIRO
- Kiro development workflow
- S
- Specs
- Requirements, design, tasks
- R
- Steering
- Persistent project rules
- P
- Powers
- Memory and context
- H
- Hooks
- Workflow automation
- Task dependencies guide the implementation order

**Speaker notes:** Specs clarify requirements before implementation. Steering captures recurring project rules. Powers and hooks support context and workflow automation. Dependency planning helps identify work that can proceed independently. These are development aids, not proof that every task executed without conflicts.

---

## Slide 8: Persistent project rules

**On-slide copy:**

- 07 · STEERING
- Consistent project rules across sessions
- architecture.md
- Lightweight hexagonal architecture
- testing.md
- TDD for critical behavior
- security.md
- SQL allowlist, bcrypt, read-only
- quality.md
- Lint and error handling
- project-conventions.md
- Stack, dependencies and commits
- design-patterns.md
- Practical patterns, limited abstraction
- Steering records the rules that guide each development session.

**Speaker notes:** The project steering files cover architecture, testing, security, quality, conventions and design patterns. They guide implementation and reduce repeated explanations. Rules are intentions and still need enforcement. Current lint findings remain, so the deck does not claim a clean lint run.

---

## Slide 9: Long-Term Memory

**On-slide copy:**

- 08 · POWER
- Long-Term Memory across sessions
- TIER 1
- Recent files
- Quick recall
- TIER 2
- Decisions
- Focused search
- TIER 3
- Full detail
- Deeper context
- “Pick up where we left off.”
- Local memory for decisions and unfinished work

**Speaker notes:** The earlier Kiro workflow used local long-term memory to retain decisions and checkpoints. Start with recent context, search for relevant decisions when necessary and retrieve the full detail before resuming. Memory supports continuity but does not replace checking the current code.

---

## Slide 10: Delivery tracking

**On-slide copy:**

- 09 · DELIVERY
- Specs and delivery tracking
- 01
- SPEC
- tasks.md
- 02
- ISSUES
- Work items
- 03
- BOARD
- Kanban V2
- 04
- INSIGHTS
- Prompts and lessons
- TASKS
- Implementation steps
- GOAL
- Shared milestone
- BOARD
- Progress tracking

**Speaker notes:** The project narrative connects spec tasks with issue tracking and a board. Demonstrate the actual artifacts available for the recording rather than claiming a particular issue count or an automatically synchronized remote board. No new remote operations are part of this update.

---

## Slide 11: Cross-device project access

**On-slide copy:**

- 10 · PROJECT ACCESS
- Development and review across devices
- D
- DESKTOP
- Build locally
- W
- WEB
- Read project files
- M
- MOBILE
- Review issues
- GIT AND SHARED PROJECT ARTIFACTS

**Speaker notes:** Local development, browser access to repository artifacts and mobile issue review are different activities. Git and shared documentation support continuity. This slide does not claim an official Kiro mobile application or automatic synchronization of every environment.

---

## Slide 12: Hexagonal architecture and AI layer

**On-slide copy:**

- 11 · ARCHITECTURE
- Hexagonal architecture and AI layer
- Dependencies point toward the domain.
- HTMX + Alpine.js + Tailwind
- Go HTTP and chi router
- Application use cases
- Domain entities and ports
- SQLite, OpenRouter and config
- DOMAIN

**Speaker notes:** The domain defines business entities and ports. Application services coordinate use cases, and infrastructure provides HTTP, persistence and AI adapters. The architecture separates business decisions from a particular database or hosting provider.

---

## Slide 13: NL-to-SQL safeguards

**On-slide copy:**

- 12 · SECURITY
- Five safeguards for AI-generated queries
- Validation and execution controls support defense in depth.
- 1
- PROMPT
- No DDL / DML
- 2
- VALIDATION
- SELECT / WITH
- 3
- CONNECTION
- Read-only
- 4
- EXECUTION
- Timeout and row cap
- 5
- AUDIT
- Query records
- GENERATED SQL STILL NEEDS VALIDATION

**Speaker notes:** The pipeline uses prompt restrictions, a Go validator, a read-only query connection, bounded execution and query audit records. Prompt instructions alone cannot guarantee safety. These controls reduce risk, and the application still needs operational monitoring and verification in its deployed environment.

---

## Slide 14: Customer management and admin settings

**On-slide copy:**

- 13 · ITERATION
- Customer management and protected admin settings
- Customers
- CRUD and validation
- Admin
- Encrypted API key
- Navigation
- HTMX no-cache
- Roles
- Role-aware sidebar
- ADMIN
- FORM
- AES-GCM
- SQLITE
- MASK ••••1234
- SHA-256 of SESSION_SECRET, decrypt on read
- ENGLISH DEFAULT · SPANISH AVAILABLE

**Speaker notes:** Customer forms validate inputs. Admin settings restrict API-key configuration by role and protect stored values with AES-GCM. The key derivation uses SESSION_SECRET, so a strong secret remains essential. The bilingual switch changes display text without translating stored business data or payment identifiers.

---

## Slide 15: Point-of-sale demo

**On-slide copy:**

- 14 · DEMO
- Questions for your POS
- A working application flow.
- 01
- PIN login
- 02
- Create a product
- 03
- Register a customer
- 04
- Record a sale
- 05
- What did I sell today?
- 06
- Dashboard + Admin
- DEMO
- What did I sell today?
- $12,480
- 37 example sales
- Prepare a recorded backup

**Speaker notes:** Start in English, briefly switch to Spanish and return to English. Show products, a customer, a sale and dashboard metrics. Ask a question using the available AI configuration. Prepare a recorded backup before the event. The slide contains example figures and does not claim the recording already exists.

---

## Slide 16: Current implementation results

**On-slide copy:**

- 15 · RESULTS
- A working MVP with local verification
- EN / ES
- UI languages
- GO
- Backend
- 5
- SQL safeguards
- CRUD
- POS workflows
- AES
- Key protection
- PASS
- Local tests
- 7
- PG repositories
- NEXT
- Render deploy
- AI chat, product and customer management, sales, encrypted settings

**Speaker notes:** Local app tests and the build pass. English is the default UI language and Spanish remains available. PostgreSQL repositories and AWS integration code also exist. The lint run still reports twelve existing findings. External AI output, hosted behavior and deployment are not verified by those local checks.

---

## Slide 17: Lessons learned

**On-slide copy:**

- 16 · LESSONS
- Strong foundations support reliable progress
- Steering
- Consistency
- LTM
- Continuity
- Specs
- Scope before code
- Waves
- Dependency planning
- NL→SQL
- Defense in depth
- HTMX
- Simple frontend
- AES-GCM
- Protected secrets
- Property tests
- Broader input coverage
- People direct the work. AI supports execution.

**Speaker notes:** Steering and memory help maintain context. Specs clarify intent, and dependency planning supports coordination. Generated SQL needs independent safeguards. Property tests explore more inputs than a few examples, but they do not prove universal correctness.

---

## Slide 18: Existing AWS integration code

**On-slide copy:**

- 17 · AWS OPTION
- Existing AWS integration code
- BROWSER
- S3 + CloudFront
- API
- API Gateway
- COMPUTE
- Lambda Go ARM64
- DATA
- RDS PostgreSQL
- AI
- Bedrock Haiku
- SECRETS
- Secrets Manager
- switch cfg.AppEnv {
- case "lambda": PostgreSQL + Bedrock
- default: SQLite + OpenRouter
- ALTERNATIVE HOSTING PATH · DEPLOYMENT NOT VERIFIED

**Speaker notes:** The repository includes an AWS deployment path using Lambda, PostgreSQL and Bedrock, with infrastructure definitions. The diagram describes that implementation path, not a verified running production environment. The current next deployment target is Render, not an AWS rollout.

---

## Slide 19: Interchangeable infrastructure adapters

**On-slide copy:**

- 18 · HEXAGONAL
- Interchangeable infrastructure adapters
- Shared ports, different implementations.
- BUSINESS LAYERS
- Domain and application
- src/domain/
- src/application/
- templates/
- ADAPTERS
- Infrastructure implementations
- Product
- Sale
- User
- Client
- Inventory
- Config
- Metrics
- Bedrock AI
- RegisterSale works through repository ports, independent of the database.

**Speaker notes:** Repository ports allow SQLite and PostgreSQL implementations to serve the same business use cases. The earlier zero-lines-changed migration claim needs historical diff evidence, so this slide explains the separation without asserting that application or templates never changed. The current language work explicitly changes both.

---

## Slide 20: Render deployment preparation

**On-slide copy:**

- 19 · DEPLOYMENT
- Render deployment preparation
- Free-plan deployment is the target. Setup and validation remain pending.
- NEXT
- Render setup
- DATA
- Storage decision
- CHECK
- Hosted behavior
- Runtime
- Go service
- Build and start
- Database
- Persistence plan
- Confirm durability
- AI provider
- API credentials
- Validate response
- Configuration
- Strong secrets
- Set environment
- Verification
- Login and sales
- Health check

**Speaker notes:** Render free-plan deployment is a requested next step, not a completed result. Confirm the supported runtime, database persistence approach, environment variables and credentials before deployment. Validate login, language switching, a sale, health and AI responses on the hosted service. This deck makes no pricing, persistent-disk or deployment-duration promise.

---

## Slide 21: Next steps

**On-slide copy:**

- 20 · ROADMAP
- Next steps after the local MVP
- NOW
- LOCAL MVP
- EN / ES
- NEXT
- RENDER SETUP
- Persistence plan
- THEN
- HOSTED DEMO
- Verify live flows
- LATER
- AI FEATURES
- Forecasting and branches
- VERIFIED STEPS BEFORE A LIVE DEMO

**Speaker notes:** Complete Render preparation, resolve durable storage and configure secrets. Then deploy within an explicitly authorized remote scope and verify the hosted demo. Conversation history, multiple branches, offline synchronization and forecasting remain future ideas, not implemented features.

---

## Slide 22: Thank you

**On-slide copy:**

- código facilito
- ×
- kiro
- Thank you
- Questions?
- POS AI-First MVP
- Hackathon 2026

**Speaker notes:** Thank you. I am Gabriel Magallon Sanchez. Questions can cover the POS workflow, bilingual UI, SQL safeguards, adapter boundaries and the remaining Render deployment work.

---
