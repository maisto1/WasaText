<script>
import ForwardModal from './ForwardModal.vue';

export default {
  components: {
    ForwardModal
  },
  
  props: {
    message: {
      type: Object,
      required: true
    },
    isMyMessage: {
      type: Boolean,
      required: true
    },
    showUsername: {
      type: Boolean,
      default: false 
    },
    availableConversations: {
      type: Array,
      default: () => []
    }
  },

  data() {
    return {
      showActions: false,
      imageLoaded: false,
      imageError: false,
      showForwardModal: false
    }
  },

  computed: {
    conversationId() {
      if (this.message.conversationId) {
        return this.message.conversationId;
      }
      
      const currentRoute = this.$route || {};
      if (currentRoute.params && currentRoute.params.id) {
        return currentRoute.params.id;
      }
      
      const path = window.location.hash;
      const match = path.match(/\/chats\/(\d+)/);
      return match ? match[1] : null;
    },
    
    groupedEmojis() {
      if (!this.comments || this.comments.length === 0) return [];
      
      const emojiGroups = {};
      this.comments.forEach(comment => {
        const emoji = comment.content;
        if (!emojiGroups[emoji]) {
          emojiGroups[emoji] = {
            emoji: emoji,
            count: 0,
            commentIds: [],
            usernames: [],
            myComment: null
          };
        }
        
        emojiGroups[emoji].count++;
        emojiGroups[emoji].commentIds.push(comment.id);
        emojiGroups[emoji].usernames.push(comment.sender.username);
        
        if (this.isMyComment(comment)) {
          emojiGroups[emoji].myComment = comment;
        }
      });
      
      return Object.values(emojiGroups);
    },

    myEmojiReactions() {
      const myUsername = sessionStorage.getItem('username');
      return this.comments
        .filter(comment => comment.sender.username === myUsername)
        .map(comment => comment.content);
    }
  },

  created() {
    this.fetchComments();
    this.startCommentsPolling();
  },
  
  beforeUnmount() {
    this.stopCommentsPolling();
  },

  methods: {
    formatTimestamp(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp);
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    },

    handleDelete() {
      if (confirm('Are you sure you want to delete this message?')) {
        this.$emit('delete', this.message.id);
      }
    },

    handleForwardClick() {
      console.log('Opening forward modal...');
      if (this.availableConversations.length === 0) {
        alert('No conversations available to forward to');
      } else {
        this.showForwardModal = true;
        console.log('Forward modal should be open now:', this.showForwardModal);
      }
    },

    handleForward(messageId, targetConversationId) {
      console.log('Forwarding message', messageId, 'to conversation', targetConversationId);
      this.$emit('forward', messageId, targetConversationId);
      this.showForwardModal = false;
    },

    closeForwardModal() {
      console.log('Closing forward modal');
      this.showForwardModal = false;
    },

    handleImageLoad() {
      this.imageLoaded = true;
    },

    handleImageError() {
      this.imageError = true;
    },

    openMediaInNewTab() {
      if (this.message.media) {
        const imgWindow = window.open();
        imgWindow.document.write(`
          <html>
            <head>
              <title>WasaText Media</title>
              <style>
                body { 
                  margin: 0; 
                  padding: 0; 
                  display: flex; 
                  justify-content: center; 
                  align-items: center; 
                  height: 100vh; 
                  background-color: #121212; 
                }
                img { 
                  max-width: 90%; 
                  max-height: 90vh; 
                  object-fit: contain; 
                }
              </style>
            </head>
            <body>
              <img src="data:image/jpeg;base64,${this.message.media}" alt="Full size media" />
            </body>
          </html>
        `);
      }
    },

    toggleEmojiPicker() {
      this.showEmojiPicker = !this.showEmojiPicker;
      this.debugInfo = null;
      

      if (this.showEmojiPicker) {
        this.showActions = false;
      }
    },

    async fetchComments() {
      if (!this.message?.id || !this.conversationId) {
        console.warn('Cannot fetch comments: message id or conversation id is missing');
        return;
      }
      
      try {
        console.log(`Fetching comments for message ${this.message.id} in conversation ${this.conversationId}`);
        const response = await this.$axios.get(
          `/conversations/${this.conversationId}/messages/${this.message.id}/comments/`
        );
        console.log('Fetched comments:', response.data);
        this.comments = response.data;
      } catch (error) {
        console.error('Error fetching comments:', error);
        
        if (error.response && error.response.status !== 404) {
          this.debugInfo = `Error fetching comments: ${error.message}`;
        }
      }
    },
    
    async addComment(emoji) {
      if (this.submittingComment) return;
      
      if (!this.message?.id || !this.conversationId) {
        console.error('Cannot add comment: message id or conversation id is missing');
        this.debugInfo = 'Cannot add comment: message id or conversation id is missing';
        return;
      }
      
      if (this.myEmojiReactions.includes(emoji)) {

        const commentToDelete = this.comments.find(
          comment => this.isMyComment(comment) && comment.content === emoji
        );
        
        if (commentToDelete) {
          this.deleteComment(commentToDelete.id);
        }
        return;
      }
      
      this.submittingComment = true;
      this.showEmojiPicker = false;
      
      try {
        console.log(`Adding emoji comment "${emoji}" to message ${this.message.id} in conversation ${this.conversationId}`);
        const payload = { content: emoji };
        console.log('Comment payload:', payload);
        
        const response = await this.$axios.post(
          `/conversations/${this.conversationId}/messages/${this.message.id}/comments/`,
          payload
        );
        
        console.log('Comment added successfully:', response.data);
        this.comments.push(response.data);
        this.debugInfo = null;
      } catch (error) {
        console.error('Error adding comment:', error);
        this.debugInfo = `Error adding comment: ${error.message}`;
        
        if (error.response) {
          console.error('Error response:', error.response.data);
          this.debugInfo += ` - Status: ${error.response.status}`;
        }
      } finally {
        this.submittingComment = false;
      }
    },
    
    async deleteComment(commentId) {
      if (!this.message?.id || !this.conversationId) {
        console.error('Cannot delete comment: message id or conversation id is missing');
        return;
      }
      
      try {
        console.log(`Deleting comment ${commentId} from message ${this.message.id} in conversation ${this.conversationId}`);
        
        await this.$axios.delete(
          `/conversations/${this.conversationId}/messages/${this.message.id}/comments/${commentId}`
        );
        
        console.log('Comment deleted successfully');
        this.comments = this.comments.filter(comment => comment.id !== commentId);
      } catch (error) {
        console.error('Error removing comment:', error);
        this.debugInfo = `Error removing comment: ${error.message}`;
      }
    },
    
    isMyComment(comment) {
      const myUsername = sessionStorage.getItem('username');
      return comment.sender && comment.sender.username === myUsername;
    },
    
    handleEmojiClick(emojiGroup) {
      if (emojiGroup.myComment) {

        this.deleteComment(emojiGroup.myComment.id);
      } else {

        this.addComment(emojiGroup.emoji);
      }
    },
    

    formatUsernames(usernames) {
      if (!usernames || usernames.length === 0) return '';
      
      if (usernames.length <= 3) {
        return usernames.join(', ');
      } else {

        return `${usernames.slice(0, 3).join(', ')} e altri ${usernames.length - 3}`;
      }
    },
    

    startCommentsPolling() {
      this.stopCommentsPolling();
      
      this.commentsPollingInterval = setInterval(() => {
        if (this.message?.id && this.conversationId) {
          this.pollForComments();
        }
      }, this.commentsPollingDelay);
    },
    
    stopCommentsPolling() {
      if (this.commentsPollingInterval) {
        clearInterval(this.commentsPollingInterval);
        this.commentsPollingInterval = null;
      }
    },
    
    async pollForComments() {
      if (!this.message?.id || !this.conversationId) return;
      
      try {
        const response = await this.$axios.get(
          `/conversations/${this.conversationId}/messages/${this.message.id}/comments/`
        );
        
        const newComments = response.data;
        if (!newComments) return;
        

        const hasChanges = this.hasCommentsChanges(this.comments, newComments);
        
        if (hasChanges) {
          console.log('Comments updated:', newComments);
          this.comments = newComments;
        }
      } catch (error) {

        if (error.response && error.response.status !== 404) {
          console.error('Error polling for comments:', error);
        }
      }
    },
    
    hasCommentsChanges(oldComments, newComments) {
      if (!oldComments || !newComments) return true;
      if (oldComments.length !== newComments.length) return true;
      

      const oldIds = new Set(oldComments.map(c => c.id));
      const newIds = new Set(newComments.map(c => c.id));
      

      return [...oldIds].some(id => !newIds.has(id)) || [...newIds].some(id => !oldIds.has(id));
    }
  }
}
</script>

<template>
  <div 
    :class="['message mb-3 d-flex', isMyMessage ? 'justify-content-end' : 'justify-content-start']"
  >
    <div 
      :class="['message-bubble p-2 rounded position-relative', 
               isMyMessage ? 'bg-primary text-white' : 'bg-secondary text-white']"
      :style="[message.type === 'media' ? {'max-width': '300px'} : {'max-width': '70%'}]"
      @mouseenter="showActions = true"
      @mouseleave="showActions = false"
    >

      <div v-if="debugInfo" class="debug-info">
        <small>{{ debugInfo }}</small>
      </div>
    

      <div v-if="message.isForwarded" class="forwarded-label">
        <i class="fas fa-share"></i> Forwarded
      </div>
      
      <div v-if="!isMyMessage && showUsername" class="message-sender mb-1">
        <small class="text-light">{{ message.sender.username }}</small>
      </div>

      <div class="message-content">
        <template v-if="message.type === 'media' && message.media">
          <div class="media-container mb-2">
            <div v-if="!imageLoaded && !imageError" class="text-center p-2">
              <div class="spinner-border spinner-border-sm text-light" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
            </div>
            
            <div v-if="imageError" class="text-center p-2">
              <i class="fas fa-exclamation-circle"></i>
              <small class="d-block">Failed to load image</small>
            </div>
            
            <img 
              v-show="!imageError"
              :src="'data:image/jpeg;base64,' + message.media"
              class="img-fluid rounded cursor-pointer"
              alt="Media content"
              @load="handleImageLoad"
              @error="handleImageError"
              @click="openMediaInNewTab"
            />
          </div>
          
          <p v-if="message.content" class="mb-1">{{ message.content }}</p>
        </template>
        
        <template v-else-if="message.type === 'text'">
          <p class="mb-1">{{ message.content }}</p>
        </template>
      </div>


      <div v-if="groupedEmojis.length > 0" class="comments-display">
        <div class="emoji-reactions">
          <div 
            v-for="emojiGroup in groupedEmojis" 
            :key="emojiGroup.emoji" 
            class="emoji-reaction-bubble"
            :class="{ 'my-reaction': emojiGroup.myComment }"
            @click="handleEmojiClick(emojiGroup)"
            :title="formatUsernames(emojiGroup.usernames)"
          >
            <span class="emoji">{{ emojiGroup.emoji }}</span>
            <span class="count">{{ emojiGroup.count }}</span>
          </div>
        </div>
      </div>

      <div class="message-meta d-flex justify-content-end align-items-center">
        <small class="text-light me-1">{{ formatTimestamp(message.timestamp) }}</small>
        <span v-if="isMyMessage" class="message-status">
          <i v-if="message.status === 'read'" class="fas fa-check-double" title="Read"></i>
          <i v-else class="fas fa-check" title="Sent"></i>
        </span>
      </div>

      <div v-if="showActions" class="message-actions position-absolute top-0 end-0 mt-1 me-1">
        <div class="btn-group">
          <button 
            @click.stop="toggleEmojiPicker"
            class="action-button emoji-action"
            title="React with emoji"
          >
            <i class="far fa-smile"></i>
          </button>
          
          <button 
            v-if="isMyMessage"
            @click.stop="handleDelete" 
            class="btn btn-sm btn-danger"
          >
            <i class="fas fa-trash"></i>
          </button>
          <button 
            @click.stop="handleForwardClick"
            class="btn btn-sm btn-info"
          >
            <i class="fas fa-share"></i>
          </button>
        </div>
      </div>


      <div v-if="showEmojiPicker" 
           :class="['emoji-picker', isMyMessage ? 'emoji-picker-user' : 'emoji-picker-other']"
           @click.stop>
        <div class="emoji-grid">
          <button 
            v-for="emoji in commonEmojis" 
            :key="emoji"
            class="emoji-item"
            :class="{ 'already-used': myEmojiReactions.includes(emoji) }"
            @click.stop="addComment(emoji)"
          >
            {{ emoji }}
          </button>
        </div>
        <div v-if="submittingComment" class="emoji-loading">
          <div class="spinner-border spinner-border-sm text-light" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
      </div>
    </div>
    

    <ForwardModal
      v-if="showForwardModal"
      :show="showForwardModal"
      :conversations="availableConversations"
      :message-id="message.id"
      @close="closeForwardModal"
      @forward="handleForward"
    />
  </div>
</template>

<style scoped>
.message-bubble {
  border-radius: 12px;
  padding: 8px 12px;
  word-break: break-word;
}

.message-actions {
  transform: translateY(-100%);
  opacity: 0.9;
  z-index: 10;
}

.message-status {
  font-size: 0.8em;
  color: rgba(255, 255, 255, 0.8);
}

.message-status .fa-check {
  color: rgba(255, 255, 255, 0.6);
}

.message-status .fa-check-double {
  color: #34B7F1;
}

.btn-group .btn {
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
}

.media-container {
  position: relative;
  min-height: 50px;
}

.cursor-pointer {
  cursor: pointer;
}

.media-container img {
  max-height: 300px;
  object-fit: contain;
  width: 100%;
}

.forwarded-label {
  color: #8696a0;
  font-size: 0.75rem;
  margin-bottom: 4px;
  font-style: italic;
}
</style>