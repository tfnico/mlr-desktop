
import React, { useState } from 'react';
import { SelectInputFile } from '../../wailsjs/go/main/App';


export default function InputSection({ onInputChange, onModeChange, mode, inputValue, options, inputFormat, ragged, headerless, fieldSeparator }) {
    // We use props for state now, but we can keep local state for immediate feedback if needed.
    // However, for controlled components, we should rely on props.

    const handleTextChange = (e) => {
        onInputChange(e.target.value, 'text', null, null, null, null, null);
    };

    const handleFileChange = (e) => {
        onInputChange(e.target.value, 'file', null, null, null, null, null);
    };

    const handleBrowseFile = async () => {
        try {
            const path = await SelectInputFile();
            if (path) {
                onInputChange(path, 'file', null, null, null, null, null);
            }
        } catch (err) {
            console.error('Error selecting file:', err);
        }
    };

    const showCsvOptions = inputFormat === '--icsv' || inputFormat === '--itsv';
    const showFieldSeparator = inputFormat === '--icsv';

    return (
        <div className="input-section" style={{ padding: '1rem', border: '1px solid #ccc', marginBottom: '1rem' }}>
            <h3>Input</h3>
            <div style={{ marginBottom: '0.5rem', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <div>
                    <button
                        onClick={() => { onModeChange('text'); }}
                        style={{ marginRight: '0.5rem', fontWeight: mode === 'text' ? 'bold' : 'normal' }}
                    >
                        Text
                    </button>
                    <button
                        onClick={() => { onModeChange('file'); }}
                        style={{ fontWeight: mode === 'file' ? 'bold' : 'normal' }}
                    >
                        File Path
                    </button>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    {showCsvOptions && (
                        <div style={{ display: 'flex', gap: '0.5rem', fontSize: '0.8rem' }}>
                            <label style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                <input
                                    type="checkbox"
                                    checked={ragged}
                                    onChange={(e) => onInputChange(null, null, null, null, e.target.checked, null, null)}
                                />
                                Ragged
                            </label>
                            <label style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                <input
                                    type="checkbox"
                                    checked={headerless}
                                    onChange={(e) => onInputChange(null, null, null, null, null, e.target.checked, null)}
                                />
                                Headerless
                            </label>
                        </div>
                    )}
                    {showFieldSeparator && (
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.25rem', fontSize: '0.8rem' }}>
                            <label>Separator:</label>
                            <input
                                type="text"
                                value={fieldSeparator}
                                onChange={(e) => onInputChange(null, null, null, null, null, null, e.target.value)}
                                style={{ width: '30px', textAlign: 'center', padding: '0.1rem' }}
                            />
                        </div>
                    )}
                    <div>
                        <label style={{ marginRight: '0.5rem', fontSize: '0.9rem' }}>Format:</label>
                        <select
                            value={inputFormat}
                            onChange={(e) => onInputChange(null, null, null, e.target.value, null, null, null)}
                            style={{ padding: '0.25rem' }}
                        >
                            <option value="">Default (Auto)</option>
                            <option value="--icsv">CSV (--icsv)</option>
                            <option value="--itsv">TSV (--itsv)</option>
                            <option value="--ijson">JSON (--ijson)</option>
                            <option value="--ijsonl">NDJSON (--ijsonl)</option>
                        </select>
                    </div>
                </div>
            </div>
            {mode === 'text' ? (
                <textarea
                    value={inputValue}
                    onChange={handleTextChange}
                    placeholder="Paste CSV/JSON data here..."
                    style={{ width: '100%', height: '150px', fontFamily: 'monospace', resize: 'vertical', whiteSpace: 'pre', overflowX: 'auto' }}
                />
            ) : (
                <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
                    <input
                        type="text"
                        value={inputValue}
                        onChange={handleFileChange}
                        placeholder="/absolute/path/to/file.csv"
                        style={{ flex: 1, padding: '0.5rem' }}
                    />
                    <button
                        onClick={handleBrowseFile}
                        style={{ padding: '0.5rem 1rem', cursor: 'pointer' }}
                    >
                        Browse...
                    </button>
                </div>
            )}
            <div style={{ marginTop: '0.5rem' }}>
                <label style={{ display: 'block', fontSize: '0.8rem', marginBottom: '0.25rem', color: '#666' }}>Additional flags</label>
                <input
                    type="text"
                    value={options}
                    placeholder="--skip-comments"
                    onChange={(e) => onInputChange(null, null, e.target.value, null, null, null, null)}
                    style={{ width: '100%', padding: '0.5rem' }}
                />
            </div>
        </div>
    );
}

