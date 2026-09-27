# POS AI-First: Kiro University Learning Showcase

**Presenter:** Gabriel Magallon Sanchez

An educational showcase of an existing project, not an eligible University final-exam submission. Original project history remains intact. English is the app default, Spanish is available and Render deployment remains pending.

## Official program context

[Kiro University Challenge](https://kiro.dev/2026/university/) describes seven required lessons and two optional bonuses, with some lessons specific to IDE, CLI or Web. The maximum possible award is 5,250 credits, not credits earned here. The final exam is due October 5, 2026 at 23:59 PDT and requires a new project whose first GitHub commit is September 21 or later. This project has earlier local history. Lesson completion and eligibility are not claimed.

Official logo source: [Kiro wordmark](https://kiro.dev/images/kiro-wordmark.png?h=0ad65a93). Brand use identifies the learning tool and does not imply endorsement.

## Slide 1: POS AI-First learning showcase

**On-slide copy:**

- SHOWCASE 2026
- POS AI-First MVP
- Your business answers questions.
- A point of sale where asking your data feels as natural as sending a message.
- código facilito
- ×
- Kiro University learning showcase
- Gabriel Magallon Sanchez

**Speaker notes:** I am Gabriel Magallon Sanchez. This is an educational showcase of an existing POS project and documented learning with Kiro. It is not an eligible final-exam submission. The original bootcamp project predates the University challenge. Official program context: https://kiro.dev/2026/university/. Official logo: https://kiro.dev/images/kiro-wordmark.png?h=0ad65a93.

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

## Slide 4: Kiro University learning framework

**On-slide copy:**

- 03 · UNIVERSITY
- Kiro University learning framework
- Official course structure. This showcase does not claim lesson completion.
- 01
- Study
- Lessons
- 02
- Apply
- POS project
- 03
- Document
- Evidence
- 04
- Verify
- Behavior
- 05
- Share
- Learning
- 7
- Required lessons
- 2
- Optional bonuses
- 3
- IDE / CLI / Web
- 5,250
- Potential credits

**Speaker notes:** The official program includes seven required lessons and two optional bonuses. Some lessons use the IDE, CLI or Web specifically. The maximum possible award is 5,250 credits, not an award this project has earned. The final-exam deadline is October 5, 2026 at 23:59 PDT. No lesson names or completion status are inferred. Source: https://kiro.dev/2026/university/.

---

## Slide 5: Educational showcase scope

**On-slide copy:**

- 04 · SHOWCASE
- An existing project, an educational showcase
- A learning demonstration with original history intact.
- HISTORY
- Existing project
- Local start
- July 21, 2026
- History
- Unchanged
- PROGRAM
- New-build rule
- First commit
- Sept 21 or later
- Context
- Challenge
- SHOWCASE
- POS AI-First
- Learning
- Documented
- Eligibility
- Not claimed
- Existing POS case study
- Educational use only. No final-exam eligibility or credit award claimed.

**Speaker notes:** The local repository begins July 21, 2026. The official exam requires a new project built during the challenge and a first GitHub commit on or after September 21. This showcase does not satisfy or claim that new-project requirement. Creating a repository or rewriting dates would not turn earlier work into a new build. No history changes are part of this presentation update. Source: https://kiro.dev/2026/university/.

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

## Slide 7: Documented Kiro workflow

**On-slide copy:**

- 06 · LEARNING
- Documented Kiro workflow
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
- Project artifacts support the learning story

**Speaker notes:** Specs clarify requirements before implementation. Steering captures recurring project rules. Powers and hooks support context and workflow automation. Dependency planning helps identify work that can proceed independently. These are development aids, not proof that every task executed without conflicts. These are documented project practices, not a mapping that proves all University lessons complete. Program source: https://kiro.dev/2026/university/.

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

## Slide 11: Kiro surfaces and available evidence

**On-slide copy:**

- 10 · SURFACES
- Kiro surfaces and available evidence
- I
- IDE
- Project artifacts
- C
- CLI
- Separate evidence
- W
- WEB
- Separate evidence
- SURFACE USE NEEDS ITS OWN EVIDENCE

**Speaker notes:** The University page distinguishes IDE, CLI and Web lessons. This project contains IDE-oriented artifacts. CLI or Web completion needs actual evidence, not a diagram. No cross-device synchronization or completed University lesson is claimed here. Source: https://kiro.dev/2026/university/.

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

## Slide 17: Lessons from the project

**On-slide copy:**

- 16 · REFLECTION
- Lessons from an existing project
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

**Speaker notes:** Steering and memory help maintain context. Specs clarify intent, and dependency planning supports coordination. Generated SQL needs independent safeguards. Property tests explore more inputs than a few examples, but they do not prove universal correctness. These are project reflections, not official University syllabus titles.

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
- Thank you
- Questions?
- POS AI-First MVP
- Educational showcase 2026

**Speaker notes:** Thank you. This educational showcase preserves the existing project history and makes no final-exam eligibility or earned-credit claim. Program source: https://kiro.dev/2026/university/. Official Kiro logo: https://kiro.dev/images/kiro-wordmark.png?h=0ad65a93.

---
