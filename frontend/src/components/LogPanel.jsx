/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 * 
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

import React, { useRef, useEffect } from 'react';
import { TrashIcon } from '@heroicons/react/20/solid';

function LogPanel({ logs = [] }) {
    const logContainerRef = useRef(null);

    // Auto-scroll to the bottom when new logs are added
    useEffect(() => {
        if (logContainerRef.current) {
            const { scrollHeight, clientHeight } = logContainerRef.current;
            logContainerRef.current.scrollTop = scrollHeight - clientHeight;
        }
    }, [logs]);

    // Function to clear logs
    const handleClearLogs = () => {
        window.runtime?.EventsEmit?.("clear-logs");
    };

    return (
        <div className="log-panel mt-2 mb-2 rounded-2xl bg-white/90 dark:bg-gray-900/90 shadow-2xl ring-1 ring-gray-200 dark:ring-gray-700 flex flex-col h-full min-h-[160px]">
            {/* Header */}
            <div className="flex items-center justify-between px-4 py-2 border-b border-gray-200 dark:border-gray-700">
                <h3 className="text-base font-semibold text-gray-700 dark:text-gray-100">Log Output</h3>
                <button
                    className="inline-flex items-center gap-1.5 rounded-md bg-cyan-600 px-2.5 py-1 text-xs font-semibold text-white shadow-sm hover:bg-cyan-500 focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-600"
                    onClick={handleClearLogs}
                    title="Clear logs"
                >
                    <TrashIcon className="size-4 -ml-0.5" />
                    Clear Logs
                </button>
            </div>

            {/* Log List */}
            <div
                className="flex-1 overflow-y-auto px-4 py-2 text-sm font-mono bg-gray-50 dark:bg-gray-950 rounded-b-2xl"
                ref={logContainerRef}
                style={{ minHeight: 96, maxHeight: 224 }}
            >
                {logs.length === 0 ? (
                    <div className="text-gray-500 italic select-none">
                        No logs yet. Start a check to see output here.
                    </div>
                ) : (
                    logs.map((log, index) => (
                        <div key={index} className="flex gap-2 mb-0.5">
                            <span className="text-gray-400">{log.timestamp && `[${log.timestamp}]`}</span>
                            <span className="text-gray-700 dark:text-gray-100">{log.message}</span>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}

export default LogPanel;