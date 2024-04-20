/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./ui/src/**/*.tsx"],
    darkMode: "class",

    plugins: [
        require("@tailwindcss/forms")({
            strategy: "base", // only generate global styles
        }),
    ],
}
