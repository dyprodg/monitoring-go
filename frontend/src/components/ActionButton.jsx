import { useState } from 'react';

/**
 * ActionButton triggers a load generation action
 * @param {string} label - Button label
 * @param {string} emoji - Emoji icon
 * @param {Function} onClick - Click handler
 * @param {boolean} disabled - Whether button is disabled
 * @param {string} color - Button color theme
 * @param {boolean} isRunning - Whether this action is currently running
 */
export default function ActionButton({
  label,
  emoji,
  onClick,
  disabled = false,
  color = "blue",
  isRunning = false
}) {
  const [isLoading, setIsLoading] = useState(false);

  const colorClasses = {
    blue: "bg-blue-600 hover:bg-blue-700 focus:ring-blue-500",
    red: "bg-red-600 hover:bg-red-700 focus:ring-red-500",
    green: "bg-green-600 hover:bg-green-700 focus:ring-green-500",
    yellow: "bg-yellow-600 hover:bg-yellow-700 focus:ring-yellow-500",
    purple: "bg-purple-600 hover:bg-purple-700 focus:ring-purple-500",
  };

  const ringClasses = {
    blue: "ring-blue-400",
    red: "ring-red-400",
    green: "ring-green-400",
    yellow: "ring-yellow-400",
    purple: "ring-purple-400",
  };

  const handleClick = async () => {
    if (disabled || isLoading) return;

    setIsLoading(true);
    try {
      await onClick();
    } catch (error) {
      console.error('Action failed:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <button
      onClick={handleClick}
      disabled={disabled || isLoading}
      className={`
        relative
        ${colorClasses[color] || colorClasses.blue}
        ${isRunning ? `ring-2 ${ringClasses[color]}` : ''}
        disabled:opacity-50 disabled:cursor-not-allowed
        px-4 py-6 rounded-lg
        transition-all duration-200
        focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-900
        flex flex-col items-center justify-center gap-2
        text-white hover:brightness-110 active:scale-95
        w-full
        shadow-lg
      `}
      aria-label={label}
    >
      {/* Running Indicator */}
      {isRunning && (
        <span className="absolute top-2 right-2">
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-white opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2 w-2 bg-white"></span>
          </span>
        </span>
      )}

      <span className="text-4xl">{emoji}</span>
      <span className="text-sm font-bold">
        {isLoading ? 'Starting...' : isRunning ? 'Running...' : label}
      </span>
    </button>
  );
}
