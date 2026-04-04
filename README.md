# PsiEncontra

Plataforma web que ajuda estudantes de Psicologia a descobrirem qual abordagem teorica e campo de atuacao mais combinam com seu perfil. Atraves de um questionario interativo com analise por inteligencia artificial, o estudante recebe um ranking personalizado de afinidade com graficos, explicacoes detalhadas e exportacao em PDF.

> Acesse a aplicacao: [psiencontra.vercel.app](https://psiencontra.vercel.app)

---

## O Problema

Estudantes de Psicologia, especialmente nos primeiros semestres, se deparam com uma grande variedade de abordagens teoricas (Psicanalise, TCC, Gestalt, Fenomenologia, entre outras) e campos de atuacao (Clinica, Organizacional, Escolar, Juridica, etc.). Essa diversidade, embora rica, pode gerar duvidas e inseguranca sobre qual caminho seguir na carreira. Faltam ferramentas praticas que ajudem o estudante a refletir sobre suas proprias inclinacoes de forma guiada e personalizada.

## A Solucao

O PsiEncontra resolve esse problema oferecendo um questionario com 15 perguntas — entre objetivas e dissertativas — que exploram a visao do estudante sobre temas como sofrimento humano, metodos terapeuticos, contextos de atuacao e papel do psicologo. As respostas sao enviadas para uma IA (Google Gemini ou Groq Llama) que analisa o vocabulario, as referencias e a forma de pensar do estudante, gerando:

- Um **ranking de afinidade** com 8 abordagens teoricas (de 0 a 100)
- Um **ranking de afinidade** com 9 campos de atuacao (de 0 a 100)
- **Graficos radar** para visualizacao intuitiva
- **Explicacoes detalhadas** sobre cada pontuacao
- Um **resumo geral** do perfil do estudante
- **Exportacao em PDF** do resultado completo

O resultado e orientativo e nao substitui acompanhamento profissional, mas funciona como uma ferramenta de autoconhecimento e reflexao para o estudante.

## Demonstracao

<!-- Adicione prints ou GIFs da aplicacao aqui -->

| Tela | Descricao |
|---|---|
| ![Landing Page](./docs/landing.png) | Pagina inicial com apresentacao e CTA |
| ![Questionario](./docs/questionario.png) | Pergunta objetiva com opcoes de resposta |
| ![Resultado](./docs/resultado.png) | Graficos radar e ranking de afinidade |

> Para adicionar as imagens, tire prints da aplicacao e salve na pasta `docs/` do projeto.

## Funcionalidades

- Questionario interativo com 15 perguntas divididas em blocos tematicos
- Perguntas objetivas (multipla escolha) e dissertativas (texto livre)
- Analise por IA com fallback automatico (Gemini -> Groq)
- Ranking de afinidade com 8 abordagens teoricas e 9 campos de atuacao
- Graficos radar interativos para visualizacao dos resultados
- Explicacoes personalizadas para cada abordagem e campo
- Exportacao do resultado completo em PDF
- Dark mode com persistencia via localStorage
- Design responsivo para desktop e mobile
- Animacoes suaves com Framer Motion

## Abordagens Teoricas Avaliadas

| Abordagem | Principais Autores |
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

## Tecnologias Utilizadas

### Frontend
- **Next.js 16** — Framework React com App Router e renderizacao server-side
- **React 19** — Biblioteca para construcao de interfaces
- **Tailwind CSS 4** — Estilizacao utilitaria com suporte a dark mode
- **Framer Motion** — Animacoes e transicoes fluidas
- **Recharts** — Graficos radar interativos para visualizacao dos resultados

### Backend
- **Go** — Linguagem do servidor, escolhida pela performance e simplicidade
- **Gin** — Framework HTTP leve e rapido
- **GORM** — ORM para comunicacao com PostgreSQL
- **gofpdf** — Geracao de documentos PDF com suporte a UTF-8
- **godotenv** — Carregamento de variaveis de ambiente

### IA
- **Google Gemini 2.0 Flash** — Provedor primario de analise
- **Groq Llama 3.3 70B** — Provedor de fallback

### Infraestrutura
- **PostgreSQL** — Banco de dados relacional
- **Docker** — Containerizacao para desenvolvimento local
- **Vercel** — Deploy do frontend
- **Railway** — Deploy da API e banco de dados

## Arquitetura do Projeto

O projeto segue uma arquitetura monorepo com separacao clara entre frontend e backend:

```
psiencontra/
├── api/                        # Backend Go
│   ├── config/                 # Configuracao (banco, env, logger)
│   ├── handler/                # Controllers HTTP (rotas)
│   ├── repository/             # Camada de acesso a dados (queries)
│   ├── router/                 # Definicao de rotas e CORS
│   ├── schemas/                # Modelos de dados (structs)
│   ├── service/                # Logica de negocio (IA, PDF, perguntas)
│   ├── Dockerfile              # Build da imagem Docker
│   ├── main.go                 # Ponto de entrada da API
│   └── go.mod                  # Dependencias Go
├── web/                        # Frontend Next.js
│   ├── app/                    # Paginas (landing, questionario, resultado)
│   ├── components/             # Componentes reutilizaveis (Button, Card, etc.)
│   ├── lib/                    # API client e constantes
│   └── package.json            # Dependencias Node.js
├── docker-compose.yml          # Orquestracao local (PostgreSQL + API)
├── .env.example                # Modelo de variaveis de ambiente
└── README.md
```

### Fluxo da Aplicacao

```
Estudante acessa o site
        |
Clica em "Comecar Questionario"
        |
Responde 15 perguntas (objetivas + dissertativas)
        |
Frontend envia respostas para a API
        |
API monta o prompt e envia para a IA (Gemini ou Groq)
        |
IA retorna JSON com scores e descricoes
        |
API salva o resultado no banco e retorna ao frontend
        |
Frontend exibe graficos radar, ranking e explicacoes
        |
Estudante pode exportar o resultado em PDF
```

## Como Executar o Projeto

### Pre-requisitos

- [Go 1.25+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker](https://www.docker.com/) (para o PostgreSQL)
- Chave de API do [Google Gemini](https://aistudio.google.com/apikey) e/ou [Groq](https://console.groq.com/)

### 1. Clone o repositorio

```bash
git clone https://github.com/lirajoaop/psiencontra.git
cd psiencontra
```

### 2. Configure as variaveis de ambiente

```bash
cp .env.example .env
```

Edite o `.env` com suas chaves:

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

Em outro terminal:

```bash
cd web
npm install
npm run dev
```

O frontend estara disponivel em `http://localhost:3000`.

## Deploy em Producao

| Servico | Plataforma | URL |
|---|---|---|
| Frontend | Vercel | [psiencontra.vercel.app](https://psiencontra.vercel.app) |
| API + Banco | Railway | psiencontra-production.up.railway.app |

### Variaveis de ambiente

**Vercel (frontend):**
- `NEXT_PUBLIC_API_URL` — URL da API + `/api/v1`

**Railway (API):**
- `DATABASE_URL` — fornecida automaticamente pelo PostgreSQL do Railway
- `GEMINI_API_KEY` — chave da API do Google Gemini
- `GROQ_API_KEY` — chave da API do Groq
- `FRONTEND_URL` — URL do frontend na Vercel
- `PORT` — 8080

## Aprendizados

Durante o desenvolvimento deste projeto, os principais aprendizados foram:

- **Integracao com IAs generativas**: como montar prompts eficientes para obter respostas em JSON estruturado, lidar com limites de tokens e implementar fallback entre provedores diferentes
- **Arquitetura fullstack**: separacao de responsabilidades entre frontend (React/Next.js) e backend (Go/Gin), comunicacao via REST API e gerenciamento de estado
- **Deploy de aplicacoes distribuidas**: configurar e conectar servicos em plataformas diferentes (Vercel + Railway), gerenciar variaveis de ambiente e lidar com CORS em producao
- **UX e acessibilidade**: implementacao de dark mode, animacoes suaves, design responsivo e feedback visual durante carregamento

## Licenca

Este projeto e de uso academico e educacional. Resultado orientativo, nao substitui acompanhamento profissional.
