const { useState, useEffect, useCallback } = React;

// Temperature Control Component
function TemperatureControl({ driverTemp, passengerTemp, onTemperatureChange, disabled }) {
    const handleTempChange = (side, delta) => {
        const newTemp = side === 'driver' 
            ? Math.max(60, Math.min(85, driverTemp + delta)) // 60°F to 85°F range
            : Math.max(60, Math.min(85, passengerTemp + delta)); // 60°F to 85°F range
        
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
                            min="60"
                            max="85"
                            value={driverTemp}
                            onChange={(e) => handleTempChange('driver', parseFloat(e.target.value) - driverTemp)}
                            disabled={disabled}
                            className="temp-slider"
                        />
                        <div className="temp-labels">
                            <span onClick={() => handleTempChange('driver', 60 - driverTemp)}>60</span>
                            <span onClick={() => handleTempChange('driver', 65 - driverTemp)}>65</span>
                            <span onClick={() => handleTempChange('driver', 70 - driverTemp)}>70</span>
                            <span onClick={() => handleTempChange('driver', 75 - driverTemp)}>75</span>
                            <span onClick={() => handleTempChange('driver', 80 - driverTemp)}>80</span>
                            <span onClick={() => handleTempChange('driver', 85 - driverTemp)}>85</span>
                        </div>
                    </div>
                </div>
                
                <div className="temp-divider"></div>
                
                <div className="temp-side">
                    <div className="temp-label">PASSENGER</div>
                    <div className="temp-value">{Math.round(passengerTemp)}°F</div>
                    <div className="temp-slider-container">
                        <input
                            type="range"
                            min="60"
                            max="85"
                            value={passengerTemp}
                            onChange={(e) => handleTempChange('passenger', parseFloat(e.target.value) - passengerTemp)}
                            disabled={disabled}
                            className="temp-slider"
                        />
                        <div className="temp-labels">
                            <span onClick={() => handleTempChange('passenger', 60 - passengerTemp)}>60</span>
                            <span onClick={() => handleTempChange('passenger', 65 - passengerTemp)}>65</span>
                            <span onClick={() => handleTempChange('passenger', 70 - passengerTemp)}>70</span>
                            <span onClick={() => handleTempChange('passenger', 75 - passengerTemp)}>75</span>
                            <span onClick={() => handleTempChange('passenger', 80 - passengerTemp)}>80</span>
                            <span onClick={() => handleTempChange('passenger', 85 - passengerTemp)}>85</span>
                        </div>
                    </div>
                </div>
            </div>
        </Card>
    );
}

// Export for use in other files
window.TemperatureControl = TemperatureControl;
