import { useState } from 'react';

/**
 * ActionButton triggers a load generation action
 * @param {string} label - Button label
 * @param {string} emoji - Emoji icon
 * @param {Function} onClick - Click handler
 * @param {boolean} disabled - Whether button is disabled
 * @param {string} color - Button color theme
 */
export default function ActionButton({
  label,
  emoji,
  onClick,
  disabled = false,
  color = "blue"
}) {
  const [isLoading, setIsLoading] = useState(false);

  const colorClasses = {
    blue: "bg-blue-600 hover:bg-blue-700 focus:ring-blue-500",
    red: "bg-red-600 hover:bg-red-700 focus:ring-red-500",
    green: "bg-green-600 hover:bg-green-700 focus:ring-green-500",
    purple: "bg-purple-600 hover:bg-purple-700 focus:ring-purple-500",
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
        ${colorClasses[color] || colorClasses.blue}
        disabled:opacity-50 disabled:cursor-not-allowed
        px-6 py-3 rounded-lg font-medium
        transition-all duration-200
        focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-900
        flex items-center gap-2
        text-white shadow-lg
      `}
      aria-label={label}
    >
      <span className="text-xl">{emoji}</span>
      <span>{isLoading ? 'Starting...' : label}</span>
    </button>
  );
}
