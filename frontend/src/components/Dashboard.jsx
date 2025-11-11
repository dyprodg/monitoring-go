import { useState, useEffect } from 'react';
import MetricCard from './MetricCard';
import ActionButton from './ActionButton';
import {
  fetchMetrics,
  fetchActiveActions,
  startCPUStress,
  startMemorySurge,
  startDiskStorm,
  startTrafficFlood,
  stopAllActions
} from '../services/api';

const POLLING_INTERVAL = 1000; // 1 second
const MAX_HISTORY_LENGTH = 60; // Keep 60 seconds of history

export default function Dashboard() {
  const [metrics, setMetrics] = useState(null);
  const [cpuHistory, setCpuHistory] = useState([]);
  const [memoryHistory, setMemoryHistory] = useState([]);
  const [diskHistory, setDiskHistory] = useState([]);
  const [networkHistory, setNetworkHistory] = useState([]);
  const [activeActions, setActiveActions] = useState([]);
  const [error, setError] = useState(null);
  const [stoppingAll, setStoppingAll] = useState(false);

  // Poll metrics and active actions every second
  useEffect(() => {
    const pollData = async () => {
      try {
        // Fetch metrics
        const data = await fetchMetrics();
        setMetrics(data);

        // Fetch active actions
        const actionsData = await fetchActiveActions();
        setActiveActions(actionsData.actions || []);

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

        // Update Disk I/O history
        setDiskHistory(prev => {
          const newHistory = [
            ...prev,
            {
              timestamp: new Date(data.timestamp).toLocaleTimeString(),
              value: data.disk_io
            }
          ];
          return newHistory.slice(-MAX_HISTORY_LENGTH);
        });

        // Update Network history
        setNetworkHistory(prev => {
          const newHistory = [
            ...prev,
            {
              timestamp: new Date(data.timestamp).toLocaleTimeString(),
              value: data.network
            }
          ];
          return newHistory.slice(-MAX_HISTORY_LENGTH);
        });

      } catch (err) {
        setError(err.message);
        console.error('Failed to poll data:', err);
      }
    };

    // Initial fetch
    pollData();

    // Set up polling interval
    const intervalId = setInterval(pollData, POLLING_INTERVAL);

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

  // Handle Memory surge button click
  const handleMemorySurge = async () => {
    try {
      const response = await startMemorySurge(500, 30);
      console.log('Memory surge started:', response);
    } catch (err) {
      console.error('Failed to start memory surge:', err);
      alert(`Failed to start memory surge: ${err.message}`);
    }
  };

  // Handle Disk storm button click
  const handleDiskStorm = async () => {
    try {
      const response = await startDiskStorm(1000, 10);
      console.log('Disk storm started:', response);
    } catch (err) {
      console.error('Failed to start disk storm:', err);
      alert(`Failed to start disk storm: ${err.message}`);
    }
  };

  // Handle Traffic flood button click
  const handleTrafficFlood = async () => {
    try {
      const response = await startTrafficFlood(100, 10);
      console.log('Traffic flood started:', response);
    } catch (err) {
      console.error('Failed to start traffic flood:', err);
      alert(`Failed to start traffic flood: ${err.message}`);
    }
  };

  // Handle Stop All button click
  const handleStopAll = async () => {
    setStoppingAll(true);
    try {
      const response = await stopAllActions();
      console.log('Stopped all actions:', response);
    } catch (err) {
      console.error('Failed to stop all actions:', err);
      alert(`Failed to stop all actions: ${err.message}`);
    } finally {
      setStoppingAll(false);
    }
  };

  // Helper function to check if an action type is running
  const isActionRunning = (actionType) => {
    return activeActions.some(action => action.type === actionType);
  };

  return (
    <div className="min-h-screen bg-slate-900">
      <div className="max-w-7xl mx-auto p-4 sm:p-6 lg:p-8">
        {/* Header */}
        <header className="mb-8">
          <div className="flex items-center justify-between flex-wrap gap-4">
            <div>
              <h1 className="text-3xl sm:text-4xl font-bold text-white mb-2">
                System Monitor
              </h1>
              <p className="text-slate-400 text-sm sm:text-base">
                Real-time system metrics and load testing
              </p>
            </div>

            {/* Active Actions Indicator */}
            {activeActions.length > 0 && (
              <div className="flex items-center gap-2 bg-blue-500/10 border border-blue-500/30 text-blue-400 px-4 py-2 rounded-lg">
                <span className="relative flex h-2 w-2">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-400"></span>
                </span>
                <span className="text-sm font-semibold">{activeActions.length} Active</span>
              </div>
            )}
          </div>
        </header>

        {/* Error Message */}
        {error && (
          <div className="mb-6 bg-red-900/50 border border-red-500 text-red-200 px-4 py-3 rounded-lg flex items-center gap-3">
            <svg className="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <span className="font-medium">{error}</span>
          </div>
        )}

        {/* Metrics Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <MetricCard
            title="CPU Usage"
            value={metrics?.cpu || 0}
            unit="%"
            history={cpuHistory}
            color="#3b82f6"
            isActive={isActionRunning('cpu-stress')}
          />
          <MetricCard
            title="Memory Usage"
            value={metrics?.memory || 0}
            unit="%"
            history={memoryHistory}
            color="#10b981"
            isActive={isActionRunning('memory-surge')}
          />
          <MetricCard
            title="Disk I/O"
            value={metrics?.disk_io || 0}
            unit="ops/s"
            history={diskHistory}
            color="#f59e0b"
            isActive={isActionRunning('disk-storm')}
          />
          <MetricCard
            title="Network"
            value={metrics?.network || 0}
            unit="MB/s"
            history={networkHistory}
            color="#8b5cf6"
            isActive={isActionRunning('traffic-flood')}
          />
        </div>

        {/* Actions Section */}
        <div className="bg-slate-800 rounded-xl p-6 border border-slate-700">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-bold text-white">Load Tests</h2>

            {/* Stop All Button */}
            {activeActions.length > 0 && (
              <button
                onClick={handleStopAll}
                disabled={stoppingAll}
                className="bg-red-600 hover:bg-red-700 disabled:opacity-50
                         px-4 py-2 rounded-lg text-sm font-semibold transition-colors
                         flex items-center gap-2 text-white"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
                </svg>
                {stoppingAll ? 'Stopping...' : 'Stop All'}
              </button>
            )}
          </div>

          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <ActionButton
              label="CPU Stress"
              emoji="🔥"
              onClick={handleCPUStress}
              color="red"
              isRunning={isActionRunning('cpu-stress')}
            />
            <ActionButton
              label="Memory Surge"
              emoji="💾"
              onClick={handleMemorySurge}
              color="green"
              isRunning={isActionRunning('memory-surge')}
            />
            <ActionButton
              label="Disk Storm"
              emoji="💿"
              onClick={handleDiskStorm}
              color="yellow"
              isRunning={isActionRunning('disk-storm')}
            />
            <ActionButton
              label="Traffic Flood"
              emoji="🌐"
              onClick={handleTrafficFlood}
              color="purple"
              isRunning={isActionRunning('traffic-flood')}
            />
          </div>
        </div>

        {/* Status Footer */}
        <footer className="mt-8 flex items-center justify-center gap-4 text-sm text-slate-500">
          <div className="flex items-center gap-2">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            <span>Live</span>
          </div>
          <span>•</span>
          <span>Updates every {POLLING_INTERVAL / 1000}s</span>
          {metrics && (
            <>
              <span className="hidden sm:inline">•</span>
              <span className="hidden sm:inline">{new Date(metrics.timestamp).toLocaleTimeString()}</span>
            </>
          )}
        </footer>
      </div>
    </div>
  );
}
