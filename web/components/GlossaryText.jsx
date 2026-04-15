"use client";

import { useState } from "react";
import MuiTooltip from "@mui/material/Tooltip";
import ClickAwayListener from "@mui/material/ClickAwayListener";
import { GLOSSARY, GLOSSARY_REGEX } from "@/lib/glossary";

function AcronymTooltip({ acronym, definition }) {
  const [open, setOpen] = useState(false);

  return (
    <ClickAwayListener onClickAway={() => setOpen(false)}>
      <MuiTooltip
        title={definition}
        placement="top"
        arrow
        open={open}
        onClose={() => setOpen(false)}
        disableFocusListener
        disableTouchListener
        slotProps={{
          tooltip: { sx: { fontSize: 12, py: 0.75, px: 1.25, maxWidth: 280 } },
        }}
      >
        <button
          type="button"
          aria-label={`Definição de ${acronym}`}
          onMouseEnter={() => setOpen(true)}
          onMouseLeave={() => setOpen(false)}
          onClick={() => setOpen((v) => !v)}
          className="no-underline border-b border-dotted border-violet-400 dark:border-violet-500 cursor-help bg-transparent p-0 text-inherit focus:outline-none focus-visible:ring-2 focus-visible:ring-violet-400 rounded-sm"
        >
          {acronym}
        </button>
      </MuiTooltip>
    </ClickAwayListener>
  );
}

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
      <AcronymTooltip
        key={`${acronym}-${match.index}`}
        acronym={acronym}
        definition={GLOSSARY[acronym]}
      />
    );
    lastIndex = match.index + acronym.length;
  }
  if (lastIndex < children.length) parts.push(children.slice(lastIndex));

  return <>{parts}</>;
}
