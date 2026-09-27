const labels = {
  queued: "Queued",
  sending: "Sending",
  sent: "Sent",
  running: "Working",
  completed: "Completed",
  failed: "Needs attention",
} as const;

export function WorkStatus({ state }: { state: keyof typeof labels }) {
  return (
    <span className="text-xs font-medium text-muted">{labels[state]}</span>
  );
}
