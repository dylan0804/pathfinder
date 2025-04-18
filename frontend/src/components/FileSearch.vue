<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="flex flex-col h-full w-full bg-white/90 dark:bg-zinc-800/90 backdrop-blur-sm shadow-lg rounded-lg overflow-hidden border border-gray-200/50 dark:border-zinc-700/30">
      <div class="p-4 border-b border-gray-200/70 dark:border-zinc-700/70 bg-white/80 dark:bg-zinc-800/60 backdrop-blur-sm">
        <div class="flex items-center w-full">
          <span class="p-input-icon-left relative flex-grow">
            <InputText
              v-model="searchPattern"
              class="w-full pl-10 py-2.5 border border-gray-300 dark:border-zinc-700/70 rounded-lg shadow-sm hover:border-gray-400 dark:hover:border-zinc-600 focus:ring-2 focus:ring-primary-500 dark:focus:ring-primary-500/70 focus:border-primary-600 dark:focus:border-primary-500 focus:bg-gray-50 focus:ring-opacity-50 bg-white dark:bg-zinc-800/80 text-gray-800 dark:text-gray-200"
              placeholder="Search files..." 
              @keyup.enter="performSearch"
            />
          </span>
          <Button
            :loading="isSearching"
            severity="primary"
            @click="performSearch"
            v-slot="slotProps"
            asChild
          >
            <button
              v-bind="slotProps.a11yAttrs"
              class="ml-3 px-4 py-2 shadow-sm bg-black rounded-lg text-white dark:text-black dark:bg-white cursor-pointer"
            >
              Search
            </button>
          </Button>
        </div>
      </div>  

      <div class="border-b border-gray-200/70 dark:border-zinc-700/70 px-4 py-2 flex justify-between items-center bg-gray-50/80 dark:bg-zinc-800/60 backdrop-blur-sm">
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
          <span class="bg-gray-700 dark:bg-white text-white dark:text-black px-2 py-0.5 rounded-full text-xs font-semibold mr-2">
            {{ tempResults.length }}
          </span>
          files found
        </span>

        <div
          v-if="isSearching"
          class="flex items-center bg-primary-100 dark:bg-primary-900/40 text-primary-700 dark:text-primary-400 px-3 py-1 rounded-full"
        >
          <span class="text-xs font-medium mr-2">Searching...</span>
          <ProgressSpinner
            style="width:16px;height:16px"
            strokeWidth="4"
          />
        </div>
      </div>

      <div class="flex flex-1 h-[calc(100%-108px)]">
        <div class="w-1/2 border-r border-gray-200 dark:border-zinc-700/70 overflow-hidden bg-white/70 dark:bg-zinc-800/50 backdrop-blur-sm">
          <div
            v-if="searchPerformed"
            class="h-full"
          >
            <DynamicScroller
              :items="tempResults"
              :min-item-size="54"
              class="h-full custom-scroller"
            >
              <template #default="{ item, index, active }">
                <DynamicScrollerItem
                  :item="item"
                  :active="active"
                  :data-index="index"
                  :data-active="active"
                >
                  <FileListItem
                    :file-info="item"
                    :class="{'bg-primary-50 dark:bg-primary-900/30': fileInfo && fileInfo.path === item.path}"
                    @on-file-selection="handleFileSelection"
                  />
                </DynamicScrollerItem>
              </template>
            </DynamicScroller>
          </div>
          <div 
            v-else 
            class="h-full flex flex-col items-center justify-center text-gray-500 dark:text-gray-500 p-6"
          >
            <i class="pi pi-search text-4xl mb-4"></i>
            <p class="text-center">Enter a search term and press Search to find files</p>
          </div>
        </div>

        <div class="w-1/2 h-full bg-gray-50/90 dark:bg-zinc-850/70 backdrop-blur-sm overflow-auto">
          <FileInfo :file-info="fileInfo" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
    import { ref, computed, onMounted } from 'vue';
    import { Search, GetHomeDirectory } from "../../wailsjs/go/main/App.js";
    import { EventsOn } from "../../wailsjs/runtime/runtime.js";

    // PrimeVue Components
    import Button from 'primevue/button';
    import InputText from 'primevue/inputtext';
    import ProgressSpinner from 'primevue/progressspinner';
    //
    import FileInfo from './FileInfo.vue';
    import { main } from 'wailsjs/go/models';
    import FileListItem from './FileListItem.vue';

    // State
    const searchPattern = ref('');
    const searchDirectory = ref('');
    const includeHidden = ref(false);
    const isSearching = ref(false);
    const searchPerformed = ref(false);
    const selectedFile = ref<main.FileResult>()
    const tempResults = ref<main.FileResult[]>([])
    const pendingBatch = ref<main.FileResult[]>([])
    const systemDarkMode = ref(false);
    const resultsFound = ref(0)

    const BATCH_TIMEOUT = 500

    let batchTimeout: number | undefined

    const processBatch = () => {
      if (pendingBatch.value.length > 0) {
        tempResults.value = [...tempResults.value, ...pendingBatch.value]
        pendingBatch.value = []

        clearTimeout(batchTimeout)
        batchTimeout = undefined
      }
    }

    const handleNewEntry = (entry: main.FileResult) => {
      pendingBatch.value.push(entry)

      if (!batchTimeout) {
        batchTimeout = setTimeout(processBatch, BATCH_TIMEOUT)
      }
    }

    EventsOn("entry:found", (entry: main.FileResult) => {
      handleNewEntry(entry)
    });

    const fileInfo = computed(() => selectedFile.value || tempResults.value[0])

    // Methods
    const performSearch = async () => {
      if (!searchPattern.value) return;

      resultsFound.value = 0
      isSearching.value = true;
      searchPerformed.value = true;
      tempResults.value = [];

      try {
        const searchOptions = {
          pattern: searchPattern.value,
          directory: searchDirectory.value || '.',
          includeHidden: includeHidden.value,
        };

        await Search(searchOptions);
      } catch (error) {
        console.error('Search error:', error);
      } finally {
        isSearching.value = false;
      }
    };

    const handleFileSelection = (file: main.FileResult) => {
      selectedFile.value = file
    }

    // Check for system dark mode
    onMounted(() => {
      // Check for system dark mode preference
      if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        systemDarkMode.value = true;
      }

      // Listen for changes to system dark mode
      window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
        systemDarkMode.value = event.matches;
      });
    });

    // Init
    async function getHomeDir() {
      searchDirectory.value = await GetHomeDirectory();
    }

    function init() {
      getHomeDir();
    }

    init()
</script>

<style>
.dark .p-inputtext {
  background-color: transparent !important;
  color: #e4e4e7 !important;
  border-color: transparent !important;
}

.p-inputtext:focus {
  box-shadow: none !important;
  border-color: transparent !important;
}

.p-button.p-button-sm {
  padding: 0.4rem !important;
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: #d4d4d8;
  border-radius: 3px;
}

.dark ::-webkit-scrollbar-thumb {
  background: #52525b;
}

.dark .bg-zinc-850 {
  background-color: #18181b;
}

.dark .hover\:bg-zinc-750:hover {
  background-color: #2d2d33;
}
</style>