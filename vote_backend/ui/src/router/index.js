import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import AdminView from '../views/AdminView.vue'
import Users from '../views/AllUsers.vue'
import Counties from '../views/AllCounties.vue'
import Constituencies from '../views/AllConstituencies.vue'
import Wards from '../views/AllWards.vue'
import PollingStations from '../views/PollingStations.vue'
import Candidates from '../views/AllCandidates.vue'
import Voters from '../views/AllVoters.vue'
import DesktopClients from '../views/DesktopClients.vue'
import TransactionPool from '../views/TransactionPool.vue'

const routes = [
  {
    path: '/',
    name: 'login',
    component: LoginView
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminView
  },
  {
    path: '/users',
    name: 'users',
    component: Users
  },
  {
    path: '/counties',
    name: 'counties',
    component: Counties
  },
  {
    path: '/constituencies',
    name: 'constituencies',
    component: Constituencies
  }, {
    path: '/wards',
    name: 'wards',
    component: Wards
  }, {
    path: '/polling-stations',
    name: 'polling-stations',
    component: PollingStations
  }, {
    path: '/candidates',
    name: 'candidates',
    component: Candidates
  }, {
    path: '/voters',
    name: 'voters',
    component: Voters
  }, {
    path: '/desktop-clients',
    name: 'desktop-clients',
    component: DesktopClients
  }, {
    path: '/transaction-pool',
    name: 'transaction-pool',
    component: TransactionPool
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }

]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
