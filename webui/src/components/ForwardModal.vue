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
      searchQuery: '',
      itemsPerPage: 15,
      currentPage: 1,
      isLoading: false
    }
  },
  
  computed: {
    filteredConversations() {
      if (!this.searchQuery.trim()) {
 
        return this.paginateResults(this.conversations);
      }
      

      const searchTerms = this.searchQuery.toLowerCase().split(/\s+/).filter(term => term.length > 0);
      

      const filteredResults = this.conversations.filter(conv => {
        const nameToSearch = conv.name.toLowerCase();
        

        return searchTerms.some(term => nameToSearch.includes(term));
      });
      

      return this.paginateResults(filteredResults);
    },
    
    totalPages() {
      if (!this.searchQuery.trim()) {
        return Math.ceil(this.conversations.length / this.itemsPerPage);
      } else {
        const searchTerms = this.searchQuery.toLowerCase().split(/\s+/).filter(term => term.length > 0);
        const totalFilteredItems = this.conversations.filter(conv => {
          const nameToSearch = conv.name.toLowerCase();
          return searchTerms.some(term => nameToSearch.includes(term));
        }).length;
        
        return Math.ceil(totalFilteredItems / this.itemsPerPage);
      }
    },
    
    showPagination() {
      return this.totalPages > 1;
    },
    
    currentPageDisplay() {
      return `${this.currentPage} / ${this.totalPages}`;
    }
  },
  
  methods: {
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
    
    closeModal() {
      this.$emit('close');
      this.selectedConversationId = null;
      this.searchQuery = '';
      this.currentPage = 1;
    },
    
    handleForward() {
      if (!this.selectedConversationId) return;
      
      this.$emit('forward', this.messageId, this.selectedConversationId);
      this.closeModal();
    },
    
    selectConversation(conversationId) {
      this.selectedConversationId = conversationId;
    },
    
    paginateResults(results) {
      const startIndex = (this.currentPage - 1) * this.itemsPerPage;
      const endIndex = startIndex + this.itemsPerPage;
      return results.slice(startIndex, endIndex);
    },
    
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
        this.scrollToTop();
      }
    },
    
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
        this.scrollToTop();
      }
    },
    
    scrollToTop() {

      this.isLoading = true;
      setTimeout(() => {
        const container = this.$refs.conversationListContainer;
        if (container) {
          container.scrollTop = 0;
        }
        this.isLoading = false;
      }, 150);
    },
    
    handleSearchInput() {

      this.currentPage = 1;
    }
  },
  
  watch: {
    searchQuery() {
      this.handleSearchInput();
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
            <button v-if="searchQuery" @click="searchQuery = ''" class="clear-search">
              <i class="fas fa-times-circle"></i>
            </button>
          </div>
        </div>
        
        <div class="modal-body" ref="conversationListContainer">
          <div v-if="isLoading" class="loading-indicator">
            <div class="spinner"></div>
          </div>
          
          <div v-else class="conversation-list">
            <div 
              v-for="conversation in filteredConversations" 
              :key="conversation.id"
              :class="['conversation-item', {'selected': selectedConversationId === conversation.id}]"
              @click="selectConversation(conversation.id)"
            >
              <div class="d-flex align-items-center">
                <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(conversation.name) }">
                  <img v-if="conversation.conversationPhoto" 
                       :src="'data:image/jpeg;base64,' + conversation.conversationPhoto"
                       class="avatar-image"
                       alt="Profile photo">
                  <div v-else class="avatar-text">
                    {{ getInitials(conversation.name) }}
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
              <small v-if="searchQuery" class="hint">Try different keywords</small>
            </div>
          </div>
          
          <div v-if="showPagination" class="pagination-controls">
            <button 
              @click="prevPage" 
              class="pagination-button"
              :disabled="currentPage === 1"
            >
              <i class="fas fa-chevron-left"></i>
            </button>
            <span class="page-indicator">{{ currentPageDisplay }}</span>
            <button 
              @click="nextPage" 
              class="pagination-button"
              :disabled="currentPage === totalPages"
            >
              <i class="fas fa-chevron-right"></i>
            </button>
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
  font-size: 0.9rem;
}

.search-input {
  width: 100%;
  padding: 10px 36px 10px 40px;
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

.clear-search {
  position: absolute;
  right: 12px;
  background: transparent;
  border: none;
  color: #8696a0;
  cursor: pointer;
  padding: 0;
  font-size: 0.9rem;
}

.clear-search:hover {
  color: #e9edef;
}

.modal-body {
  padding: 0;
  max-height: 350px;
  overflow-y: auto;
  position: relative;
}

.loading-indicator {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 3px solid rgba(0, 168, 132, 0.3);
  border-radius: 50%;
  border-top-color: #00a884;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
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

.avatar-text {
  color: white;
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

.hint {
  display: block;
  margin-top: 8px;
  font-size: 0.8rem;
  opacity: 0.8;
}

.pagination-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 12px;
  background-color: #111b21;
  border-top: 1px solid #2a3942;
}

.pagination-button {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: transparent;
  border: none;
  color: #e9edef;
  cursor: pointer;
  transition: background-color 0.2s;
}

.pagination-button:hover:not(:disabled) {
  background-color: #2a3942;
}

.pagination-button:disabled {
  color: #8696a0;
  cursor: not-allowed;
}

.page-indicator {
  margin: 0 12px;
  color: #8696a0;
  font-size: 0.9rem;
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

/* Scrollbar styling */
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