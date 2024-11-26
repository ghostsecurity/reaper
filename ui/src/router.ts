import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

import ScanMain from './components/ScanMain.vue'
import ExploreMain from './components/ExploreMain.vue'
import ReplayMain from './components/ReplayMain.vue'
import AttackMain from './components/AttackMain.vue'
import AgentMain from './components/AgentMain.vue'
import ReportsMain from './components/ReportsMain.vue'
import SettingsMain from './components/SettingsMain.vue'
import CollaborateMain from './components/CollaborateMain.vue'

const routes: Array<RouteRecordRaw> = [
  { path: '/', name: 'root', component: ScanMain },
  { path: '/scan', name: 'scan', component: ScanMain },
  { path: '/explore', name: 'explore', component: ExploreMain },
  { path: '/replay', name: 'replay', component: ReplayMain },
  { path: '/test', name: 'test', component: AttackMain },
  { path: '/agent', name: 'agent', component: AgentMain },
  { path: '/collaborate', name: 'collaborate', component: CollaborateMain },
  { path: '/settings', name: 'settings', component: SettingsMain },
  { path: '/reports', name: 'reports', component: ReportsMain },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router