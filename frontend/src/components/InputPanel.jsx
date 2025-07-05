/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 * 
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */


import { useState, useEffect } from 'react';
import { PlayCircleIcon, StopCircleIcon, PauseCircleIcon, ArrowDownTrayIcon, TrashIcon, SunIcon, MoonIcon } from '@heroicons/react/20/solid';

export default function InputPanel({
    onStart, onStop, onPauseResume, onExport, onClear,
    isChecking, isPaused, isPausing, isStopping,
    canStart, canStop, canPause, canResume, canClear, canExport,
    proxyList, setProxyList  // Add proxyList as a prop
}) {
    const [proxyType, setProxyType] = useState('Auto');
    const [threads, setThreads] = useState(20);

    // Dark mode state
    const [darkMode, setDarkMode] = useState(true);

    useEffect(() => {
        document.documentElement.classList.add("dark");
    }, []);

    function toggleDarkMode() {
        setDarkMode(prev => {
            if (prev) {
                document.documentElement.classList.remove("dark");
            } else {
                document.documentElement.classList.add("dark");
            }
            return !prev;
        });
    }

    const handleStart = () => {
        const proxies = proxyList.split('\n').map(line => line.trim()).filter(Boolean);
        if (proxies.length === 0) {
            if (window.runtime?.EventsEmit) {
                window.runtime.EventsEmit("log", "No proxies entered");
            } else {
                alert("No proxies entered");
            }
            return;
        }
        const request = {
            ProxyList: proxies,
            ProxyType: proxyType.toLowerCase(),
            Threads: threads
        };
        onStart(request);
    };


    return (
        <div className="input-panel w-full rounded-2xl bg-white/90 dark:bg-gray-900/90 p-6 shadow-2xl ring-1 ring-gray-200 dark:ring-gray-700">
            {/* Mode Switcher */}
            <div className="flex items-center justify-end mb-4">
                <button
                    type="button"
                    className="inline-flex items-center gap-1 rounded-full bg-gray-200/70 dark:bg-gray-800/70 p-2 shadow hover:bg-gray-300 dark:hover:bg-gray-700 transition"
                    onClick={toggleDarkMode}
                >
                    {darkMode ? (
                        <>
                            <SunIcon className="size-5 text-yellow-400" />
                            <span className="text-sm text-gray-700 dark:text-gray-200">Light</span>
                        </>
                    ) : (
                        <>
                            <MoonIcon className="size-5 text-blue-500" />
                            <span className="text-sm text-gray-700 dark:text-gray-200">Dark</span>
                        </>
                    )}
                </button>
            </div>
            {/* Proxy Type */}
            <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Proxy Type</label>
                <select
                    value={proxyType}
                    onChange={e => setProxyType(e.target.value)}
                    className="block w-full rounded-md border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800 py-2 px-3 text-gray-900 dark:text-gray-100 shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 dark:focus:ring-indigo-900 transition"
                >
                    <option>Auto</option>
                    <option>HTTP</option>
                    <option>HTTPS</option>
                    <option>SOCKS4</option>
                    <option>SOCKS5</option>
                </select>
            </div>
            {/* Threads */}
            <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Threads</label>
                <input
                    type="number"
                    value={threads}
                    onChange={e => setThreads(parseInt(e.target.value) || 1)}
                    min="1"
                    max="100"
                    className="block w-full rounded-md border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800 py-2 px-3 text-gray-900 dark:text-gray-100 shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 dark:focus:ring-indigo-900 transition"
                />
            </div>
            {/* Proxies */}
            <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Proxies</label>
                <textarea
                    value={proxyList}
                    onChange={(e) => setProxyList(e.target.value)}
                    placeholder="Paste proxies here, one per line (ip:port)"
                    rows={8}
                    className="
    block w-full rounded-md border border-gray-300 dark:border-gray-700
    bg-gray-50 dark:bg-gray-800 py-2 px-3 text-gray-900 dark:text-gray-100
    shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 dark:focus:ring-indigo-900
    transition font-mono resize-y
    min-h-[120px] max-h-[260px]
    md:min-h-[160px] md:max-h-[320px]
    lg:min-h-[200px] lg:max-h-[400px]
  "
                />
            </div>
            {/* Buttons */}
            <div className="flex flex-col gap-2 w-full mt-2">

                <button
                    type="button"
                    className={`inline-flex items-center justify-center gap-x-2 rounded-xl 
                        ${isChecking
                            ? canStop ? "bg-orange-600 hover:bg-orange-700" : "bg-orange-400 cursor-not-allowed"
                            : canStart ? "bg-indigo-600 hover:bg-indigo-700" : "bg-indigo-400 cursor-not-allowed"} 
                        px-4 py-2 text-base font-semibold text-white shadow-md 
                        focus-visible:outline-none focus-visible:ring-2 
                        focus-visible:ring-${isChecking ? "orange" : "indigo"}-400 
                        transition w-full`}
                    onClick={isChecking ? onStop : handleStart}
                    disabled={isChecking ? !canStop : (!canStart || !proxyList.length)}
                    title={isChecking ? "Stop checking" : "Start checking proxies"}
                >
                    {isChecking ? (
                        <>
                            <StopCircleIcon aria-hidden="true" className="size-6" />
                            {isStopping ? "Stopping..." : "Stop Check"}
                        </>
                    ) : (
                        <>
                            <PlayCircleIcon aria-hidden="true" className="size-6" />
                            Start Check
                        </>
                    )}
                </button>
                <div className="grid grid-cols-2 gap-2">


                    <button
                        type="button"
                        className={`inline-flex items-center justify-center gap-x-2 rounded-xl 
                                ${isPaused
                                ? canResume ? "bg-green-600 hover:bg-green-700" : "bg-green-400 cursor-not-allowed"
                                : isPausing
                                    ? "bg-yellow-600 hover:bg-yellow-700"
                                    : canPause ? "bg-red-600 hover:bg-red-700" : "bg-red-400 cursor-not-allowed"} 
                                                px-4 py-2 text-base font-semibold text-white shadow-md 
                                                focus-visible:outline-none focus-visible:ring-2 
                                                focus-visible:ring-${isPaused ? "green" : isPausing ? "yellow" : "red"}-400 
                                                transition`}
                        onClick={onPauseResume}
                        disabled={isPausing || !(canPause || canResume)}
                        title={isPaused ? "Resume checking" : isPausing ? "Pausing..." : "Pause checking"}
                    >
                        {isPaused ? (
                            <PlayCircleIcon aria-hidden="true" className="size-6" />
                        ) : isPausing ? (
                            <PauseCircleIcon aria-hidden="true" className="size-6 animate-pulse" />
                        ) : (
                            <PauseCircleIcon aria-hidden="true" className="size-6" />
                        )}
                        {isPaused ? "Resume" : isPausing ? "Pausing..." : "Pause"}
                    </button>
                    <button
                        type="button"
                        className={`inline-flex items-center justify-center gap-x-2 rounded-xl 
                            ${canClear ? "bg-cyan-600 hover:bg-cyan-700" : "bg-cyan-400 cursor-not-allowed"}
                            px-4 py-2 text-base font-semibold text-white shadow-md 
                            focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400 
                            transition`}
                        onClick={onClear}
                        disabled={!canClear}
                        title="Clear all results"
                    >
                        <TrashIcon aria-hidden="true" className="size-6" />
                        Clear
                    </button>
                </div>
                <button
                    type="button"
                    className={`mt-2 w-full inline-flex items-center justify-center gap-x-2 rounded-xl 
                        ${canExport ? "bg-lime-600 hover:bg-lime-700" : "bg-lime-400 cursor-not-allowed"}
                        px-4 py-2 text-base font-semibold text-white shadow-md 
                        focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-lime-400 
                        transition`}
                    onClick={onExport}
                    disabled={!canExport}
                    title="Export working proxies"
                >
                    <ArrowDownTrayIcon aria-hidden="true" className="size-6" />
                    Export Working
                </button>
            </div>
        </div>
    );
}