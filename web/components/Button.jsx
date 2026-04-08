"use client";

export default function Button({ children, onClick, variant = "primary", disabled = false, className = "", type = "button" }) {
  const base = "px-6 py-3 rounded-xl font-semibold transition-all duration-200 cursor-pointer text-center";
  const variants = {
    primary: "bg-violet-600 text-white hover:bg-violet-700 shadow-lg shadow-violet-200 hover:shadow-violet-300 dark:shadow-violet-900/50 dark:hover:shadow-violet-800/50",
    secondary: "bg-white text-violet-700 border-2 border-violet-200 hover:border-violet-400 hover:bg-violet-50 dark:bg-gray-800 dark:text-violet-300 dark:border-violet-700 dark:hover:border-violet-500 dark:hover:bg-gray-700",
    ghost: "text-violet-600 hover:bg-violet-50 dark:text-violet-400 dark:hover:bg-violet-900/30",
  };

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={`${base} ${variants[variant]} ${disabled ? "opacity-50 cursor-not-allowed" : ""} ${className}`}
    >
      {children}
    </button>
  );
}
