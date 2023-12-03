<template>
    <div class="ml-96 mr-12 mt-32 flex flex-col justify-start  bg-[#ffffff] h-full">

        <section class="">
            <div class="py-8 px-4 mx-auto max-w-2xl lg:py-16">
                <h2 class="mb-4 text-xl font-bold text-[#1E1E1E]">Add a new ward</h2>
                <form ref="form" @submit.prevent="newWard">
                    <div class="grid gap-4 sm:grid-cols-2 sm:gap-6">
                        <div class="sm:col-span-2">
                            <label for="name" class="block mb-2 text-sm font-medium text-[#1E1E1E]">Ward
                                </label>
                            <input type="text" name="ward_name" id="ward_name"
                                class="bg-gray-50 border border-gray-300 text-gray-900 rounded-2xl"
                                placeholder="e.g Ngara" required="" v-model="formFields.ward_name">
                        </div>
                    </div>

                    <div>
                        <label for="constituency"
                            class="block mb-2 text-sm font-medium text-gray-900 dark:text-black pt-2">Constituency</label>
                        <select required id="constituency" v-model="formFields.constituency_name"
                            class="bg-gray-50 border border-gray-300 text-gray-900 rounded-2xl py-2">
                            <option selected="" v-for="(constituency, index) in this.allConstituencies" :key="index">{{constituency.Constituency }}</option>
                        </select>
                    </div>

                    <div>
                        <label for="county"
                            class="block mb-2 text-sm font-medium text-gray-900 dark:text-black pt-2">County</label>
                        <select required id="county" v-model="formFields.county_name"
                            class="bg-gray-50 border border-gray-300 text-gray-900 rounded-2xl py-2">
                            <option selected="" v-for="(county, index) in this.allCounties" :key="index">{{county.name }}</option>
                        </select>
                    </div>
                    <button type="submit"
                        class="inline-flex items-center px-5 py-2.5 mt-4 bg-[#1E1E1E] text-white rounded-2xl">
                        Save
                    </button>
                </form>
            </div>
        </section>
    </div>
</template>

<script>
import axios from 'axios'
import SecureLS from 'secure-ls'
import Toastify from 'toastify-js'
import 'toastify-js/src/toastify.css'

export default {
  data () {
    this.getCounties()
    this.getConstituencies()

    return {
      componentKey: 0,

      allCounties: {},
      allConstituencies: {},
      formFields: {
        county_name: '',
        constituency_name: '',
        ward_name: ''
      }
    }
  },

  methods: {
    newWard () {
      let countyID = ''
      let constituencyID = ''

      for (let i = 0; i < this.allCounties.length; i++) {
        if (this.allCounties[i].name === this.formFields.county_name) {
          console.log(this.allCounties[i].ID)
          countyID = this.allCounties[i].ID
        }
      }

      for (let i = 0; i < this.allConstituencies.length; i++) {
        if (this.allConstituencies[i].Constituency === this.formFields.constituency_name) {
          console.log(',,,,' + this.allConstituencies[i].ConstituencyID)
          constituencyID = this.allConstituencies[i].ConstituencyID
        }
      }

      const ls = new SecureLS()
      const axiosConfig = {
        headers: { Authorization: `Bearer ${ls.get('user').token}` }
      }
      try {
        axios.post('http://127.0.0.1:3500/api/secured/new-ward', [{
          name: this.formFields.ward_name,
          CountyID: parseInt(countyID),
          ConstituencyID: parseInt(constituencyID)

        }], axiosConfig).then(result => {
          console.log(result.data)
          this.$refs.form.reset()
          if (result.status === 201) {
            this.toast('Ward added successfuly')
          }
        })
      } catch (error) {
        console.error(error.response.data)
      }
    },

    getCounties () {
      const ls = new SecureLS()
      const config = {
        headers: { Authorization: `Bearer ${ls.get('user').token}` }
      }
      axios.get(
        'http://127.0.0.1:3500/api/secured/get-all-counties', config
      ).then((response) => {
        console.log(response.data)
        this.allCounties = response.data
      }).catch(function (error) {
        if (error.response.status === 401) {
          ls.removeAll()
          window.location.href = '/'
        } if (error.toJSON().message === 'Network Error') {
          this.toast('no internet connection')
        }
      })
    },

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
          this.toast('no internet connection')
        }
      })
    },

    toast (message) {
      Toastify({
        text: message,
        duration: 3000,
        destination: '',
        newWindow: true,
        close: true,
        gravity: 'top', // `top` or `bottom`
        position: 'center', // `left`, `center` or `right`
        stopOnFocus: true, // Prevents dismissing of toast on hover
        style: {
          background: '#000000'
        },
        onClick: function () { } // Callback after click
      }).showToast()
    }
  }
}
</script>
