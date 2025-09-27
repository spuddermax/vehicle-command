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
    const [isExpanded, setIsExpanded] = useState(false);

    const handleTitleClick = () => {
        setIsExpanded(!isExpanded);
    };

    const cardStyle = {
        borderColor: borderColor,
        backgroundColor: backgroundColor
    };

    return (
        <div 
            className={`card ${className} ${isExpanded ? 'expanded' : ''}`}
            style={cardStyle}
            {...props}
        >
            {title && (
                <h2 className="card-title" onClick={handleTitleClick}>{title}</h2>
            )}
            <div className="card-content">
                {children}
            </div>
        </div>
    );
}

// Export for use in other files
window.Card = Card;
