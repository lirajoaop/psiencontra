# PsiEncontra

A web platform that helps Psychology students discover which theoretical approach and field of practice best match their profile. The student answers a structured questionnaire and receives a personalized affinity ranking with charts, detailed explanations and PDF export.

> Live app: [psiencontra.vercel.app](https://psiencontra.vercel.app)

---

## The Problem

Psychology students, especially in their first semesters, are confronted with a wide variety of theoretical approaches (Psychoanalysis, CBT, Gestalt, Phenomenology, among others) and fields of practice (Clinical, Organizational, School, Forensic, etc.). This diversity, while rich, can generate doubts and insecurity about which path to follow in their career. There is a lack of practical tools that help students reflect on their own inclinations in a guided and personalized way.

## The Solution

PsiEncontra offers two questionnaire modes that explore the student's view on topics such as human suffering, therapeutic methods, contexts of practice and the role of the psychologist:

- **Simple mode** — 15 questions combining multiple-choice and open-ended items. Scores are produced by an AI (Google Gemini, with Groq Llama as fallback) that interprets vocabulary, references and reasoning.
- **Detailed mode** — ~76 Likert items (1–5) distributed across 8 approaches (5 items each, 4 pro-construct + 1 reverse-coded) and 9 fields (4 items each). Scores are computed by a **deterministic ipsative algorithm** on the backend; the AI is used only to write personalized descriptions for each result.

Both modes generate:

- An **affinity ranking** with 8 theoretical approaches (0–100)
- An **affinity ranking** with 9 fields of practice (0–100)
- **Radar charts** for intuitive visualization
- **Personalized explanations** for each score
- A **general summary** of the student's profile
- **PDF export** of the complete result

The result is informative and does not replace professional guidance, but works as a self-knowledge and reflection tool for the student.

## Demo

| Screen | Description |
|---|---|
| ![Landing Page](./docs/landing.png) | Landing page with introduction and CTA |
| ![Questionnaire](./docs/questionnaire.png) | Question with answer options |
| ![Result](./docs/result.png) | Radar charts and affinity ranking |

## Features

- Two questionnaire modes (simple AI-scored and detailed Likert with deterministic scoring)
- Items shuffled at runtime to reduce halo and carry-over effects
- Reverse-coded items (1 per approach) to reduce construct transparency and flag inattentive responders
- AI analysis with automatic fallback (Gemini → Groq)
- Deterministic ipsative scoring (reproducible, free of acquiescence bias) for the detailed mode
- Affinity ranking with 8 theoretical approaches and 9 fields of practice
- Interactive radar charts for visualizing the results
- Full result export to PDF
- **Authentication**: email/password and Google OAuth (JWT delivered via cookie or `Authorization: Bearer` header)
- **Anonymous flow**: sessions can be created without an account
- **Questionnaire history**: authenticated users can view all past sessions with top approach/field badges and compare results over time
- **Rate limiting**: token-bucket per-IP rate limiting on all API routes (stricter limits on auth and response-submission endpoints)
- **Input validation**: 500-character limit on open-ended answers (enforced on both frontend and backend)
- **Keyboard navigation**: Likert scale questions support 1–5 key selection and Enter to advance
- **Contextual tooltips**: Material-UI tooltips with glossary hints for approaches and fields
- Dark mode with persistence via localStorage
- Responsive design for desktop and mobile
- Smooth animations with Framer Motion

## Theoretical Approaches Evaluated

| Approach | Main Authors |
|---|---|
| Psychoanalysis | Freud, Lacan, Winnicott |
| Existential-Phenomenology | Husserl, Heidegger, Sartre |
| Behavior Analysis | Skinner |
| Cognitive Behavioral Therapy | Beck, Ellis |
| Analytical Psychology | Jung |
| Gestalt Therapy | Perls |
| Socio-Historical Psychology | Vygotsky |
| Humanism (Person-Centered Approach) | Rogers, Maslow |

## Fields of Practice Evaluated

Clinical Psychology, Organizational, School/Educational, Social and Community, Health/Hospital, Forensic, Sports, Neuropsychology and Psychometrics.

## Scoring Methodology (Detailed Mode)

The detailed mode uses **ipsative normalization**, computing each respondent's deviation from their own personal mean instead of comparing absolute ratings between people. This cancels the acquiescence bias: two students with different response styles but the same relative preference produce the same profile.

Algorithm per block (approaches / fields):
1. Group ratings by mapping (e.g. `psicanalise`, `clinica`).
2. Compute mean rating per mapping.
3. Compute the grand mean across all ratings in the block.
4. Compute each mapping's deviation from the grand mean.
5. Linearly normalize deviations into `[5, 95]` and round to integer.
6. If all deviations are equal (flat profile), every mapping receives `50`.

References: Holland (1997), Harmon et al. (1994), Cattell & Brennan (1994).

See [`api/service/scoring.go`](api/service/scoring.go) for the implementation.

## Tech Stack

### Frontend
- **Next.js 16** — React framework with App Router
- **React 19** — UI library
- **Tailwind CSS 4** — utility-first styling with dark mode
- **Material-UI (MUI) 9** — Tooltip, ClickAwayListener
- **Emotion** — CSS-in-JS engine (required by MUI)
- **Framer Motion** — animations and transitions
- **Recharts** — interactive radar charts

### Backend
- **Go** — server language
- **Gin** — HTTP framework
- **GORM** — ORM for PostgreSQL
- **JWT** (`golang-jwt`) — auth tokens
- **bcrypt** — password hashing
- **Google OAuth 2.0** — third-party login
- **gofpdf** — PDF generation with UTF-8 support
- **godotenv** — environment variables

### AI
- **Google Gemini 2.0 Flash** — primary analysis provider
- **Groq Llama 3.3 70B** — fallback provider

### Infrastructure
- **PostgreSQL** — relational database (users, sessions, responses, results)
- **Docker** — local development
- **Vercel** — frontend deployment
- **Railway** — API and database deployment

## Project Architecture

```
psiencontra/
├── api/                        # Go backend
│   ├── config/                 # Configuration (database, env, logger)
│   ├── handler/                # HTTP controllers (auth, sessions, questions, pdf, history)
│   ├── repository/             # Data access (users, sessions, responses, results)
│   ├── router/                 # Route definitions, CORS, auth and rate-limit middleware
│   ├── schemas/                # GORM models (User, Session, Response, Result)
│   ├── service/                # Business logic (AI, PDF, questions, scoring, auth, OAuth)
│   ├── Dockerfile
│   ├── main.go
│   └── go.mod
├── web/                        # Next.js frontend
│   ├── app/                    # App Router pages (auth, questionnaire, results, history, privacy)
│   ├── components/             # UI (LikertScale, RadarChartResult, Tooltip, AuthProvider, ...)
│   ├── lib/                    # API client, auth-token, constants
│   └── package.json
├── docker-compose.yml          # Local orchestration (PostgreSQL + API)
├── .env.example
└── README.md
```

### API Endpoints

```
GET    /api/v1/health
GET    /api/v1/questions?type=simple|detailed

POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/auth/me
GET    /api/v1/auth/google
GET    /api/v1/auth/google/callback

POST   /api/v1/sessions                    # optional auth (anonymous-friendly)
POST   /api/v1/sessions/:id/responses
GET    /api/v1/sessions/:id/result
GET    /api/v1/sessions/:id/pdf

GET    /api/v1/user/sessions               # requires auth — questionnaire history
```

#### Rate Limiting

All routes are protected by per-IP token-bucket rate limiting:

| Scope | Limit | Burst |
|---|---|---|
| Global (default) | 60 req/min | 10 |
| Auth endpoints | 10 req/min | 5 |
| Response submission | 5 req/min | 2 |

### Application Flow

```
Student visits the site
        |
(Optional) signs in with email/password or Google
        |
Chooses simple or detailed questionnaire
        |
Answers the questions (multiple-choice, open-ended or Likert)
        |
Frontend submits answers to the API
        |
   ┌────────────────────┴────────────────────┐
   |                                          |
Simple mode:                              Detailed mode:
AI returns scores + descriptions          Backend computes ipsative scores;
                                          AI only writes descriptions
   └────────────────────┬────────────────────┘
                        |
API persists the session/result and returns it to the frontend
                        |
Radar charts, ranking and explanations are displayed
                        |
Student can export the full result as a PDF
        |
(If authenticated) can view history at /historico and compare past results
```

## How to Run the Project

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker](https://www.docker.com/) (for PostgreSQL)
- An API key from [Google Gemini](https://aistudio.google.com/apikey) and/or [Groq](https://console.groq.com/)
- (Optional) Google OAuth credentials for "Sign in with Google"

### 1. Clone the repository

```bash
git clone https://github.com/lirajoaop/psiencontra.git
cd psiencontra
```

### 2. Configure environment variables

```bash
cp .env.example .env
```

Edit `.env` with your keys:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/psiencontra?sslmode=disable

GEMINI_API_KEY=your_gemini_key
GROQ_API_KEY=your_groq_key

PORT=8080
FRONTEND_URL=http://localhost:3000

# Local HTTP dev: unlocks an insecure JWT placeholder and disables Secure cookies
APP_ENV=development
JWT_SECRET=                       # required in production

# Google OAuth (optional — leave empty to disable Google login)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback

# Cookie behavior (defaults derived from APP_ENV)
# COOKIE_SECURE=true
# COOKIE_DOMAIN=
```

### 3. Start the database

```bash
docker compose up -d postgres
```

### 4. Start the API

```bash
cd api
go run .
```

Available at `http://localhost:8080`.

### 5. Start the frontend

```bash
cd web
npm install
npm run dev
```

Available at `http://localhost:3000`.

## Production Deployment

| Service | Platform | URL |
|---|---|---|
| Frontend | Vercel | [psiencontra.vercel.app](https://psiencontra.vercel.app) |
| API + Database | Railway | psiencontra-production.up.railway.app |

### Environment variables

**Vercel (frontend):**
- `NEXT_PUBLIC_API_URL` — API URL + `/api/v1`

**Railway (API):**
- `DATABASE_URL` — provided automatically by Railway PostgreSQL
- `GEMINI_API_KEY`
- `GROQ_API_KEY`
- `FRONTEND_URL`
- `PORT` — `8080`
- `APP_ENV` — leave unset (or set to `production`) so Secure/SameSite=None cookies are emitted
- `JWT_SECRET` — strong random secret (required)
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL` — required for Google login

### Auth & cookies in production

The API issues JWTs both as an `HttpOnly` cookie (`psiencontra_auth`) and in the JSON response body. Browsers that allow cross-site cookies use the cookie path; Safari and incognito modes fall back to `Authorization: Bearer` from `localStorage`. The OAuth callback returns the token via the URL fragment (`#token=...`) so the JWT never appears in any access log along the redirect chain.

## Lessons Learned

- **Generative AIs**: crafting prompts for structured JSON, dealing with token limits, and falling back across providers
- **Psychometrics in code**: separating deterministic scoring (ipsative normalization, reproducible, audit-friendly) from AI-generated descriptions, so numeric results never depend on model variability
- **Authentication across origins**: combining cookie-based and Bearer-token flows to survive Safari/WebKit cross-site cookie restrictions; passing the OAuth token via URL fragment to keep it out of server logs
- **Rate limiting**: implementing per-IP token-bucket rate limiting with automatic stale-entry cleanup to protect AI-backed endpoints from abuse
- **Fullstack architecture**: clear separation between Next.js (App Router) and Go/Gin, communication via REST, GORM models and migrations
- **Distributed deployment**: connecting Vercel and Railway, managing environment variables, CORS and Secure/SameSite cookies in production

## Disclaimer & Privacy

PsiEncontra is a **self-reflection tool for Psychology students**, not a psychological test under [CFP Resolution nº 09/2018](https://satepsi.cfp.org.br/). Results are informative and do not replace professional guidance, academic supervision, or psychotherapy. A plain-language privacy notice is served at [`/privacidade`](https://psiencontra.vercel.app/privacidade); the submit step in the questionnaire asks for explicit consent before responses are persisted.

**Known gaps (portfolio-stage, not production-grade):**
- No LGPD data-processing agreement with third-party AI providers; no automated data-erasure endpoint (requests handled manually via the e-mail listed in the privacy page).
- No formal validation study for the detailed questionnaire (no published α, test-retest, or factor analysis on a Brazilian sample). Scores are face-valid but not psychometrically validated.

## License

This project is for academic and educational use. The result is informative and does not replace professional guidance.
