// API client for monitoring dashboard backend
const API_BASE_URL = 'http://localhost:8080/api';

/**
 * Fetch current system metrics
 * @returns {Promise<Object>} Metrics data
 */
export async function fetchMetrics() {
  try {
    const response = await fetch(`${API_BASE_URL}/metrics`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch metrics:', error);
    throw error;
  }
}

/**
 * Start a CPU stress action
 * @param {number} targetPercent - Target CPU percentage (0-95)
 * @param {number} durationSeconds - Duration in seconds (1-30)
 * @returns {Promise<Object>} Action response
 */
export async function startCPUStress(targetPercent, durationSeconds) {
  try {
    const response = await fetch(`${API_BASE_URL}/actions/cpu-stress`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        target_percent: targetPercent,
        duration_seconds: durationSeconds,
      }),
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to start CPU stress:', error);
    throw error;
  }
}

/**
 * Start a memory surge action
 * @param {number} sizeMB - Memory size in MB
 * @param {number} durationSeconds - Duration in seconds (1-60)
 * @returns {Promise<Object>} Action response
 */
export async function startMemorySurge(sizeMB, durationSeconds) {
  try {
    const response = await fetch(`${API_BASE_URL}/actions/memory-surge`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        size_mb: sizeMB,
        duration_seconds: durationSeconds,
      }),
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to start memory surge:', error);
    throw error;
  }
}

/**
 * Start a disk storm action
 * @param {number} operations - Number of file operations
 * @param {number} fileSizeKB - File size in KB
 * @returns {Promise<Object>} Action response
 */
export async function startDiskStorm(operations, fileSizeKB) {
  try {
    const response = await fetch(`${API_BASE_URL}/actions/disk-storm`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        operations: operations,
        file_size_kb: fileSizeKB,
      }),
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to start disk storm:', error);
    throw error;
  }
}

/**
 * Start a traffic flood action
 * @param {number} requestsPerSec - Requests per second
 * @param {number} durationSeconds - Duration in seconds (1-60)
 * @returns {Promise<Object>} Action response
 */
export async function startTrafficFlood(requestsPerSec, durationSeconds) {
  try {
    const response = await fetch(`${API_BASE_URL}/actions/traffic-flood`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        requests_per_sec: requestsPerSec,
        duration_seconds: durationSeconds,
        target_url: '', // Will use default
      }),
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to start traffic flood:', error);
    throw error;
  }
}

/**
 * Get all active actions
 * @returns {Promise<Object>} Active actions
 */
export async function fetchActiveActions() {
  try {
    const response = await fetch(`${API_BASE_URL}/actions/active`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch active actions:', error);
    throw error;
  }
}

/**
 * Stop a running action
 * @param {string} actionId - Action ID to stop
 * @returns {Promise<Object>} Stop response
 */
export async function stopAction(actionId) {
  try {
    const response = await fetch(`${API_BASE_URL}/actions/${actionId}/stop`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to stop action:', error);
    throw error;
  }
}

/**
 * Stop all running actions
 * @returns {Promise<Object>} Stop all response
 */
export async function stopAllActions() {
  try {
    const response = await fetch(`${API_BASE_URL}/actions/stop-all`, {
      method: 'POST',
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to stop all actions:', error);
    throw error;
  }
}

/**
 * Check backend health
 * @returns {Promise<Object>} Health status
 */
export async function fetchHealth() {
  try {
    const response = await fetch(`${API_BASE_URL}/health`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch health:', error);
    throw error;
  }
}
