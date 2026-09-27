export function Notice({ children }: { children: string }) {
  return (
    <p
      role="alert"
      className="rounded-md border border-danger p-3 text-sm text-danger"
    >
      {children}
    </p>
  );
}
