/** @type {import('tailwindcss').Config} */
module.exports = {
  variants: {
    extend: {
      textColor: ['hover', 'group-hover'],
    },
  },
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx,vue}'],
  theme: {
    extend: {
      /*

      These aren't part of the Nord palette, but they are the bg colors
      used in the Nord website.
      #f2f4f8 - light bg
      #242933 - dark bg
      
      Nord - published palette
      
      Polar Night: Backgrounds (dark to light
      #2e3440
      #3b4252
      #434c5e
      #4c566a

      Snow Storm: Bright text/ui colours (dark to light)
      #d8dee9
      #e5e9f0
      #eceff4

      Frost: Bluish colours for primary and secondary UI components (light to dark)
      #8fbcbb
      #88c0d0
      #81a1c1
      #5e81ac

      Aurora: Accent colours (light to dark)
      #bf616a red
      #d08770 orange
      #ebcb8b yellow
      #a3be8c green
      #b48ead pink

      */
      colors: {
        reaper: {
          'bg-dark': '#242933',
          'bg-light': '#f2f4f8',
        },
        'polar-night': {
          DEFAULT: '#2e3440',
          1: '#2e3440',
          '1a': '#353b47',
          2: '#3b4252',
          3: '#434c5e',
          4: '#4c566a',
        },

        'snow-storm': {
          DEFAULT: '#eceff4',
          1: '#d8dee9',
          2: '#e5e9f0',
          3: '#eceff4',
        },

        frost: {
          DEFAULT: '#8fbcbb',
          1: '#8fbcbb',
          2: '#88c0d0',
          3: '#81a1c1',
          4: '#5e81ac',
        },

        aurora: {
          DEFAULT: '#bf616a',
          1: '#bf616a',
          2: '#d08770',
          3: '#ebcb8b',
          4: '#a3be8c',
          5: '#b48ead',
        },
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
}
