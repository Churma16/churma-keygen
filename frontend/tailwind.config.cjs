/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{svelte,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
      },
      colors: {
        forest: {
          dark: '#003D29',
          light: '#0A4D38',
          accent: '#E8F5F1',
        },
        coral: {
          DEFAULT: '#FF6B5A',
          soft: '#FEE2E2',
        }
      }
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        light: {
          ...require("daisyui/src/theming/themes")["light"],
          primary: "#003D29",      // Forest Green
          secondary: "#FF6B5A",    // Coral/Salmon Red
          accent: "#0D9488",       // Teal Accent
          neutral: "#1e293b",      // Slate
          "base-100": "#F8F9FA",   // Off-white canvas background
          "base-200": "#FFFFFF",   // Pure white container background
          "base-300": "#F1F5F9",   // Border/divider line
          info: "#0284c7",
          success: "#E8F5F1",      // Pastel mint background
          warning: "#f59e0b",
          error: "#FEE2E2",        // Pastel rose/coral background
        },
      },
    ],
    darkTheme: "light", // Disable dark theme fallback, keep it bright
  },
}
