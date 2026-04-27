"use client";

import { TaskEvent } from "../lib/api";

interface Props {
  events: TaskEvent[] | null;
}

const badgeConfig: Record<string, { bg: string; text: string }> = {
  completed: { bg: "rgba(63, 185, 80, 0.15)", text: "#3fb950" },
  started: { bg: "rgba(88, 166, 255, 0.15)", text: "#58a6ff" },
  failed: { bg: "rgba(248, 81, 73, 0.15)", text: "#f85149" },
  retrying: { bg: "rgba(210, 153, 34, 0.15)", text: "#d29922" },
  promoted: { bg: "rgba(163, 113, 247, 0.15)", text: "#a371f7" },
  submitted: { bg: "rgba(57, 211, 83, 0.15)", text: "#39d353" },
  dead_lettered: { bg: "rgba(248, 81, 73, 0.15)", text: "#f85149" },
};

export default function ActivityLog({ events }: Props) {
  const list = events ?? [];

  return (
    <section className="rounded-xl border border-[#30363d] bg-[#161b22] p-6 flex flex-col">
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className="font-mono text-sm text-[#8b949e]">&gt;_</span>
          <h2 className="text-sm font-semibold uppercase tracking-wider text-[#8b949e]">
            Activity Log
          </h2>
        </div>
        <span className="rounded-full bg-[#21262d] px-2.5 py-0.5 text-[10px] font-mono text-[#8b949e]">
          {list.length} events
        </span>
      </div>

      <div className="flex-1 overflow-y-auto max-h-[400px] scrollbar-thin">
        {list.length === 0 ? (
          <p className="text-xs italic text-[#6e7681]">No events yet</p>
        ) : (
          <div className="flex flex-col gap-0.5">
            {list.map((ev) => (
              <EventRow key={ev.id} event={ev} />
            ))}
          </div>
        )}
      </div>
    </section>
  );
}

function EventRow({ event }: { event: TaskEvent }) {
  const badge = badgeConfig[event.type] ?? badgeConfig.submitted;
  const time = new Date(event.timestamp).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: true,
  });

  // Parse the worker ID display
  const workerLabel =
    event.worker_id >= 0 ? `W${event.worker_id}` : "";

  return (
    <div className="group flex items-start gap-3 rounded px-2 py-1.5 transition-colors hover:bg-[#1c2333]">
      {/* Timestamp */}
      <span className="shrink-0 text-[11px] font-mono text-[#6e7681] pt-0.5 min-w-[85px]">
        {time}
      </span>

      {/* Badge */}
      <span
        className="shrink-0 mt-0.5 inline-flex items-center rounded px-2 py-0.5 text-[10px] font-bold"
        style={{ backgroundColor: badge.bg, color: badge.text }}
      >
        {event.type}
      </span>

      {/* Detail */}
      <span className="text-xs text-[#e6edf3] font-mono leading-relaxed">
        <span className="text-[#58a6ff]">{event.task_id}</span>
        {workerLabel && (
          <span className="text-[#d29922]"> {workerLabel}</span>
        )}
        {" — "}
        <span className="text-[#8b949e]">
          {event.detail.replace(event.task_id, "").replace(/^[\s\-]+/, "")}
        </span>
      </span>
    </div>
  );
}
