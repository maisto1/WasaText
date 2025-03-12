<script>
import ForwardModal from './ForwardModal.vue';
import ReplyModal from './ReplyModal.vue';

export default {
  components: {
    ForwardModal,
    ReplyModal
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
      showForwardModal: false,
      showReplyModal: false,
      showEmojiPicker: false,
      comments: [],
      loading: false,
      submittingComment: false,
      debugInfo: null,
      commonEmojis: ['👍', '❤️', '😂', '😮', '😢', '🙏', '🔥', '👏', '🎉', '🤔'],
      commentsPollingInterval: null,
      commentsPollingDelay: 2000
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
    
    handleReplyClick() {
      console.log('Opening reply modal...');
      this.showReplyModal = true;
    },
    
    handleSendReply(originalMessage, replyData) {
      console.log('Sending reply to message:', originalMessage.id, 'with data:', replyData);
      this.$emit('send-reply', originalMessage, replyData);
      this.showReplyModal = false;
    },

    closeForwardModal() {
      console.log('Closing forward modal');
      this.showForwardModal = false;
    },
    
    closeReplyModal() {
      console.log('Closing reply modal');
      this.showReplyModal = false;
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
      

      <div v-if="message.replyTo" class="reply-info mb-2">
        <div class="reply-preview">
          <div class="reply-sender">
            {{ message.replyTo.sender }}
          </div>
          <div class="reply-content">
            <template v-if="message.replyTo.content">
              {{ message.replyTo.content }}
            </template>
            <template v-else>
              <i class="fas fa-camera me-1"></i> Photo
            </template>
          </div>
        </div>
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


      <div v-if="showActions" :class="['message-actions-below', isMyMessage ? 'user-message-actions' : 'other-message-actions']">
        <div class="action-button-container">
          <button 
            @click.stop="handleReplyClick" 
            class="action-button reply-action"
            title="Reply"
          >
            <i class="fas fa-reply"></i>
          </button>
          
          <button 
            @click.stop="handleForwardClick"
            class="action-button forward-action"
            title="Forward"
          >
            <i class="fas fa-share"></i>
          </button>
          
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
            class="action-button delete-action"
            title="Delete"
          >
            <i class="fas fa-trash"></i>
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
    

    <ReplyModal
      v-if="showReplyModal"
      :show="showReplyModal"
      :message="message"
      @close="closeReplyModal"
      @send-reply="handleSendReply"
    />
  </div>
</template>

<style scoped>
.message-bubble {
  border-radius: 12px;
  padding: 8px 12px;
  word-break: break-word;
}


.message-actions-below {
  position: absolute;
  bottom: -45px;
  z-index: 10;
}


.user-message-actions {
  left: 40%; 
  transform: translateX(-50%);
}


.other-message-actions {
  left: 50%;
  transform: translateX(-50%);
}

.action-button-container {
  display: flex;
  background-color: #202c33;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
  padding: 2px;
}

.action-button {
  background: transparent;
  border: none;
  color: #fff;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 2px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.action-button:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.reply-action {
  color: #00a884;
}

.forward-action {
  color: #34B7F1;
}

.emoji-action {
  color: #ffd700;
}

.delete-action {
  color: #f15c6d;
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

.reply-info {
  padding: 8px;
  border-radius: 6px;
  background-color: rgba(0, 0, 0, 0.2);
  position: relative;
  border-left: 3px solid #34B7F1;
  margin-top: 4px;
}

.reply-preview {
  font-size: 0.85rem;
  overflow: hidden;
}

.reply-sender {
  font-weight: 600;
  color: #34B7F1;
  margin-bottom: 2px;
}

.reply-content {
  color: rgba(255, 255, 255, 0.8);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
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
  cursor: pointer;
  position: relative;
  transition: background-color 0.2s;
}

.emoji-reaction-bubble:hover {
  background-color: #2a3942;
}

.emoji-reaction-bubble.my-reaction {
  background-color: #36434a;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.emoji-reaction-bubble.my-reaction:hover {
  background-color: #445258;
}

.emoji {
  font-size: 1.1rem;
  margin-right: 4px;
}

.count {
  font-size: 0.85rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
}


.emoji-picker {
  position: absolute;
  bottom: -65px;
  background-color: #2a3942;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  padding: 8px;
  z-index: 20;
}


.emoji-picker-user {
  left: 40%;
  transform: translateX(-50%);
  animation: fade-in-user 0.2s ease;
}


.emoji-picker-other {
  left: 50%;
  transform: translateX(-50%);
  animation: fade-in-other 0.2s ease;
}

@keyframes fade-in-user {
  from {
    opacity: 0;
    transform: translateY(10px) translateX(-50%);
  }
  to {
    opacity: 1;
    transform: translateY(0) translateX(-50%);
  }
}

@keyframes fade-in-other {
  from {
    opacity: 0;
    transform: translateY(10px) translateX(-50%);
  }
  to {
    opacity: 1;
    transform: translateY(0) translateX(-50%);
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

.emoji-item.already-used {
  position: relative;
  background-color: #36434a;
}

.emoji-item.already-used:hover {
  background-color: #445258;
}

.emoji-loading {
  display: flex;
  justify-content: center;
  margin-top: 8px;
}

.debug-info {
  position: absolute;
  top: -20px;
  left: 0;
  right: 0;
  color: #ff6b6b;
  font-size: 0.7rem;
  background-color: rgba(0, 0, 0, 0.6);
  padding: 2px 5px;
  border-radius: 3px;
  z-index: 5;
  max-width: 100%;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>