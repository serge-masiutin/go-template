import type { ButtonHTMLAttributes } from "react";

const variants = {
  primary: "bg-accent text-white hover:bg-accent-hover",
  secondary: "border border-border bg-panel text-ink hover:bg-surface",
} as const;

type Props = Omit<ButtonHTMLAttributes<HTMLButtonElement>, "className"> & {
  variant?: keyof typeof variants;
};

export function Button({ variant = "primary", ...props }: Props) {
  return (
    <button
      type="button"
      {...props}
      className={`rounded-md px-4 py-2 text-sm transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-wait disabled:opacity-50 ${variants[variant]}`}
    />
  );
}
