const swapTheme = (dark: boolean) => {
    if (dark) {
        document.documentElement.classList.add('dark')
    } else {
        document.documentElement.classList.remove('dark')
    }
    localStorage.darkMode = dark;
}

export default swapTheme;