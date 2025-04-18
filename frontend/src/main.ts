import {createApp} from 'vue'
import App from './App.vue'
import router from './router';
import './style.css';
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura';
import 'primeicons/primeicons.css'
import VueVirtualScroller from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import ToastService from 'primevue/toastservice';
import { definePreset } from '@primeuix/styled'

const CustomAura = definePreset(Aura, {
    components: {
      toast: {
        colorScheme: {
          dark: {
            error: {
              background: '{red.500}',
              color: '{white}',
            },
            success: {
              background: '{green.500}',
              color: '{white}',
            },
            warn: {
              background: '{yellow.600}',
              color: '{white}',
            },
            info: {
              background: '{blue.500}',
              color: '{white}',
            },
          },
        },
      },
    },
  })
  

createApp(App)
    .use(router)
    .use(PrimeVue, {
        theme: {
            preset: CustomAura,
            options: {
                darkModeSelector: 'system'
            }
        }
    })
    .use(VueVirtualScroller)
    .use(ToastService)
    .mount('#app')

