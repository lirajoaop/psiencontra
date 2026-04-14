import Link from "next/link";
import ThemeToggle from "@/components/ThemeToggle";

export const metadata = {
  title: "Privacidade · PsiEncontra",
  description: "Política de privacidade e uso de dados do PsiEncontra.",
};

export default function Privacidade() {
  return (
    <main className="flex-1 bg-gradient-to-br from-violet-50 to-purple-50 dark:from-gray-900 dark:to-gray-950">
      <div className="absolute top-4 right-4 z-20">
        <ThemeToggle />
      </div>
      <div className="max-w-3xl mx-auto px-6 py-16 text-gray-700 dark:text-gray-300">
        <h1 className="text-3xl md:text-4xl font-bold text-violet-900 dark:text-violet-200 mb-10">
          Privacidade e uso dos seus dados
        </h1>

        <section className="space-y-4 mb-8">
          <h2 className="text-xl font-semibold text-violet-800 dark:text-violet-300">
            O que o PsiEncontra é — e o que não é
          </h2>
          <p>
            O PsiEncontra é uma <strong>ferramenta de autorreflexão</strong> voltada a
            estudantes de Psicologia. Ele <strong>não é um teste psicológico</strong>,
            não fornece diagnóstico e não substitui orientação profissional, supervisão
            acadêmica ou psicoterapia. Os resultados são uma leitura aproximada das suas
            respostas no momento em que foram dadas.
          </p>
        </section>

        <section className="space-y-4 mb-8">
          <h2 className="text-xl font-semibold text-violet-800 dark:text-violet-300">
            Que dados coletamos
          </h2>
          <ul className="list-disc pl-6 space-y-2">
            <li>
              <strong>Respostas ao questionário</strong> (opções marcadas e textos
              livres), vinculadas à sessão da sua resposta.
            </li>
            <li>
              <strong>Dados de conta</strong>, apenas se você optar por criar uma: nome,
              e-mail e, no caso do login com Google, informações básicas do perfil
              (nome, e-mail, foto).
            </li>
            <li>
              <strong>Metadados técnicos</strong> mínimos (data/hora da sessão, tipo de
              questionário escolhido).
            </li>
          </ul>
          <p>
            Não coletamos dados sensíveis fora do próprio conteúdo das suas respostas.
            Você pode usar o PsiEncontra <strong>de forma anônima</strong>, sem criar
            conta.
          </p>
        </section>

        <section className="space-y-4 mb-8">
          <h2 className="text-xl font-semibold text-violet-800 dark:text-violet-300">
            Para que usamos
          </h2>
          <ul className="list-disc pl-6 space-y-2">
            <li>Gerar o seu resultado (ranking de afinidade, gráficos e PDF).</li>
            <li>Permitir que você recupere resultados anteriores, se tiver conta.</li>
            <li>
              Aprimorar a ferramenta de forma agregada e anonimizada (nunca
              identificando respostas individuais).
            </li>
          </ul>
        </section>

        <section className="space-y-4 mb-8">
          <h2 className="text-xl font-semibold text-violet-800 dark:text-violet-300">
            Com quem compartilhamos
          </h2>
          <p>
            Suas respostas do modo rápido são enviadas a provedores de IA (Google Gemini
            e, como fallback, Groq) para gerar a análise interpretativa. Nenhum outro
            terceiro recebe seus dados. Não vendemos nem cedemos dados a anunciantes.
          </p>
        </section>

        <section className="space-y-4 mb-8">
          <h2 className="text-xl font-semibold text-violet-800 dark:text-violet-300">
            Base legal (LGPD)
          </h2>
          <p>
            O tratamento dos seus dados ocorre com base no <strong>consentimento</strong>{" "}
            que você dá ao enviar o questionário, e no <strong>legítimo interesse</strong>{" "}
            em oferecer o serviço. Você pode, a qualquer momento, solicitar exclusão das
            suas respostas e da sua conta pelo e-mail abaixo.
          </p>
        </section>

        <section className="space-y-4 mb-8">
          <h2 className="text-xl font-semibold text-violet-800 dark:text-violet-300">
            Seus direitos
          </h2>
          <p>
            Conforme a Lei Geral de Proteção de Dados (Lei nº 13.709/2018), você tem
            direito a confirmar o tratamento, acessar, corrigir, anonimizar, portar ou
            eliminar seus dados, bem como revogar o consentimento. Para exercer qualquer
            um deles, escreva para{" "}
            <a
              href="mailto:joaopedroababa132@gmail.com"
              className="text-violet-700 dark:text-violet-300 underline"
            >
              joaopedroababa132@gmail.com
            </a>
            .
          </p>
        </section>

        <section className="space-y-4 mb-8">
          <h2 className="text-xl font-semibold text-violet-800 dark:text-violet-300">
            Segurança
          </h2>
          <p>
            Senhas são armazenadas com hash bcrypt, tráfego é servido sob HTTPS e tokens
            de autenticação são emitidos como cookies HttpOnly. Ainda assim, nenhum
            sistema é 100% à prova de incidentes — por isso recomendamos não inserir no
            questionário informações identificáveis de terceiros.
          </p>
        </section>

        <div className="mt-10">
          <Link
            href="/"
            className="text-violet-700 dark:text-violet-300 underline text-sm"
          >
            ← Voltar ao início
          </Link>
        </div>
      </div>
    </main>
  );
}
