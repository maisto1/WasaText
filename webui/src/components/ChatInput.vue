// components/ChatInput.vue
<script>
export default {
  props: {
    disabled: {
      type: Boolean,
      default: false
    }
  },

  data() {
    return {
      message: '',
      isSubmitting: false,
      selectedMedia: null,
      mediaPreview: null,
      showMediaPreview: false
    }
  },

  methods: {
    async handleSubmit() {
      if ((!this.message.trim() && !this.selectedMedia) || this.isSubmitting || this.disabled) return;

      this.isSubmitting = true;
      try {
        const messageData = {
          content: this.message,
          type: this.selectedMedia ? 'media' : 'text'
        };

        if (this.selectedMedia) {
          messageData.media = await this.getBase64FromFile(this.selectedMedia);
        }

        await this.$emit('send', messageData);
        this.clearForm();
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

    clearForm() {
      this.message = '';
      this.removeMedia();
    }
  }
}
</script>

<template>
  <div class="chat-input p-3 bg-dark border-top">
    <div v-if="showMediaPreview" class="media-preview mb-2">
      <div class="position-relative d-inline-block">
        <img :src="mediaPreview" class="preview-image rounded" alt="Media preview" />
        <button 
          type="button" 
          class="btn btn-sm btn-danger position-absolute top-0 end-0 rounded-circle" 
          @click="removeMedia"
        >
          <i class="fas fa-times"></i>
        </button>
      </div>
    </div>
    
    <form @submit.prevent="handleSubmit" class="d-flex">
      <input 
        v-model="message"
        type="text"
        class="form-control form-control-lg bg-secondary text-white me-2"
        placeholder="Type a message..."
        :disabled="isSubmitting || disabled"
      />
      
      <button 
        type="button" 
        class="btn btn-secondary btn-lg me-2"
        @click="$refs.fileInput.click()"
        :disabled="isSubmitting || disabled"
      >
        <i class="fas fa-paperclip"></i>
      </button>
      
      <button 
        type="submit" 
        class="btn btn-primary btn-lg"
        :disabled="(!message.trim() && !selectedMedia) || isSubmitting || disabled"
      >
        <i class="fas" :class="isSubmitting ? 'fa-spinner fa-spin' : 'fa-paper-plane'"></i>
      </button>
      
      <input 
        ref="fileInput" 
        type="file" 
        accept="image/*" 
        class="d-none" 
        @change="handleFileSelect" 
      />
    </form>
  </div>
</template>

<style>
.chat-input {
  background-color: #202c33;
}

.form-control {
  background-color: #2a3942 !important;
  border: none;
}

.form-control:focus {
  box-shadow: none;
  border-color: #00a884;
}

.btn-primary {
  background-color: #00a884;
  border-color: #00a884;
}

.btn-primary:hover:not(:disabled) {
  background-color: #008f6f;
  border-color: #008f6f;
}

.btn-primary:disabled {
  background-color: #2a3942;
  border-color: #2a3942;
}

.btn-outline-secondary {
  color: #8696a0;
  border-color: #2a3942;
  background-color: #2a3942;
}

.btn-outline-secondary:hover:not(:disabled) {
  background-color: #36434a;
  color: #e9edef;
}

.media-preview {
  display: flex;
  justify-content: center;
}

.preview-image {
  max-height: 150px;
  max-width: 100%;
  object-fit: contain;
}
</style>