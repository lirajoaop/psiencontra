"use client";

import { useState } from "react";

function buildShareText(topApproach, topField) {
  const parts = [];
  if (topApproach) parts.push(`minha principal abordagem é **${topApproach}**`);
  if (topField) parts.push(`meu campo de maior afinidade é **${topField}**`);
  const intro = parts.length
    ? `Fiz o PsiEncontra e descobri que ${parts.join(" e ")}.`
    : "Fiz o PsiEncontra e descobri minha afinidade em Psicologia.";
  return `${intro} Descubra a sua:`;
}

function stripMarkdown(text) {
  return text.replace(/\*\*(.+?)\*\*/g, "$1");
}

export default function ShareButtons({ url, topApproach, topField }) {
  const [copied, setCopied] = useState(false);
  const [nativeShared, setNativeShared] = useState(false);

  const plainText = stripMarkdown(buildShareText(topApproach, topField));
  const whatsappText = `${plainText} ${url}`;
  const xText = plainText;

  const whatsappHref = `https://wa.me/?text=${encodeURIComponent(whatsappText)}`;
  const xHref = `https://x.com/intent/post?text=${encodeURIComponent(xText)}&url=${encodeURIComponent(url)}`;

  async function copyLink() {
    try {
      await navigator.clipboard.writeText(`${plainText} ${url}`);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      setCopied(false);
    }
  }

  async function nativeShare() {
    if (typeof navigator === "undefined" || !navigator.share) {
      copyLink();
      return;
    }
    try {
      await navigator.share({
        title: "Meu resultado no PsiEncontra",
        text: plainText,
        url,
      });
      setNativeShared(true);
      setTimeout(() => setNativeShared(false), 2000);
    } catch {
      // user cancelled — no-op
    }
  }

  return (
    <section className="rounded-2xl border border-violet-100 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 md:p-8 mb-10">
      <h2 className="text-xl font-bold text-violet-900 dark:text-violet-200 mb-2">
        Compartilhar resultado
      </h2>
      <p className="text-sm text-gray-600 dark:text-gray-400 mb-5">
        Este resultado é público via link — qualquer pessoa com a URL pode visualizá-lo.
      </p>
      <div className="flex flex-wrap gap-3">
        <a
          href={whatsappHref}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl font-semibold text-sm bg-[#25D366] text-white hover:brightness-110 transition-all"
          aria-label="Compartilhar no WhatsApp"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path d="M20.52 3.48A11.93 11.93 0 0 0 12.05 0C5.5 0 .16 5.34.16 11.9c0 2.1.55 4.15 1.6 5.96L0 24l6.3-1.65a11.87 11.87 0 0 0 5.74 1.46h.01c6.55 0 11.89-5.34 11.89-11.9 0-3.18-1.24-6.17-3.42-8.43zM12.05 21.8h-.01a9.9 9.9 0 0 1-5.04-1.38l-.36-.21-3.74.98 1-3.65-.24-.37a9.88 9.88 0 0 1-1.52-5.27c0-5.46 4.45-9.9 9.91-9.9 2.64 0 5.13 1.03 7 2.9a9.85 9.85 0 0 1 2.9 7c0 5.46-4.44 9.9-9.9 9.9zm5.43-7.42c-.3-.15-1.76-.87-2.03-.97-.27-.1-.47-.15-.67.15-.2.3-.77.97-.94 1.17-.17.2-.35.22-.65.07-.3-.15-1.26-.46-2.4-1.48-.88-.79-1.48-1.76-1.66-2.06-.17-.3-.02-.46.13-.6.14-.13.3-.35.45-.52.15-.17.2-.3.3-.5.1-.2.05-.37-.03-.52-.07-.15-.67-1.62-.92-2.22-.24-.58-.48-.5-.66-.51l-.57-.01c-.2 0-.52.07-.8.37-.27.3-1.04 1.02-1.04 2.48s1.07 2.88 1.22 3.08c.15.2 2.1 3.2 5.08 4.49.71.3 1.26.48 1.7.62.71.23 1.36.2 1.87.12.57-.08 1.76-.72 2.01-1.41.25-.7.25-1.29.17-1.41-.07-.12-.27-.2-.57-.35z"/>
          </svg>
          WhatsApp
        </a>
        <a
          href={xHref}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl font-semibold text-sm bg-black text-white hover:bg-gray-800 transition-all"
          aria-label="Compartilhar no X"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
          </svg>
          X (Twitter)
        </a>
        <button
          type="button"
          onClick={nativeShare}
          className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl font-semibold text-sm bg-gradient-to-tr from-[#feda75] via-[#fa7e1e] to-[#d62976] text-white hover:brightness-110 transition-all cursor-pointer"
          aria-label="Compartilhar no Instagram ou app nativo"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path d="M12 2.163c3.204 0 3.584.012 4.85.07 1.17.053 1.805.248 2.227.413.56.217.96.477 1.38.896.42.42.68.82.896 1.38.165.422.36 1.057.413 2.227.058 1.266.07 1.646.07 4.85s-.012 3.584-.07 4.85c-.053 1.17-.248 1.805-.413 2.227-.217.56-.477.96-.896 1.38-.42.42-.82.68-1.38.896-.422.165-1.057.36-2.227.413-1.266.058-1.646.07-4.85.07s-3.584-.012-4.85-.07c-1.17-.053-1.805-.248-2.227-.413a3.72 3.72 0 0 1-1.38-.896 3.72 3.72 0 0 1-.896-1.38c-.165-.422-.36-1.057-.413-2.227-.058-1.266-.07-1.646-.07-4.85s.012-3.584.07-4.85c.053-1.17.248-1.805.413-2.227.217-.56.477-.96.896-1.38.42-.42.82-.68 1.38-.896.422-.165 1.057-.36 2.227-.413 1.266-.058 1.646-.07 4.85-.07zm0 3.838a5.999 5.999 0 1 0 0 11.998 5.999 5.999 0 0 0 0-11.998zm0 9.897a3.898 3.898 0 1 1 0-7.796 3.898 3.898 0 0 1 0 7.796zm6.406-10.845a1.44 1.44 0 1 1-2.88 0 1.44 1.44 0 0 1 2.88 0z"/>
          </svg>
          {nativeShared ? "Compartilhado!" : "Instagram / Mais"}
        </button>
        <button
          type="button"
          onClick={copyLink}
          className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl font-semibold text-sm border-2 border-violet-300 dark:border-violet-700 text-violet-700 dark:text-violet-300 hover:border-violet-500 dark:hover:border-violet-500 transition-colors cursor-pointer"
          aria-label="Copiar link"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
          </svg>
          {copied ? "Link copiado!" : "Copiar link"}
        </button>
      </div>
    </section>
  );
}
