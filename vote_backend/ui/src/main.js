import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import AdminBase from './admin-base.vue'
import SecureLS from 'secure-ls'

let auth = false

const ls = new SecureLS()
if (ls.get('user').token != null) {
  auth = true
}
if (auth) {
  createApp(AdminBase).use(router).mount('#app')
} else {
  createApp(App).use(router).mount('#app')
}
