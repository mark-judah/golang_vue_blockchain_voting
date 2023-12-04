<template>
    <div class="flex justify-between">
        <div class="flex flex-col bg-white  m-2 rounded-2xl">
            <div class="bg-[#4D4D4D] rounded-t-2xl">
                <p class="text-white font-bold p-3 text-center">Presidential Election Results</p>
            </div>

            <div>
                <table class="table-auto space-y-6">
                    <thead class="">

                        <tr class="rounded-3xl">
                            <th class="px-6 py-2 border">Candidate</th>
                            <th class="px-6 py-2 border">Party</th>
                            <th class="px-6 py-2 border">Total Votes</th>
                            <th class="px-6 py-2 border">Percentage</th>
                        </tr>
                    </thead>
                    <tr v-for="(data, i) in this.results" :key="i" class="bg-white border-b">
                        <td class="px-6 py-2   border ">{{ data.Candidate }}</td>
                        <td class="px-6 py-2   border ">{{ data.Party }}</td>
                        <td class="px-6 py-2   border ">{{ data.TotalVotes }}</td>
                        <td class="px-6 py-2   border ">{{(data.TotalVotes/totalvotes)*100 }}%</td>
                    </tr>
                </table>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'

export default {
  data () {
    this.getElectionResults()
    return {
      results: {}

    }
  },
  computed: {
    totalvotes: function () {
      let x = 0
      const y = this.results
      for (const index in y) {
        x = x + parseInt(y[index].TotalVotes)
        console.log(x)
      }
      return x
    }
  },

  methods: {
    getElectionResults () {
      axios.get(
        'http://127.0.0.1:3500/api/tally-votes').then((response) => {
        console.log(response.data)
        this.results = response.data
      }).catch(function (error) {
        if (error.toJSON().message === 'Network Error') {
          alert('no internet connection')
        }
      })
    }
  }
}
</script>
