const { useState, useEffect, useCallback } = React;

// Standardized Card Component
function Card({ 
    title, 
    children, 
    borderColor = 'var(--primary-blue)', 
    backgroundColor = 'var(--bg-panel)',
    className = '',
    ...props 
}) {
    const cardStyle = {
        borderColor: borderColor,
        backgroundColor: backgroundColor
    };

    return (
        <div 
            className={`card ${className}`}
            style={cardStyle}
            {...props}
        >
            {title && (
                <h2 className="card-title">{title}</h2>
            )}
            <div className="card-content">
                {children}
            </div>
        </div>
    );
}

// Export for use in other files
window.Card = Card;
