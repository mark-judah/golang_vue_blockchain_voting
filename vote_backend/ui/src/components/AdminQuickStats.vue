<template>
    <div>
        <div class="text-3xl font-bold p-6 text-center text-[#4D4D4D]">Quick Statistics</div>

    <div class="flex justify-between mx-24 mb-6">
                <div class="flex flex-col bg-white  mb-4 rounded-2xl pb-4 p-3">
                    <div class="px-3 mb-1">
                        <h1 class="text-center font-bold">Voting Stats</h1>
                    </div>
                    <div class="p-1">
                        <p>Total Registered Voters: {{ stats.TotalRegisteredVoters }}</p>
                    </div>
                    <div class="p-1">
                        <p>Total Votes: {{ stats.TotalVotes }}</p>
                    </div>
                    <div class="p-1">
                        <p>Total Polling Stations: {{ stats.TotalPollingStations }}</p>
                    </div>
                </div>

                <div class="flex flex-col bg-white  mb-4 rounded-2xl pb-4 p-3">
                    <div class="px-3 mb-1">
                        <h1 class="text-center font-bold">Candidates Stats</h1>
                    </div>
                    <div class="p-1">
                        <p>Presidential Candidates: {{ stats.PresidentialCandidates }}</p>
                    </div>
                    <div class="p-1">
                        <p>Gubernatorial Candidates: 0</p>
                    </div>
                    <div class="p-1">
                        <p>Senatorial Candidates: 0</p>
                    </div>
                    <div class="p-1">
                        <p>Nat Assembly Candidates: 0</p>
                    </div>
                    <div class="p-1">
                        <p>Woman Rep Candidates: 0</p>
                    </div>
                    <div class="p-1">
                        <p>MCA Candidates: 0</p>
                    </div>
                </div>

                <div class="flex flex-col bg-white  mb-4 rounded-2xl pb-4 p-3">
                    <div class="px-3 mb-1">
                        <h1 class="text-center font-bold">Desktop Clients Stats</h1>
                    </div>
                    <div class="p-1">
                        <p>Total DesktopClients {{ stats.TotalDesktopClients }}</p>
                    </div>
                    <div class="p-1">
                        <p>Total Online Clients: {{ stats.OnlineClients }}</p>
                    </div>
                </div>

                <div class="flex flex-col bg-white  mb-4 rounded-2xl pb-4 p-3">
                    <div class="px-3 mb-1">
                        <h1 class="text-center font-bold">Transaction Pool Stats</h1>
                    </div>
                    <div class="p-1">
                        <p>Total Transactions: {{ stats.TransactionPoolSize }}</p>
                    </div>
                    <div class="p-1">
                        <p>Total Processed Transactions: {{ stats.TotalProcessedTransactions }}</p>
                    </div>
                </div>
            </div>
            </div>
</template>
<script>
import axios from 'axios'
import SecureLS from 'secure-ls'

export default {

  data () {
    this.getStats()
    return {
      stats: {}

    }
  },

  methods: {
    getStats () {
      const ls = new SecureLS()
      const config = {
        headers: { Authorization: `Bearer ${ls.get('user').token}` }
      }
      axios.get(
        'http://127.0.0.1:3500/api/secured/get-quick-stats', config
      ).then((response) => {
        console.log(response.data)
        this.stats = response.data
      }).catch(function (error) {
        if (error.response.status === 401) {
          ls.removeAll()
          window.location.href = '/'
        } if (error.toJSON().message === 'Network Error') {
          alert('no internet connection')
        }
      })
    }
  }

}
</script>

<style scoped></style>
