<template>
  <div class="ml-96 mr-12 mt-32 flex flex-col justify-start  bg-[#ffffff] h-full">

    <section class="">
      <div class="py-8 px-4 mx-auto max-w-2xl lg:py-16">
        <h2 class="mb-4 text-xl font-bold text-[#1E1E1E]">Add a new county</h2>
        <form ref="form" @submit.prevent="newCounty">
          <div class="grid gap-4 sm:grid-cols-2 sm:gap-6">
            <div class="sm:col-span-2">
              <label for="name" class="block mb-2 text-sm font-medium text-[#1E1E1E]">County
                Name</label>
              <input type="text" name="county_name" id="county_name"
                class="bg-gray-50 border border-gray-300 text-gray-900 rounded-2xl" placeholder="e.g Nairobi" required=""
                v-model="formFields.county_name">
            </div>
          </div>
          <button type="submit" class="inline-flex items-center px-5 py-2.5 mt-4 bg-[#1E1E1E] text-white rounded-2xl">
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
    return {
      componentKey: 0,

      formFields: {
        county_name: ''
      }
    }
  },

  methods: {
    newCounty () {
      const ls = new SecureLS()
      const axiosConfig = {
        headers: { Authorization: `Bearer ${ls.get('user').token}` }
      }
      try {
        axios.post('http://127.0.0.1:3500/api/secured/new-county', [{ name: this.formFields.county_name }], axiosConfig).then(result => {
          console.log(result.data)
          this.$refs.form.reset()
          if (result.status === 201) {
            Toastify({
              text: 'County added successfuly',
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
        })
      } catch (error) {
        console.error(error.response.data)
      }
    }
  }
}
</script>
