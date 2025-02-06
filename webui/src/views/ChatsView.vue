// views/ChatsView.vue
<script>
import ConversationItem from '../components/ConversationItem.vue'
import UserSearchResult from '../components/UserSearchResult.vue'
import SearchBar from '../components/SearchBar.vue'
import ChatSection from '../components/ChatSection.vue'

export default {
  components: {
    ConversationItem,
    UserSearchResult,
    SearchBar,
    ChatSection
  },
  
  data() {
    return {
      conversations: [],
      loading: false,
      error: null,
      selectedConversation: null,
      searchQuery: "",
      searchResults: [],
      isSearching: false,
      searchLoading: false,
      showCreateGroup: false,
      newGroupData: {
        name: '',
        participants: []
      }
    }
  },

  async created() {
    await this.fetchConversations();
  },

  methods: {
    async fetchConversations() {
      this.loading = true;
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

    async startPrivateChat(user) {
      try {
        await this.$axios.post('/conversations/', {
          conversationType: 'private',
          partecipant: user.username
        });
        await this.fetchConversations();
        this.isSearching = false;
        this.searchQuery = '';
      } catch (error) {
        console.error('Error starting private chat:', error);
      }
    },

    async createGroup() {
      if (!this.newGroupData.name || this.newGroupData.participants.length === 0) return;
      
      try {
        await this.$axios.post('/conversations/', {
          groupName: this.newGroupData.name,
          conversationType: 'group',
          participants: this.newGroupData.participants
        });
        this.showCreateGroup = false;
        this.newGroupData = { name: '', participants: [] };
        await this.fetchConversations();
      } catch (error) {
        console.error('Error creating group:', error);
      }
    },

    selectConversation(conversation) {
      this.selectedConversation = conversation;
      this.isSearching = false;
    },

    startNewChat() {
      this.selectedConversation = null;
      this.isSearching = true;
      this.searchQuery = '';
    }
  }
}
</script>

<template>
  <div class="chat-container">
    <!-- Sidebar -->
    <div class="sidebar">
      <!-- Search and New Chat Section -->
      <div class="search-header p-2 bg-dark border-bottom">
        <div class="d-flex justify-content-between align-items-center mb-2">
          <h5 class="text-white mb-0">Chats</h5>
          <button @click="startNewChat" class="btn btn-outline-light btn-sm">
            <i class="fas fa-plus"></i> New Chat
          </button>
        </div>
        <SearchBar 
          v-model="searchQuery"
          @search="searchUsers"
        />
      </div>

      <!-- Search Results or Conversations List -->
      <div v-if="loading" class="text-center p-4">
        <div class="spinner-border text-light" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>

      <div v-else-if="isSearching" class="search-results">
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
            @click="startPrivateChat(user)"
          />
          <div v-if="searchResults.length === 0 && searchQuery" class="text-center text-light p-4">
            <p>No users found</p>
          </div>
        </div>
      </div>

      <div v-else class="conversations-list">
        <ConversationItem
          v-for="conv in conversations"
          :key="conv.id"
          :conversation="conv"
          :isActive="selectedConversation?.id === conv.id"
          @click="selectConversation(conv)"
        />
        <div v-if="conversations.length === 0" class="text-center text-light p-4">
          <p>No conversations yet</p>
          <button @click="startNewChat" class="btn btn-outline-light">
            Start a new chat
          </button>
        </div>
      </div>
    </div>

    <!-- Main Chat Area -->
    <div class="main-chat">
      <div v-if="!selectedConversation" class="empty-chat-container">
        <div class="text-center">
          <i class="fas fa-comments fa-3x text-secondary mb-3"></i>
          <p class="text-secondary">Select a conversation to start chatting</p>
        </div>
      </div>
      <ChatSection 
        v-else 
        :conversation="selectedConversation"
        @refresh-conversations="fetchConversations"
      />
    </div>

    <!-- Create Group Modal -->
    <div v-if="showCreateGroup" class="modal fade show" style="display: block;">
      <div class="modal-dialog">
        <div class="modal-content bg-dark text-white">
          <div class="modal-header">
            <h5 class="modal-title">Create New Group</h5>
            <button @click="showCreateGroup = false" class="btn-close btn-close-white"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label class="form-label">Group Name</label>
              <input v-model="newGroupData.name" type="text" class="form-control bg-secondary text-white">
            </div>
            <!-- Add participant selection here -->
          </div>
          <div class="modal-footer">
            <button @click="showCreateGroup = false" class="btn btn-secondary">Cancel</button>
            <button @click="createGroup" class="btn btn-primary">Create Group</button>
          </div>
        </div>
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

.search-header {
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
  flex-direction: column;
}

.empty-chat-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.conversations-list::-webkit-scrollbar {
  width: 6px;
}

.conversations-list::-webkit-scrollbar-track {
  background: #202c33;
}

.conversations-list::-webkit-scrollbar-thumb {
  background-color: #374045;
  border-radius: 6px;
}

.modal {
  background-color: rgba(0, 0, 0, 0.5);
}

.search-results {
  flex: 1;
  overflow-y: auto;
}

.user-list {
  padding: 8px;
}

.btn-outline-light {
  border-color: #36434a;
}

.btn-outline-light:hover {
  background-color: #36434a;
  border-color: #36434a;
}
</style>