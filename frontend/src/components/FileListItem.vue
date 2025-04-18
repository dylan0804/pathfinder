<template>
  <div class="file-item w-full px-2">
    <div 
      class="flex items-center cursor-pointer px-3 py-2 rounded-md transition-all duration-200 w-full hover:bg-gray-200 dark:hover:bg-gray-600/60 group"
      :class="{'bg-gray-200/60 dark:bg-gray-600/60': isExpanded}"
      @click="handleFileClick"
    >     
      <div class="w-6 flex justify-center items-center">
        <i
          v-if="fileInfo.isDir"
          class="pi transition-transform duration-200 text-gray-400"
          :class="isExpanded ? 'pi-angle-down' : 'pi-angle-right'"
        ></i>
        <span
          v-else
          class="w-6"
        ></span>
      </div>
      
      <div class="flex items-center justify-center w-8 h-8 mx-2">
        <img
          :src="fileIcon"
          class="w-6 h-6 drop-shadow-sm transition-all duration-200 group-hover:drop-shadow-md"
          :alt="fileInfo.extension || 'file'"
        />
      </div>
      
      <div class="flex flex-col min-w-0 flex-1">
        <p class="text-black dark:text-white font-medium truncate">{{ fileInfo.name }}</p>
        <div class="flex items-center">
          <span
            v-if="!fileInfo.isDir"
            class="text-xs text-gray-400 group-hover:text-gray-400 dark:group-hover:text-gray-400"
          >
            {{ fileInfo.extension || 'file' }}
          </span>
          <span 
            v-if="!fileInfo.isDir && fileInfo.size" 
            class="text-xs text-gray-500 ml-2 group-hover:text-gray-400"
          >
            · {{ formatFileSize(fileInfo.size) }}
          </span>
        </div>
      </div>
      
      <div class="opacity-0 group-hover:opacity-100 transition-opacity duration-200 flex items-center gap-1">
        <button
          class="py-2 px-3 text-gray-400 hover:text-white hover:bg-gray-600/70 rounded-full cursor-pointer"
          @click.stop="copyToClipboard"
        >
          <i class="pi pi-copy text-xs"></i>
        </button>
        
        <button
          v-if="!fileInfo.isDir"
          class="py-2 px-3 text-gray-400 hover:text-white hover:bg-gray-600/70 rounded-full cursor-pointer"
          @click.stop="openFile"
        >
          <i class="pi pi-external-link text-xs"></i>
        </button>
      </div>
    </div>
    
    <div 
      v-if="isExpanded"
      class="overflow-hidden transition-all duration-300"
    >
      <div
        v-if="files?.length > 0"
        class="flex flex-col ml-6 pl-2 border-l border-gray-700/50 mt-1"
      >
        <FileListItem
          v-for="file in files"
          :key="file.path" 
          :file-info="file"
          class="my-0.5"
          @on-file-selection="$emit('on-file-selection', $event)"
        />
      </div>
      
      <div
        v-if="isLoading"
        class="ml-12 py-2 text-gray-400 flex items-center"
      >
        <i class="pi pi-spinner animate-spin mr-2 text-blue-400"></i>
        <span class="text-sm">Loading contents...</span>
      </div>
      
      <div
        v-if="!isLoading && files?.length === 0"
        class="ml-12 py-3 text-gray-500 flex items-center bg-gray-800/30 rounded-md px-3 my-1.5"
      >
        <i class="pi pi-inbox mr-2 text-gray-400"></i>
        <span class="text-sm">This folder is empty</span>
      </div>
    </div>
  </div>
</template>
  
<script setup lang="ts">
    import { ref, computed, watch, } from 'vue';
    import { SearchSubdirectory, OpenFile } from "../../wailsjs/go/main/App.js";
    import { main } from 'wailsjs/go/models';

    import pdfIcon from "../assets/images/pdf.svg"
    import xlsIcon from "../assets/images/xls.svg"
    import unknownIcon from "../assets/images/unknown.svg"

    import { useToast } from 'primevue/usetoast';
    
    const props = defineProps<{
        fileInfo: {
            id: number
            path: string;
            name: string;
            extension: string;
            isDir: boolean;
            size: number;
        }
    }>();

    const emit = defineEmits<{
      (e: 'on-file-selection', file: main.FileResult): void
    }>()

    const toast = useToast()

    watch(() => props.fileInfo, () => {
      files.value = []
      isExpanded.value = false
      isLoading.value = false
    })

    const isExpanded = ref(false);
    const isLoading = ref(false);
    const files = ref<main.FileResult[]>([]);

    const fileIcon = computed(() => {        
        switch(props.fileInfo.extension.toLowerCase()) {
            case ".pdf":
                return pdfIcon;
            case ".xls":
            case ".xlsx":
                return xlsIcon;
            default:
                return unknownIcon;
        }
    });

    const formatFileSize = (bytes: number) => {
      if (!bytes) return '';
      if (bytes === 0) return '0 B';
      
      const k = 1024;
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
    };

    const handleFileClick = async () => {
        emit('on-file-selection', props.fileInfo)

        if (!props.fileInfo.isDir) return;

        isExpanded.value = !isExpanded.value;
        
        if (!isExpanded.value) return;
        
        try {
          isLoading.value = true;
          const result = await SearchSubdirectory(props.fileInfo.path);
          files.value = result === null ? [] : result
        } catch (e) {
          console.log("Error loading directory:", e);
        } finally {
          isLoading.value = false;
        }
    };

    const openFile = async () => {
      try {
        await OpenFile(props.fileInfo.path)
      } catch(e) {
        console.log(e)
      }
    }

    const copyToClipboard = async () => {
      try {
        await navigator.clipboard.writeText(props.fileInfo.path)
        toast.add({ 
          severity: 'info', 
          summary: 'File path copied', 
          detail: 'Path copied to clipboard',
          life: 2750 
        });
      } catch (e) {
        console.log(e)
      }
    }
</script>