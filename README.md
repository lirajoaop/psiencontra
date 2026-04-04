# PsiEncontra

Plataforma web que ajuda estudantes de Psicologia a descobrirem sua afinidade com diferentes abordagens teoricas e campos de atuacao. O estudante responde 15 perguntas (objetivas e dissertativas) e uma IA analisa as respostas, gerando um ranking personalizado com graficos e explicacoes detalhadas.

## Funcionalidades

- **Questionario interativo** com 15 perguntas divididas em blocos tematicos (objetivas e dissertativas)
- **Analise por IA** usando Gemini 2.0 Flash (primario) e Groq Llama 3.3 70B (fallback)
- **Ranking de afinidade** com 8 abordagens teoricas e 9 campos de atuacao
- **Graficos radar** para visualizacao intuitiva dos resultados
- **Exportacao em PDF** do resultado completo
- **Dark mode** com persistencia via localStorage
- **Design responsivo** para desktop e mobile

## Abordagens Teoricas Avaliadas

| Abordagem | Referencia |
|---|---|
| Psicanalise | Freud, Lacan, Winnicott |
| Fenomenologia-Existencial | Husserl, Heidegger, Rogers |
| Analise do Comportamento | Skinner |
| Terapia Cognitivo-Comportamental | Beck, Ellis |
| Psicologia Analitica | Jung |
| Gestalt-terapia | Perls |
| Psicologia Socio-Historica | Vigotski |
| Sistemica | Bateson, Minuchin |

## Campos de Atuacao Avaliados

Psicologia Clinica, Organizacional, Escolar/Educacional, Social e Comunitaria, da Saude/Hospitalar, Juridica, do Esporte, Neuropsicologia e Psicometria.

## Tecnologias

### Frontend
- **Next.js 16** (App Router)
- **React 19**
- **Tailwind CSS 4**
- **Framer Motion** (animacoes)
- **Recharts** (graficos radar)

### Backend
- **Go** (Gin framework)
- **GORM** (ORM para PostgreSQL)
- **gofpdf** (geracao de PDF)
- **godotenv** (variaveis de ambiente)

### Infraestrutura
- **PostgreSQL** (banco de dados)
- **Docker** (desenvolvimento local)
- **Vercel** (deploy do frontend)
- **Railway** (deploy da API + banco)

## Arquitetura

```
psiencontra/
├── api/                    # Backend Go
│   ├── config/             # Configuracao (banco, env, logger)
│   ├── handler/            # Controllers HTTP
│   ├── repository/         # Camada de acesso a dados
│   ├── router/             # Rotas e CORS
│   ├── schemas/            # Modelos de dados
│   ├── service/            # Logica de negocio (IA, PDF, perguntas)
│   ├── Dockerfile
│   ├── main.go
│   └── go.mod
├── web/                    # Frontend Next.js
│   ├── app/                # Paginas (landing, questionario, resultado)
│   ├── components/         # Componentes reutilizaveis
│   ├── lib/                # API client e constantes
│   └── package.json
├── docker-compose.yml
└── .env.example
```

## API Endpoints

| Metodo | Rota | Descricao |
|---|---|---|
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/questions` | Lista todas as perguntas |
| POST | `/api/v1/sessions` | Cria uma nova sessao |
| POST | `/api/v1/sessions/:id/responses` | Envia respostas do questionario |
| GET | `/api/v1/sessions/:id/result` | Retorna o resultado da analise |
| GET | `/api/v1/sessions/:id/pdf` | Download do resultado em PDF |

## Como Rodar Localmente

### Pre-requisitos

- [Go 1.25+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker](https://www.docker.com/) (para o PostgreSQL)

### 1. Clone o repositorio

```bash
git clone https://github.com/lirajoaop/psiencontra.git
cd psiencontra
```

### 2. Configure as variaveis de ambiente

```bash
cp .env.example .env
```

Edite o `.env` com suas chaves de API:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/psiencontra?sslmode=disable
GEMINI_API_KEY=sua_chave_gemini
GROQ_API_KEY=sua_chave_groq
PORT=8080
FRONTEND_URL=http://localhost:3000
```

### 3. Suba o banco de dados

```bash
docker compose up -d postgres
```

### 4. Inicie a API

```bash
cd api
go run .
```

A API estara disponivel em `http://localhost:8080`.

### 5. Inicie o frontend

```bash
cd web
npm install
npm run dev
```

O frontend estara disponivel em `http://localhost:3000`.

## Deploy

| Servico | Plataforma | URL |
|---|---|---|
| Frontend | Vercel | [psiencontra.vercel.app](https://psiencontra.vercel.app) |
| API + Banco | Railway | psiencontra-production.up.railway.app |

### Variaveis de ambiente em producao

**Vercel (frontend):**
- `NEXT_PUBLIC_API_URL` = URL da API no Railway + `/api/v1`

**Railway (API):**
- `DATABASE_URL` = fornecida automaticamente pelo PostgreSQL do Railway
- `GEMINI_API_KEY` = chave da API do Google Gemini
- `GROQ_API_KEY` = chave da API do Groq
- `FRONTEND_URL` = URL do frontend na Vercel
- `PORT` = 8080

## Licenca

Este projeto e de uso academico e educacional.
