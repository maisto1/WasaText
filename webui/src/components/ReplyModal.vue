<script>
export default {
  props: {
    show: {
      type: Boolean,
      required: true
    },
    message: {
      type: Object,
      default: null
    }
  },
  
  data() {
    return {
      replyText: '',
      isSubmitting: false,
      selectedMedia: null,
      mediaPreview: null,
      showMediaPreview: false
    }
  },
  
  methods: {
    closeModal() {
      this.resetForm();
      this.$emit('close');
    },
    
    async handleSubmit() {
      if ((!this.replyText.trim() && !this.selectedMedia) || this.isSubmitting) return;
      
      this.isSubmitting = true;
      
      try {
        const messageData = {
          content: this.replyText,
          type: this.selectedMedia ? 'media' : 'text'
        };
        
        if (this.selectedMedia) {
          messageData.media = await this.getBase64FromFile(this.selectedMedia);
        }
        
        this.$emit('send-reply', this.message, messageData);
        this.resetForm();
        this.closeModal();
      } finally {
        this.isSubmitting = false;
      }
    },
    
    async getBase64FromFile(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => {
          const base64String = reader.result.split(',')[1];
          resolve(base64String);
        };
        reader.onerror = error => reject(error);
      });
    },
    
    handleFileSelect(event) {
      const file = event.target.files[0];
      if (!file) return;
      
      if (!file.type.match('image.*')) {
        alert('Please select an image file');
        return;
      }
      
      if (file.size > 5 * 1024 * 1024) {
        alert('File size exceeds 5MB limit');
        return;
      }
      
      this.selectedMedia = file;
      this.createMediaPreview(file);
    },
    
    createMediaPreview(file) {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        this.mediaPreview = reader.result;
        this.showMediaPreview = true;
      };
    },
    
    removeMedia() {
      this.selectedMedia = null;
      this.mediaPreview = null;
      this.showMediaPreview = false;
      
      const fileInput = this.$refs.fileInput;
      if (fileInput) {
        fileInput.value = '';
      }
    },
    
    resetForm() {
      this.replyText = '';
      this.removeMedia();
    },
    
    formatMessagePreview(message) {
      if (!message) return '';
      
      if (message.type === 'media') {
        if (message.content) {
          return `[Photo] ${message.content}`;
        } else {
          return '[Photo]';
        }
      } else {
        return message.content;
      }
    }
  },
  
  mounted() {
    // Focus the input field when modal opens
    this.$nextTick(() => {
      if (this.$refs.replyInput) {
        this.$refs.replyInput.focus();
      }
    });
  }
}
</script>

<template>
  <div v-if="show" class="modal-overlay" @click="closeModal">
    <div class="modal-dialog" @click.stop>
      <div class="modal-content">
        <div class="modal-header">
          <div class="d-flex align-items-center">
            <i class="fas fa-reply modal-icon"></i>
            <h5 class="modal-title ms-2">Reply to message</h5>
          </div>
          <button @click="closeModal" class="close-button">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <!-- Original message preview -->
          <div v-if="message" class="message-preview mb-3">
            <div class="sender-info mb-2">
              <strong>{{ message.sender.username }}</strong>
            </div>
            
            <div class="message-content">
              <template v-if="message.type === 'media' && message.media">
                <div class="media-preview mb-2">
                  <img 
                    :src="'data:image/jpeg;base64,' + message.media"
                    class="preview-image rounded"
                    alt="Media preview"
                  />
                </div>
              </template>
              
              <p class="message-text">{{ message.content }}</p>
            </div>
          </div>
          
          <!-- Media preview -->
          <div v-if="showMediaPreview" class="reply-media-preview mb-3">
            <div class="position-relative d-inline-block">
              <img :src="mediaPreview" class="reply-preview-image rounded" alt="Media preview" />
              <button 
                type="button" 
                class="btn btn-sm btn-danger position-absolute top-0 end-0 rounded-circle remove-media-btn" 
                @click="removeMedia"
              >
                <i class="fas fa-times"></i>
              </button>
            </div>
          </div>
          
          <!-- Reply form -->
          <form @submit.prevent="handleSubmit" class="reply-form">
            <div class="d-flex mb-2">
              <input 
                ref="replyInput"
                v-model="replyText"
                type="text"
                class="form-control form-control-lg bg-secondary text-white me-2"
                placeholder="Type your reply..."
                :disabled="isSubmitting"
              />
              
              <button 
                type="button" 
                class="btn btn-secondary btn-lg me-2"
                @click="$refs.fileInput.click()"
                :disabled="isSubmitting"
              >
                <i class="fas fa-paperclip"></i>
              </button>
            </div>
            
            <input 
              ref="fileInput" 
              type="file" 
              accept="image/*" 
              class="d-none" 
              @change="handleFileSelect" 
            />
          </form>
        </div>
        
        <div class="modal-footer">
          <button @click="closeModal" class="cancel-button">
            Cancel
          </button>
          <button 
            @click="handleSubmit" 
            class="reply-button"
            :disabled="(!replyText.trim() && !selectedMedia) || isSubmitting"
          >
            <i v-if="!isSubmitting" class="fas fa-paper-plane me-2"></i>
            <span v-else class="spinner-border spinner-border-sm me-2"></span>
            Send Reply
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

.modal-body {
  padding: 20px;
  max-height: 350px;
  overflow-y: auto;
}

.message-preview {
  padding: 15px;
  background-color: rgba(42, 57, 66, 0.6);
  border-radius: 12px;
  border-left: 4px solid #00a884;
}

.sender-info {
  color: #00a884;
}

.message-text {
  color: #e9edef;
  margin-bottom: 0;
  word-break: break-word;
}

.media-preview {
  display: flex;
  justify-content: center;
}

.preview-image {
  max-height: 120px;
  max-width: 100%;
  object-fit: contain;
}

.reply-media-preview {
  display: flex;
  justify-content: center;
}

.reply-preview-image {
  max-height: 150px;
  max-width: 100%;
  object-fit: contain;
}

.remove-media-btn {
  width: 24px;
  height: 24px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 4px;
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

.reply-button {
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

.reply-button:hover:not(:disabled) {
  background-color: #008f6f;
}

.reply-button:disabled {
  background-color: rgba(0, 168, 132, 0.5);
  cursor: not-allowed;
}

.form-control {
  background-color: #2a3942 !important;
  border: none;
}

.form-control:focus {
  box-shadow: none;
  border-color: #00a884;
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