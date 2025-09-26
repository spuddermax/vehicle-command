const { useState, useEffect, useCallback } = React;

// Status Display Component
function StatusDisplay({ hvacState, connectionStatus }) {
    return (
        <div className="status-display">
            <div className="status-grid">
                <div className="status-item">
                    <div className="status-label">INSIDE</div>
                    <div className="status-value">{Math.round(hvacState.insideTemp)}°F</div>
                </div>
                <div className="status-item">
                    <div className="status-label">OUTSIDE</div>
                    <div className="status-value">{Math.round(hvacState.outsideTemp)}°F</div>
                </div>
                <div className="status-item">
                    <div className="status-label">FAN</div>
                    <div className="status-value">
                        {hvacState.fanSpeed === -1 ? 'AUTO' : hvacState.fanSpeed}
                    </div>
                </div>
                <div className="status-item">
                    <div className="status-label">MODE</div>
                    <div className="status-value">
                        {hvacState.autoMode ? 'AUTO' : 'MANUAL'}
                    </div>
                </div>
            </div>
        </div>
    );
}

// Export for use in other files
window.StatusDisplay = StatusDisplay;
