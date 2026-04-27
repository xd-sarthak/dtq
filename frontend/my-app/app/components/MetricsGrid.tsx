"use client";

import { EnhancedMetrics } from "../lib/api";

interface Props {
  metrics: EnhancedMetrics | null;
}

const cards: {
  key: keyof EnhancedMetrics;
  label: string;
  color: string;
  icon: string;
}[] = [
  { key: "total_processed", label: "PROCESSED", color: "#3fb950", icon: "✅" },
  { key: "total_failed", label: "FAILED", color: "#f85149", icon: "❌" },
  { key: "total_retries", label: "RETRIES", color: "#d29922", icon: "🔄" },
  { key: "queue_size", label: "QUEUE SIZE", color: "#58a6ff", icon: "📋" },
  { key: "active_workers", label: "ACTIVE WORKERS", color: "#db6d28", icon: "👷" },
  { key: "delayed_queue_size", label: "DELAYED", color: "#d29922", icon: "⏳" },
  { key: "dead_letter_size", label: "DEAD LETTER", color: "#f85149", icon: "💀" },
  { key: "total_submitted", label: "SUBMITTED", color: "#a371f7", icon: "📊" },
];

export default function MetricsGrid({ metrics }: Props) {
  const successRate = metrics?.success_rate ?? 0;

  return (
    <section className="rounded-xl border border-[#30363d] bg-[#161b22] p-6">
      <h2 className="mb-4 text-sm font-semibold uppercase tracking-wider text-[#8b949e]">
        Metrics
      </h2>

      <div className="grid grid-cols-2 gap-3">
        {cards.map(({ key, label, color, icon }) => (
          <div
            key={key}
            className="rounded-lg border border-[#30363d] bg-[#0d1117] p-4 transition-colors hover:border-[#484f58]"
          >
            <div className="flex items-center gap-2">
              <span className="text-base">{icon}</span>
              <span className="text-2xl font-bold" style={{ color }}>
                {(metrics?.[key] ?? 0).toLocaleString()}
              </span>
            </div>
            <span className="mt-1 block text-[10px] font-semibold uppercase tracking-widest text-[#8b949e]">
              {label}
            </span>
          </div>
        ))}
      </div>

      {/* Success rate bar */}
      <div className="mt-4">
        <div className="mb-1 flex items-center justify-between">
          <span className="text-xs text-[#8b949e]">Success Rate</span>
          <span
            className="text-xs font-bold"
            style={{ color: successRate > 90 ? "#3fb950" : successRate > 50 ? "#d29922" : "#f85149" }}
          >
            {successRate.toFixed(1)}%
          </span>
        </div>
        <div className="h-2 w-full overflow-hidden rounded-full bg-[#21262d]">
          <div
            className="h-full rounded-full transition-all duration-500 ease-out"
            style={{
              width: `${Math.min(successRate, 100)}%`,
              backgroundColor:
                successRate > 90
                  ? "#58a6ff"
                  : successRate > 50
                  ? "#d29922"
                  : "#f85149",
            }}
          />
        </div>
      </div>
    </section>
  );
}
