import antfu from "@antfu/eslint-config";

export default antfu(
  {
    vue: true,
    typescript: {
      overrides: {
        "ts/no-explicit-any": "error",
      },

      overridesTypeAware: {
        "ts/no-unsafe-argument": "error",
        "ts/no-unsafe-assignment": "error",
        "ts/no-unsafe-call": "error",
        "ts/no-unsafe-member-access": "error",
        "ts/no-unsafe-type-assertion": "error",
      },
    },

    stylistic: false,
    lessOpinionated: true,
    ignores: ["**/components/ui"],
  },
  {
    rules: {
      "jsonc/sort-keys": "off",
      "perfectionist/sort-imports": "off",
    },
  },
  {
    files: ["**/*.vue"],
    rules: {
      "vue/first-attribute-linebreak": "off",
      "vue/html-closing-bracket-newline": "off",
      "vue/html-indent": "off",
      "vue/html-self-closing": "off",
    },
  },
);
