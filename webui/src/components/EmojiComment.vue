<script>
export default {
  props: {
    message: {
      type: Object,
      required: true
    },
    isMyMessage: {
      type: Boolean,
      required: true
    }
  },
  
  data() {
    return {
      showEmojiPicker: false,
      comments: [],
      loading: false,
      error: null,
      selectedEmoji: '',
      submittingComment: false,
      showToast: false,
      toastMessage: '',
      toastType: 'success',
      commonEmojis: ['👍', '❤️', '😂', '😮', '😢', '🙏', '🔥', '👏', '🎉', '🤔']
    }
  },
  
  created() {
    this.fetchComments();
  },
  
  methods: {
    async fetchComments() {
      if (!this.message?.id) return;
      
      this.loading = true;
      
      try {
        const response = await this.$axios.get(
          `/conversations/${this.message.conversationId}/messages/${this.message.id}/comments/`
        );
        this.comments = response.data;
      } catch (error) {
        console.error('Error fetching comments:', error);
        this.error = 'Failed to load comments';
      } finally {
        this.loading = false;
      }
    },
    
    toggleEmojiPicker() {
      this.showEmojiPicker = !this.showEmojiPicker;
    },
    
    async addComment(emoji) {
      if (this.submittingComment) return;
      
      this.submittingComment = true;
      this.showEmojiPicker = false;
      
      try {
        const response = await this.$axios.post(
          `/conversations/${this.message.conversationId}/messages/${this.message.id}/comments/`,
          {
            content: emoji
          }
        );
        
        this.comments.push(response.data);
        this.showNotification('Comment added successfully', 'success');
      } catch (error) {
        console.error('Error adding comment:', error);
        this.showNotification('Failed to add comment', 'error');
      } finally {
        this.submittingComment = false;
      }
    },
    
    async deleteComment(commentId) {
      try {
        await this.$axios.delete(
          `/conversations/${this.message.conversationId}/messages/${this.message.id}/comments/${commentId}`
        );
        
        this.comments = this.comments.filter(comment => comment.id !== commentId);
        this.showNotification('Comment removed successfully', 'success');
      } catch (error) {
        console.error('Error removing comment:', error);
        this.showNotification('Failed to remove comment', 'error');
      }
    },
    
    canDeleteComment(comment) {
      const myUsername = sessionStorage.getItem('username');
      return comment.sender.username === myUsername;
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
  <div class="emoji-comments-container">

    <div v-if="comments.length > 0" class="comments-display">
      <div class="emoji-reactions">
        <div 
          v-for="comment in comments" 
          :key="comment.id" 
          class="emoji-reaction-bubble"
          :class="{ 'my-reaction': canDeleteComment(comment) }"
          @click="canDeleteComment(comment) ? deleteComment(comment.id) : null"
        >
          <span class="emoji">{{ comment.content }}</span>
          <span v-if="canDeleteComment(comment)" class="remove-hint">
            <i class="fas fa-times"></i>
          </span>
        </div>
      </div>
    </div>
    

    <div class="emoji-comment-actions">
      <button 
        @click="toggleEmojiPicker" 
        class="emoji-comment-button"
        v-if="!submittingComment"
      >
        <i class="far fa-smile"></i>
      </button>
      
      <div v-if="submittingComment" class="emoji-loading">
        <div class="spinner-border spinner-border-sm text-light" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      
      <div v-if="showEmojiPicker" class="emoji-picker" :class="{ 'right-aligned': isMyMessage }">
        <div class="emoji-grid">
          <button 
            v-for="emoji in commonEmojis" 
            :key="emoji"
            class="emoji-item"
            @click="addComment(emoji)"
          >
            {{ emoji }}
          </button>
        </div>
      </div>
    </div>
    

    <transition name="toast">
      <div v-if="showToast" class="comment-toast-container">
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
  </div>
</template>

<style scoped>
.emoji-comments-container {
  position: relative;
  margin-top: 4px;
}

.comments-display {
  margin-top: 6px;
  margin-bottom: 6px;
}

.emoji-reactions {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.emoji-reaction-bubble {
  display: inline-flex;
  align-items: center;
  background-color: #202c33;
  border-radius: 15px;
  padding: 3px 8px;
  cursor: default;
  position: relative;
  transition: background-color 0.2s;
}

.emoji-reaction-bubble.my-reaction {
  background-color: #2a3942;
  cursor: pointer;
}

.emoji-reaction-bubble.my-reaction:hover {
  background-color: #36434a;
}

.emoji {
  font-size: 1.1rem;
}

.remove-hint {
  font-size: 0.6rem;
  margin-left: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.my-reaction:hover .remove-hint {
  opacity: 1;
}

.emoji-comment-actions {
  display: flex;
  margin-top: 5px;
  position: relative;
}

.emoji-comment-button {
  background: transparent;
  border: none;
  color: #8696a0;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0;
  transition: color 0.2s;
}

.emoji-comment-button:hover {
  color: #e9edef;
}

.emoji-loading {
  font-size: 0.8rem;
}

.emoji-picker {
  position: absolute;
  bottom: 100%;
  left: 0;
  background-color: #2a3942;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  padding: 8px;
  margin-bottom: 8px;
  z-index: 10;
  animation: fade-in 0.2s ease;
}

.emoji-picker.right-aligned {
  left: auto;
  right: 0;
}

@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.emoji-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 6px;
}

.emoji-item {
  background: transparent;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 6px;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.emoji-item:hover {
  background-color: #36434a;
}

.comment-toast-container {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-bottom: 8px;
  z-index: 1000;
}

.custom-toast {
  min-width: 180px;
  background: white;
  padding: 10px 12px;
  border-radius: 8px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.15);
  position: relative;
  overflow: hidden;
}

.toast-content {
  display: flex;
  align-items: center;
}

.toast-icon {
  font-size: 18px;
  margin-right: 8px;
}

.toast-message {
  font-size: 12px;
  font-weight: 500;
}

.toast-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 3px;
  width: 100%;
  animation: progress 3s linear forwards;
}

.success {
  border-left: 3px solid #2ecc71;
}

.error {
  border-left: 3px solid #e74c3c;
}

.success .toast-icon {
  color: #2ecc71;
}

.error .toast-icon {
  color: #e74c3c;
}

.toast-progress.success {
  background: #2ecc71;
}

.toast-progress.error {
  background: #e74c3c;
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
    opacity: 0;
    transform: translateY(-10px) translateX(-50%);
  }
  to {
    opacity: 1;
    transform: translateY(0) translateX(-50%);
  }
}

@keyframes slideOut {
  from {
    opacity: 1;
    transform: translateY(0) translateX(-50%);
  }
  to {
    opacity: 0;
    transform: translateY(-10px) translateX(-50%);
  }
}
</style>