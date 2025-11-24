import React from 'react';

class ErrorBoundary extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            hasError: false,
            error: null,
            errorInfo: null
        };
    }

    static getDerivedStateFromError(error) {
        // Update state so the next render will show the fallback UI
        return { hasError: true };
    }

    componentDidCatch(error, errorInfo) {
        // Log error details
        console.error('ErrorBoundary caught an error:', error, errorInfo);

        this.setState({
            error: error,
            errorInfo: errorInfo
        });

        // Send error to backend for logging (future enhancement)
        // This could be implemented to call a backend logging endpoint
    }

    handleReload = () => {
        // Reload the application
        window.location.reload();
    }

    render() {
        if (this.state.hasError) {
            return (
                <div style={{
                    padding: '2rem',
                    maxWidth: '800px',
                    margin: '2rem auto',
                    backgroundColor: '#fee',
                    border: '2px solid #c00',
                    borderRadius: '8px',
                    fontFamily: 'system-ui, -apple-system, sans-serif'
                }}>
                    <h1 style={{ color: '#c00', marginTop: 0 }}>
                        ⚠️ Something Went Wrong
                    </h1>
                    <p style={{ fontSize: '1.1rem', marginBottom: '1.5rem' }}>
                        The application encountered an unexpected error. This has been logged for debugging.
                    </p>

                    <details style={{
                        marginBottom: '1.5rem',
                        padding: '1rem',
                        backgroundColor: '#fff',
                        borderRadius: '4px',
                        cursor: 'pointer'
                    }}>
                        <summary style={{ fontWeight: 'bold', marginBottom: '0.5rem' }}>
                            Error Details
                        </summary>
                        <pre style={{
                            overflow: 'auto',
                            fontSize: '0.85rem',
                            padding: '1rem',
                            backgroundColor: '#f5f5f5',
                            borderRadius: '4px',
                            marginTop: '0.5rem'
                        }}>
                            {this.state.error && this.state.error.toString()}
                            {'\n\n'}
                            {this.state.errorInfo && this.state.errorInfo.componentStack}
                        </pre>
                    </details>

                    <div style={{ display: 'flex', gap: '1rem' }}>
                        <button
                            onClick={this.handleReload}
                            style={{
                                padding: '0.75rem 1.5rem',
                                fontSize: '1rem',
                                backgroundColor: '#0066cc',
                                color: 'white',
                                border: 'none',
                                borderRadius: '4px',
                                cursor: 'pointer',
                                fontWeight: 'bold'
                            }}
                            onMouseOver={(e) => e.target.style.backgroundColor = '#0052a3'}
                            onMouseOut={(e) => e.target.style.backgroundColor = '#0066cc'}
                        >
                            Reload Application
                        </button>
                    </div>

                    <p style={{
                        marginTop: '1.5rem',
                        fontSize: '0.9rem',
                        color: '#666'
                    }}>
                        If this problem persists, check the application logs in <code>~/.mlr-desktop/logs/</code>
                    </p>
                </div>
            );
        }

        return this.props.children;
    }
}

export default ErrorBoundary;
