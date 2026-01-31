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
    'text-dashboard-hero',
    'text-dashboard-secondary',
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
        'text-primary': '#ffffff',     // White text
        'text-muted': '#94a3b8',       // Slate-400 for labels
        'border-subtle': '#1f2937',    // Gray-800 for borders
      },
      fontSize: {
        'dashboard-hero': ['80px', { lineHeight: '1.0', fontWeight: '700' }],
        'dashboard-secondary': ['60px', { lineHeight: '1.0', fontWeight: '700' }],
      },
      spacing: {
        'dashboard-gap': '32px',
      }
    },
  },
  plugins: [],
}
