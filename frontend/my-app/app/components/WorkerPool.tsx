"use client";

import { WorkerState } from "../lib/api";

interface Props {
  workers: WorkerState[] | null;
}

export default function WorkerPool({ workers }: Props) {
  const sorted = [...(workers ?? [])].sort((a, b) => a.worker_id - b.worker_id);

  return (
    <section className="rounded-xl border border-[#30363d] bg-[#161b22] p-6">
      <h2 className="mb-4 text-sm font-semibold uppercase tracking-wider text-[#8b949e]">
        Worker Pool
      </h2>

      <div className="flex flex-wrap gap-2">
        {sorted.length === 0 ? (
          <span className="text-xs text-[#6e7681]">No workers registered</span>
        ) : (
          sorted.map((w) => (
            <div
              key={w.worker_id}
              className="flex items-center gap-2 rounded-lg border border-[#30363d] bg-[#0d1117] px-4 py-2.5 min-w-[140px] transition-colors hover:border-[#484f58]"
            >
              <svg className="h-4 w-4 text-[#8b949e]" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M10.343 3.94c.09-.542.56-.94 1.11-.94h1.093c.55 0 1.02.398 1.11.94l.149.894c.07.424.384.764.78.93.398.164.855.142 1.205-.108l.737-.527a1.125 1.125 0 011.45.12l.773.774c.39.389.44 1.002.12 1.45l-.527.737c-.25.35-.272.806-.107 1.204.165.397.505.71.93.78l.893.15c.543.09.94.56.94 1.109v1.094c0 .55-.397 1.02-.94 1.11l-.893.149c-.425.07-.765.383-.93.78-.165.398-.143.854.107 1.204l.527.738c.32.447.269 1.06-.12 1.45l-.774.773a1.125 1.125 0 01-1.449.12l-.738-.527c-.35-.25-.806-.272-1.203-.107-.397.165-.71.505-.781.929l-.149.894c-.09.542-.56.94-1.11.94h-1.094c-.55 0-1.019-.398-1.11-.94l-.148-.894c-.071-.424-.384-.764-.781-.93-.398-.164-.854-.142-1.204.108l-.738.527c-.447.32-1.06.269-1.45-.12l-.773-.774a1.125 1.125 0 01-.12-1.45l.527-.737c.25-.35.273-.806.108-1.204-.165-.397-.506-.71-.93-.78l-.894-.15c-.542-.09-.94-.56-.94-1.109v-1.094c0-.55.398-1.02.94-1.11l.894-.149c.424-.07.765-.383.93-.78.165-.398.143-.854-.107-1.204l-.527-.738a1.125 1.125 0 01.12-1.45l.773-.773a1.125 1.125 0 011.45-.12l.737.527c.35.25.807.272 1.204.107.397-.165.71-.505.78-.929l.15-.894z" />
                <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span className="text-sm font-semibold text-[#e6edf3]">
                W{w.worker_id}
              </span>
              <StatusBadge status={w.status} />
            </div>
          ))
        )}
      </div>
    </section>
  );
}

function StatusBadge({ status }: { status: string }) {
  const config: Record<string, { bg: string; text: string; label: string }> = {
    idle: { bg: "rgba(210, 153, 34, 0.15)", text: "#d29922", label: "idle" },
    busy: { bg: "rgba(63, 185, 80, 0.15)", text: "#3fb950", label: "busy" },
    offline: { bg: "rgba(248, 81, 73, 0.15)", text: "#f85149", label: "offline" },
  };
  const c = config[status] ?? config.offline;

  return (
    <span
      className="ml-auto rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase"
      style={{ backgroundColor: c.bg, color: c.text }}
    >
      {c.label}
    </span>
  );
}
