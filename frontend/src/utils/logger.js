/**
 * Frontend logging utility
 * Provides structured logging for frontend errors and events
 */

class Logger {
    constructor() {
        this.isProduction = import.meta.env.PROD || false;
    }

    _formatMessage(level, message, context = {}) {
        const timestamp = new Date().toISOString();
        return {
            timestamp,
            level,
            message,
            context,
            userAgent: navigator.userAgent,
            url: window.location.href
        };
    }

    _log(level, message, context = {}) {
        const logEntry = this._formatMessage(level, message, context);

        // Always log to console in development
        if (!this.isProduction) {
            const consoleMethod = level === 'error' ? console.error :
                level === 'warn' ? console.warn :
                    console.log;
            consoleMethod(`[${logEntry.timestamp}] ${level.toUpperCase()}: ${message}`, context);
        }

        // In production, only log warnings and errors
        if (this.isProduction && (level === 'error' || level === 'warn')) {
            console[level](`[${logEntry.timestamp}] ${level.toUpperCase()}: ${message}`, context);
        }

        // Future: Send to backend logging endpoint
        // this._sendToBackend(logEntry);
    }

    info(message, context) {
        this._log('info', message, context);
    }

    warn(message, context) {
        this._log('warn', message, context);
    }

    error(message, context) {
        this._log('error', message, context);
    }

    // Helper to log caught errors
    logError(error, context = {}) {
        const errorContext = {
            ...context,
            errorMessage: error.message,
            errorStack: error.stack,
            errorName: error.name
        };
        this.error(error.message || 'Unknown error', errorContext);
    }
}

// Create singleton instance
const logger = new Logger();

// Set up global error handler
window.addEventListener('error', (event) => {
    logger.error('Uncaught error', {
        message: event.message,
        filename: event.filename,
        lineno: event.lineno,
        colno: event.colno,
        error: event.error?.stack
    });
});

// Set up unhandled promise rejection handler
window.addEventListener('unhandledrejection', (event) => {
    logger.error('Unhandled promise rejection', {
        reason: event.reason,
        promise: String(event.promise)
    });
});

export default logger;
