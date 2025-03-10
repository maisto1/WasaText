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
      },
      tempChatUser: null,
      isTempChat: false,
      creatingConversation: false,
      pollingInterval: null,
      pollingDelay: 1000
    }
  },

  async created() {
    await this.fetchConversations();
    this.startPolling();
  },
  
  beforeUnmount() {
    this.stopPolling();
  },

  methods: {
    async fetchConversations() {

      const isInitialLoad = this.conversations.length === 0;
      if (isInitialLoad) {
        this.loading = true;
      }
      
      try {
        const response = await this.$axios.get('/conversations/');

        if (this.selectedConversation) {
          const currentConvId = this.selectedConversation.id;
          this.conversations = response.data;
          if (currentConvId) {
            const updatedConv = this.conversations.find(c => c.id === currentConvId);
            if (updatedConv) {
              this.selectedConversation = updatedConv;
            }
          }
        } else {
          this.conversations = response.data;
        }
      } catch (error) {
        console.error('Error fetching conversations:', error);
        this.error = 'Failed to load conversations';
      } finally {
        if (isInitialLoad) {
          this.loading = false;
        }
      }
    },
    
    startPolling() {
      this.stopPolling();

      this.pollingInterval = setInterval(() => {
        if (!this.isSearching) {
          this.fetchConversations();
        }
      }, this.pollingDelay);
    },
    
    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null;
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

    selectUserForChat(user) {
      console.log('Setting up temporary chat with user:', user.username);
      
      this.tempChatUser = user;
      this.isTempChat = true;
      

      this.selectedConversation = {
        id: null,
        name: user.username,
        conversationType: 'private',
        conversationPhoto: user.profilePhoto
      };
      
      this.isSearching = false;
      this.searchQuery = '';
    },

    async createConversationAndSendMessage(messageContent) {
      console.log('createConversationAndSendMessage called with message:', messageContent);
      
      if (this.creatingConversation) {
        console.log('Already creating conversation, ignoring duplicate request');
        return false;
      }
      
      this.creatingConversation = true;
      
      try {
        if (!this.tempChatUser) {
          console.error('No temporary user selected');
          this.creatingConversation = false;
          return false;
        }
        
        console.log('Creating new conversation with user:', this.tempChatUser.username);
        
        const conversationResponse = await this.$axios.post('/conversations/', {
          conversationType: 'private',
          partecipant: this.tempChatUser.username,
          groupName: ''
        });
        
        const conversationId = conversationResponse.data.id;
        console.log('Conversation created with ID:', conversationId);
        
        const messageResponse = await this.$axios.post(`/conversations/${conversationId}/messages/`, {
          type: "text",
          content: messageContent
        });
        console.log('Message sent successfully:', messageResponse.data);
        
        await this.fetchConversations();
        
        const newConversation = this.conversations.find(c => c.id === conversationId);
        if (newConversation) {
          console.log('Selecting new conversation:', newConversation);
          this.selectedConversation = newConversation;
          this.isTempChat = false;
          this.tempChatUser = null;
        } else {
          console.error('Could not find newly created conversation in the list');
        }
        
        this.creatingConversation = false;
        return true;
      } catch (error) {
        console.error('Error creating conversation and sending message:', error);
        this.creatingConversation = false;
        return false;
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
      console.log('Selecting existing conversation:', conversation);
      this.selectedConversation = conversation;
      this.isSearching = false;
      this.isTempChat = false;
      this.tempChatUser = null;
    },

    startNewChat() {
      this.selectedConversation = null;
      this.isSearching = true;
      this.searchQuery = '';
      this.isTempChat = false;
      this.tempChatUser = null;
    }
  }
}
</script>

<template>
  <div class="chats-view">
    <div class="sidebar-container">
      <div class="search-header p-2 bg-dark border-bottom">
        <div class="d-flex justify-content-between align-items-center px-3 py-2">
          <h5 class="text-white mb-0">Chats</h5>
        </div>
        <SearchBar 
          v-model="searchQuery"
          @search="searchUsers"
        />
      </div>

      <div class="sidebar-content">
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
              @click="selectUserForChat(user)"
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
            <p class="text-secondary">Search for a user to start chatting</p>
          </div>
        </div>
      </div>
    </div>

    <div class="main-chat-container">
      <div v-if="!selectedConversation" class="empty-chat-container">
        <div class="text-center">
          <i class="fas fa-comments fa-3x text-secondary mb-3"></i>
          <p class="text-secondary">Select a conversation or search for a user to start chatting</p>
        </div>
      </div>
      <ChatSection 
        v-else 
        :conversation="selectedConversation"
        :is-temp-chat="isTempChat"
        @refresh-conversations="fetchConversations"
        @send-first-message="createConversationAndSendMessage"
      />
    </div>

    <div v-if="showCreateGroup" class="modal-overlay">
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

<style>
.chats-view {
  display: flex;
  height: calc(100vh - 48px);
  background-color: #111b21;
  position: relative;
}

.sidebar-container {
  width: 400px;
  min-width: 400px;
  background-color: #202c33;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #36434a;
  height: 100%;
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
}

.search-header {
  padding: 16px;
  background-color: #202c33;
  border-bottom: 1px solid #36434a;
}

.main-chat-container {
  flex: 1;
  background-color: #111b21;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.empty-chat-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.conversations-list, .search-results {
  height: 100%;
  overflow-y: auto;
}

.conversations-list::-webkit-scrollbar,
.search-results::-webkit-scrollbar {
  width: 6px;
}

.conversations-list::-webkit-scrollbar-track,
.search-results::-webkit-scrollbar-track {
  background: #202c33;
}

.conversations-list::-webkit-scrollbar-thumb,
.search-results::-webkit-scrollbar-thumb {
  background-color: #374045;
  border-radius: 6px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.btn-outline-light {
  border-color: #36434a;
}

.btn-outline-light:hover {
  background-color: #36434a;
  border-color: #36434a;
}

.user-list {
  padding: 8px;
}
</style>