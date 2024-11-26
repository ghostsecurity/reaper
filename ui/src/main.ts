import './assets/index.css'
import './assets/themes.css'
import { createApp } from 'vue'
import router from './router'
import { createPinia } from 'pinia'
import App from './App.vue'
import { useNavStore } from './stores/nav'

const app = createApp(App)

router.afterEach((to, from) => {
  const navStore = useNavStore()
  navStore.recordNavigation(to.path, from.path)
})

app.use(router)
app.use(createPinia())
app.mount('#app')
