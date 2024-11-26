import vue from "eslint-plugin-vue";
import typescriptEslint from "@typescript-eslint/eslint-plugin";
import tailwindcss from "eslint-plugin-tailwindcss";
import parser from "vue-eslint-parser";
import path from "node:path";
import { fileURLToPath } from "node:url";
import js from "@eslint/js";
import { FlatCompat } from "@eslint/eslintrc";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
    allConfig: js.configs.all
});

export default [...compat.extends(
    "eslint:recommended",
    "plugin:vue/vue3-essential",
    "plugin:@typescript-eslint/recommended",
), {
    plugins: {
        vue,
        "@typescript-eslint": typescriptEslint,
        tailwindcss,
    },

    languageOptions: {
        globals: {
            NodeJS: true,
        },

        parser: parser,
        ecmaVersion: "latest",
        sourceType: "script",

        parserOptions: {
            parser: "@typescript-eslint/parser",
            extraFileExtensions: [".vue"],
        },
    },

    settings: {
        tailwindcss: {
            config: "./tailwind.config.js",
        },
    },

    rules: {
        "tailwindcss/classnames-order": "warn",
        "tailwindcss/enforces-negative-arbitrary-values": "warn",
        "tailwindcss/enforces-shorthand": "warn",
        "tailwindcss/no-contradicting-classname": "warn",
    },
}, {
    files: ["**/*.vue"],

    languageOptions: {
        parser: parser,
    },
}];