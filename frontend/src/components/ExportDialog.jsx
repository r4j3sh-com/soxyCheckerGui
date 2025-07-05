/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 * 
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

import React, { useState } from 'react';
import { XMarkIcon, ClipboardDocumentIcon, ArrowDownTrayIcon } from '@heroicons/react/24/outline';

function ExportDialog({ isOpen, onClose, workingProxies }) {
    const [format, setFormat] = useState('plain'); // plain, json, csv

    if (!isOpen) return null;

    const getFormattedProxies = () => {
        switch (format) {
            case 'json':
                return JSON.stringify(workingProxies, null, 2);
            case 'csv':
                return workingProxies.join('\n');
            case 'plain':
            default:
                return workingProxies.join('\n');
        }
    };

    const handleCopy = () => {
        const formattedProxies = getFormattedProxies();
        navigator.clipboard.writeText(formattedProxies)
            .then(() => {
                window.runtime.EventsEmit("log", "Proxies copied to clipboard");
            })
            .catch(err => {
                window.runtime.EventsEmit("log", "Failed to copy proxies to clipboard");
            });
    };

    const handleSave = () => {
        // Implement save functionality (can use Wails backend or download as file)
        window.runtime.EventsEmit("log", "Save functionality not implemented yet");
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm transition-all">
            <div className="relative w-full max-w-lg mx-4 md:mx-auto bg-gray-900 dark:bg-gray-900 rounded-2xl shadow-2xl ring-1 ring-gray-700 flex flex-col p-6 space-y-5 animate-fadeIn">
                {/* Close Button */}
                <button
                    className="absolute top-4 right-4 rounded-full p-1.5 bg-gray-700/50 hover:bg-gray-700 text-gray-300 hover:text-white transition"
                    onClick={onClose}
                    title="Close"
                >
                    <XMarkIcon className="h-6 w-6" />
                </button>
                {/* Title */}
                <h2 className="text-xl md:text-2xl font-bold text-center text-white mb-1">Export Working Proxies</h2>

                {/* Format Selector */}
                <div className="flex flex-wrap justify-center gap-3 mb-2">
                    {[
                        { val: 'plain', label: 'Plain Text' },
                        { val: 'json', label: 'JSON' },
                        { val: 'csv', label: 'CSV' },
                    ].map(opt => (
                        <label
                            key={opt.val}
                            className={`cursor-pointer flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium
                ${format === opt.val
                                    ? 'bg-indigo-600 text-white shadow-md'
                                    : 'bg-gray-800 text-gray-200 hover:bg-gray-700 transition'
                                }`}
                        >
                            <input
                                type="radio"
                                name="format"
                                value={opt.val}
                                checked={format === opt.val}
                                onChange={() => setFormat(opt.val)}
                                className="accent-indigo-600"
                            />
                            {opt.label}
                        </label>
                    ))}
                </div>

                {/* Proxies Preview */}
                <div>
                    <textarea
                        readOnly
                        value={getFormattedProxies()}
                        rows={10}
                        className="w-full rounded-xl border border-gray-700 bg-gray-800 text-gray-100 font-mono p-3 text-sm shadow-inner resize-y min-h-[180px] max-h-[360px] focus:outline-none"
                    />
                </div>

                {/* Buttons */}
                <div className="flex flex-wrap gap-3 justify-center mt-2">
                    <button
                        onClick={handleCopy}
                        className="inline-flex items-center gap-2 rounded-xl bg-indigo-600 px-4 py-2 font-semibold text-white shadow hover:bg-indigo-700 transition"
                        title="Copy to clipboard"
                    >
                        <ClipboardDocumentIcon className="h-5 w-5" />
                        Copy
                    </button>
                    <button
                        onClick={handleSave}
                        className="inline-flex items-center gap-2 rounded-xl bg-lime-600 px-4 py-2 font-semibold text-white shadow hover:bg-lime-700 transition"
                        title="Save to file"
                    >
                        <ArrowDownTrayIcon className="h-5 w-5" />
                        Save
                    </button>
                    <button
                        onClick={onClose}
                        className="inline-flex items-center gap-2 rounded-xl bg-gray-700 px-4 py-2 font-semibold text-gray-200 shadow hover:bg-gray-600 transition"
                        title="Close"
                    >
                        <XMarkIcon className="h-5 w-5" />
                        Close
                    </button>
                </div>
            </div>
        </div>
    );
}

export default ExportDialog;