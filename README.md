# PsiEncontra

A web platform that helps Psychology students discover which theoretical approach and field of practice best match their profile. Through an interactive questionnaire analyzed by artificial intelligence, the student receives a personalized affinity ranking with charts, detailed explanations and PDF export.

> Live app: [psiencontra.vercel.app](https://psiencontra.vercel.app)

---

## The Problem

Psychology students, especially in their first semesters, are confronted with a wide variety of theoretical approaches (Psychoanalysis, CBT, Gestalt, Phenomenology, among others) and fields of practice (Clinical, Organizational, School, Forensic, etc.). This diversity, while rich, can generate doubts and insecurity about which path to follow in their career. There is a lack of practical tools that help students reflect on their own inclinations in a guided and personalized way.

## The Solution

PsiEncontra solves this problem by offering a 15-question questionnaire — combining multiple-choice and open-ended questions — that explores the student's view on topics such as human suffering, therapeutic methods, contexts of practice and the role of the psychologist. The answers are sent to an AI (Google Gemini or Groq Llama) that analyzes the student's vocabulary, references and way of thinking, generating:

- An **affinity ranking** with 8 theoretical approaches (from 0 to 100)
- An **affinity ranking** with 9 fields of practice (from 0 to 100)
- **Radar charts** for intuitive visualization
- **Detailed explanations** for each score
- A **general summary** of the student's profile
- **PDF export** of the complete result

The result is informative and does not replace professional guidance, but works as a self-knowledge and reflection tool for the student.

## Demo

<!-- Add screenshots or GIFs of the application here -->

| Screen | Description |
|---|---|
| ![Landing Page](./docs/landing.png) | Landing page with introduction and CTA |
| ![Questionnaire](./docs/questionario.png) | Multiple-choice question with answer options |
| ![Result](./docs/resultado.png) | Radar charts and affinity ranking |

> To add the images, take screenshots of the application and save them in the `docs/` folder of the project.

## Features

- Interactive questionnaire with 15 questions divided into thematic blocks
- Multiple-choice and open-ended questions
- AI analysis with automatic fallback (Gemini -> Groq)
- Affinity ranking with 8 theoretical approaches and 9 fields of practice
- Interactive radar charts for visualizing the results
- Personalized explanations for each approach and field
- Full result export to PDF
- Dark mode with persistence via localStorage
- Responsive design for desktop and mobile
- Smooth animations with Framer Motion

## Theoretical Approaches Evaluated

| Approach | Main Authors |
|---|---|
| Psychoanalysis | Freud, Lacan, Winnicott |
| Existential-Phenomenology | Husserl, Heidegger, Rogers |
| Behavior Analysis | Skinner |
| Cognitive Behavioral Therapy | Beck, Ellis |
| Analytical Psychology | Jung |
| Gestalt Therapy | Perls |
| Socio-Historical Psychology | Vygotsky |
| Systemic | Bateson, Minuchin |

## Fields of Practice Evaluated

Clinical Psychology, Organizational, School/Educational, Social and Community, Health/Hospital, Forensic, Sports, Neuropsychology and Psychometrics.

## Tech Stack

### Frontend
- **Next.js 16** — React framework with App Router and server-side rendering
- **React 19** — Library for building user interfaces
- **Tailwind CSS 4** — Utility-first styling with dark mode support
- **Framer Motion** — Smooth animations and transitions
- **Recharts** — Interactive radar charts for visualizing results

### Backend
- **Go** — Server language, chosen for its performance and simplicity
- **Gin** — Lightweight and fast HTTP framework
- **GORM** — ORM for communication with PostgreSQL
- **gofpdf** — PDF document generation with UTF-8 support
- **godotenv** — Loading of environment variables

### AI
- **Google Gemini 2.0 Flash** — Primary analysis provider
- **Groq Llama 3.3 70B** — Fallback provider

### Infrastructure
- **PostgreSQL** — Relational database
- **Docker** — Containerization for local development
- **Vercel** — Frontend deployment
- **Railway** — API and database deployment

## Project Architecture

The project follows a monorepo architecture with a clear separation between frontend and backend:

```
psiencontra/
├── api/                        # Go backend
│   ├── config/                 # Configuration (database, env, logger)
│   ├── handler/                # HTTP controllers (routes)
│   ├── repository/             # Data access layer (queries)
│   ├── router/                 # Route definitions and CORS
│   ├── schemas/                # Data models (structs)
│   ├── service/                # Business logic (AI, PDF, questions)
│   ├── Dockerfile              # Docker image build
│   ├── main.go                 # API entry point
│   └── go.mod                  # Go dependencies
├── web/                        # Next.js frontend
│   ├── app/                    # Pages (landing, questionnaire, result)
│   ├── components/             # Reusable components (Button, Card, etc.)
│   ├── lib/                    # API client and constants
│   └── package.json            # Node.js dependencies
├── docker-compose.yml          # Local orchestration (PostgreSQL + API)
├── .env.example                # Environment variables template
└── README.md
```

### Application Flow

```
Student visits the site
        |
Clicks "Start Questionnaire"
        |
Answers 15 questions (multiple-choice + open-ended)
        |
Frontend sends answers to the API
        |
API builds the prompt and sends it to the AI (Gemini or Groq)
        |
AI returns JSON with scores and descriptions
        |
API saves the result in the database and returns it to the frontend
        |
Frontend displays radar charts, ranking and explanations
        |
Student can export the result as a PDF
```

## How to Run the Project

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker](https://www.docker.com/) (for PostgreSQL)
- An API key from [Google Gemini](https://aistudio.google.com/apikey) and/or [Groq](https://console.groq.com/)

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

The API will be available at `http://localhost:8080`.

### 5. Start the frontend

In another terminal:

```bash
cd web
npm install
npm run dev
```

The frontend will be available at `http://localhost:3000`.

## Production Deployment

| Service | Platform | URL |
|---|---|---|
| Frontend | Vercel | [psiencontra.vercel.app](https://psiencontra.vercel.app) |
| API + Database | Railway | psiencontra-production.up.railway.app |

### Environment variables

**Vercel (frontend):**
- `NEXT_PUBLIC_API_URL` — API URL + `/api/v1`

**Railway (API):**
- `DATABASE_URL` — provided automatically by Railway's PostgreSQL
- `GEMINI_API_KEY` — Google Gemini API key
- `GROQ_API_KEY` — Groq API key
- `FRONTEND_URL` — frontend URL on Vercel
- `PORT` — 8080

## Lessons Learned

While building this project, the main takeaways were:

- **Integration with generative AIs**: how to craft efficient prompts to obtain structured JSON responses, deal with token limits and implement fallback between different providers
- **Fullstack architecture**: separation of responsibilities between frontend (React/Next.js) and backend (Go/Gin), communication via REST API and state management
- **Deploying distributed applications**: configuring and connecting services across different platforms (Vercel + Railway), managing environment variables and handling CORS in production
- **UX and accessibility**: dark mode, smooth animations, responsive design and visual feedback during loading

## License

This project is for academic and educational use. The result is informative and does not replace professional guidance.
