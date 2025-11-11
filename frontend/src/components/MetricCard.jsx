import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

/**
 * MetricCard displays a single metric with a line chart
 * @param {string} title - Metric title (e.g., "CPU Usage")
 * @param {number} value - Current value
 * @param {string} unit - Unit of measurement (e.g., "%")
 * @param {Array} history - Historical data points for chart
 * @param {string} color - Line color
 */
export default function MetricCard({ title, value, unit, history, color = "#3b82f6" }) {
  // Format value to 1 decimal place
  const formattedValue = typeof value === 'number' ? value.toFixed(1) : '0.0';

  // Calculate dynamic Y-axis domain based on data
  const calculateYDomain = () => {
    if (!history || history.length === 0) {
      return [0, 100]; // Default for percentage metrics
    }

    // For percentage metrics (CPU, Memory), use fixed 0-100 scale
    if (unit === '%') {
      return [0, 100];
    }

    // For other metrics (Disk I/O, Network), use dynamic scale
    const values = history.map(item => item.value);
    const maxValue = Math.max(...values);
    const minValue = Math.min(...values);

    // Add 10% padding to top and bottom for better visualization
    const padding = (maxValue - minValue) * 0.1 || 10;
    const yMin = Math.max(0, minValue - padding);
    const yMax = maxValue + padding;

    return [yMin, yMax];
  };

  const yDomain = calculateYDomain();

  return (
    <div className="bg-slate-800 rounded-lg p-6 shadow-lg border border-slate-700">
      {/* Header */}
      <div className="mb-4">
        <h3 className="text-sm font-medium text-slate-400 uppercase tracking-wide">
          {title}
        </h3>
        <div className="mt-2 flex items-baseline">
          <p className="text-4xl font-semibold text-white">
            {formattedValue}
          </p>
          <p className="ml-2 text-xl text-slate-400">
            {unit}
          </p>
        </div>
      </div>

      {/* Chart */}
      <div className="h-32">
        {history && history.length > 0 ? (
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={history}>
              <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
              <XAxis
                dataKey="timestamp"
                tick={{ fill: '#64748b', fontSize: 10 }}
                stroke="#475569"
                tickCount={5}
              />
              <YAxis
                domain={yDomain}
                tick={{ fill: '#64748b', fontSize: 10 }}
                stroke="#475569"
                tickFormatter={(value) => value.toFixed(0)}
                width={40}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: '#1e293b',
                  border: '1px solid #475569',
                  borderRadius: '0.375rem'
                }}
                labelStyle={{ color: '#94a3b8' }}
                formatter={(value) => [`${value.toFixed(1)} ${unit}`, title]}
              />
              <Line
                type="monotone"
                dataKey="value"
                stroke={color}
                strokeWidth={2}
                dot={false}
                isAnimationActive={false}
              />
            </LineChart>
          </ResponsiveContainer>
        ) : (
          <div className="flex items-center justify-center h-full text-slate-500">
            No data
          </div>
        )}
      </div>
    </div>
  );
}
