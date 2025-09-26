const { useState, useEffect, useCallback } = React;

// Temperature Control Component
function TemperatureControl({ driverTemp, passengerTemp, onTemperatureChange, disabled }) {
    const handleTempChange = (side, delta) => {
        const newTemp = side === 'driver' 
            ? Math.max(59, Math.min(86, driverTemp + delta)) // 59°F to 86°F range
            : Math.max(59, Math.min(86, passengerTemp + delta)); // 59°F to 86°F range
        
        if (side === 'driver') {
            onTemperatureChange(newTemp, passengerTemp);
        } else {
            onTemperatureChange(driverTemp, newTemp);
        }
    };

    return (
        <Card title="TEMPERATURE" className="temperature-control">
            <div className="temperature-display">
                <div className="temp-side">
                    <div className="temp-label">DRIVER</div>
                    <div className="temp-value">{Math.round(driverTemp)}°F</div>
                    <div className="temp-slider-container">
                        <input
                            type="range"
                            min="59"
                            max="86"
                            value={driverTemp}
                            onChange={(e) => handleTempChange('driver', parseFloat(e.target.value) - driverTemp)}
                            disabled={disabled}
                            className="temp-slider"
                        />
                    </div>
                </div>
                
                <div className="temp-divider"></div>
                
                <div className="temp-side">
                    <div className="temp-label">PASSENGER</div>
                    <div className="temp-value">{Math.round(passengerTemp)}°F</div>
                    <div className="temp-slider-container">
                        <input
                            type="range"
                            min="59"
                            max="86"
                            value={passengerTemp}
                            onChange={(e) => handleTempChange('passenger', parseFloat(e.target.value) - passengerTemp)}
                            disabled={disabled}
                            className="temp-slider"
                        />
                    </div>
                </div>
            </div>
        </Card>
    );
}

// Export for use in other files
window.TemperatureControl = TemperatureControl;
