import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import AdminHome from '../views/AdminHome.vue'
import Users from '../views/AllUsers.vue'
import NewUser from '../views/NewUser.vue'
import Counties from '../views/AllCounties.vue'
import NewCounty from '../views/NewCounty.vue'
import NewConstituency from '../views/NewConstituency.vue'
import Constituencies from '../views/AllConstituencies.vue'
import Wards from '../views/AllWards.vue'
import NewWard from '../views/NewWard.vue'
import PollingStations from '../views/AllPollingStations.vue'
import NewPollingStation from '../views/NewPollingStation.vue'
import Candidates from '../views/AllCandidates.vue'
import NewCandidate from '../views/NewCandidate.vue'
import Voters from '../views/AllVoters.vue'
import NewVoter from '../views/NewVoter.vue'
import DesktopClients from '../views/DesktopClients.vue'
import TransactionPool from '../views/TransactionPool.vue'
import ConnectedNodes from '../views/ConnectedNodes.vue'
import AllResults from '../views/AllResults.vue'
import BlockExplorer from '../views/BlockExplorer.vue'

const routes = [
  {
    path: '/',
    name: 'login',
    component: LoginView
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminHome
  },
  {
    path: '/users',
    name: 'users',
    component: Users
  },
  {
    path: '/new-user',
    name: 'new-user',
    component: NewUser
  },
  {
    path: '/counties',
    name: 'counties',
    component: Counties
  },
  {
    path: '/new-county',
    name: 'new-county',
    component: NewCounty
  },
  {
    path: '/new-constituency',
    name: 'new-constituency',
    component: NewConstituency
  },
  {
    path: '/constituencies',
    name: 'constituencies',
    component: Constituencies
  }, {
    path: '/wards',
    name: 'wards',
    component: Wards
  },
  {
    path: '/new-ward',
    name: 'new-ward',
    component: NewWard
  },
  {
    path: '/polling-stations',
    name: 'polling-stations',
    component: PollingStations
  },
  {
    path: '/new-polling-station',
    name: 'new-polling-station',
    component: NewPollingStation
  },
  {
    path: '/candidates',
    name: 'candidates',
    component: Candidates
  },
  {
    path: '/new-candidate',
    name: 'new-candidate',
    component: NewCandidate
  },
  {
    path: '/voters',
    name: 'voters',
    component: Voters
  },
  {
    path: '/new-voter',
    name: 'new-voter',
    component: NewVoter
  },
  {
    path: '/desktop-clients',
    name: 'desktop-clients',
    component: DesktopClients
  },
  {
    path: '/transaction-pool',
    name: 'transaction-pool',
    component: TransactionPool
  },
  {
    path: '/connected-nodes',
    name: 'connected-nodes',
    component: ConnectedNodes
  },
  {
    path: '/results',
    name: 'results',
    component: AllResults
  },
  {
    path: '/block-explorer',
    name: 'block-explorer',
    component: BlockExplorer
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
