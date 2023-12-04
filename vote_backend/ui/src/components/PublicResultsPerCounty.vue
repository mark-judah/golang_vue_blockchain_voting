<template>
  <div class="flex justify-between">
    <div class="flex flex-col bg-white  m-2 rounded-2xl">
      <div class="bg-[#4D4D4D] rounded-t-2xl">
        <p class="text-white font-bold p-3 text-center">Results Per County</p>
      </div>

      <div>
        <canvas ref="myChart"></canvas>
      </div>
    </div>

  </div>
</template>

<script>
import axios from 'axios'
import { Chart as ChartJS } from 'chart.js/auto'

export default {
  name: 'BarChart',
  data () {
    this.getStats()

    return {
      stats: {},
      totalVotes: [],
      counties: []

    }
  },
  // computed: {
  //   totalvotes: function () {
  //     let x = 0
  //     const y = this.stats.ResultsPerCounty
  //     for (const index in y) {
  //       x = x + parseInt(y[index].TotalVotes)
  //       console.log(x)
  //     }
  //     return x
  //   }
  // },

  methods: {
    getStats () {
      axios.get(
        'http://127.0.0.1:3500/api/get-quick-stats').then((response) => {
        console.log(response.data.ResultsPerCounty[0].TotalVotes)
        this.stats = response.data.ResultsPerCounty
        for (const i in this.stats) {
          this.totalVotes.push(this.stats[i].TotalVotes)
          this.counties.push(this.stats[i].County)

          console.log(this.stats[i].TotalVotes)
        }
        this.renderChat()
      }).catch(function (error) {
        console.log(error)
      })
    },

    renderChat () {
      // Process the data to aggregate total votes per candidate for each county
      const aggregatedData = {}
      this.stats.forEach(entry => {
        if (!aggregatedData[entry.County]) {
          aggregatedData[entry.County] = {}
        }
        if (!aggregatedData[entry.County][entry.Candidate]) {
          aggregatedData[entry.County][entry.Candidate] = 0
        }
        aggregatedData[entry.County][entry.Candidate] += entry.TotalVotes
      })

      // Prepare data for Chart.js
      const labels = Object.keys(aggregatedData) // County names
      const datasets = Object.keys(this.stats.reduce((acc, entry) => {
        acc[entry.Candidate] = true
        return acc
      }, {})).map(candidate => ({
        label: candidate,
        data: labels.map(county => aggregatedData[county][candidate] || 0)
      }))

      const chartData = {
        labels,
        datasets
      }

      // Get the drawing context on the canvas
      const ctx = this.$refs.myChart.getContext('2d')

      const mychart = new ChartJS(ctx, {
        type: 'bar',
        data: chartData,
        options: {

          scales: {
            x: {
              stacked: true
            },
            y: {
              stacked: true
            }
          }
        }
      })
      return mychart
    }
  }
}

</script>
