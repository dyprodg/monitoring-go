import { useState, useEffect } from 'react';
import MetricCard from './MetricCard';
import ActionButton from './ActionButton';
import { fetchMetrics, startCPUStress } from '../services/api';

const POLLING_INTERVAL = 1000; // 1 second
const MAX_HISTORY_LENGTH = 60; // Keep 60 seconds of history

export default function Dashboard() {
  const [metrics, setMetrics] = useState(null);
  const [cpuHistory, setCpuHistory] = useState([]);
  const [memoryHistory, setMemoryHistory] = useState([]);
  const [error, setError] = useState(null);

  // Poll metrics every second
  useEffect(() => {
    const pollMetrics = async () => {
      try {
        const data = await fetchMetrics();
        setMetrics(data);
        setError(null);

        // Update CPU history
        setCpuHistory(prev => {
          const newHistory = [
            ...prev,
            {
              timestamp: new Date(data.timestamp).toLocaleTimeString(),
              value: data.cpu
            }
          ];
          // Keep only last 60 data points
          return newHistory.slice(-MAX_HISTORY_LENGTH);
        });

        // Update Memory history
        setMemoryHistory(prev => {
          const newHistory = [
            ...prev,
            {
              timestamp: new Date(data.timestamp).toLocaleTimeString(),
              value: data.memory
            }
          ];
          return newHistory.slice(-MAX_HISTORY_LENGTH);
        });

      } catch (err) {
        setError(err.message);
        console.error('Failed to poll metrics:', err);
      }
    };

    // Initial fetch
    pollMetrics();

    // Set up polling interval
    const intervalId = setInterval(pollMetrics, POLLING_INTERVAL);

    // Cleanup on unmount
    return () => clearInterval(intervalId);
  }, []);

  // Handle CPU stress button click
  const handleCPUStress = async () => {
    try {
      const response = await startCPUStress(80, 10);
      console.log('CPU stress started:', response);
    } catch (err) {
      console.error('Failed to start CPU stress:', err);
      alert(`Failed to start CPU stress: ${err.message}`);
    }
  };

  return (
    <div className="min-h-screen bg-slate-900 p-8">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <header className="mb-8">
          <h1 className="text-4xl font-bold text-white mb-2">
            Interactive System Monitoring Dashboard
          </h1>
          <p className="text-slate-400">
            Click buttons to generate load and watch metrics spike in real-time
          </p>
        </header>

        {/* Error Message */}
        {error && (
          <div className="mb-6 bg-red-900/50 border border-red-700 text-red-200 px-4 py-3 rounded-lg">
            <strong>Error:</strong> {error}
          </div>
        )}

        {/* Metrics Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <MetricCard
            title="CPU Usage"
            value={metrics?.cpu || 0}
            unit="%"
            history={cpuHistory}
            color="#3b82f6"
          />
          <MetricCard
            title="Memory Usage"
            value={metrics?.memory || 0}
            unit="%"
            history={memoryHistory}
            color="#10b981"
          />
        </div>

        {/* Actions */}
        <div className="bg-slate-800 rounded-lg p-6 shadow-lg border border-slate-700">
          <h2 className="text-2xl font-semibold text-white mb-4">
            Load Actions
          </h2>
          <div className="flex flex-wrap gap-4">
            <ActionButton
              label="CPU Stress"
              emoji="🔥"
              onClick={handleCPUStress}
              color="red"
            />
          </div>
        </div>

        {/* Status Footer */}
        <footer className="mt-8 text-center text-slate-500 text-sm">
          <p>
            Polling every {POLLING_INTERVAL / 1000}s •
            {metrics && ` Last update: ${new Date(metrics.timestamp).toLocaleTimeString()}`}
          </p>
        </footer>
      </div>
    </div>
  );
}
