import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'

export default [
  // add more generic rulesets here, such as:
  // js.configs.recommended,
  ...pluginVue.configs['flat/recommended'],
  // ...pluginVue.configs['flat/vue2-recommended'], // Use this if you are using Vue.js 2.x.
  {
    rules: {
        'vue/multi-word-component-names': 'error',
        'vue/component-name-in-template-casing': ['error', 'PascalCase'],
        'vue/no-v-html': 'warn',
        'vue/require-default-prop': 'error',
        'vue/require-prop-types': 'error',
        'vue/v-bind-style': ['error', 'shorthand'],
        'vue/v-on-style': ['error', 'shorthand'],
        'vue/no-unused-vars': 'error',
        'vue/html-closing-bracket-newline': ['error', {
          'singleline': 'never',
          'multiline': 'always'
        }],
        'vue/html-self-closing': ['error', {
          'html': {
            'void': 'always',
            'normal': 'never',
            'component': 'always'
          }
        }],
        'vue/max-attributes-per-line': ['error', {
          'singleline': 3,
          'multiline': 1
        }],
        'vue/no-multiple-template-root': 'off', // For Vue 3
    
        // TypeScript rules
        '@typescript-eslint/explicit-module-boundary-types': 'off',
        '@typescript-eslint/no-explicit-any': 'warn',
        '@typescript-eslint/no-unused-vars': ['error', {
          'argsIgnorePattern': '^_',
          'varsIgnorePattern': '^_'
        }],
    
        // General rules
        'eqeqeq': ['error', 'always'],
        'curly': ['error', 'all'],
        'prefer-const': 'error',
        'no-var': 'error',
        'object-shorthand': 'error',
        'array-callback-return': 'error',
        'prefer-template': 'error',
        'quote-props': ['error', 'as-needed'],
        'vue/script-indent': ['error', 'tab', {
            'baseIndent': 1,  // This enforces 1 tab indentation inside script tags
            'switchCase': 1,
            'ignores': []
        }],
    },
    languageOptions: {
      sourceType: 'module',
      globals: {
        ...globals.browser
      }
    }
  }
]