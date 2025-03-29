import { createMemoryHistory, createRouter } from 'vue-router'

import FileSearch from './components/FileSearch.vue'

const routes = [
  { path: '/', component: FileSearch },
]

const router = createRouter({
  history: createMemoryHistory(),
  routes,
})

export default router