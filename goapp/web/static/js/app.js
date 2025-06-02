// Initialize Alpine.js for dropdown interactions
document.addEventListener('alpine:init', () => {
    Alpine.data('dropdown', () => ({
        open: false,
        toggle() {
            this.open = !this.open;
        },
        close() {
            this.open = false;
        }
    }));
});

// HTMX event handlers
document.body.addEventListener('htmx:configRequest', (event) => {
    // Add CSRF token if needed
    // event.detail.headers['X-CSRF-Token'] = getCsrfToken();
});

document.body.addEventListener('htmx:afterSwap', (event) => {
    // Re-initialize any JavaScript components after HTMX swap
    if (typeof Alpine !== 'undefined') {
        Alpine.initTree(event.detail.target);
    }
});

// Handle HTMX errors
document.body.addEventListener('htmx:responseError', (event) => {
    console.error('HTMX request failed:', event.detail);
    // Show error notification
    showNotification('error', 'Request failed. Please try again.');
});

// Notification system
function showNotification(type, message) {
    const notification = document.createElement('div');
    notification.className = `fixed top-4 right-4 px-6 py-3 rounded-lg shadow-lg text-white z-50 transition-all transform translate-x-0 ${
        type === 'success' ? 'bg-green-500' :
        type === 'error' ? 'bg-red-500' :
        type === 'warning' ? 'bg-yellow-500' :
        'bg-blue-500'
    }`;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    // Auto-remove after 5 seconds
    setTimeout(() => {
        notification.classList.add('translate-x-full', 'opacity-0');
        setTimeout(() => notification.remove(), 300);
    }, 5000);
}

// Sidebar toggle for mobile
function toggleSidebar() {
    const sidebar = document.getElementById('sidebar');
    sidebar.classList.toggle('hidden');
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    // Add Alpine.js for interactivity
    const alpineScript = document.createElement('script');
    alpineScript.src = 'https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js';
    alpineScript.defer = true;
    document.head.appendChild(alpineScript);
});