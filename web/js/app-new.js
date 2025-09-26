// Main application entry point
// This file loads all components and renders the app

// Load all component files
// Note: In a real React app, these would be proper ES6 imports
// For this vanilla React setup, we're using script tags in HTML

// Render the app once all components are loaded
document.addEventListener('DOMContentLoaded', function() {
    // Wait for all components to be available
    if (window.TeslaHVACApp) {
        ReactDOM.render(<TeslaHVACApp />, document.getElementById('root'));
    } else {
        // Fallback: retry after a short delay
        setTimeout(() => {
            if (window.TeslaHVACApp) {
                ReactDOM.render(<TeslaHVACApp />, document.getElementById('root'));
            } else {
                console.error('TeslaHVACApp component not found');
            }
        }, 100);
    }
});
