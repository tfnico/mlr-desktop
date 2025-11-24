import React from 'react';

export default function OutputPreview({ output, error, outputFormat, onOutputFormatChange, command, onSave }) {
    return (
        <div className="output-preview" style={{ padding: '1rem', border: '1px solid #ccc' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
                <h3>Output Preview</h3>
                <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
                    <div>
                        <label style={{ marginRight: '0.5rem', fontSize: '0.9rem' }}>Format:</label>
                        <select
                            value={outputFormat}
                            onChange={(e) => onOutputFormatChange(e.target.value)}
                            style={{ padding: '0.25rem' }}
                        >
                            <option value="">Default (Auto)</option>
                            <option value="--opprint">Pretty Print (--opprint)</option>
                            <option value="--ocsv">CSV (--ocsv)</option>
                            <option value="--otsv">TSV (--otsv)</option>
                            <option value="--ojson">JSON (--ojson)</option>
                            <option value="--ojsonl">NDJSON (--ojsonl)</option>
                        </select>
                    </div>
                    <button onClick={onSave} disabled={!output} style={{ padding: '0.25rem 0.5rem', cursor: 'pointer' }}>
                        Save to File
                    </button>
                </div>
            </div>
            {error ? (
                <div style={{ color: 'red', whiteSpace: 'pre-wrap' }}>
                    <strong>Error:</strong>
                    <pre>{error}</pre>
                </div>
            ) : (
                <textarea
                    readOnly
                    value={output}
                    style={{
                        width: '100%',
                        height: '300px',
                        fontFamily: 'monospace',
                        resize: 'vertical',
                        whiteSpace: 'pre',
                        overflowX: 'auto',
                        backgroundColor: '#f5f5f5',
                        border: '1px solid #ddd'
                    }}
                />
            )}
            {command && (
                <div style={{ marginTop: '1rem' }}>
                    <h4 style={{ margin: '0 0 0.5rem 0', textAlign: 'left' }}>Generated Command</h4>
                    <textarea
                        readOnly
                        value={command}
                        onClick={(e) => e.target.select()}
                        style={{
                            width: '100%',
                            height: '80px',
                            fontFamily: 'monospace',
                            padding: '0.5rem',
                            border: '1px solid #ddd',
                            borderRadius: '4px',
                            backgroundColor: '#f5f5f5',
                            color: '#333',
                            resize: 'vertical'
                        }}
                    />
                </div>
            )}
        </div>
    );
}
