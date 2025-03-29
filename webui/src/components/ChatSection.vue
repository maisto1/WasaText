<script>
import MessageBubble from './MessageBubble.vue';
import ChatInput from './ChatInput.vue';

export default {
  components: {
    MessageBubble,
    ChatInput,
  },

  props: {
    conversation: {
      type: Object,
      required: true
    },
    isTempChat: {
      type: Boolean,
      default: false
    },
    conversations: {
      type: Array,
      default: () => []
    }
  },

  data() {
    return {
      messages: [],
      loading: false,
      error: null,
      showToast: false,
      toastMessage: '',
      toastType: 'success',
      sendingMessage: false,
      isGroupChat: false,
      pollingInterval: null,
      pollingDelay: 6000,
      lastMessageId: null,
      replyToMessage: null
    }
  },

  computed: {
    forwardingConversations() {
      return this.conversations; 
    }
  },

  created() {
    this.checkIfGroupChat();
    this.startPolling();
  },

  beforeUnmount() {
    this.stopPolling();
  },

  watch: {
    'conversation': {
      immediate: true,
      deep: true,
      handler() {
        this.checkIfGroupChat();
        if (this.conversation?.id && !this.isTempChat) {
          this.fetchMessages();
        } else if (this.isTempChat) {
          this.messages = [];
        }
        

        this.replyToMessage = null;
      }
    },
    'messages': {
      deep: true,
      handler() {
        this.scrollToBottom();
      }
    }
  },

  methods: {
    checkIfGroupChat() {
      this.isGroupChat = this.conversation?.conversationType === 'group';
    },

    async fetchMessages() {
        if (!this.conversation?.id || this.isTempChat) return;
        
        const isInitialLoad = this.messages.length === 0;
        if (isInitialLoad) {
          this.loading = true;
        }
        
        try {
          console.log('Fetching messages for conversation:', this.conversation.id);
          const response = await this.$axios.get(`/conversations/${this.conversation.id}`);
          
       
          this.messages = response.data.map(msg => {

            if (msg.ReplyTo && !msg.ReplyTo.Content) {
              msg.ReplyTo.Content = "Messaggio non disponibile";
            }
            return msg;
          });
          
          this.$emit('refresh-conversations');
          this.scrollToBottom();
        } catch (error) {
          console.error('Error fetching messages:', error);
          this.showNotification('Failed to load messages', 'error');
        } finally {
          if (isInitialLoad) {
            this.loading = false;
          }
        }
      },

    async sendMessage(messageData) {
      if (this.sendingMessage) return;
      
      if (messageData.type === 'text' && !messageData.content.trim()) return;
      
      this.sendingMessage = true;
      
      try {
        if (this.isTempChat) {
        
          const textContent = messageData.content || 'Hello!';
          this.$emit('send-first-message', textContent);
        } else {
          
          
        
          if (this.replyToMessage && this.replyToMessage.id) {
           
            
            const payload = {
              type: messageData.type,
              content: messageData.content || ''
            };
            
            if (messageData.type === 'media' && messageData.media) {
              payload.media = messageData.media;
            }
            
            const response = await this.$axios.post(
              `/conversations/${this.conversation.id}/messages/${this.replyToMessage.id}/reply`,
              payload
            );
            
            this.messages.push(response.data);
            this.replyToMessage = null;
          } else {
          
            const payload = {
              type: messageData.type,
              content: messageData.content || ''
            };
            
            if (messageData.type === 'media' && messageData.media) {
              payload.media = messageData.media;
            }
            
            const response = await this.$axios.post(
              `/conversations/${this.conversation.id}/messages/`, 
              payload
            );
            
            this.messages.push(response.data);
          }
          
          this.scrollToBottom();
        }
      } catch (error) {
        console.error('Error sending message:', error);
        this.showNotification('Failed to send message', 'error');
      } finally {
        this.sendingMessage = false;
      }
    },
    async deleteMessage(messageId) {
      try {
        await this.$axios.delete(`/conversations/${this.conversation.id}/messages/${messageId}`);
        
       
        this.messages = this.messages.filter(msg => msg.id !== messageId);
        
        
        this.messages.forEach(msg => {
          if (msg.ReplyTo && msg.ReplyTo.ID === messageId) {
          
            msg.ReplyTo.Content = null;
          }
        });
        
        this.showNotification('Message deleted successfully', 'success');
      } catch (error) {
        console.error('Error deleting message:', error);
        this.showNotification('Failed to delete message', 'error');
      }
    },
    async forwardMessage(messageId, targetConversationId) {
      try {
        
        this.showNotification('Forwarding message...', 'success');
        
        await this.$axios.post(
          `/conversations/${this.conversation.id}/messages/${messageId}`,
          { conversationId: targetConversationId }
        );
        
        this.showNotification('Message forwarded successfully', 'success');
        this.$emit('refresh-conversations');
      } catch (error) {
        console.error('Error forwarding message:', error);
        this.showNotification('Failed to forward message', 'error');
      }
    },
    
    async handleSendReply(originalMessage, replyData) {
      
      
      if (this.sendingMessage) return;
      this.sendingMessage = true;
      
      try {
        const payload = {
          type: replyData.type,
          content: replyData.content || ''
        };
        
        if (replyData.type === 'media' && replyData.media) {
          payload.media = replyData.media;
        }
        
        const response = await this.$axios.post(
          `/conversations/${this.conversation.id}/messages/${originalMessage.id}/reply`,
          payload
        );
        
        this.messages.push(response.data);
        this.scrollToBottom();
        this.showNotification('Reply sent successfully', 'success');
      } catch (error) {
        console.error('Error sending reply:', error);
        this.showNotification('Failed to send reply', 'error');
      } finally {
        this.sendingMessage = false;
      }
    },
    
    
    async handleReactionAdded(reaction) {
      console.log('Reaction added:', reaction);
      
    },

    async handleReactionUpdated(reaction) {
      console.log('Reaction updated:', reaction);
     
    },

    async handleReactionRemoved(reaction) {
      console.log('Reaction removed:', reaction);
      
    },
    
    cancelReply() {
      this.replyToMessage = null;
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },

    showNotification(message, type) {
      this.toastMessage = message;
      this.toastType = type;
      this.showToast = true;
      setTimeout(() => {
        this.showToast = false;
      }, 3000);
    },

    isMyMessage(message) {
      const myUsername = sessionStorage.getItem('username');
      return message.sender.username === myUsername;
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
    
    startPolling() {
      this.stopPolling();

      this.pollingInterval = setInterval(() => {
        if (this.conversation?.id && !this.isTempChat) {
          this.pollForNewMessages();
        }
      }, this.pollingDelay);
    },
    
    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null;
      }
    },
    
    async pollForNewMessages() {
      if (!this.conversation?.id || this.isTempChat) return;
      
      try {
        const response = await this.$axios.get(`/conversations/${this.conversation.id}`);
        const newMessages = response.data;
        
        if (newMessages.length === 0) return;
        
        const hasChanges = this.checkForMessageChanges(this.messages, newMessages);
        
        if (hasChanges) {
          this.messages = newMessages;
          this.$emit('refresh-conversations');
        }
      } catch (error) {
        console.error('Error polling for new messages:', error);
      }
    },

    checkForMessageChanges(oldMessages, newMessages) {
      if (oldMessages.length !== newMessages.length) {
        return true;
      }
      
      if (oldMessages.length > 0 && newMessages.length > 0) {
        const lastOldMsg = oldMessages[oldMessages.length - 1];
        const lastNewMsg = newMessages[newMessages.length - 1];
        
        if (lastOldMsg.id !== lastNewMsg.id) {
          return true;
        }
        
        for (let i = 0; i < oldMessages.length; i++) {
          if (oldMessages[i].status !== newMessages[i].status) {
            return true;
          }
        }
      }
      
      return false;
    }
  }
}
</script>

<template>
  <div class="chat-section d-flex flex-column h-100">
    <div class="chat-header p-3 d-flex align-items-center justify-content-between border-bottom">
      <div class="d-flex align-items-center">
        <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(conversation.name) }">
          <img 
            v-if="conversation.conversationPhoto" 
            :src="'data:image/jpeg;base64,' + conversation.conversationPhoto"
            class="avatar-image"
            alt="Profile"
          />
          <div v-else class="avatar-text">
            {{ getInitials(conversation.name) }}
          </div>
        </div>
        <div class="ms-3">
          <h6 class="mb-0 text-white">{{ conversation.name }}</h6>
          <small v-if="isTempChat" class="text-light-grey">New chat</small>
        </div>
      </div>

      <button 
        v-if="isGroupChat"
        @click="$emit('open-group-settings')" 
        class="btn-group-settings"
        title="Manage Group"
      >
        <i class="fas fa-cog"></i>
      </button>
    </div>


    <div ref="messagesContainer" class="messages-container flex-grow-1 p-3 overflow-auto">
      <div v-if="loading && !isGroupChat" class="d-flex justify-content-center align-items-center h-100">
        <div class="spinner-border text-light" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>

    <template v-else-if="messages.length === 0">
      <div class="text-center text-light-grey py-5">
        <i class="fas fa-comment-dots fa-3x mb-3"></i>
        <p>Send a message to start the conversation</p>
      </div>
      </template>

      <template v-else>
        <MessageBubble
          v-for="message in messages"
          :key="message.id"
          :message="message"
          :isMyMessage="isMyMessage(message)"
          :showUsername="isGroupChat"
          :availableConversations="forwardingConversations"
          :conversationId="conversation.id"
          @delete="deleteMessage"
          @forward="forwardMessage"
          @send-reply="handleSendReply"
          @reaction-added="handleReactionAdded"
          @reaction-updated="handleReactionUpdated"
          @reaction-removed="handleReactionRemoved"
        />
      </template>
    </div>

    <ChatInput 
      @send="sendMessage" 
      @cancel-reply="cancelReply"
      :disabled="sendingMessage" 
      :reply-to-message="replyToMessage"
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
  </div>
</template>

<style>
.chat-section {
  background-color: #111b21;
}

.chat-header {
  background-color: #202c33;
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


.messages-container {
  background-color: #0b141a;
}

.messages-container::-webkit-scrollbar {
  width: 6px;
}

.messages-container::-webkit-scrollbar-track {
  background: #0b141a;
}

.messages-container::-webkit-scrollbar-thumb {
  background-color: #374045;
  border-radius: 6px;
}

.text-light-grey {
  color: #8696a0;
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

.btn-group-settings {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2a3942;
  color: #e9edef;
  border: none;
  transition: all 0.2s ease;
}

.btn-group-settings:hover {
  background-color: #36434a;
  color: #00a884;
  transform: rotate(30deg);
}
</style>