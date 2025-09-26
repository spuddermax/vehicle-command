const { useState, useEffect, useCallback } = React;

// Audio Control Component
function AudioControl({ volume, onVolumeChange, disabled }) {
    const handleVolumeChange = (delta) => {
        const newVolume = Math.max(0, Math.min(100, volume + delta));
        onVolumeChange(newVolume);
    };

    return (
        <div className="control-panel audio-control">
            <h1 className="control-title">VOLUME</h1>
            <div className="volume-display">
                <div className="volume-value">{volume}%</div>
                <div className="volume-buttons">
                    <button 
                        className="volume-button volume-down"
                        onClick={() => handleVolumeChange(-5)}
                        disabled={disabled}
                    >
                        −
                    </button>
                    <button 
                        className="volume-button volume-up"
                        onClick={() => handleVolumeChange(5)}
                        disabled={disabled}
                    >
                        +
                    </button>
                </div>
            </div>
        </div>
    );
}

// Export for use in other files
window.AudioControl = AudioControl;
