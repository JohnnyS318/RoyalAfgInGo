const colors = require("tailwindcss/colors");

module.exports = {
    future: {
        removeDeprecatedGapUtilities: true,
        purgeLayersByDefault: true
    },
    purge: [],
    theme: {
        fontFamily: { sans: ["Poppins", "sans-serif"] },
        colors: {
            transparent: "transparent",
            current: "currentColor",
            black: colors.black,
            white: colors.white,
            gray: {
                50: "#fafafa",
                100: "#f5f5f5",
                200: "#efefef",
                300: "#e0e0e0",
                400: "#d0d0d0",
                500: "#737373",
                600: "#525252",
                700: "#404040",
                800: "#262626",
                900: "#171717"
            },
            indigo: colors.indigo,
            red: colors.rose,
            blue: colors.blue,
            yellow: colors.amber
        },
        extend: {
            width: {
                fit: "fit-content",
                min: "min-content"
            },
            minWidth: {
                36: "9rem",
                40: "10rem",
                44: "11rem",
                48: "12rem",
                52: "13rem",
                56: "14rem",
                60: "15rem"
            }
        }
    },
    variants: {},
    plugins: [require("@tailwindcss/custom-forms")]
};
