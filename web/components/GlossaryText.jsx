"use client";

import { GLOSSARY, GLOSSARY_REGEX } from "@/lib/glossary";

// Wraps known Brazilian institutional acronyms (SUAS, CRAS, UBS, AVC, ...) in
// <abbr> with a native tooltip so students unfamiliar with the terminology get
// context on hover (desktop) or long-press (mobile).
export default function GlossaryText({ children }) {
  if (typeof children !== "string") return children;

  const parts = [];
  let lastIndex = 0;
  let match;
  const re = new RegExp(GLOSSARY_REGEX.source, GLOSSARY_REGEX.flags);

  while ((match = re.exec(children)) !== null) {
    if (match.index > lastIndex) {
      parts.push(children.slice(lastIndex, match.index));
    }
    const acronym = match[1];
    parts.push(
      <abbr
        key={`${acronym}-${match.index}`}
        title={GLOSSARY[acronym]}
        className="no-underline border-b border-dotted border-violet-400 dark:border-violet-500 cursor-help decoration-0"
      >
        {acronym}
      </abbr>
    );
    lastIndex = match.index + acronym.length;
  }
  if (lastIndex < children.length) parts.push(children.slice(lastIndex));

  return <>{parts}</>;
}
