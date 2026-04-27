"use client";

import { useCallback } from "react";
import {
  getMetrics,
  getEvents,
  getWorkers,
  getQueues,
  getFailedTasks,
} from "./lib/api";
import { usePolling } from "./hooks/usePolling";

import TaskFlowPipeline from "./components/TaskFlowPipeline";
import SubmitTask from "./components/SubmitTask";
import MetricsGrid from "./components/MetricsGrid";
import WorkerPool from "./components/WorkerPool";
import QueueContents from "./components/QueueContents";
import ActivityLog from "./components/ActivityLog";
import FailedTasks from "./components/FailedTasks";

export default function Dashboard() {
  const metrics = usePolling(useCallback(() => getMetrics(), []), 2000);
  const events = usePolling(useCallback(() => getEvents(80), []), 2000);
  const workers = usePolling(useCallback(() => getWorkers(), []), 2000);
  const queues = usePolling(useCallback(() => getQueues(), []), 2000);
  const failed = usePolling(useCallback(() => getFailedTasks(), []), 2000);

  function refreshAll() {
    metrics.refresh();
    events.refresh();
    workers.refresh();
    queues.refresh();
    failed.refresh();
  }

  return (
    <div className="mx-auto max-w-[1400px] px-4 py-6 sm:px-6 lg:px-8">
      {/* Header */}
      <header className="mb-6">
        <h1 className="text-2xl font-bold text-[#e6edf3]">
          Task Queue Dashboard
        </h1>
        <p className="text-sm text-[#8b949e]">
          An educational view into distributed task processing
        </p>
      </header>

      {/* Pipeline */}
      <div className="mb-6">
        <TaskFlowPipeline metrics={metrics.data} />
      </div>

      {/* Submit + Metrics side by side */}
      <div className="mb-6 grid grid-cols-1 gap-6 lg:grid-cols-[320px_1fr]">
        <SubmitTask onRefresh={refreshAll} />
        <MetricsGrid metrics={metrics.data} />
      </div>

      {/* Worker Pool */}
      <div className="mb-6">
        <WorkerPool workers={workers.data} />
      </div>

      {/* Queue Contents + Activity Log side by side */}
      <div className="mb-6 grid grid-cols-1 gap-6 lg:grid-cols-[280px_1fr]">
        <QueueContents queues={queues.data} />
        <ActivityLog events={events.data} />
      </div>

      {/* Failed Tasks */}
      <div className="mb-6">
        <FailedTasks tasks={failed.data} />
      </div>
    </div>
  );
}
