import type { InputHTMLAttributes, TextareaHTMLAttributes } from "react";
import { Notice } from "./Notice";

type NativeFieldProps =
  | "id"
  | "name"
  | "className"
  | "style"
  | "children"
  | "aria-invalid"
  | "aria-label"
  | "aria-labelledby";
type FieldContent = {
  id: string;
  name: string;
  label: string;
  hint?: string;
  error?: string;
};
type Props = FieldContent &
  (
    | (Omit<
        InputHTMLAttributes<HTMLInputElement>,
        NativeFieldProps | "type"
      > & {
        multiline?: false;
        type?: "text" | "email" | "password";
      })
    | (Omit<TextareaHTMLAttributes<HTMLTextAreaElement>, NativeFieldProps> & {
        multiline: true;
      })
  );

export function Field({ id, name, label, hint, error, ...props }: Props) {
  const describedBy =
    [props["aria-describedby"], hint && `${id}-hint`, error && `${id}-error`]
      .filter(Boolean)
      .join(" ") || undefined;
  const accessibility = {
    id,
    name,
    "aria-describedby": describedBy,
    "aria-invalid": error ? true : undefined,
  } as const;
  const focus =
    "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50";
  let control;
  if (props.multiline) {
    const { multiline: _, ...nativeProps } = props;
    control = (
      <textarea
        {...nativeProps}
        {...accessibility}
        className={`w-full rounded-md border border-border bg-panel p-3 ${focus}`}
      />
    );
  } else {
    const { multiline: _, type = "text", ...nativeProps } = props;
    control = (
      <input
        {...nativeProps}
        {...accessibility}
        type={type}
        className={`w-full rounded-md border border-border bg-panel px-3 py-2 ${focus}`}
      />
    );
  }
  return (
    <div className={props.multiline ? "space-y-3" : "space-y-2"}>
      <label
        htmlFor={id}
        className={
          props.multiline ? "block text-sm font-medium" : "block text-sm"
        }
      >
        {label}
      </label>
      {control}
      {hint && (
        <p id={`${id}-hint`} className="text-sm text-muted">
          {hint}
        </p>
      )}
      {error && <Notice id={`${id}-error`}>{error}</Notice>}
    </div>
  );
}
