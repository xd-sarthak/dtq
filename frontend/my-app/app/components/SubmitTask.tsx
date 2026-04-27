"use client";

import { useState } from "react";
import { submitTask, flushData } from "../lib/api";

interface Props {
  onRefresh: () => void;
}

export default function SubmitTask({ onRefresh }: Props) {
  const [id, setId] = useState("my-task-1");
  const [priority, setPriority] = useState(5);
  const [delay, setDelay] = useState(0);
  const [maxRetries, setMaxRetries] = useState(3);
  const [submitting, setSubmitting] = useState(false);
  const [batchLabel, setBatchLabel] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    try {
      await submitTask({
        id: id || `task-${Date.now()}`,
        priority,
        delay,
        max_retries: maxRetries,
      });
      onRefresh();
    } catch {
      // silently handle — backend may be down
    } finally {
      setSubmitting(false);
    }
  }

  async function handleBatch(count: number, label: string) {
    setBatchLabel(label);
    const batchId = `batch-${Date.now()}`;
    try {
      const promises = [];
      for (let i = 0; i < count; i++) {
        promises.push(
          submitTask({
            id: `${batchId}-${i}`,
            priority: Math.floor(Math.random() * 10) + 1,
            delay,
            max_retries: maxRetries,
          })
        );
        // batch in groups of 50 to avoid overwhelming
        if (promises.length >= 50) {
          await Promise.all(promises);
          promises.length = 0;
        }
      }
      if (promises.length > 0) await Promise.all(promises);
      onRefresh();
    } catch {
      // silently handle
    } finally {
      setBatchLabel(null);
    }
  }

  async function handleFlush() {
    if (!confirm("Delete all queues, metrics, events & dead letter data from Redis?")) return;
    try {
      await flushData();
      onRefresh();
    } catch {
      // silently handle
    }
  }

  return (
    <section className="rounded-xl border border-[#30363d] bg-[#161b22] p-6">
      <h2 className="mb-2 text-sm font-semibold uppercase tracking-wider text-[#8b949e]">
        Submit Task
      </h2>
      <p className="mb-4 text-xs text-[#6e7681] leading-relaxed">
        Submit a task to see it flow through the pipeline: submit → queue → worker → outcome.
      </p>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        {/* Task ID */}
        <div>
          <label className="mb-1 block text-xs text-[#8b949e]">
            Task ID (auto-generated if empty)
          </label>
          <input
            type="text"
            value={id}
            onChange={(e) => setId(e.target.value)}
            placeholder="my-task-1"
            className="w-full rounded-lg border border-[#30363d] bg-[#0d1117] px-3 py-2 text-sm text-[#e6edf3] placeholder-[#6e7681] outline-none focus:border-[#58a6ff] transition-colors"
          />
        </div>

        {/* Priority */}
        <div>
          <label className="mb-1 block text-xs text-[#8b949e]">
            Priority (1-10, higher = first)
          </label>
          <input
            type="number"
            min={1}
            max={10}
            value={priority}
            onChange={(e) => setPriority(Number(e.target.value))}
            className="w-full rounded-lg border border-[#30363d] bg-[#0d1117] px-3 py-2 text-sm text-[#e6edf3] outline-none focus:border-[#58a6ff] transition-colors"
          />
        </div>

        {/* Delay */}
        <div>
          <label className="mb-1 block text-xs text-[#8b949e]">
            Delay (seconds)
          </label>
          <input
            type="number"
            min={0}
            value={delay}
            onChange={(e) => setDelay(Number(e.target.value))}
            className="w-full rounded-lg border border-[#30363d] bg-[#0d1117] px-3 py-2 text-sm text-[#e6edf3] outline-none focus:border-[#58a6ff] transition-colors"
          />
        </div>

        {/* Max Retries */}
        <div>
          <label className="mb-1 block text-xs text-[#8b949e]">
            Max Retries
          </label>
          <input
            type="number"
            min={0}
            value={maxRetries}
            onChange={(e) => setMaxRetries(Number(e.target.value))}
            className="w-full rounded-lg border border-[#30363d] bg-[#0d1117] px-3 py-2 text-sm text-[#e6edf3] outline-none focus:border-[#58a6ff] transition-colors"
          />
        </div>

        {/* Submit button */}
        <button
          type="submit"
          disabled={submitting}
          className="flex w-full items-center justify-center gap-2 rounded-lg bg-[#238636] px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-[#2ea043] disabled:opacity-50"
        >
          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5" />
          </svg>
          {submitting ? "Submitting…" : "Submit Task"}
        </button>
      </form>

      {/* Batch submit */}
      <div className="mt-5">
        <p className="mb-2 text-xs text-[#8b949e] font-semibold">Batch Submit</p>
        <div className="flex gap-2">
          {[
            { count: 1000, label: "1k" },
            { count: 5000, label: "5k" },
            { count: 10000, label: "10k" },
          ].map(({ count, label }) => (
            <button
              key={label}
              onClick={() => handleBatch(count, label)}
              disabled={batchLabel !== null}
              className="flex flex-1 items-center justify-center gap-1.5 rounded-lg border border-[#30363d] bg-[#0d1117] px-3 py-2 text-xs font-medium text-[#8b949e] transition-colors hover:border-[#58a6ff] hover:text-[#58a6ff] disabled:opacity-50"
            >
              <svg className="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 12c0-1.232-.046-2.453-.138-3.662a4.006 4.006 0 00-3.7-3.7 48.678 48.678 0 00-7.324 0 4.006 4.006 0 00-3.7 3.7c-.017.22-.032.441-.046.662M19.5 12l3-3m-3 3l-3-3m-12 3c0 1.232.046 2.453.138 3.662a4.006 4.006 0 003.7 3.7 48.656 48.656 0 007.324 0 4.006 4.006 0 003.7-3.7c.017-.22.032-.441.046-.662M4.5 12l3 3m-3-3l-3 3" />
              </svg>
              {batchLabel === label ? "Sending…" : label}
            </button>
          ))}
        </div>
      </div>

      {/* Flush */}
      <button
        onClick={handleFlush}
        className="mt-4 w-full rounded-lg border border-[#f85149]/30 bg-[#f85149]/10 px-4 py-2 text-xs font-semibold text-[#f85149] transition-colors hover:bg-[#f85149]/20"
      >
        🗑 Flush All Data
      </button>
      <p className="mt-1 text-[10px] text-[#6e7681]">
        Deletes all queues, metrics, events & dead letter data from Redis
      </p>
    </section>
  );
}
