/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*_templ.go"],
  blocklist: ["./vendor/**"],
  safelist: [],
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["cupcake"],
  },
};
