<template>
    <nav class="fixed top-0 z-50 w-full bg-white">
        <div class="flex items-center justify-start ml-60">

            <BlockExplorerNav></BlockExplorerNav>

            <router-link :to="{ name: 'results' }"  class="flex items-center p-2 ml-32">
                <div class="px-1">
                    <button class="bg-transparent  font-semibold text-black py-2 px-4 ">
                        Results</button>

                </div>
            </router-link>

            <router-link :to="{ name: 'block-explorer' }" class="flex items-center p-2">
                <div class="px-1">
                    <button class="bg-transparent  font-semibold text-black py-2 px-4 ">
                        Public Ledger</button>

                </div>
            </router-link>
        </div>
    </nav>
    <div class="">
        <div class="mt-32 flex justify-center space-x-10   h-full">

            <NetworkState></NetworkState>
            <div class="w-1/2">
                <h1 class="text-center  pb-3">Blocks</h1>
            <div class="flex justify-center  bg-white mr-6">
                <table class="table-auto space-y-6">
                    <thead class="">

                        <tr class="rounded-3xl">
                            <th class="px-1.5 py-2 border">Block Height</th>
                            <th class="px-1.5 py-2 border">Block Hash</th>
                            <th class="px-1.5 py-2 border">Previous BlockHash</th>
                            <th class="px-1.5 py-2 border">No of Transactions</th>
                            <th class="px-1.5 py-2 border">Block Size</th>

                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(block, index) in this.allBlocks" :key="index" class="bg-white border-b">
                            <td class="px-1.5 py-2   border ">{{ block.BlockHeight }}</td>
                            <td class="px-1.5 py-2   border break-all">{{ block.BlockHash }}</td>
                            <td class="px-1.5 py-2   border break-all">{{ block.PreviousBlockHash }}</td>
                            <td class="px-1.5 py-2   border ">{{  block.NoOfTransactions}}</td>
                            <td class="px-1.5 py-2   border ">{{  block.BlockSize}} bytes</td>

                        </tr>
                    </tbody>
                </table>
            </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import BlockExplorerNav from '../components/BlockExplorerNav.vue'
import NetworkState from '../components/NetworkState.vue'

export default {
  components: {
    BlockExplorerNav,
    NetworkState
  },

  data () {
    this.getUsers()
    return {
      allBlocks: {}

    }
  },

  methods: {
    getUsers () {
      axios.get(
        'http://127.0.0.1:3500/api/get-blockchain').then((response) => {
        console.log(response.data)
        this.allBlocks = response.data
      }).catch(function (error) {
        if (error.toJSON().message === 'Network Error') {
          alert('no internet connection')
        }
      })
    }
  }

}
</script>

<style scoped></style>
