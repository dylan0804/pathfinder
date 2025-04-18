<template>
  <div
    v-if="!fileInfo"
    class="h-full flex items-center justify-center p-8 bg-white/40 dark:bg-zinc-800 backdrop-blur-sm"
  >
    <div class="text-center p-6 rounded-lg">
      <i class="pi pi-file-o text-5xl text-gray-300 dark:text-gray-600 mb-4"></i>
      <p class="text-gray-400 dark:text-gray-500">Select a file to view details</p>
    </div>
  </div>
  <div
    v-else
    class="h-full flex flex-col bg-white/40 dark:bg-zinc-800"
  >
    <!-- Main file preview area -->
    <div class="flex-1 flex flex-col items-stretch justify-between">
      <div class="relative w-64 h-80 m-auto transition-all duration-300 hover:scale-105">
        <!-- File icon/preview -->
        <div class="bg-white/90 dark:bg-zinc-100/90 rounded-lg shadow-md w-full h-full flex flex-col items-center justify-center p-4 relative overflow-hidden">
          <!-- The blurred content inside the file -->
          <div class="w-4/5 h-10 bg-gray-200 dark:bg-zinc-600 mb-4 blur-sm"></div>
          <div class="w-4/5 h-24 bg-gray-200 dark:bg-zinc-600 mb-8 blur-sm"></div>
            
          <!-- The folded corner -->
          <div
            class="absolute top-0 right-0 w-16 h-16 bg-gray-100/90 dark:bg-zinc-200/90"
            style="clip-path: polygon(100% 0, 100% 100%, 0 100%)"
          ></div>
            
          <!-- File type label -->
          <div class="text-4xl font-bold text-gray-800 mt-6">
            {{ fileExtension }}
          </div>
        </div>
        
        <!-- Shadow effect under the file -->
        <div class="absolute -bottom-2 left-4 right-4 h-8 bg-black/20 dark:bg-black/40 blur-md rounded-full z-[-1]"></div>
      </div>

      <!-- Additional file info -->
      <div class="px-6 pb-4 pt-2 flex justify-between items-center">
        <div class="flex items-center">
          <span
            v-if="fileInfo.isDir"
            class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-blue-100/90 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300"
          >
            <i class="pi pi-folder mr-1"></i> Directory
          </span>
          <span 
            v-else
            class="inline-flex items-center px-2.5 py-2 rounded-full text-xs font-medium bg-indigo-100/90 dark:bg-indigo-900/40 text-indigo-800 dark:text-indigo-300"
          >
            <i class="pi pi-file mr-1"></i> File
          </span>
        </div>
        
        <button 
          v-if="!fileInfo.isDir" 
          class="inline-flex items-center px-3 py-1.5 bg-black dark:bg-white text-white dark:text-black cursor-pointer hover:bg-primary-600/80 dark:text-black rounded-full text-sm"
          @click="openFile"
        >
          <i class="pi pi-external-link mr-1"></i> Open
        </button>
      </div>
    </div>
    
    <!-- File information area -->
    <div class="mt-auto bg-white/80 dark:bg-zinc-800/80 backdrop-blur-sm rounded-t-xl border-t border-gray-200/70 dark:border-gray-700/70">
      <!-- File name row -->
      <div class="py-4 px-6 flex items-center border-b border-gray-200/50 dark:border-gray-700/50">
        <div class="text-gray-500 dark:text-gray-400 w-20 flex-shrink-0">Name</div>
        <div class="text-gray-800 dark:text-gray-200 font-medium truncate">{{ fileInfo.name }}</div>
      </div>
        
      <!-- File location row -->
      <div class="py-4 px-6 flex items-center border-b border-gray-200/50 dark:border-gray-700/50">
        <div class="text-gray-500 dark:text-gray-400 w-20 flex-shrink-0">Location</div>
        <div class="text-gray-800 dark:text-gray-200 truncate text-sm">{{ fileInfo.path }}</div>
      </div>
        
      <!-- File size row -->
      <div class="py-4 px-6 flex items-center">
        <div class="text-gray-500 dark:text-gray-400 w-20 flex-shrink-0">Size</div>
        <div class="text-gray-800 dark:text-gray-200">{{ formatFileSize(fileInfo.size) }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { OpenFile } from "../../wailsjs/go/main/App.js";

const props = defineProps<{
  fileInfo: {
    id: number;
    path: string;
    name: string;
    extension: string;
    isDir: boolean;
    size: number;
  }
}>()

// Extract file extension from name or extension property
const fileExtension = computed(() => {
  if (!props.fileInfo) return '';
  
  // If it's a directory
  if (props.fileInfo.isDir) return 'DIR';
  
  // Get extension from the extension property or from filename
  const ext = props.fileInfo.extension || '';
  if (ext.startsWith('.')) {
    return ext.substring(1).toUpperCase();
  }
  return ext.toUpperCase() || 'FILE';
});

// Format file size with appropriate units
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B';
  if (!bytes) return 'Unknown';
  
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const openFile = async () => {
  try {
    await OpenFile(props.fileInfo.path)
  } catch(e) {
    console.log(e)
  }
}
</script>