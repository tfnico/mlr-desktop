import { useState, useEffect, useCallback } from 'react';
import './App.css';
import InputSection from './components/InputSection';
import VerbBuilder from './components/VerbBuilder';
import OutputPreview from './components/OutputPreview';
import ErrorBoundary from './components/ErrorBoundary';
import logger from './utils/logger';
import { Preview, SaveConfig, LoadConfig, ExportGoProgram, ReadFileHead, SaveLastState, LoadLastState, GetCommand, SaveOutput, ParseCommand } from '../wailsjs/go/main/App';

const DEFAULT_INPUT_CONTENT = `SKU,Product Name,Price,Barcode
FRO-010,Organic Free-Range Eggs (Dozen),5.99,5012345678901
PNC-025,Fresh Baked Sourdough Loaf,4.25,5023456789012
DRY-050,Basmati Rice (5kg),12.50,5034567890123
PRO-005,Roma Tomatoes (Per Lb),2.99,5045678901234
CNS-012,Black Beans (Canned),0.89,5056789012345
DAI-033,Whole Milk (Gallon),3.75,5067890123456
BEV-070,Sparkling Water (12-Pack),6.50,5078901234567
CLE-045,All-Purpose Kitchen Cleaner,4.99,5089012345678
PET-099,Adult Dog Food (10kg),25.99,5090123456789
BAK-001,All-Purpose Flour (2kg),3.20,5001234567890`;

function App() {
    const [inputMode, setInputMode] = useState('text'); // 'text' or 'file'
    const [inputValue, setInputValue] = useState(DEFAULT_INPUT_CONTENT); // text content or file path
    const [options, setOptions] = useState('');
    const [inputFormat, setInputFormat] = useState('');
    const [ragged, setRagged] = useState(false);
    const [headerless, setHeaderless] = useState(false);
    const [fieldSeparator, setFieldSeparator] = useState(',');
    const [outputFormat, setOutputFormat] = useState('');
    const [verbs, setVerbs] = useState([]);
    const [output, setOutput] = useState('');
    const [error, setError] = useState('');
    const [command, setCommand] = useState('');

    // Load last state on startup
    useEffect(() => {
        const loadState = async () => {
            try {
                const config = await LoadLastState();
                if (config) {
                    setInputValue(config.inputPath || '');
                    setVerbs(config.verbs || []);
                    setOptions(config.options || '');
                    setInputMode(config.inputMode || 'text');
                    setInputFormat(config.inputFormat || '');
                    setRagged(config.ragged || false);
                    setHeaderless(config.headerless || false);
                    setFieldSeparator(config.fieldSeparator || ',');
                    setOutputFormat(config.outputFormat || '');
                }
            } catch (err) {
                logger.logError(err, { context: 'LoadLastState' });
                console.log("No previous state found or error loading:", err);
            }
        };
        loadState();
    }, []);

    const updatePreview = useCallback(async () => {
        setError('');
        if (!inputValue) {
            setOutput('');
            setCommand('');
            return;
        }

        // Check if there are any enabled verbs
        const hasEnabledVerbs = verbs.some(verb => verb.enabled);
        if (!hasEnabledVerbs) {
            setOutput('');
            setCommand('');
            return;
        }

        let inputContent = inputValue;
        if (inputMode === 'file') {
            try {
                // Read first 100 lines for preview
                // Check if path is not empty
                if (!inputValue.trim()) return;
                inputContent = await ReadFileHead(inputValue, 100);
            } catch (err) {
                const errorMsg = `Error reading file: ${err}`;
                logger.error(errorMsg, { context: 'ReadFileHead', path: inputValue });
                setError(errorMsg);
                return;
            }
        }

        try {
            const result = await Preview(inputContent, verbs, options, inputFormat, ragged, headerless, fieldSeparator, outputFormat);
            setOutput(result);
            const cmd = await GetCommand(verbs, options, inputFormat, ragged, headerless, fieldSeparator, outputFormat, inputMode, inputValue);
            setCommand(cmd);
            // Auto-save state on success
            SaveLastState({ inputPath: inputValue, inputMode, inputFormat, ragged, headerless, fieldSeparator, outputFormat, verbs, options });
        } catch (err) {
            logger.logError(err, { context: 'Preview', verbs, inputFormat, outputFormat });
            setError(String(err));
        }
    }, [inputValue, inputMode, verbs, options, inputFormat, ragged, headerless, fieldSeparator, outputFormat]);

    useEffect(() => {
        const timer = setTimeout(() => {
            updatePreview();
        }, 500); // Debounce
        return () => clearTimeout(timer);
    }, [updatePreview]);

    const handleSaveOutput = async () => {
        if (!output) return;
        try {
            await SaveOutput(output);
        } catch (err) {
            logger.logError(err, { context: 'SaveOutput' });
            alert("Error saving output: " + err);
        }
    };

    const handleImportCommand = async () => {
        const commandStr = prompt("Paste your mlr command:");
        if (!commandStr) return;

        try {
            const config = await ParseCommand(commandStr);

            // Update all state from the parsed config
            setInputFormat(config.inputFormat || '');
            setOutputFormat(config.outputFormat || '');
            setOptions(config.options || '');
            setRagged(config.ragged || false);
            setHeaderless(config.headerless || false);
            setFieldSeparator(config.fieldSeparator || ',');
            setVerbs(config.verbs || []);

            // Update input mode and path if present
            if (config.inputPath) {
                setInputMode('file');
                setInputValue(config.inputPath);
            }

            logger.info("Command imported successfully", {
                verbs_count: config.verbs?.length,
                input_mode: config.inputMode,
                has_input_path: !!config.inputPath
            });
        } catch (err) {
            logger.logError(err, { context: 'ImportCommand', command: commandStr });
            alert("Error parsing command: " + err);
        }
    };

    return (
        <ErrorBoundary>
            <div id="app" className="App">
                <header style={{ padding: '1rem', background: '#282c34', color: 'white', marginBottom: '1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h1 style={{ margin: 0 }}>MLR Desktop Tool</h1>
                    <button
                        onClick={handleImportCommand}
                        style={{
                            padding: '0.5rem 1rem',
                            backgroundColor: '#4CAF50',
                            color: 'white',
                            border: 'none',
                            borderRadius: '4px',
                            cursor: 'pointer',
                            fontWeight: 'bold'
                        }}
                        onMouseOver={(e) => e.target.style.backgroundColor = '#45a049'}
                        onMouseOut={(e) => e.target.style.backgroundColor = '#4CAF50'}
                    >
                        Import Command
                    </button>
                </header>
                <main style={{ padding: '1rem', width: 'calc(100% - 2rem)', margin: '0 auto' }}>
                    <InputSection
                        mode={inputMode}
                        inputValue={inputValue}
                        options={options}
                        inputFormat={inputFormat}
                        ragged={ragged}
                        headerless={headerless}
                        fieldSeparator={fieldSeparator}
                        onInputChange={(val, mode, opts, fmt, rag, hdl, fs) => {
                            if (val !== null) setInputValue(val);
                            if (mode !== null) setInputMode(mode);
                            if (opts !== null) setOptions(opts);
                            if (fmt !== null) setInputFormat(fmt);
                            if (rag !== null) setRagged(rag);
                            if (hdl !== null) setHeaderless(hdl);
                            if (fs !== null) setFieldSeparator(fs);
                        }}
                        onModeChange={(mode) => setInputMode(mode)}
                    />
                    <VerbBuilder verbs={verbs} setVerbs={setVerbs} />
                    <OutputPreview
                        output={output}
                        error={error}
                        outputFormat={outputFormat}
                        onOutputFormatChange={setOutputFormat}
                        command={command}
                        onSave={handleSaveOutput}
                    />
                </main>
            </div>
        </ErrorBoundary>
    );
}

export default App;
