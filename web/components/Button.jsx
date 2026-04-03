"use client";

export default function Button({ children, onClick, variant = "primary", disabled = false, className = "" }) {
  const base = "px-6 py-3 rounded-xl font-semibold transition-all duration-200 cursor-pointer text-center";
  const variants = {
    primary: "bg-violet-600 text-white hover:bg-violet-700 shadow-lg shadow-violet-200 hover:shadow-violet-300",
    secondary: "bg-white text-violet-700 border-2 border-violet-200 hover:border-violet-400 hover:bg-violet-50",
    ghost: "text-violet-600 hover:bg-violet-50",
  };

  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`${base} ${variants[variant]} ${disabled ? "opacity-50 cursor-not-allowed" : ""} ${className}`}
    >
      {children}
    </button>
  );
}
