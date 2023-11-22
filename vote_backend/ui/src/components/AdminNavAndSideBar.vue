<template>
    <nav class="fixed top-0 z-50 w-full bg-white border-b border-gray-200 dark:bg-gray-800 dark:border-gray-700">
        <div class="px-3 py-3 lg:px-5 lg:pl-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center justify-start">
                    <!-- <button type="button" class="" aria-expanded="false" data-dropdown-toggle="logo-sidebar">
                        <span class="sr-only">Open side bar</span>
                        <img class="w-8 h-8 rounded-full" src="/images/menu.svg" alt="menu">
                    </button> -->

                    <a href="/" class="flex ml-2 md:mr-24">
                        <span class="self-center text-xl font-semibold sm:text-xl whitespace-nowrap dark:text-white"> ChainVote Admin
                            Dashboard</span>
                    </a>
                </div>
                <div class="flex items-center">
                    <div class="flex items-center mr-6">
                        <div>
                            <button type="button"
                                class="flex text-sm bg-gray-800 rounded-full focus:ring-4 focus:ring-gray-300 dark:focus:ring-gray-600"
                                aria-expanded="false" data-dropdown-toggle="dropdown-user">
                                <span class="sr-only">Open user menu</span>
                                <img class="w-8 h-8 rounded-full"
                                    src="https://flowbite.com/docs/images/people/profile-picture-5.jpg" alt="user photo">
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
                                    <a href="#"
                                        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white"
                                        role="menuitem">Dashboard</a>
                                </li>
                                <li>
                                    <a href="#"
                                        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white"
                                        role="menuitem">Settings</a>
                                </li>
                                <li>
                                    <a href="#"
                                        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white"
                                        role="menuitem">Sign out</a>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </nav>

    <div>
        <aside id="logo-sidebar"
            class="fixed top-0 left-0 z-40 w-64 h-screen pt-5 transition-transform -translate-x-full bg-white border-r border-gray-200 sm:translate-x-0 dark:bg-gray-800 dark:border-gray-700"
            aria-label="Sidebar">
            <div class="h-full px-3 pb-4 overflow-y-auto bg-white dark:bg-gray-800">
                <ul class="space-y-2">
                    <li>
                        <router-link :to="{ name: 'admin' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/images/home.svg" class="w-6 h-6">
                            <span class="ml-3">Home</span>
                        </router-link>
                    </li>

                    <li>
                        <router-link :to="{ name: 'admin' }"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/app/ui/src/assets/images/election_officials.svg" class="w-6 h-6">
                            <span class="flex-1 ml-3 whitespace-nowrap">Election Officers</span>
                        </router-link>

                    </li>

                    <li>
                        <a href="/admin-people"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/images/people.svg" class="w-6 h-6">
                            <span class="flex-1 ml-3 whitespace-nowrap">Counties</span>
                        </a>
                    </li>

                    <li>
                        <a href="/admin-manage-account"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/images/admin.svg" class="w-6 h-6">
                            <span class="flex-1 ml-3 text-left whitespace-nowrap">Constituencies</span>
                        </a>

                    </li>
                    <li>
                        <a href="/logout-auth"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/images/logout.svg" class="w-6 h-6">
                            <span class="flex-1 ml-3 whitespace-nowrap">Wards</span>
                        </a>
                    </li>
                    <li>
                        <a href="/logout-auth"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/images/logout.svg" class="w-6 h-6">
                            <span class="flex-1 ml-3 whitespace-nowrap">Polling Stations</span>
                        </a>
                    </li>
                    <li>
                        <a href="/logout-auth"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/images/logout.svg" class="w-6 h-6">
                            <span class="flex-1 ml-3 whitespace-nowrap">Candidates</span>
                        </a>
                    </li>
                    <li>
                        <a href="/logout-auth"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/images/logout.svg" class="w-6 h-6">
                            <span class="flex-1 ml-3 whitespace-nowrap">Voters</span>
                        </a>
                    </li>
                    <li>
                        <a href="/logout-auth"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/images/logout.svg" class="w-6 h-6">
                            <span class="flex-1 ml-3 whitespace-nowrap">Desktop Clients</span>
                        </a>
                    </li>
                    <li>
                        <a href="/logout-auth"
                            class="flex items-center p-2 text-base font-normal text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                            <img src="/images/logout.svg" class="w-6 h-6">
                            <span class="flex-1 ml-3 whitespace-nowrap">Transaction Pool</span>
                        </a>
                    </li>
                </ul>
            </div>
        </aside>

    </div>
</template>

<script>
import SecureLS from 'secure-ls'

export default {
  data () {
    const ls = new SecureLS()
    return {
      username: ls.get('user').username,
      email: ls.get('user').email,
      token: ls.get('user').token
    }
  },
  mounted () {
    this.guardMyRoute()
  },
  methods: {
    guardMyRoute () {
      const ls = new SecureLS()
      if (ls.get('user') != null) {
        console.log('authenticated')
      } else {
        window.location.href = '/'
      }
    }
  }
}
</script>

<style scoped></style>
