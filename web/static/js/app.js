class FotoBooApp {
    constructor() {
        this.currentScreen = 'welcome';
        this.stream = null;
        this.capturedPhoto = null;
        this.currentFilter = 'none';
        this.currentFrame = 'none';
        this.photoId = null;
        
        this.initializeElements();
        this.bindEvents();
        this.startCamera();
    }

    initializeElements() {
        // Screens
        this.screens = {
            welcome: document.getElementById('welcome-screen'),
            capture: document.getElementById('capture-screen'),
            preview: document.getElementById('preview-screen'),
            result: document.getElementById('result-screen')
        };

        // Camera elements
        this.camera = document.getElementById('camera');
        this.photoCanvas = document.getElementById('photo-canvas');
        this.previewCanvas = document.getElementById('preview-canvas');
        this.frameOverlay = document.getElementById('frame-overlay');
        this.previewFrame = document.getElementById('preview-frame');
        this.countdown = document.getElementById('countdown');
        this.flash = document.getElementById('flash');

        // Buttons
        this.startCaptureBtn = document.getElementById('start-capture');
        this.captureBtn = document.getElementById('capture-btn');
        this.backBtn = document.getElementById('back-btn');
        this.retakeBtn = document.getElementById('retake-btn');
        this.saveBtn = document.getElementById('save-btn');
        this.downloadBtn = document.getElementById('download-btn');
        this.newPhotoBtn = document.getElementById('new-photo-btn');

        // Filter controls
        this.filterButtons = document.querySelectorAll('.filter-btn');
        this.frameButtons = document.querySelectorAll('.frame-btn');
        this.brightnessSlider = document.getElementById('brightness-slider');
        this.contrastSlider = document.getElementById('contrast-slider');
        this.brightnessValue = document.getElementById('brightness-value');
        this.contrastValue = document.getElementById('contrast-value');

        // Result elements
        this.finalPhoto = document.getElementById('final-photo');
        this.qrCode = document.getElementById('qr-code');
    }

    bindEvents() {
        // Navigation
        this.startCaptureBtn.addEventListener('click', () => this.showScreen('capture'));
        this.backBtn.addEventListener('click', () => this.showScreen('welcome'));
        this.retakeBtn.addEventListener('click', () => this.showScreen('capture'));
        this.newPhotoBtn.addEventListener('click', () => this.showScreen('welcome'));

        // Photo capture
        this.captureBtn.addEventListener('click', () => this.capturePhoto());
        this.saveBtn.addEventListener('click', () => this.savePhoto());
        this.downloadBtn.addEventListener('click', () => this.downloadPhoto());

        // Filter controls
        this.filterButtons.forEach(btn => {
            btn.addEventListener('click', (e) => this.selectFilter(e.target.dataset.filter));
        });

        this.frameButtons.forEach(btn => {
            btn.addEventListener('click', (e) => this.selectFrame(e.target.dataset.frame));
        });

        this.brightnessSlider.addEventListener('input', (e) => {
            this.brightnessValue.textContent = e.target.value;
            this.applyFilters();
        });

        this.contrastSlider.addEventListener('input', (e) => {
            this.contrastValue.textContent = e.target.value;
            this.applyFilters();
        });
    }

    showScreen(screenName) {
        // Hide all screens
        Object.values(this.screens).forEach(screen => {
            screen.classList.remove('active');
        });

        // Show selected screen
        this.screens[screenName].classList.add('active');
        this.currentScreen = screenName;

        // Special handling for preview screen
        if (screenName === 'preview' && this.capturedPhoto) {
            this.displayPreview();
        }
    }

    async startCamera() {
        try {
            this.stream = await navigator.mediaDevices.getUserMedia({
                video: {
                    width: { ideal: 1280 },
                    height: { ideal: 720 },
                    facingMode: 'user'
                }
            });
            this.camera.srcObject = this.stream;
        } catch (error) {
            console.error('Error accessing camera:', error);
            alert('Unable to access camera. Please check permissions.');
        }
    }

    async capturePhoto() {
        // Disable capture button during countdown
        this.captureBtn.disabled = true;

        // Start countdown
        await this.startCountdown();

        // Capture photo
        this.flash.classList.add('active');
        
        const canvas = this.photoCanvas;
        const context = canvas.getContext('2d');
        canvas.width = this.camera.videoWidth;
        canvas.height = this.camera.videoHeight;
        context.drawImage(this.camera, 0, 0);

        // Store photo data
        this.capturedPhoto = canvas.toDataURL('image/jpeg', 0.9);

        setTimeout(() => {
            this.flash.classList.remove('active');
            this.captureBtn.disabled = false;
            this.showScreen('preview');
        }, 300);
    }

    startCountdown() {
        return new Promise(resolve => {
            let count = 3;
            this.countdown.textContent = count;
            this.countdown.classList.add('active');

            const countInterval = setInterval(() => {
                count--;
                if (count > 0) {
                    this.countdown.textContent = count;
                    this.countdown.classList.remove('active');
                    void this.countdown.offsetWidth; // Force reflow
                    this.countdown.classList.add('active');
                } else {
                    clearInterval(countInterval);
                    this.countdown.classList.remove('active');
                    resolve();
                }
            }, 1000);
        });
    }

    selectFilter(filter) {
        this.currentFilter = filter;
        
        // Update UI
        this.filterButtons.forEach(btn => {
            btn.classList.toggle('active', btn.dataset.filter === filter);
        });

        this.applyFilters();
    }

    selectFrame(frame) {
        this.currentFrame = frame;
        
        // Update UI
        this.frameButtons.forEach(btn => {
            btn.classList.toggle('active', btn.dataset.frame === frame);
        });

        this.applyFrame();
    }

    applyFilters() {
        if (!this.capturedPhoto) return;

        const canvas = this.previewCanvas;
        const context = canvas.getContext('2d');
        const img = new Image();

        img.onload = () => {
            canvas.width = img.width;
            canvas.height = img.height;

            // Apply CSS filters
            context.filter = this.getFilterString();
            context.drawImage(img, 0, 0);
        };

        img.src = this.capturedPhoto;
    }

    getFilterString() {
        let filterString = '';
        
        const brightness = this.brightnessSlider.value / 100;
        const contrast = this.contrastSlider.value / 100;

        filterString += `brightness(${brightness}) `;
        filterString += `contrast(${contrast}) `;

        switch (this.currentFilter) {
            case 'grayscale':
                filterString += 'grayscale(100%) ';
                break;
            case 'vintage':
                filterString += 'sepia(50%) contrast(1.2) brightness(0.9) ';
                break;
            case 'brightness':
                filterString += 'brightness(1.3) ';
                break;
            case 'contrast':
                filterString += 'contrast(1.5) ';
                break;
        }

        return filterString.trim();
    }

    applyFrame() {
        // Clear existing frames
        this.previewFrame.className = '';
        this.frameOverlay.className = '';

        const frameClass = `frame-${this.currentFrame}`;
        this.previewFrame.classList.add(frameClass);
        this.frameOverlay.classList.add(frameClass);
    }

    displayPreview() {
        this.applyFilters();
        this.applyFrame();
    }

    async savePhoto() {
        if (!this.capturedPhoto) return;

        try {
            // Convert to blob
            const response = await fetch(this.capturedPhoto);
            const blob = await response.blob();
            
            // Upload to backend
            const uploadResponse = await fetch('/photos', {
                method: 'POST',
                body: blob
            });

            if (!uploadResponse.ok) {
                throw new Error('Failed to upload photo');
            }

            const result = await uploadResponse.json();
            this.photoId = result.id;

            // Generate the final photo URL
            const photoUrl = `/photos/${this.photoId}`;
            
            // Show result screen
            this.showResult(photoUrl);

        } catch (error) {
            console.error('Error saving photo:', error);
            alert('Failed to save photo. Please try again.');
        }
    }

    showResult(photoUrl) {
        this.finalPhoto.src = photoUrl;
        
        // Generate QR code (simple implementation - you'd want to use a proper QR code library)
        this.qrCode.innerHTML = `
            <p><strong>Photo URL:</strong></p>
            <p>${window.location.origin}${photoUrl}</p>
            <p>Scan to view photo</p>
        `;

        this.showScreen('result');
    }

    downloadPhoto() {
        if (!this.photoId) return;

        const link = document.createElement('a');
        link.href = `/photos/${this.photoId}`;
        link.download = `fotoboo-${this.photoId}.jpg`;
        link.click();
    }
}

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new FotoBooApp();
});