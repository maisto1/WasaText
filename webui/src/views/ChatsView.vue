<script>
import ConversationItem from '../components/ConversationItem.vue'
import UserSearchResult from '../components/UserSearchResult.vue'
import SearchBar from '../components/SearchBar.vue'
import ChatSection from '../components/ChatSection.vue'
import ForwardModal from '../components/ForwardModal.vue'
import CreateGroupModal from '../components/CreateGroupModal.vue'
import GroupManagementModal from '../components/GroupManagementModal.vue'

export default {
  components: {
    ConversationItem,
    UserSearchResult,
    SearchBar,
    ChatSection,
    ForwardModal,
    CreateGroupModal,
    GroupManagementModal
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
      tempChatUser: null,
      isTempChat: false,
      creatingConversation: false,
      pollingInterval: null,
      pollingDelay: 1000,
      showGroupManagement: false,
      currentUser: {
        username: '',
        profilePhoto: null
      },
      showToast: false,
      toastMessage: '',
      toastType: 'success'
    }
  },

  async created() {
    this.loadCurrentUser();
    await this.fetchConversations();
    this.startPolling();
  },
  
  beforeUnmount() {
    this.stopPolling();
  },
  
  beforeRouteEnter(to, from, next) {
    next(vm => {
      if (from.path === '/me') {
        vm.loadCurrentUser();
      }
    });
  },

  methods: {
    loadCurrentUser() {
      this.currentUser.username = sessionStorage.getItem('username') || '';
    
      const profilePhoto = sessionStorage.getItem('profilePhoto');
      if (profilePhoto) {
        this.currentUser.profilePhoto = profilePhoto;
      }
    },
    
    goToProfile() {
      this.$router.push('/me');
    },
    
    async fetchConversations() {
      const isInitialLoad = this.conversations.length === 0;
      if (isInitialLoad) {
        this.loading = true;
      }
      
      try {
        const response = await this.$axios.get('/conversations/');

        if (this.selectedConversation) {
          const currentConvId = this.selectedConversation.id;
          
          const newConversations = response.data;
          
          if (currentConvId) {
            const updatedConv = newConversations.find(c => c.id === currentConvId);
            if (updatedConv) {
              this.selectedConversation = updatedConv;
            }
          }
          
          this.conversations = newConversations;
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
    openGroupManagement() {
      this.showGroupManagement = true;
    },

    handleGroupUpdated() {
      this.fetchConversations();
      this.showNotification('Group updated successfully', 'success');
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
    },
    
    getInitials(name) {
      if (!name) return '';
      return name.charAt(0).toUpperCase();
    },
    
    getAvatarColor(name) {
      if (!name) return '#202c33';
      
      let hash = 0;
      for (let i = 0; i < name.length; i++) {
        hash = name.charCodeAt(i) + ((hash << 5) - hash);
      }

      const h = hash % 360;
      const s = 65 + (hash % 20);
      const l = 45 + (hash % 10);
      
      return `hsl(${h}, ${s}%, ${l}%)`;
    },

    handleGroupCreated() {
      this.fetchConversations();
      this.showNotification('Group created successfully', 'success');
    },

    logout() {
      sessionStorage.clear();
      
      this.showToast = true;
      this.toastMessage = "Logout effettuato con successo!";
      this.toastType = "success";
      
      setTimeout(() => {
        this.$router.push('/login');
      }, 2000);
    },
    
    showNotification(message, type) {
      this.toastMessage = message;
      this.toastType = type;
      this.showToast = true;
      setTimeout(() => { 
        this.showToast = false; 
      }, 3000);
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
          <div class="d-flex align-items-center">
            <button @click="showCreateGroup = true" class="btn btn-new-group me-2">
              <i class="fas fa-users me-2"></i>
              New Group
            </button>
            <button @click="logout" class="btn btn-logout me-3">
              <i class="fas fa-sign-out-alt me-2"></i>
              Logout
            </button>
            <div class="profile-icon" @click="goToProfile">
              <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(currentUser.username) }">
                <img 
                  v-if="currentUser.profilePhoto" 
                  :src="'data:image/jpeg;base64,' + currentUser.profilePhoto"
                  class="avatar-image"
                  alt="Profile"
                />
                <div v-else class="avatar-text">
                  {{ getInitials(currentUser.username) }}
                </div>
              </div>
            </div>
          </div>
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
      :conversations="conversations"
      @refresh-conversations="fetchConversations"
      @send-first-message="createConversationAndSendMessage"
      @open-group-settings="openGroupManagement"
    />
    </div>

    <CreateGroupModal 
      :show="showCreateGroup" 
      @close="showCreateGroup = false"
      @created="handleGroupCreated"
    />

    <transition name="toast">
      <div v-if="showToast" class="toast-container position-fixed p-3">
        <div class="custom-toast" :class="toastType">
          <div class="toast-content">
            <i v-if="toastType === 'success'" class="fas fa-check-circle toast-icon"></i>
            <i v-else class="fas fa-exclamation-circle toast-icon"></i>
            <div class="toast-message">{{ toastMessage }}</div>
          </div>
          <div class="toast-progress" :class="toastType"></div>
        </div>
      </div>
    </transition>

    <GroupManagementModal 
      v-if="selectedConversation && selectedConversation.conversationType === 'group'"
      :show="showGroupManagement" 
      :conversation="selectedConversation"
      @close="showGroupManagement = false"
      @group-updated="handleGroupUpdated"
    />
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

.profile-icon {
  cursor: pointer;
  position: relative;
}

.profile-icon:hover .avatar-container {
  box-shadow: 0 0 0 2px #00a884;
}

.btn-logout {
  background-color: transparent;
  color: #e9edef;
  border: 1px solid #e74c3c;
  border-radius: 20px;
  padding: 0.25rem 0.75rem;
  font-size: 0.85rem;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
}

.btn-logout:hover {
  background-color: #e74c3c;
  color: white;
  box-shadow: 0 0 8px rgba(231, 76, 60, 0.7);
}

.btn-logout:active {
  transform: scale(0.95);
}

.btn-new-group {
  background-color: transparent;
  color: #e9edef;
  border: 1px solid #00a884;
  border-radius: 20px;
  padding: 0.25rem 0.75rem;
  font-size: 0.85rem;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
}

.btn-new-group:hover {
  background-color: #00a884;
  color: white;
  box-shadow: 0 0 8px rgba(0, 168, 132, 0.7);
}

.btn-new-group:active {
  transform: scale(0.95);
}

.toast-container {
  bottom: 20px;
  right: 20px;
  z-index: 9999;
}

.custom-toast {
  min-width: 300px;
  background: white;
  padding: 15px 20px;
  border-radius: 12px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.15);
  position: relative;
  overflow: hidden;
}

.toast-content {
  display: flex;
  align-items: center;
}

.toast-icon {
  font-size: 24px;
  margin-right: 12px;
}

.toast-message {
  font-size: 14px;
  font-weight: 500;
}

.success {
  border-left: 4px solid #2ecc71;
}

.error {
  border-left: 4px solid #e74c3c;
}

.success .toast-icon {
  color: #2ecc71;
}

.error .toast-icon {
  color: #e74c3c;
}

.toast-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 3px;
  width: 100%;
  animation: progress 3s linear forwards;
}

.toast-progress.success {
  background: #2ecc71;
}

.toast-progress.error {
  background: #e74c3c;
}

.avatar-container {
  width: 40px;
  height: 40px;
  min-width: 40px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-text {
  color: white;
  font-size: 1.2rem;
  font-weight: 600;
  text-align: center;
}

@keyframes progress {
  from { width: 100%; }
  to { width: 0%; }
}

.toast-enter-active {
  animation: slideIn 0.3s ease-out;
}

.toast-leave-active {
  animation: slideOut 0.3s ease-in;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes slideOut {
  from {
    transform: translateX(0);
    opacity: 1;
  }
  to {
    transform: translateX(100%);
    opacity: 0;
  }
}
</style>