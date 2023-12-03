<template>
  <div class="">
    <div class="ml-96 mr-12 mt-32 flex flex-col justify-start  bg-[#ffffff] h-full">
      <div class="flex flex-row justify-between items-center w-full">
        <div>
          <p class="px-12">Counties</p>
        </div>

        <div class="flex flex-row px-12">
          <router-link :to="{ name: 'candidates' }" @click="$router.go()" class="flex items-center p-2">
            <img src="../assets/images/refresh.svg" class="w-6 h-8">
          </router-link>

          <router-link :to="{ name: 'new-candidate' }" class="flex items-center p-2">
            <img src="../assets/images/add.svg" class="w-6 h-8">
          </router-link>
        </div>
      </div>

           <div class="px-12 flex justify-center w-fit">

        <table class="table-auto space-y-6">
          <thead class="">

            <tr class="rounded-3xl">
              <th class="px-1.5 py-2 border">Index</th>
              <th class="px-1.5 py-2 border">Name</th>
              <th class="px-1.5 py-2 border">Position</th>
              <th class="px-1.5 py-2 border">Party</th>
              <th class="px-1.5 py-2 border">Slogan</th>
              <th class="px-1.5 py-2 border">Statement</th>
              <th class="px-1.5 py-2 border">County</th>
              <th class="px-1.5 py-2 border">Constituency</th>
              <th class="px-1.5 py-2 border">Ward</th>
              <th class="px-1.5 py-2 border">Polling Station</th>
              <th class="px-1.5 py-2 border">Photo</th>
            </tr>
          </thead>
          <tr v-for="(data, i) in this.allCandidates" :key="data" class="bg-white border-b">
            <td class="px-1.5 py-2   border ">{{i+1 }}</td>
            <td class="px-1.5 py-2   border ">{{ data.Candidate }}</td>
            <td class="px-1.5 py-2   border ">{{ data.Position }}</td>
            <td class="px-1.5 py-2   border ">{{ data.Party }}</td>
            <td class="px-1.5 py-2   border ">{{ data.Slogan }}</td>
            <td class="px-1.5 py-2   border ">{{ data.Statement }}</td>
            <td class="px-1.5 py-2   border ">{{ data.County }}</td>
            <td class="px-1.5 py-2   border ">{{ data.Constituency }}</td>
            <td class="px-1.5 py-2   border ">{{ data.Ward }}</td>
            <td class="px-1.5 py-2   border ">{{ data.PollingStation }}</td>
            <td class="px-1.5 py-2   border "><img :src=data.Photo></td>

          </tr>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import SecureLS from 'secure-ls'

export default {

  data () {
    this.getCandidates()
    return {
      allCandidates: {}

    }
  },

  methods: {
    getCandidates () {
      const ls = new SecureLS()
      const config = {
        headers: { Authorization: `Bearer ${ls.get('user').token}` }
      }
      axios.get(
        'http://127.0.0.1:3500/api/get-all-candidates', config
      ).then((response) => {
        console.log(response.data)
        this.allCandidates = response.data
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
