"use client";

import { useState } from "react";
import MuiTooltip from "@mui/material/Tooltip";
import ClickAwayListener from "@mui/material/ClickAwayListener";
import { Info } from "lucide-react";

export default function Tooltip({ label, side = "top", className = "" }) {
  const [open, setOpen] = useState(false);

  return (
    <ClickAwayListener onClickAway={() => setOpen(false)}>
      <MuiTooltip
        title={label}
        placement={side}
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
          aria-label="Mais informações"
          onMouseEnter={() => setOpen(true)}
          onMouseLeave={() => setOpen(false)}
          onClick={() => setOpen((v) => !v)}
          className={`inline-flex items-center justify-center rounded-full text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-1 focus-visible:ring-gray-500 ${className}`}
        >
          <Info className="h-4 w-4" strokeWidth={2} />
        </button>
      </MuiTooltip>
    </ClickAwayListener>
  );
}
