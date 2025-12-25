import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { VueQueryPlugin } from '@tanstack/vue-query'
import Login from "./components/login.vue"
import { createRouter, createWebHashHistory } from "vue-router"

const routes = [
  { path: "/", redirect: "/login" },
  { path: '/login', name: "login", component: Login }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

createApp(App).use(router).use(VueQueryPlugin).mount('#app')
