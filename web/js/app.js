// Main application entry point
// This file loads all components and renders the app

// Render the app once all components are loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, checking for components...');
    console.log('Available components:', Object.keys(window).filter(key => key.includes('Control') || key.includes('App') || key.includes('Header')));
    
    // Wait for all components to be available
    if (window.TeslaHVACApp) {
        console.log('TeslaHVACApp found, rendering...');
        ReactDOM.render(<TeslaHVACApp />, document.getElementById('root'));
    } else {
        console.log('TeslaHVACApp not found, retrying...');
        // Fallback: retry after a short delay
        setTimeout(() => {
            if (window.TeslaHVACApp) {
                console.log('TeslaHVACApp found on retry, rendering...');
                ReactDOM.render(<TeslaHVACApp />, document.getElementById('root'));
            } else {
                console.error('TeslaHVACApp component not found after retry');
                console.log('Available window properties:', Object.keys(window).filter(key => key.includes('Control') || key.includes('App') || key.includes('Header')));
            }
        }, 500);
    }
});
