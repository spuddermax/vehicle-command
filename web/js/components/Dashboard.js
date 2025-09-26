const { useState, useEffect, useCallback } = React;

// Speedometer Component
function Speedometer({ speed, maxSpeed = 120 }) {
    const percentage = (speed / maxSpeed) * 100;
    const rotation = (percentage / 100) * 180 - 90; // -90 to 90 degrees
    
    return (
        <Card title="SPEED" className="speedometer">
            <div className="speedometer-container">
                <div className="speedometer-gauge">
                    <div className="speedometer-needle" style={{ transform: `rotate(${rotation}deg)` }}></div>
                    <div className="speedometer-center"></div>
                </div>
                <div className="speed-display">
                    <span className="speed-value">{speed}</span>
                    <span className="speed-unit">MPH</span>
                </div>
            </div>
        </Card>
    );
}

// Battery Gauge Component
function BatteryGauge({ batteryLevel, range = 250 }) {
    const percentage = Math.max(0, Math.min(100, batteryLevel));
    const batteryColor = percentage > 50 ? '#00ff88' : percentage > 20 ? '#ffaa00' : '#ff4444';
    
    return (
        <Card title="BATTERY" className="battery-gauge">
            <div className="battery-container">
                <div className="battery-level">
                    <div 
                        className="battery-fill" 
                        style={{ 
                            width: `${percentage}%`,
                            backgroundColor: batteryColor
                        }}
                    ></div>
                </div>
                <div className="battery-info">
                    <span className="battery-percentage">{Math.round(percentage)}%</span>
                    <span className="battery-range">{range} mi</span>
                </div>
            </div>
        </Card>
    );
}

// Power Meter Component
function PowerMeter({ power, maxPower = 100 }) {
    const percentage = Math.abs(power / maxPower) * 100;
    const isRegen = power < 0;
    
    return (
        <Card title="POWER" className="power-meter">
            <div className="power-container">
                <div className="power-bar">
                    <div 
                        className={`power-fill ${isRegen ? 'regen' : 'accel'}`}
                        style={{ 
                            width: `${percentage}%`,
                            right: isRegen ? '0' : 'auto',
                            left: isRegen ? 'auto' : '0'
                        }}
                    ></div>
                </div>
                <div className="power-info">
                    <span className="power-value">{power > 0 ? '+' : ''}{power}%</span>
                    <span className="power-label">{isRegen ? 'REGEN' : 'ACCEL'}</span>
                </div>
            </div>
        </Card>
    );
}

// Temperature Gauge Component
function TemperatureGauge({ temperature, label = "MOTOR" }) {
    const tempColor = temperature > 80 ? '#ff4444' : temperature > 60 ? '#ffaa00' : '#00ff88';
    
    return (
        <Card title={label} className="temperature-gauge">
            <div className="temp-container">
                <div className="temp-display">
                    <span className="temp-value" style={{ color: tempColor }}>{temperature}°</span>
                    <span className="temp-unit">F</span>
                </div>
                <div className="temp-bar">
                    <div 
                        className="temp-fill" 
                        style={{ 
                            width: `${Math.min(100, (temperature / 100) * 100)}%`,
                            backgroundColor: tempColor
                        }}
                    ></div>
                </div>
            </div>
        </Card>
    );
}

// Main Dashboard Component
function Dashboard({ dashboardState }) {
    return (
        <div className="controls-grid">
            <Speedometer speed={dashboardState.speed} />
            <BatteryGauge batteryLevel={dashboardState.batteryLevel} range={dashboardState.range} />
            <PowerMeter power={dashboardState.power} />
            <TemperatureGauge temperature={dashboardState.motorTemp} label="MOTOR" />
            <TemperatureGauge temperature={dashboardState.batteryTemp} label="BATTERY" />
        </div>
    );
}

// Export for use in other files
window.Dashboard = Dashboard;
