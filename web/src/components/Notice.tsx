export function Notice({ children, id }: { children: string; id?: string }) {
  return (
    <p
      id={id}
      role="alert"
      className="rounded-md border border-danger p-3 text-sm text-danger"
    >
      {children}
    </p>
  );
}
