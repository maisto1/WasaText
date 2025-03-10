<script>
export default {
  props: {
    show: {
      type: Boolean,
      required: true
    },
    conversations: {
      type: Array,
      required: true
    },
    messageId: {
      type: Number,
      default: null
    }
  },
  
  data() {
    return {
      selectedConversationId: null,
      searchQuery: ''
    }
  },
  
  computed: {
    filteredConversations() {
      if (!this.searchQuery) return this.conversations;
      const query = this.searchQuery.toLowerCase();
      return this.conversations.filter(conv => 
        conv.name.toLowerCase().includes(query)
      );
    }
  },
  
  methods: {
    closeModal() {
      this.$emit('close');
      this.selectedConversationId = null;
      this.searchQuery = '';
    },
    
    handleForward() {
      if (!this.selectedConversationId) return;
      
      this.$emit('forward', this.messageId, this.selectedConversationId);
      this.closeModal();
    },
    
    selectConversation(conversationId) {
      this.selectedConversationId = conversationId;
    }
  }
}
</script>

<template>
  <div v-if="show" class="modal-overlay">
    <div class="modal-dialog" @click.stop>
      <div class="modal-content">
        <div class="modal-header">
          <div class="d-flex align-items-center">
            <i class="fas fa-share modal-icon"></i>
            <h5 class="modal-title ms-2">Forward message</h5>
          </div>
          <button @click="closeModal" class="close-button">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="search-container">
          <div class="search-wrapper">
            <i class="fas fa-search search-icon"></i>
            <input 
              v-model="searchQuery"
              type="text"
              class="search-input"
              placeholder="Search conversations..."
            >
          </div>
        </div>
        
        <div class="modal-body">
          <div class="conversation-list">
            <div 
              v-for="conversation in filteredConversations" 
              :key="conversation.id"
              :class="['conversation-item', {'selected': selectedConversationId === conversation.id}]"
              @click="selectConversation(conversation.id)"
            >
              <div class="d-flex align-items-center">
                <div class="avatar-container">
                  <img v-if="conversation.conversationPhoto" 
                       :src="'data:image/jpeg;base64,' + conversation.conversationPhoto"
                       class="avatar-image"
                       alt="Profile photo">
                  <div v-else class="avatar-placeholder">
                    {{ conversation.name.charAt(0).toUpperCase() }}
                  </div>
                </div>
                
                <div class="conversation-info">
                  <h6 class="conversation-name">{{ conversation.name }}</h6>
                  <span v-if="selectedConversationId === conversation.id" class="selected-indicator">
                    <i class="fas fa-check"></i>
                  </span>
                </div>
              </div>
            </div>
            
            <div v-if="filteredConversations.length === 0" class="empty-state">
              <i class="fas fa-search fa-2x mb-2"></i>
              <p>No conversations found</p>
            </div>
          </div>
        </div>
        
        <div class="modal-footer">
          <button @click="closeModal" class="cancel-button">
            Cancel
          </button>
          <button 
            @click="handleForward" 
            class="forward-button"
            :class="{'disabled': !selectedConversationId}"
            :disabled="!selectedConversationId"
          >
            <i class="fas fa-paper-plane me-2"></i>
            Forward
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1050;
  backdrop-filter: blur(4px);
}

.modal-dialog {
  width: 420px;
  max-width: 90%;
  max-height: 90vh;
  animation: fadeInUp 0.3s ease;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-content {
  border-radius: 16px;
  background-color: #111b21;
  border: 1px solid #2a3942;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}

.modal-header {
  padding: 16px 20px;
  background-color: #202c33;
  border-bottom: 1px solid #2a3942;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-icon {
  color: #00a884;
  font-size: 1.2rem;
}

.modal-title {
  font-size: 1.1rem;
  font-weight: 600;
  margin: 0;
  color: #e9edef;
}

.close-button {
  background: transparent;
  border: none;
  color: #8696a0;
  font-size: 1.1rem;
  cursor: pointer;
  transition: color 0.2s;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
}

.close-button:hover {
  color: #e9edef;
  background-color: #2a3942;
}

.search-container {
  padding: 12px 16px;
  background-color: #202c33;
  border-bottom: 1px solid #2a3942;
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: #8696a0;
  z-index: 1;
}

.search-input {
  width: 100%;
  padding: 10px 10px 10px 40px;
  background-color: #2a3942;
  border: none;
  border-radius: 8px;
  color: #e9edef;
  font-size: 0.9rem;
}

.search-input:focus {
  outline: none;
  box-shadow: 0 0 0 2px #00a884;
}

.search-input::placeholder {
  color: #8696a0;
}

.modal-body {
  padding: 0;
  max-height: 350px;
  overflow-y: auto;
}

.conversation-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.conversation-item {
  padding: 12px 20px;
  transition: background-color 0.2s ease;
  cursor: pointer;
  border-bottom: 1px solid #2a3942;
  position: relative;
}

.conversation-item:hover {
  background-color: #2a3942;
}

.conversation-item.selected {
  background-color: rgba(0, 168, 132, 0.1);
}

.avatar-container {
  width: 48px;
  height: 48px;
  min-width: 48px;
  border-radius: 50%;
  overflow: hidden;
  background-color: #2a3942;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  color: #e9edef;
  font-size: 1.2rem;
  font-weight: 600;
}

.conversation-info {
  margin-left: 12px;
  flex-grow: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.conversation-name {
  margin: 0;
  font-size: 1rem;
  color: #e9edef;
  font-weight: 500;
}

.selected-indicator {
  color: #00a884;
  font-size: 1.2rem;
}

.empty-state {
  padding: 24px;
  text-align: center;
  color: #8696a0;
}

.modal-footer {
  padding: 16px 20px;
  background-color: #202c33;
  border-top: 1px solid #2a3942;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.cancel-button {
  padding: 8px 16px;
  background-color: transparent;
  border: 1px solid #2a3942;
  border-radius: 4px;
  color: #e9edef;
  font-size: 0.9rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.cancel-button:hover {
  background-color: #2a3942;
}

.forward-button {
  padding: 8px 16px;
  background-color: #00a884;
  border: none;
  border-radius: 4px;
  color: #fff;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
}

.forward-button:hover {
  background-color: #008f6f;
}

.forward-button.disabled {
  background-color: rgba(0, 168, 132, 0.5);
  cursor: not-allowed;
}

.modal-body::-webkit-scrollbar {
  width: 6px;
}

.modal-body::-webkit-scrollbar-track {
  background: transparent;
}

.modal-body::-webkit-scrollbar-thumb {
  background-color: #374045;
  border-radius: 6px;
}
</style>