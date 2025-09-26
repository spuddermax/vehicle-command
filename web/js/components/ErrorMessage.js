const { useState, useEffect, useCallback } = React;

// Error Message Component
function ErrorMessage({ message, onDismiss }) {
    return (
        <div className="error-message">
            <div className="error-content">
                <span className="error-icon">⚠</span>
                <span className="error-text">{message}</span>
                <button className="error-dismiss" onClick={onDismiss}>×</button>
            </div>
        </div>
    );
}

// Export for use in other files
window.ErrorMessage = ErrorMessage;
