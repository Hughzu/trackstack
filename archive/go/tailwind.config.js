/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/**/*.templ", "./internal/**/*_templ.go"],
  safelist: [
    'bg-status-safe',
    'bg-status-warn',
    'bg-status-danger',
    'text-status-safe',
    'text-status-warn',
    'text-status-danger',
    'text-text-muted',
  ],
  theme: {
    extend: {
      colors: {
        // Semantic Status Colors (Traffic Light)
        'status-safe': '#00d08a',      // Green (>20% remaining)
        'status-warn': '#ffb800',      // Yellow (5-20% remaining)
        'status-danger': '#ff3b3f',    // Red (<5% or over)
        
        // Base Colors (Dark Mode)
        'bg-primary': '#0a0a0a',       // Near-black background
        'bg-secondary': '#1a1a1a',     // Slightly lighter background for cards
        'text-primary': '#ffffff',     // White text
        'text-muted': '#94a3b8',       // Slate-400 for labels
        'border-subtle': '#1f2937',    // Gray-800 for borders
        'interactive': '#ffffff',
        'interactive-hover': '#e5e7eb',
      },
      fontSize: {
        // Removed custom dashboard sizes in favor of standard utility classes
      },
      spacing: {
        'dashboard-gap': '32px',
        'safe': 'env(safe-area-inset-bottom)',
      }
    },
  },
  plugins: [],
}
