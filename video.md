# POS AI-First MVP: English Video Script

## Bootcamp Kiro × Código Facilito, Hackathon 2026

**Presenter:** Gabriel Magallon Sanchez

**Target length:** approximately five minutes, including demo pauses. Rehearse and trim to the event limit.
**Current status:** local app verified. Render deployment remains pending. Use a local recording until hosted checks pass.

## Recording plan

| Section | Time | Visual |
|---------|------|--------|
| Problem | 0:00–0:30 | Title and business-owner scenario |
| Solution | 0:30–1:10 | App and language switch |
| Kiro workflow | 1:10–2:15 | Actual project artifacts in the IDE |
| App demo | 2:15–3:40 | Local POS screencast |
| Architecture and safeguards | 3:40–4:20 | Architecture and security slides |
| Results and next steps | 4:20–5:00 | Results, Render preparation and closing |

## 1. Problem (0:00–0:30)

**Visual:** Project title. Cut to a small-business owner reviewing sales records.

**Narration:**
> Hi, I am Gabriel Magallon Sanchez. Imagine Lupita, a taquería owner, finishing a busy day. She has recorded sales, but answering a simple question still takes work: what did I sell today?
>
> POS AI-First combines everyday point-of-sale tasks with a conversational way to explore business data.

**Production note:** Lupita is a fictional user scenario. Do not present her as an interviewed customer.

## 2. Solution (0:30–1:10)

**Visual:** Open the English UI, show the Spanish switch, then return to English.

**Narration:**
> The app starts in English and also supports Spanish. The language switch changes the interface while preserving product names and stored business data.
>
> An owner can manage products and customers, register sales and view dashboard metrics. They can also ask a question such as: which product sold the most this week?
>
> The AI provider proposes a SQL query. The application checks that query before execution. The selected language also guides the provider's response.

**Production note:** Local tests verify language selection and prompt instructions. Verify an actual external AI response before recording it. Do not claim a response-time benchmark.

## 3. Kiro workflow (1:10–2:15)

**Visual:** Show the actual specs, steering files, custom agents, hooks and Power materials that exist in the project.

**Narration:**
> I used Kiro to support a structured development workflow. Specs describe requirements, design decisions and implementation tasks before code changes.
>
> Steering files capture project rules for architecture, testing, security and conventions. They help maintain the same expectations across sessions.
>
> The project also includes specialist agent instructions, workflow hooks and long-term memory materials. These support security checks, realistic demo data and continuity between sessions.
>
> The useful lesson is that the person defines the goal and constraints. AI helps execute within them. I still need to inspect the result, run tests and distinguish verified behavior from assumptions.

**Production note:** Show only configured tools and artifacts you can verify. Do not imply that a historical MCP connection is currently active or that every hook has run successfully. Avoid displaying API keys, PINs or other credentials.

## 4. App demo (2:15–3:40)

**Visual:** Local screencast with prepared, clearly identified demo data. Leave short pauses for each action.

**Narration:**
> Here is the working app. First, I sign in with a PIN. Authentication checks the PIN and handles failed attempts.

*[Show login without exposing a production credential.]*

> The dashboard shows sales metrics and stock information. I can switch to Spanish and back to English. The data stays the same.

*[Show dashboard and both languages.]*

> Next, I open the product list and the customer form. These are part of the same sales workflow.

*[Show a product and a demo customer.]*

> I select a product, add it to the cart and complete a sale. The application records the sale and updates inventory through its business logic.

*[Complete the sale and show the resulting state.]*

> Now I ask: what did I sell today? With a configured provider, the app validates the generated query before reading the data. The answer should match the records from this demo.

*[Show a verified real response. If unavailable, show the localized provider error and state the limitation honestly.]*

**Production note:** Do not narrate example slide amounts as real results. Record the actual amount and sale count. Prepare a recorded backup before the event. A Render URL belongs here only after deployment and hosted verification succeed.

## 5. Architecture and safeguards (3:40–4:20)

**Visual:** Slides 12, 13 and 19.

**Narration:**
> The architecture separates domain rules, application use cases and infrastructure adapters. Repository ports allow different database implementations without tying a sale to a hosting provider.
>
> AI-generated SQL needs more than a prompt. The pipeline combines prompt restrictions with validation, a read-only query connection, bounded execution and audit records.
>
> The repository includes SQLite and PostgreSQL adapters and an AWS integration path. That code is different from proof of a running cloud deployment.

## 6. Results and next steps (4:20–5:00)

**Visual:** Results slide, Render preparation slide, then closing slide.

**Narration:**
> The current result is a working local POS with an English default, Spanish support and an AI-query workflow. Local tests and the build pass. Existing lint findings remain, and live provider behavior still needs verification.
>
> The next target is Render's free plan. Deployment is pending. Before going live, I need to confirm database durability, configure secrets and test login, sales and AI responses on the hosted service.
>
> This MVP shows how a focused business workflow, clear architecture and guided AI development can work together. Thank you.

## Production checklist

- [ ] Prepare local demo data and identify it as illustrative.
- [ ] Verify the live AI response and record its actual result.
- [ ] Show the English default and Spanish switch.
- [ ] Keep secrets and real customer data out of the recording.
- [ ] Show only project artifacts and integrations that are available.
- [ ] Prepare a demo recording as a backup.
- [ ] Rehearse narration and pauses within the event's video limit.
- [ ] Export at 1080p or higher and test playback.
- [ ] Replace local footage with hosted footage only after Render deployment and validation.
