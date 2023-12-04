<template>
    <div class="w-1/4">
                <h1 class="text-center  pb-3">Network State</h1>
                <div class="ml-10 flex flex-col justify-center items-center bg-white px-10">
                    <div class="py-4">
                        <p>Total Blocks: {{ state.BlockHeight }} </p>
                    </div>
                    <div class="py-4">
                        <p>Processed Transactions: {{ state.NoOfTransactions }} </p>
                    </div>

                </div>
            </div>
</template>

<script>
import axios from 'axios'

export default {
  data () {
    this.getNetworkState()
    return {
      state: {}

    }
  },

  methods: {
    getNetworkState () {
      axios.get(
        'http://127.0.0.1:3500/api/get-network-state').then((response) => {
        console.log(response.data)
        this.state = response.data
      }).catch(function (error) {
        if (error.toJSON().message === 'Network Error') {
          alert('no internet connection')
        }
      })
    }
  }
}
</script>
