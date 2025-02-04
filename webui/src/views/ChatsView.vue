<script>
import ConversationItem from '../components/ConversationItem.vue'
import UserSearchResult from '../components/UserSearchResult.vue'
import SearchBar from '../components/SearchBar.vue'

export default {
  components: {
    ConversationItem,
    UserSearchResult,
    SearchBar
  },
  
  data() {
    return {
      conversations: [],
      loading: true,
      error: null,
      selectedConversation: null,
      searchQuery: "",
      searchResults: [],
      isSearching: false,
      searchLoading: false
    }
  },

  async created() {
    await this.fetchConversations();
  },

  methods: {
    async fetchConversations() {
      try {
        const response = await this.$axios.get('/conversations/');
        this.conversations = response.data;
      } catch (error) {
        console.error('Error fetching conversations:', error);
        this.error = 'Failed to load conversations';
      } finally {
        this.loading = false;
      }
    },

    async searchUsers() {
      if (!this.searchQuery) {
        this.searchResults = [];
        this.isSearching = false;
        return;
      }
      
      this.isSearching = true;
      this.searchLoading = true;
      try {
        const response = await this.$axios.get(`/users/?username=${this.searchQuery}`);
        this.searchResults = response.data;
      } catch (error) {
        console.error('Error searching users:', error);
      } finally {
        this.searchLoading = false;
      }
    },

    selectConversation(conversation) {
      this.selectedConversation = conversation;
    }
  }
}
</script>

<template>
  <div class="chat-container">
    <div class="sidebar">
      <SearchBar 
        v-model="searchQuery"
        @search="searchUsers"
      />

      <div class="sidebar-header" v-if="!isSearching">
        <h5 class="text-white mb-0">Chats</h5>
      </div>

      <div v-if="loading && !isSearching" class="text-center p-4">
        <div class="spinner-border text-light" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>

      <div v-if="isSearching" class="search-results">
        <div v-if="searchLoading" class="text-center p-4">
          <div class="spinner-border text-light spinner-border-sm" role="status">
            <span class="visually-hidden">Searching...</span>
          </div>
        </div>
        <div v-else class="user-list">
          <UserSearchResult 
            v-for="user in searchResults"
            :key="user.id"
            :user="user"
          />
          <div v-if="searchResults.length === 0" class="text-center text-light p-4">
            <p>No users found</p>
          </div>
        </div>
      </div>

      <div v-else-if="!loading && !error" class="conversations-list">
        <ConversationItem
          v-for="conv in conversations"
          :key="conv.id"
          :conversation="conv"
          :isActive="selectedConversation?.id === conv.id"
          @click="selectConversation(conv)"
        />
      </div>
    </div>

    <div class="main-chat">
      <div v-if="!selectedConversation" class="empty-chat-container">
        <p class="text-secondary">Select a conversation to start chatting</p>
      </div>
      <div v-else class="chat-content">
        <!-- Chat messages view implementation -->
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  background-color: #111b21;
  position: fixed;
  top: 0;
  left: 0;
}

.sidebar {
  width: 400px;
  background-color: #202c33;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #36434a;
}

.sidebar-header {
  padding: 16px;
  background-color: #202c33;
  border-bottom: 1px solid #36434a;
}

.conversations-list {
  flex: 1;
  overflow-y: auto;
}

.main-chat {
  flex: 1;
  background-color: #111b21;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-chat-container {
  text-align: center;
}

.empty-chat-container p {
  margin: 0;
  color: #8696a0;
}

.conversations-list::-webkit-scrollbar {
  width: 0px;
}

.conversations-list {
  scrollbar-width: none;
}
</style>