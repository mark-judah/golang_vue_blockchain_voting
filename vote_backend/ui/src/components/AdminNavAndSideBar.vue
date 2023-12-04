<template>
    <nav class="fixed top-0 z-50 w-full bg-[#4D4D4D]">
        <div class="">
            <div class="flex items-center justify-between">
                <div class="flex items-center justify-start">
                    <div class="bg-[#1E1E1E] flex">
                        <router-link :to="{ name: 'admin' }">
                            <p class="p-2.5 text-xl font-semibold dark:text-white mr-5">
                                ChainVote Admin
                                Dashboard</p>
                        </router-link>
                    </div>
                    <div class="px-1">
                        <button type="button"
                            class="flex text-sm  rounded-full focus:ring-4 focus:ring-gray-300 dark:focus:ring-gray-600"
                            aria-expanded="false" data-dropdown-toggle="dropdown-user">
                            <img class="w-8 h-8 rounded-full" src="../assets/images/user.svg" alt="user photo">

                        </button>
                    </div>
                    <div>
                        <p class="text-white ml-2"> {{ username }}</p>
                    </div>
                </div>
                <div class="flex items-center">
                    <div class="flex items-center">
                        <router-link :to="{ name: 'block-explorer' }" target="_blank" class="flex items-center p-2">
                            <div class="px-1">
                                <button
                                    class="bg-[#1E1E1E]  font-semibold text-white py-1 px-4 border border-ffffff hover:border-transparent rounded">
                                    Block Explorer</button>

                            </div>
                        </router-link>
                        <div class="px-1">
                            <button type="button" class="flex text-sm  rounded-full" aria-expanded="false"
                                data-dropdown-toggle="dropdown-user">
                                <span class="sr-only">Open user menu</span>
                                <img class="w-8 h-8 rounded-full" src="../assets/images/admin.svg" alt="user photo">
                            </button>
                        </div>
                        <div class="z-50 hidden my-4 text-base list-none bg-white divide-y divide-gray-100 rounded shadow dark:bg-gray-700 dark:divide-gray-600"
                            id="dropdown-user">
                            <div class="px-4 py-3" role="none">
                                <p class="text-sm text-gray-900 dark:text-white" role="none">
                                    {{ username }}
                                </p>
                                <p class="text-sm font-medium text-gray-900 truncate dark:text-gray-300" role="none">
                                    {{ email }}
                                </p>
                            </div>
                            <ul class="py-1" role="none">
                                <li>
                                    <!-- only visible to superuser and hr-admin -->
                                    <router-link :to="{ name: 'users' }"
                                        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white"
                                        role="menuitem">Users </router-link>
                                </li>
                                <li>
                                    <router-link :to="{ name: 'admin' }"
                                        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white"
                                        role="menuitem">Profile </router-link>
                                </li>
                                <li>
                                    <p v-on:click="logout"
                                        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white"
                                        role="menuitem">Sign out </p>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </nav>

    <div class="">
        <aside id="logo-sidebar"
            class="fixed  top-0 left-0 z-40 w-80 h-screen pt-5 transition-transform -translate-x-full bg-white border-r border-gray-200 sm:translate-x-0 dark: dark:border-gray-700"
            aria-label="Sidebar">
            <div class="h-full pb-4 overflow-y-auto bg-[#4D4D4D]">
                <ul class="mt-7">
                    <li>
                        <router-link :to="{ name: 'admin' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/home.svg" class="w-6 h-8">
                            <span class="ml-3">Home</span>
                        </router-link>
                    </li>

                    <li class="bg-[#1E1E1E] w-full">
                        <router-link :to="{ name: 'counties' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/county.svg" class="w-6 h-8">
                            <span class="flex-1 ml-3 whitespace-nowrap">Counties</span>
                        </router-link>
                    </li>

                    <li>
                        <router-link :to="{ name: 'constituencies' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/parliament.svg" class="w-6 h-8">
                            <span class="flex-1 ml-3 text-left whitespace-nowrap">Constituencies</span>
                        </router-link>

                    </li>

                    <li class="bg-[#1E1E1E] w-full">
                        <router-link :to="{ name: 'wards' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/ward.svg" class="w-6 h-8">
                            <span class="flex-1 ml-3 whitespace-nowrap">Wards</span>
                        </router-link>
                    </li>

                    <li>
                        <router-link :to="{ name: 'polling-stations' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/queue.svg" class="w-6 h-8">
                            <span class="flex-1 ml-3 whitespace-nowrap">Polling Stations</span>
                        </router-link>
                    </li>

                    <li class="bg-[#1E1E1E] w-full">
                        <router-link :to="{ name: 'candidates' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/candidates.svg" class="w-6 h-8">
                            <span class="flex-1 ml-3 whitespace-nowrap">Candidates</span>
                        </router-link>
                    </li>

                    <li>
                        <router-link :to="{ name: 'voters' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/voters.svg" class="w-6 h-8">
                            <span class="flex-1 ml-3 whitespace-nowrap">Voters</span>
                        </router-link>
                    </li>

                    <li class="bg-[#1E1E1E] w-full">
                        <router-link :to="{ name: 'desktop-clients' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/desktop_clients.svg" class="w-6 h-8">
                            <span class="flex-1 ml-3 whitespace-nowrap">Desktop Clients</span>
                        </router-link>
                    </li>

                    <li>
                        <router-link :to="{ name: 'transaction-pool' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/transactions.svg" class="w-6 h-8">
                            <span class="flex-1 ml-3 whitespace-nowrap">Transaction Pool</span>
                        </router-link>
                    </li>

                    <li class="bg-[#1E1E1E] w-full">
                        <router-link :to="{ name: 'connected-nodes' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="../assets/images/node.svg" class="w-6 h-8">
                            <span class="flex-1 ml-3 whitespace-nowrap">Connected Nodes</span>
                        </router-link>
                    </li>
                </ul>
            </div>
        </aside>

    </div>
</template>

<script>
import SecureLS from 'secure-ls'
import axios from 'axios'

export default {
  data () {
    const ls = new SecureLS()
    return {
      username: ls.get('user').username,
      email: ls.get('user').email
    }
  },
  mounted () {
    this.guardMyRoute()
  },
  methods: {
    guardMyRoute () {
      const ls = new SecureLS()

      const config = {
        headers: { Authorization: `Bearer ${ls.get('user').token}` }
      }

      axios.get(
        'http://127.0.0.1:3500/api/secured/check-auth', config
      ).then((response) => { }).catch(function (error) {
        if (error.response.status === 401) {
          ls.removeAll()
          window.location.href = '/'
        } if (error.toJSON().message === 'Network Error') {
          alert('no internet connection')
        }
      })
    },

    logout () {
      const ls = new SecureLS()
      ls.removeAll()
      this.guardMyRoute()
    }
  }
}
</script>

<style scoped></style>
