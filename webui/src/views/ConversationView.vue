<script>
export default {
 data: function() {
   return {
     searchQuery: "",
     users: [],
     loading: false,
     defaultImage: '../public/default-icon.png'
   }
 },
 methods: {
   async searchUsers() {
     if (!this.searchQuery) {
       this.users = []
       return
     }
     this.loading = true
     try {
       const token = sessionStorage.getItem('authToken')
       const response = await this.$axios.get(`/users/`, {
         params: {
           username: this.searchQuery
         },
         headers: {
           'Authorization': token
         }
       })
       console.log("Search results:", response.data)
       this.users = response.data
     } catch (e) {
       console.error("Error while searching users:", e)
       this.users = []
     } finally {
       this.loading = false
     }
   },
   handleImageError(e) {
     e.target.src = this.defaultImage
   }
 }
}
</script>

<template>
<div class="conversations-container">
 <div class="search-bar">
   <div class="search-wrapper">
     <svg
       class="search-icon"
       xmlns="http://www.w3.org/2000/svg"
       width="20"
       height="20"
       viewBox="0 0 24 24"
       fill="none"
       stroke="#808080"
       stroke-width="2"
       stroke-linecap="round"
       stroke-linejoin="round"
     >
       <circle cx="11" cy="11" r="8"></circle>
       <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
     </svg>
     <input
       v-model="searchQuery"
       @input="searchUsers"
       type="text"
       placeholder="Search or start a new conversation"
       class="search-input"
     >
   </div>
 </div>

 <div v-if="loading" class="loading-indicator">
   Searching...
 </div>

 <div v-else-if="users.length > 0" class="search-results">
   <div v-for="user in users"
        :key="user.id"
        class="user-item"
   >
     <img 
       :src="user.photo ? `data:image/jpeg;base64,${user.photo}` : defaultImage"
       :alt="user.username"
       class="user-avatar"
       @error="handleImageError"
     >
     <div class="user-info">
       <div class="username">{{ user.username }}</div>
     </div>
   </div>
 </div>
</div>
</template>

<style>
.loading-indicator {
 color: #808080;
 text-align: center;
 padding: 20px;
}

.conversations-container {
 width: 100%;
 max-width: 400px;
 height: 100vh;
 border-right: 1px solid #2c2c2c;
 background-color: #111;
}

.search-bar {
 padding: 10px;
 background-color: #202020;
}

.search-wrapper {
 position: relative;
 display: flex;
 align-items: center;
 background-color: #323232;
 border-radius: 8px;
 padding: 0 12px;
}

.search-icon {
 min-width: 20px;
 margin-right: 8px;
}

.search-input {
 width: 100%;
 padding: 8px 0;
 border: none;
 background-color: transparent;
 color: white;
 outline: none;
}

.search-input::placeholder {
 color: #808080;
}

.search-results {
 overflow-y: auto;
 max-height: calc(100vh - 60px);
}

.user-item {
 display: flex;
 align-items: center;
 padding: 12px;
 cursor: pointer;
 border-bottom: 1px solid #2c2c2c;
}

.user-item:hover {
 background-color: #202020;
}

.user-avatar {
 width: 40px;
 height: 40px;
 border-radius: 50%;
 margin-right: 12px;
 object-fit: cover;
}

.user-info {
 display: flex;
 flex-direction: column;
}

.username {
 color: white;
 font-weight: 500;
}

</style>