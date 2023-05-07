/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{gohtml,js}"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
    require('@tailwindcss/aspect-ratio'),
    require('@tailwindcss/container-queries'),
  ],
}