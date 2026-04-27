"use client";

import { EnhancedMetrics } from "../lib/api";

interface Props {
  metrics: EnhancedMetrics | null;
}

export default function TaskFlowPipeline({ metrics }: Props) {
  const submitted = metrics?.total_submitted ?? 0;
  const delayed = metrics?.delayed_queue_size ?? 0;
  const ready = metrics?.queue_size ?? 0;
  const active = metrics?.active_workers ?? 0;
  const completed = metrics?.total_processed ?? 0;
  const retried = metrics?.total_retries ?? 0;
  const deadLetter = metrics?.dead_letter_size ?? 0;

  return (
    <section className="rounded-xl border border-[#30363d] bg-[#161b22] p-6">
      <h2 className="mb-6 text-sm font-semibold uppercase tracking-wider text-[#8b949e]">
        Task Flow Pipeline
      </h2>

      <div className="flex items-center justify-center gap-3 flex-wrap">
        {/* Submitted */}
        <PipelineStage value={submitted} label="Submitted" color="#58a6ff" />
        <Arrow />
        {/* Delayed Queue */}
        <PipelineStage value={delayed} label="Delayed Queue" color="#d29922" />
        <Arrow />
        {/* Ready Queue */}
        <PipelineStage value={ready} label="Ready Queue" color="#58a6ff" />
        <Arrow />
        {/* Workers Active */}
        <PipelineStage value={active} label="Workers Active" color="#58a6ff" />
        <Arrow />
        {/* Outcomes */}
        <div className="flex flex-col gap-1 rounded-lg border border-[#30363d] bg-[#0d1117] px-4 py-3 min-w-[140px]">
          <OutcomeLine label="Completed" value={completed} color="#3fb950" />
          <OutcomeLine label="Retried" value={retried} color="#d29922" />
          <OutcomeLine label="Dead Letter" value={deadLetter} color="#f85149" />
        </div>
      </div>
    </section>
  );
}

function PipelineStage({
  value,
  label,
  color,
}: {
  value: number;
  label: string;
  color: string;
}) {
  return (
    <div className="flex flex-col items-center gap-1 rounded-lg border border-[#30363d] bg-[#0d1117] px-6 py-3 min-w-[100px]">
      <span className="text-2xl font-bold" style={{ color }}>
        {value.toLocaleString()}
      </span>
      <span className="text-xs text-[#8b949e]">{label}</span>
    </div>
  );
}

function Arrow() {
  return (
    <span className="text-[#30363d] text-lg font-bold select-none">→</span>
  );
}

function OutcomeLine({
  label,
  value,
  color,
}: {
  label: string;
  value: number;
  color: string;
}) {
  return (
    <div className="flex items-center gap-2 text-xs">
      <span
        className="inline-block h-2 w-2 rounded-full"
        style={{ backgroundColor: color }}
      />
      <span style={{ color }}>
        {label}: {value.toLocaleString()}
      </span>
    </div>
  );
}
