/**
 * Theme configuration store
 */
import { Ref, computed, ref, onMounted } from 'vue'
import { themes } from '@/lib/registry/themes'

interface ConfigStore {
  // theme
  isDark: Ref<boolean>
  theme: Ref<string>
  themeColors: string[]
  themePrimary: Ref<string>
  themeSecondary: Ref<string>
  themeBackground: Ref<string>
  // helpers
  setThemeColor: (color: string) => void
  setThemeData: () => void
  toggleThemeDark: () => void
}

interface Theme {
  name: string
  cssVars: {
    dark: {
      primary: string
      secondary: string
      background: string
    }
    light: {
      primary: string
      secondary: string
      background: string
    }
  }

}

export function useConfigStore(): ConfigStore {
  const isDark = ref(true)
  const theme = ref('')
  const themeColors = [
    'zinc',
    'rose',
    'blue',
    'green',
    'orange',
    'red',
    'violet',
    'slate',
    'stone',
    'gray',
    'neutral',
    'yellow',
  ]

  // toggl dark mode
  const toggleThemeDark = () => {
    console.log('config store toggle dark')
    isDark.value = !isDark.value
    document.documentElement.classList.toggle('dark')
    localStorage.setItem('ghost.reaper.dark', isDark.value.toString())
  }

  // set dark mode
  const setThemeDark = () => {
    isDark.value = true
    document.documentElement.classList.add('dark')
    localStorage.setItem('ghost.reaper.dark', isDark.value.toString())
  }

  // set light mode
  const setThemeLight = () => {
    isDark.value = false
    document.documentElement.classList.remove('dark')
    localStorage.setItem('ghost.reaper.dark', isDark.value.toString())
  }

  // set theme color
  const setThemeColor = (color: string) => {
    // valid color
    if (!themeColors.includes(color)) {
      return
    }

    // remove any existing theme color
    document.documentElement.classList.remove(...themeColors.map((color) => `theme-${color}`))

    // add new theme color
    document.documentElement.classList.add(`theme-${color}`)

    theme.value = color

    // write it to localStorage
    localStorage.setItem('ghost.reaper.theme', color)
  }

  const setThemeData = () => {
    // set dark mode
    const dark = localStorage.getItem('ghost.reaper.dark')
    if (dark == 'true') {
      setThemeDark()
    } else {
      setThemeLight()
    }

    // set theme color
    const theme = localStorage.getItem('ghost.reaper.theme')
    if (theme) {
      setThemeColor(theme)
    }
  }

  // get theme colors for charts
  // TODO: refactor to a combined function/map and make it reactive, so color updates on theme change
  const themePrimary = computed(() => {
    const t = themes.find((t: Theme) => t.name === theme.value)
    return `hsl(${t?.cssVars[isDark.value ? 'dark' : 'light'].primary})`
  })
  const themeSecondary = computed(() => {
    const t = themes.find((t: Theme) => t.name === theme.value)
    return `hsl(${t?.cssVars[isDark.value ? 'dark' : 'light'].secondary})`
  })
  const themeBackground = computed(() => {
    const t = themes.find((t: Theme) => t.name === theme.value)
    return `hsl(${t?.cssVars[isDark.value ? 'dark' : 'light'].background})`
  })

  // TODO: remove this; replace with a function and call it directly from components that need it
  onMounted(() => {
    // set default theme color if none is set
    if (localStorage.getItem('ghost.reaper.dark') === null) {
      // setThemeDark()
      setThemeLight()
    }
    if (localStorage.getItem('ghost.reaper.theme') === null) {
      setThemeColor('stone')
    }

    setThemeData()
  })

  return {
    themeColors,
    themePrimary,
    themeSecondary,
    themeBackground,
    isDark,
    theme,
    toggleThemeDark,
    setThemeColor,
    setThemeData,
  }
}
