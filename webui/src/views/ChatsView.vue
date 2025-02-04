<script>
export default {
  data() {
    return {
      conversations: [],
      loading: true,
      error: null,
      selectedConversation: null
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
    formatTimestamp(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp * 1000);
      return date.toLocaleString();
    },
    getMessagePreview(message) {
      if (!message) return '';
      return message.type === 'text' ? message.content : '🖼️ Photo';
    },
    selectConversation(conversation) {
      this.selectedConversation = conversation;
    }
  }
}
</script>

<template>
  <div class="chat-container">
    <!-- Sidebar with conversations -->
    <div class="sidebar">
      <!-- Sidebar Header -->
      <div class="sidebar-header">
        <h5 class="text-white mb-0">Chats</h5>
      </div>

      <!-- Loading state -->
      <div v-if="loading" class="text-center p-4">
        <div class="spinner-border text-light" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="p-3">
        <div class="alert alert-danger" role="alert">
          {{ error }}
        </div>
      </div>

      <!-- Conversations list -->
      <div v-else class="conversations-list">
        <div v-for="conv in conversations" 
             :key="conv.id" 
             @click="selectConversation(conv)"
             :class="['conversation-item p-3', {'active': selectedConversation?.id === conv.id}]">
          <div class="d-flex align-items-center">
            <!-- Profile/Group photo -->
            <div class="rounded-circle bg-secondary flex-shrink-0" 
                 style="width: 40px; height: 40px;">
              <img v-if="conv.conversationPhoto" 
                   :src="'data:image/jpeg;base64,' + conv.conversationPhoto"
                   class="rounded-circle w-100 h-100"
                   alt="Profile photo">
            </div>
            
            <div class="ms-3 overflow-hidden">
              <div class="d-flex justify-content-between align-items-center">
                <h6 class="mb-1 text-truncate text-white">{{ conv.name }}</h6>
                <small class="text-light ms-2" v-if="conv.latestMessage">
                  {{ formatTimestamp(conv.latestMessage.timestamp) }}
                </small>
              </div>
              <p class="mb-0 text-light-grey small text-truncate" v-if="conv.latestMessage">
                {{ getMessagePreview(conv.latestMessage) }}
              </p>
            </div>
          </div>
        </div>

        <!-- Empty state -->
        <div v-if="!loading && !error && conversations.length === 0" 
             class="text-center text-light p-4">
          <p>No conversations yet</p>
        </div>
      </div>
    </div>

    <!-- Main chat area -->
    <div class="main-chat">
      <div v-if="!selectedConversation" class="empty-chat-container">
        <p class="text-secondary">Select a conversation to start chatting</p>
      </div>
      <div v-else class="chat-content">
        <!-- TODO: Implement chat messages view -->
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

.conversation-item {
  cursor: pointer;
  transition: background-color 0.2s ease;
  border-bottom: 1px solid #36434a;
}

.conversation-item:hover {
  background-color: #2a3942;
}

.conversation-item.active {
  background-color: #2a3942;
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

.text-light-grey {
  color: #8696a0;
}

/* Rimuove la scrollbar ma mantiene la funzionalità di scorrimento */
.conversations-list::-webkit-scrollbar {
  width: 0px;
}

.conversations-list {
  scrollbar-width: none;
}

/* Assicura che il contenuto non vada sotto la navbar */
.empty-chat-container p {
  margin: 0;
  color: #8696a0;
}
</style>