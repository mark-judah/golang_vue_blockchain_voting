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
                    <tr v-for="(data, i) in this.allConstituencies"  :key="data" class="bg-white border-b">
                            <td class="px-1.5 py-2   border ">{{++i }}</td>
                            <td class="px-1.5 py-2   border ">{{ data.Constituency }}</td>
                            <td class="px-1.5 py-2   border "  >{{ data.County   }}</td>
                            <td class="px-1.5 py-2   border ">Edit|Delete</td>
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
