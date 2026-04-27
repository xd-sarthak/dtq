"use client";

import { QueueData } from "../lib/api";

interface Props {
  queues: QueueData | null;
}

export default function QueueContents({ queues }: Props) {
  const ready = queues?.ready ?? [];
  const delayed = queues?.delayed ?? [];

  return (
    <section className="rounded-xl border border-[#30363d] bg-[#161b22] p-6">
      <h2 className="mb-4 text-sm font-semibold uppercase tracking-wider text-[#8b949e]">
        Queue Contents
      </h2>

      {/* Ready Queue */}
      <div className="mb-4">
        <div className="mb-2 flex items-center gap-2">
          <span className="text-sm">⚡</span>
          <h3 className="text-xs font-semibold text-[#e6edf3]">
            Ready Queue ({ready.length})
          </h3>
        </div>
        {ready.length === 0 ? (
          <p className="pl-5 text-xs italic text-[#6e7681]">Empty</p>
        ) : (
          <div className="flex flex-col gap-1 pl-5">
            {ready.map((task) => (
              <div
                key={task.id}
                className="flex items-center justify-between rounded border border-[#21262d] bg-[#0d1117] px-3 py-1.5 text-xs"
              >
                <span className="font-mono text-[#58a6ff]">{task.id}</span>
                <span className="text-[#6e7681]">p:{task.priority}</span>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Delayed Queue */}
      <div>
        <div className="mb-2 flex items-center gap-2">
          <span className="text-sm">⏳</span>
          <h3 className="text-xs font-semibold text-[#e6edf3]">
            Delayed Queue ({delayed.length})
          </h3>
        </div>
        {delayed.length === 0 ? (
          <p className="pl-5 text-xs italic text-[#6e7681]">Empty</p>
        ) : (
          <div className="flex flex-col gap-1 pl-5">
            {delayed.map((entry) => (
              <div
                key={entry.task.id}
                className="flex items-center justify-between rounded border border-[#21262d] bg-[#0d1117] px-3 py-1.5 text-xs"
              >
                <span className="font-mono text-[#d29922]">
                  {entry.task.id}
                </span>
                <span className="text-[#6e7681]">
                  {new Date(entry.execute_at * 1000).toLocaleTimeString()}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
