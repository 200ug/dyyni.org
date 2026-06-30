import js from "@eslint/js"
import prettier from "eslint-config-prettier"
import globals from "globals"

export default [
    js.configs.recommended,
    prettier,
    {
        ignores: ["dist/**", "node_modules/**", ".stormdrain/**"],
        languageOptions: {
            globals: {
                ...globals.browser,
                ENDPOINT_URL: "readonly"
            }
        }
    }
]
