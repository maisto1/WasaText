<script>
export default {
  data() {
    return {
      username: sessionStorage.getItem('username') || '',
      newUsername: '',
      loading: false,
      showToast: false,
      toastMessage: '',
      toastType: 'success',
      profilePhoto: null,
      newProfilePhoto: null,
      photoPreview: null,
      updatingUsername: false,
      updatingPhoto: false
    }
  },
  
  created() {
    this.loadUserData();
  },
  
  methods: {
    loadUserData() {
      this.loading = true;
      try {
        const username = sessionStorage.getItem('username');
        if (username) {
          this.username = username;
          this.newUsername = username;
        }

        const profilePhoto = sessionStorage.getItem('profilePhoto');
        if (profilePhoto) {
          this.profilePhoto = profilePhoto;
        }
        
      } catch (error) {
        console.error('Error loading user data:', error);
        this.showNotification('Failed to load user data', 'error');
      } finally {
        this.loading = false;
      }
    },
    
    async updateUsername() {
      if (!this.newUsername.trim() || this.newUsername === this.username || this.updatingUsername) {
        return;
      }
      
      if (this.newUsername.length < 3 || this.newUsername.length > 16) {
        this.showNotification('Username must be between 3 and 16 characters', 'error');
        return;
      }
      
      this.updatingUsername = true;
      
      try {
        await this.$axios.put('/users/profile/username', {
          username: this.newUsername
        });
        
        sessionStorage.setItem('username', this.newUsername);
        this.username = this.newUsername;
        
        this.showNotification('Username updated successfully', 'success');
      } catch (error) {
        console.error('Error updating username:', error);
        
        if (error.response) {
          switch (error.response.status) {
            case 409:
              this.showNotification('Username already exists, please try another one', 'error');
              break;
            case 400:
              this.showNotification('Invalid username format', 'error');
              break;
            case 404:
              this.showNotification('User not found', 'error');
              break;
            default:
              this.showNotification('Failed to update username', 'error');
          }
        } else {
          this.showNotification('Failed to update username', 'error');
        }

        this.newUsername = this.username;
      } finally {
        this.updatingUsername = false;
      }
    },
    
    handleFileSelect(event) {
      const file = event.target.files[0];
      if (!file) return;
      
      if (!file.type.match('image.*')) {
        this.showNotification('Please select an image file', 'error');
        return;
      }
      
      if (file.size > 5 * 1024 * 1024) {
        this.showNotification('File size exceeds 5MB limit', 'error');
        return;
      }
      
      this.newProfilePhoto = file;
      this.createPhotoPreview(file);
    },
    
    createPhotoPreview(file) {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        this.photoPreview = reader.result;
      };
    },
    
    async updateProfilePhoto() {
      if (!this.newProfilePhoto || this.updatingPhoto) {
        return;
      }
      
      this.updatingPhoto = true;
      
      try {
        const base64String = await this.getBase64FromFile(this.newProfilePhoto);
        
        await this.$axios.put('/users/profile/photo', {
          photo: base64String
        });
        
        this.profilePhoto = base64String;
        
        sessionStorage.setItem('profilePhoto', base64String);
        
        this.resetPhotoSelection();
        
        this.showNotification('Profile photo updated successfully', 'success');
      } catch (error) {
        console.error('Error updating profile photo:', error);
        if (error.response && error.response.status === 400) {
          this.showNotification('Invalid image format', 'error');
        } else {
          this.showNotification('Failed to update profile photo', 'error');
        }
      } finally {
        this.updatingPhoto = false;
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
    
    resetPhotoSelection() {
      this.newProfilePhoto = null;
      this.photoPreview = null;
      
      const fileInput = this.$refs.fileInput;
      if (fileInput) {
        fileInput.value = '';
      }
    },
    
    goBack() {
      this.$router.push('/chats');
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
  <div class="profile-view">
    <div class="profile-header p-3 d-flex align-items-center border-bottom">
      <button @click="goBack" class="btn btn-link text-light me-2">
        <i class="fas fa-arrow-left"></i>
      </button>
      <h5 class="mb-0 text-white">Profile Settings</h5>
    </div>
    
    <div class="profile-content p-4">
      <div v-if="loading" class="text-center p-4">
        <div class="spinner-border text-light" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      
      <div v-else>
        <div class="profile-photo-section text-center mb-4">
          <div class="profile-photo-container mx-auto mb-3">
            <img 
              v-if="photoPreview || profilePhoto" 
              :src="photoPreview || `data:image/jpeg;base64,${profilePhoto}`"
              class="profile-photo"
              alt="Profile Photo"
            />
            <div v-else class="profile-photo-placeholder">
              <span>{{ username.charAt(0).toUpperCase() }}</span>
            </div>
          </div>
          
          <div class="mb-3">
            <button @click="$refs.fileInput.click()" class="btn btn-outline-light me-2">
              <i class="fas fa-camera me-2"></i>Change Photo
            </button>
            <input 
              ref="fileInput" 
              type="file" 
              accept="image/*" 
              class="d-none" 
              @change="handleFileSelect" 
            />
            
            <button 
              v-if="photoPreview" 
              @click="updateProfilePhoto" 
              class="btn btn-success me-2"
              :disabled="updatingPhoto"
            >
              <span v-if="!updatingPhoto">
                <i class="fas fa-save me-2"></i>Save
              </span>
              <span v-else class="spinner-border spinner-border-sm"></span>
            </button>
            
            <button 
              v-if="photoPreview" 
              @click="resetPhotoSelection" 
              class="btn btn-outline-danger"
              :disabled="updatingPhoto"
            >
              <i class="fas fa-times me-2"></i>Cancel
            </button>
          </div>
        </div>
        
    <div class="username-section mb-4">
    <div class="card bg-dark border-secondary" style="border-radius: 15px; overflow: hidden;">
        <div class="card-header bg-secondary d-flex align-items-center">
        <i class="fas fa-user me-2 text-light"></i>
        <h6 class="mb-0 text-white">Account Information</h6>
        </div>
        <div class="card-body">
        <div class="mb-3">
            <label for="username" class="form-label text-white d-flex align-items-center">
            <i class="fas fa-signature me-2 text-light-grey"></i>
            Username
            </label>
            <div class="input-group input-group-lg">
            <input 
                id="username"
                v-model="newUsername"
                type="text"
                class="form-control bg-secondary text-white border-end-0"
                placeholder="Enter a new username"
                minlength="3"
                maxlength="16"
                :disabled="updatingUsername"
                style="border-radius: 10px 0 0 10px;"
            />
            <button 
                @click="updateUsername" 
                class="btn btn-primary"
                :disabled="!newUsername.trim() || newUsername === username || updatingUsername"
                style="border-radius: 0 10px 10px 0; padding-left: 20px; padding-right: 20px;"
            >
                <span v-if="!updatingUsername">
                <i class="fas fa-save me-2"></i>Update
                </span>
                <span v-else class="spinner-border spinner-border-sm"></span>
            </button>
            </div>
            <div class="d-flex align-items-center mt-2">
            <i class="fas fa-info-circle text-light-grey me-2"></i>
            <small class="text-light-grey">
                Username must be between 3 and 16 characters
            </small>
            </div>
        </div>
        </div>
    </div>
    </div>
      </div>
    </div>
    
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

<style scoped>
.profile-view {
  height: calc(100vh - 48px);
  background-color: #111b21;
  overflow-y: auto;
}

.profile-header {
  background-color: #202c33;
}

.profile-content {
  max-width: 600px;
  margin: 0 auto;
}

.profile-photo-container {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
  background-color: #2a3942;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.profile-photo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-photo-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 4rem;
  font-weight: bold;
  color: white;
  background-color: #00a884;
}

.text-light-grey {
  color: #8696a0;
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

.toast-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 3px;
  width: 100%;
  animation: progress 3s linear forwards;
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