<template>
	<div class="p-6 max-w-6xl mx-auto">
		<!-- Search Header -->
		<div class="mb-8">
			<h1 class="text-2xl font-bold text-gray-800 mb-2">File Finder</h1>
			<p class="text-gray-600">Fast file search powered by fd</p>
		</div>
		
		<!-- Search Form -->
		<div class="bg-white rounded-lg shadow-md p-6 mb-6">
			<div class="space-y-4">
				<!-- Search Pattern Input -->
				<div>
					<div class="flex">
						<input
							v-model="searchPattern"
							type="text"
							class="flex-1 px-4 py-2 border border-gray-300 rounded-l-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
							placeholder="Search for files..."
							@keyup.enter="performSearch"
						/>
						<button
							@click="performSearch"
							:disabled="isSearching"
							class="px-4 py-2 bg-blue-600 text-white rounded-r-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-blue-400"
						>
							<span v-if="isSearching">
								<svg class="animate-spin h-5 w-5 mr-2 inline-block" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
								Searching...
						    </span>
							<span v-else>Search</span>
						</button>
					</div>
					<p class="mt-1 text-xs text-gray-500">
						Supports regex patterns like "\.js$" for JavaScript files
					</p>
				</div>
				
				<!-- Directory Selector -->
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Search Directory</label>
					<div class="flex">
						<input
							v-model="searchDirectory"
							type="text"
							class="flex-1 px-4 py-2 border border-gray-300 rounded-l-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
							placeholder="/path/to/search"
						/>
						<button
							@click="browseDirectory"
							class="px-4 py-2 bg-gray-200 text-gray-800 rounded-r-lg hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-400"
						>
							Browse
						</button>
					</div>
				</div>
				
				<!-- Search Options -->
				<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
					<div class="flex items-center">
						<input
							id="includeHidden"
							v-model="includeHidden"
							type="checkbox"
							class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
						/>
						<label for="includeHidden" class="ml-2 block text-sm text-gray-700">
							Include hidden files
						</label>
					</div>
					
					<div class="flex items-center">
						<input
							id="caseSensitive"
							v-model="caseSensitive"
							type="checkbox"
							class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
						/>
						<label for="caseSensitive" class="ml-2 block text-sm text-gray-700">
							Case sensitive
						</label>
					</div>
					
					<div>
						<label for="fileType" class="block text-sm text-gray-700">File type</label>
						<select
							id="fileType"
							v-model="fileType"
							class="mt-1 block w-full pl-3 pr-10 py-2 text-base border border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
						>
							<option value="">All</option>
							<option value="f">Files only</option>
							<option value="d">Directories only</option>
						</select>
					</div>
				</div>
			</div>
		</div>
		
		<!-- Results Section -->
		<div v-if="searchPerformed" class="bg-white rounded-lg shadow-md p-6">
			<div class="flex justify-between items-center mb-4">
				<h2 class="text-lg font-medium text-gray-800">
					{{ results.length }} Results
				</h2>
				
				<div class="flex space-x-2">
					<select 
						v-model="sortField" 
						@change="sortResults"
						class="block pl-3 pr-10 py-1 text-sm border border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 rounded-md"
					>
						<option value="path">Path</option>
						<option value="name">Name</option>
						<option value="modified">Modified</option>
						<option value="size">Size</option>
					</select>
					
					<button 
						@click="sortDirection = sortDirection === 'asc' ? 'desc' : 'asc'; sortResults()"
						class="p-1 rounded border border-gray-300 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
					>
						<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" :class="{ 'rotate-180': sortDirection === 'desc' }">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
						</svg>
					</button>
				</div>
			</div>
			
			<!-- Loading Indicator -->
			<div v-if="isSearching" class="py-20 flex flex-col items-center justify-center text-gray-500">
				<svg class="animate-spin h-10 w-10 mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				<p>Searching your files...</p>
			</div>
			
			<!-- No Results -->
			<div v-else-if="results.length === 0" class="py-20 flex flex-col items-center justify-center text-gray-500">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
				<p>No files matching your search criteria</p>
			</div>
			
			<!-- Results Table -->
			<div v-else class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Name
							</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Path
							</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Size
							</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Modified
							</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-white divide-y divide-gray-200">
						<tr v-for="result in results" :key="result.path" class="hover:bg-gray-50">
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="flex items-center">
									<svg v-if="result.isDir" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
									</svg>
									<svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
									</svg>
									<span class="text-sm font-medium text-gray-900">{{ result.name }}</span>
								</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="text-sm text-gray-500 truncate max-w-xs">{{ result.directory }}</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="text-sm text-gray-500">{{ formatSize(result.size) }}</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="text-sm text-gray-500">{{ formatDate(result.modTime) }}</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
								<button 
									@click="openFile(result.path)" 
									class="text-blue-600 hover:text-blue-900 mr-2"
								>
									Open
								</button>
								<button 
									@click="showInFolder(result.path)" 
									class="text-gray-600 hover:text-gray-900"
								>
									Show in folder
								</button>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
	</div>
</template>

<script setup>
	import { ref, computed, onMounted } from 'vue';
    import { Search } from "../../wailsjs/go/main/App.js"
	
	// State
	const searchPattern = ref('');
	const searchDirectory = ref('');
	const includeHidden = ref(false);
	const caseSensitive = ref(false);
	const fileType = ref('');
	const isSearching = ref(false);
	const searchPerformed = ref(false);
	const rawResults = ref([]);
	const results = computed(() => sortedResults.value);
	
	// Sorting
	const sortField = ref('name');
	const sortDirection = ref('asc');
	const sortedResults = computed(() => {
		const sorted = [...rawResults.value].sort((a, b) => {
			let comparison = 0;
			switch (sortField.value) {
				case 'name':
					comparison = a.name.localeCompare(b.name);
					break;
				case 'path':
					comparison = a.path.localeCompare(b.path);
					break;
				case 'size':
					comparison = a.size - b.size;
					break;
				case 'modified':
					comparison = new Date(a.modTime) - new Date(b.modTime);
					break;
				default:
					comparison = a.name.localeCompare(b.name);
			}
			
			return sortDirection.value === 'asc' ? comparison : -comparison;
		});
		
		return sorted;
	});
	
	// Initialization
	onMounted(async () => {
		// Get user's home directory as default search location
		searchDirectory.value = await window.go.main.App.GetHomeDirectory();
	});
	
	// Methods
	const performSearch = async () => {
		if (!searchPattern.value) return;
		
		isSearching.value = true;
		searchPerformed.value = true;
		
		try {
			const searchOptions = {
				pattern: searchPattern.value,
				directory: searchDirectory.value,
				caseSensitive: caseSensitive.value,
				includeHidden: includeHidden.value,
				fileTypes: fileType.value ? [fileType.value] : [],
			};
			
			// Call to Go backend via Wails
			rawResults.value = await Search(searchOptions)
		} catch (error) {
			console.error('Search error:', error);
			// TODO: Show error to user
		} finally {
			isSearching.value = false;
		}
	};
	
	const browseDirectory = async () => {
		try {
			const selectedDir = await window.go.main.App.BrowseDirectory();
			if (selectedDir) {
				searchDirectory.value = selectedDir;
			}
		} catch (error) {
			console.error('Directory browse error:', error);
		}
	};
	
	const sortResults = () => {
		// Sorting is handled by the computed property
	};
	
	const openFile = async (path) => {
		try {
			await window.go.main.App.OpenFile(path);
		} catch (error) {
			console.error('Error opening file:', error);
		}
	};
	
	const showInFolder = async (path) => {
		try {
			await window.go.main.App.ShowInFolder(path);
		} catch (error) {
			console.error('Error showing in folder:', error);
		}
	};
	
	// Utility functions
	const formatSize = (bytes) => {
		if (bytes === 0) return '0 B';
		
		const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
		const i = Math.floor(Math.log(bytes) / Math.log(1024));
		return parseFloat((bytes / Math.pow(1024, i)).toFixed(2)) + ' ' + sizes[i];
	};
	
	const formatDate = (dateString) => {
		const date = new Date(dateString);
		return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
	};
</script>