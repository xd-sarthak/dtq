"use client";

import { FailedTask } from "../lib/api";

interface Props {
  tasks: FailedTask[] | null;
}

export default function FailedTasks({ tasks }: Props) {
  const list = tasks ?? [];

  return (
    <section className="rounded-xl border border-[#30363d] bg-[#161b22] p-6">
      <div className="mb-4 flex items-center gap-2">
        <span className="text-sm">⚠</span>
        <h2 className="text-sm font-semibold uppercase tracking-wider text-[#8b949e]">
          Failed Tasks (Dead Letter)
        </h2>
      </div>

      {list.length === 0 ? (
        <p className="text-xs italic text-[#6e7681]">No failed tasks</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead>
              <tr className="border-b border-[#21262d] text-[#8b949e]">
                <th className="pb-2 pr-4 font-medium">Task ID</th>
                <th className="pb-2 pr-4 font-medium">Priority</th>
                <th className="pb-2 pr-4 font-medium">Attempts</th>
                <th className="pb-2 pr-4 font-medium">Reason</th>
                <th className="pb-2 font-medium text-right">Failed At</th>
              </tr>
            </thead>
            <tbody>
              {list.map((ft, i) => (
                <tr
                  key={`${ft.task.id}-${i}`}
                  className="border-b border-[#21262d]/50 transition-colors hover:bg-[#1c2333]"
                >
                  <td className="py-2.5 pr-4">
                    <span className="font-mono text-[#58a6ff]">
                      {ft.task.id}
                    </span>
                  </td>
                  <td className="py-2.5 pr-4 text-[#e6edf3]">
                    {ft.task.priority}
                  </td>
                  <td className="py-2.5 pr-4 text-[#e6edf3]">
                    {ft.task.retries}/{ft.task.max_retries}
                  </td>
                  <td className="py-2.5 pr-4">
                    <span className="text-[#f85149]">
                      {ft.error}
                    </span>
                  </td>
                  <td className="py-2.5 text-right text-[#8b949e]">
                    {new Date(ft.failed_at).toLocaleTimeString([], {
                      hour: "2-digit",
                      minute: "2-digit",
                      second: "2-digit",
                      hour12: true,
                    })}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
