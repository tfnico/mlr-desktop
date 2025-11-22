import React, { useState } from 'react';

export default function VerbBuilder({ verbs, setVerbs }) {
    const [newVerb, setNewVerb] = useState('');
    const [editingIndex, setEditingIndex] = useState(-1);
    const [editingValue, setEditingValue] = useState('');

    const addVerb = () => {
        if (newVerb.trim()) {
            setVerbs([...verbs, { value: newVerb.trim(), enabled: true }]);
            setNewVerb('');
        }
    };

    const removeVerb = (index) => {
        const newVerbs = [...verbs];
        newVerbs.splice(index, 1);
        setVerbs(newVerbs);
    };

    const moveVerb = (index, direction) => {
        if (direction === 'up' && index > 0) {
            const newVerbs = [...verbs];
            const temp = newVerbs[index];
            newVerbs[index] = newVerbs[index - 1];
            newVerbs[index - 1] = temp;
            setVerbs(newVerbs);
        } else if (direction === 'down' && index < verbs.length - 1) {
            const newVerbs = [...verbs];
            const temp = newVerbs[index];
            newVerbs[index] = newVerbs[index + 1];
            newVerbs[index + 1] = temp;
            setVerbs(newVerbs);
        }
    };

    const toggleEnabled = (index) => {
        const newVerbs = [...verbs];
        newVerbs[index].enabled = !newVerbs[index].enabled;
        setVerbs(newVerbs);
    };

    const startEditing = (index, verb) => {
        setEditingIndex(index);
        setEditingValue(verb.value);
    };

    const saveEdit = (index) => {
        if (editingValue.trim()) {
            const newVerbs = [...verbs];
            newVerbs[index].value = editingValue.trim();
            setVerbs(newVerbs);
        }
        setEditingIndex(-1);
        setEditingValue('');
    };

    const cancelEdit = () => {
        setEditingIndex(-1);
        setEditingValue('');
    };

    return (
        <div className="verb-builder" style={{ padding: '1rem', border: '1px solid #ccc', marginBottom: '1rem' }}>
            <h3>Transformations</h3>
            <ul style={{ listStyle: 'none', padding: 0 }}>
                {verbs.map((verb, index) => (
                    <li
                        key={index}
                        style={{
                            display: 'flex',
                            alignItems: 'center',
                            background: '#f0f0f0',
                            padding: '0.5rem',
                            marginBottom: '0.5rem',
                            borderRadius: '4px',
                            opacity: verb.enabled ? 1 : 0.6
                        }}
                    >
                        <div style={{ display: 'flex', flexDirection: 'column', marginRight: '0.5rem' }}>
                            <button
                                onClick={() => moveVerb(index, 'up')}
                                disabled={index === 0}
                                style={{
                                    fontSize: '0.6rem',
                                    padding: '0 0.2rem',
                                    cursor: index === 0 ? 'default' : 'pointer',
                                    opacity: index === 0 ? 0.3 : 1
                                }}
                            >
                                ▲
                            </button>
                            <button
                                onClick={() => moveVerb(index, 'down')}
                                disabled={index === verbs.length - 1}
                                style={{
                                    fontSize: '0.6rem',
                                    padding: '0 0.2rem',
                                    cursor: index === verbs.length - 1 ? 'default' : 'pointer',
                                    opacity: index === verbs.length - 1 ? 0.3 : 1
                                }}
                            >
                                ▼
                            </button>
                        </div>
                        <input
                            type="checkbox"
                            checked={verb.enabled}
                            onChange={() => toggleEnabled(index)}
                            style={{ marginRight: '0.5rem', cursor: 'pointer' }}
                            title="Enable/Disable verb"
                        />
                        <span style={{ marginRight: '0.5rem', fontWeight: 'bold', color: '#555' }}>{index + 1}.</span>

                        {editingIndex === index ? (
                            <div style={{ flex: 1, display: 'flex', gap: '0.5rem' }}>
                                <input
                                    type="text"
                                    value={editingValue}
                                    onChange={(e) => setEditingValue(e.target.value)}
                                    style={{ flex: 1, padding: '0.25rem' }}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') saveEdit(index);
                                        if (e.key === 'Escape') cancelEdit();
                                    }}
                                    autoFocus
                                />
                                <button onClick={() => saveEdit(index)} style={{ color: 'green', cursor: 'pointer' }}>✓</button>
                                <button onClick={cancelEdit} style={{ color: 'gray', cursor: 'pointer' }}>✕</button>
                            </div>
                        ) : (
                            <>
                                <code
                                    style={{ flex: 1, color: '#333', cursor: 'pointer', borderBottom: '1px dashed #ccc', textDecoration: verb.enabled ? 'none' : 'line-through' }}
                                    onClick={() => startEditing(index, verb)}
                                    title="Click to edit"
                                >
                                    {verb.value}
                                </code>
                                <button onClick={() => removeVerb(index)} style={{ color: 'red', cursor: 'pointer', marginLeft: '0.5rem' }}>X</button>
                            </>
                        )}
                    </li>
                ))}
            </ul>
            {verbs.length === 0 && <p style={{ color: '#888' }}>No transformations added.</p>}
            <div style={{ display: 'flex', gap: '0.5rem', marginTop: '1rem' }}>
                <input
                    type="text"
                    value={newVerb}
                    onChange={(e) => setNewVerb(e.target.value)}
                    placeholder="e.g. sort -f field1"
                    style={{ flex: 1, padding: '0.5rem' }}
                    onKeyDown={(e) => e.key === 'Enter' && addVerb()}
                />
                <button onClick={addVerb}>Add Verb</button>
            </div>
            <div style={{ marginTop: '1rem', borderTop: '1px solid #eee', paddingTop: '1rem' }}>
                <h4 style={{ margin: '0 0 0.5rem 0', fontSize: '0.9rem', color: '#666' }}>Quick Add:</h4>
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem' }}>
                    {[
                        { label: 'Head 5', value: 'head -n 5' },
                        { label: 'Clean Headers', value: "rename -g -r ' ,_'" },
                        { label: 'Filter SKU', value: 'filter \'$SKU == "DAI-033"\'' },
                        { label: 'Label Cols', value: 'label sku,name,price,code' },
                        { label: 'Cut Cols', value: 'cut -f SKU,Price' },
                        { label: 'Add Split Col', value: 'put \'$first_word = splitax($Product_Name, " ")[1]\'' }
                    ].map((shortcut, i) => (
                        <button
                            key={i}
                            onClick={() => setVerbs([...verbs, { value: shortcut.value, enabled: true }])}
                            style={{
                                fontSize: '0.8rem',
                                padding: '0.25rem 0.5rem',
                                background: '#e1e1e1',
                                border: 'none',
                                borderRadius: '3px',
                                cursor: 'pointer'
                            }}
                            title={shortcut.value}
                        >
                            {shortcut.label}
                        </button>
                    ))}
                </div>
            </div>
        </div>
    );
}
