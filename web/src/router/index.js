import { createRouter, createWebHistory } from 'vue-router'

import BoothLayout from '../views/booth/BoothLayout.vue'
import WelcomeScreen from '../views/booth/WelcomeScreen.vue'
import CaptureScreen from '../views/booth/CaptureScreen.vue'
import CollageScreen from '../views/booth/CollageScreen.vue'
import PreviewScreen from '../views/booth/PreviewScreen.vue'
import ResultScreen from '../views/booth/ResultScreen.vue'
import AnimalScreen from '../views/booth/AnimalScreen.vue'

import AdminLayout from '../views/admin/AdminLayout.vue'
import AdminDashboard from '../views/admin/AdminDashboard.vue'
import AdminPhotos from '../views/admin/AdminPhotos.vue'
import AdminSessions from '../views/admin/AdminSessions.vue'
import AdminDevices from '../views/admin/AdminDevices.vue'
import AdminConfig from '../views/admin/AdminConfig.vue'

const routes = [
  {
    path: '/',
    component: BoothLayout,
    children: [
      { path: '', name: 'welcome', component: WelcomeScreen },
      { path: 'capture', name: 'capture', component: CaptureScreen },
      { path: 'collage', name: 'collage', component: CollageScreen },
      { path: 'preview', name: 'preview', component: PreviewScreen },
      { path: 'result', name: 'result', component: ResultScreen },
      { path: 'animal', name: 'animal', component: AnimalScreen },
    ],
  },
  {
    path: '/admin',
    component: AdminLayout,
    children: [
      { path: '', name: 'admin-dashboard', component: AdminDashboard },
      { path: 'photos', name: 'admin-photos', component: AdminPhotos },
      { path: 'sessions', name: 'admin-sessions', component: AdminSessions },
      { path: 'devices', name: 'admin-devices', component: AdminDevices },
      { path: 'config', name: 'admin-config', component: AdminConfig },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
