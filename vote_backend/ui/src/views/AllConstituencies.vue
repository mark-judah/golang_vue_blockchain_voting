<template>
    <div class="">
        <div class="ml-96 mr-12 mt-32 flex justify-start  bg-[#ffffff] h-full">

            <div class="p-12 flex justify-center">
                <table class="table-auto space-y-6">
                    <thead class="">

                        <tr class="rounded-3xl">
                            <th class="px-1.5 py-2 border">Index</th>
                            <th class="px-1.5 py-2 border">Name</th>
                            <th class="px-1.5 py-2 border">County</th>
                            <th class="px-1.5 py-2 border">Action</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(constituency, index) in this.allConstituencies" :key="index"
                            class="bg-white border-b">
                            <td class="px-1.5 py-2   border ">{{ index+1 }} </td>
                            <td class="px-1.5 py-2   border ">{{ constituency.name }}</td>
                            <td class="px-1.5 py-2   border ">{{ constituency.county }}</td>

                            <td class="px-1.5 py-2   border ">Edit|Delete</td>

                        </tr>
                    </tbody>
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
    this.getConstituencies()
    return {
      allConstituencies: {}

    }
  },

  methods: {
    getConstituencies () {
      const ls = new SecureLS()
      const config = {
        headers: { Authorization: `Bearer ${ls.get('user').token}` }
      }
      axios.get(
        'http://127.0.0.1:3500/api/secured/get-all-constituencies', config
      ).then((response) => {
        console.log(response.data)
        this.allConstituencies = response.data
      })
    }
  }

}
</script>

<style scoped></style>
