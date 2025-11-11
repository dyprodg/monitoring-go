import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

/**
 * MetricCard displays a single metric with a line chart
 * @param {string} title - Metric title (e.g., "CPU Usage")
 * @param {number} value - Current value
 * @param {string} unit - Unit of measurement (e.g., "%")
 * @param {Array} history - Historical data points for chart
 * @param {string} color - Line color
 * @param {boolean} isActive - Whether this metric is currently under load
 */
export default function MetricCard({ title, value, unit, history, color = "#3b82f6", isActive = false }) {
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
    <div className={`
      relative bg-slate-800 rounded-xl p-5 border transition-all duration-300
      ${isActive
        ? 'border-orange-500 shadow-lg shadow-orange-500/20'
        : 'border-slate-700 hover:border-slate-600'
      }
    `}>
      {/* Header */}
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <h3 className="text-xs font-semibold text-slate-400 uppercase tracking-wide">
              {title}
            </h3>
            {isActive && (
              <span className="flex items-center gap-1 text-[10px] text-orange-400 font-bold px-1.5 py-0.5 rounded bg-orange-500/10 border border-orange-500/30">
                <span className="relative flex h-1.5 w-1.5">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-orange-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-1.5 w-1.5 bg-orange-500"></span>
                </span>
                LIVE
              </span>
            )}
          </div>
          <div className="flex items-baseline gap-1.5">
            <span className="text-3xl font-bold text-white tabular-nums">
              {formattedValue}
            </span>
            <span className="text-base font-medium text-slate-400">
              {unit}
            </span>
          </div>
        </div>
      </div>

      {/* Chart */}
      <div className="h-24 -mx-2">
        {history && history.length > 0 ? (
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={history}>
              <defs>
                <linearGradient id={`area-${title}`} x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor={color} stopOpacity={0.4}/>
                  <stop offset="95%" stopColor={color} stopOpacity={0.05}/>
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#334155" vertical={false} strokeOpacity={0.3} />
              <XAxis
                dataKey="timestamp"
                hide={true}
              />
              <YAxis
                domain={yDomain}
                hide={true}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: '#1e293b',
                  border: '1px solid #475569',
                  borderRadius: '0.5rem',
                  padding: '8px 12px'
                }}
                labelStyle={{ color: '#94a3b8', fontSize: '12px' }}
                itemStyle={{ color: color, fontSize: '14px', fontWeight: 'bold' }}
                formatter={(value) => [`${value.toFixed(1)} ${unit}`]}
              />
              <Area
                type="monotone"
                dataKey="value"
                stroke={color}
                strokeWidth={2}
                fill={`url(#area-${title})`}
                isAnimationActive={true}
                animationDuration={300}
              />
            </AreaChart>
          </ResponsiveContainer>
        ) : (
          <div className="flex items-center justify-center h-full text-slate-600 text-sm">
            Waiting for data...
          </div>
        )}
      </div>
    </div>
  );
}
