import { preset } from "@rebass/preset";

const theme = {
  ...preset,
  colors: {
    ...preset.colors,
    primary: "#34C363",
    dark: "#2b2b2b",
    lightGrey: "#6d6d6d",
    fadeGrey: "#7a7a7a",
  },
  fonts: {
    body: "system-ui, sans-serif",
    heading: "system-ui, sans-serif",
    monospace: "Menlo, monospace",
  },
  fontSizes: [12, 14, 16, 20, 24, 32, 48, 64],
  fontWeights: {
    body: 400,
    heading: 700,
    bold: 700,
  },
  space: [0, 4, 8, 16, 32, 64, 128, 256, 512],
  radii: {
    default: 4,
    small: 2,
    medium: 6,
    large: 8,
  },
  buttons: {
    primary: {
      color: "white",
      bg: "primary",
      "&:hover": {
        bg: "primary",
        opacity: 0.9,
      },
    },
    secondary: {
      color: "primary",
      bg: "secondary",
      "&:hover": {
        bg: "secondary",
        opacity: 0.8,
      },
    },
    outline: {
      color: "primary",
      bg: "transparent",
      border: "1px solid",
      borderColor: "primary",
      "&:hover": {
        bg: "#fff",
      },
    },
  },
};

export default theme;
