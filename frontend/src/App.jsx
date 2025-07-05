/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 * 
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

import { useState, useEffect, useRef } from 'react';
import './App.css';
import InputPanel from './components/InputPanel';
import ResultsTable from './components/ResultsTable';
import LogPanel from './components/LogPanel';
import ExportDialog from './components/ExportDialog';
//import StatsPanel from './components/StatsPanel';
import { StartCheck, StopCheck, PauseCheck, ResumeCheck, GetWorkingProxies, ClearResults } from '../wailsjs/go/backend/App';

export default function App() {
    const [results, setResults] = useState([]);
    const [proxyList, setProxyList] = useState('');
    const [stats, setStats] = useState({
        Total: 0,
        Pending: 0,
        Live: 0,
        Dead: 0,
        Errors: 0,
        TypeCounts: {}
    });
    const [logs, setLogs] = useState([]);
    const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
    const [workingProxies, setWorkingProxies] = useState([]);
    const [isChecking, setIsChecking] = useState(false);
    const [isPaused, setIsPaused] = useState(false);
    const [isPausing, setIsPausing] = useState(false);
    const [isStopping, setIsStopping] = useState(false);

    // Top-wide controls state
    const [upstreamProxy, setUpstreamProxy] = useState('');
    const [upstreamType, setUpstreamType] = useState('HTTP');
    const [endpoint, setEndpoint] = useState('https://api.ipify.org');
    const pauseTimeoutRef = useRef(null);


    const handlePauseResumeCheck = async () => {
        try {
            if (isPaused) {
                await ResumeCheck();
                // When resuming, we should immediately update the UI state
                setIsPaused(false);
                setIsChecking(true);
            } else {
                setIsPausing(true);
                // Clear any existing timeout
                if (pauseTimeoutRef.current) {
                    clearTimeout(pauseTimeoutRef.current);
                }

                // Set a timeout to force transition to paused state after 5 seconds
                pauseTimeoutRef.current = setTimeout(() => {
                    if (isPausing) {
                        console.log("Pause timeout reached, forcing paused state");
                        setIsPausing(false);
                        setIsPaused(true);
                        // Important: Keep isChecking true when paused
                        setIsChecking(true);
                        window.runtime.EventsEmit("log", "Check paused (timeout reached)");
                    }
                }, 5000); // 5 seconds timeout
                await PauseCheck();
            }
        } catch (error) {
            if (!isPaused) setIsPausing(false);
            window.runtime.EventsEmit("log", `Error ${isPaused ? "resuming" : "pausing"} check: ${error.message}`);
        }
    };

    // Make sure the stop buttons are calling the correct functions
    const handleStopCheck = async () => {
        try {
            setIsStopping(true);
            await StopCheck();
        }
        catch (error) {
            setIsStopping(false);
            window.runtime.EventsEmit("log", `Error stopping check: ${error.message}`);
        }
    };


    const handleExport = async () => {
        try {
            const proxies = await GetWorkingProxies();
            if (proxies && proxies.length > 0) {
                setWorkingProxies(proxies);
                setIsExportDialogOpen(true);
            } else {
                window.runtime.EventsEmit("log", "No working proxies to export");
            }
        } catch (error) {
            window.runtime.EventsEmit("log", `Error getting working proxies: ${error.message}`);
        }
    };

    const handleClearResults = () => {
        setResults([]);
        setStats({
            Total: 0,
            Pending: 0,
            Live: 0,
            Dead: 0,
            Errors: 0,
            TypeCounts: {}
        });
        // Clear proxy list
        setProxyList('');
        // Also tell the backend to clear its results
        // This will prevent any pending results from reappearing
        if (window.runtime?.EventsEmit) {
            window.runtime.EventsEmit("clear-results");
            window.runtime.EventsEmit("log", "Results and proxy list cleared");
        }
    };

    useEffect(() => {
        window.runtime.EventsOn("log", (message) => {
            setLogs(prevLogs => [...prevLogs, {
                timestamp: new Date().toLocaleTimeString(),
                message
            }]);
        });
        window.runtime.EventsOn("results-update", setResults);
        window.runtime.EventsOn("stats-update", setStats);

        window.runtime.EventsOn("check-status", (status) => {
            console.log("Check status:", status);
            // Clear the pause timeout if we get a definitive status
            if (status !== "pausing" && pauseTimeoutRef.current) {
                clearTimeout(pauseTimeoutRef.current);
                pauseTimeoutRef.current = null;
            }

            // Important: Keep isChecking true when paused
            setIsChecking(status === "running" || status === "pausing" || status === "paused");
            setIsPaused(status === "paused");
            setIsPausing(status === "pausing");

            if (status === "stopped") {
                setIsStopping(false);
            }
        });
        window.runtime.EventsOn("pause-progress", (progress) => {
            console.log("Pause progress:", progress);
            // Handle the progress safely to avoid NaN issues
            if (progress && progress.total > 0) {
                // Only process valid progress data
                console.log(`Paused ${progress.paused}/${progress.total} workers (${Math.round(progress.percent)}%)`);
            }
        });
        window.runtime.EventsOn("clear-results", async () => {
            try {
                await ClearResults();
            } catch (error) {
                console.error("Error clearing results:", error);
            }
        });
        window.runtime.EventsOn("clear-logs", () => setLogs([]));

        return () => {
            window.runtime.EventsOff("log");
            window.runtime.EventsOff("results-update");
            window.runtime.EventsOff("stats-update");
            window.runtime.EventsOff("check-status");
            window.runtime.EventsOff("pause-progress");
            window.runtime.EventsOff("clear-logs");
            if (pauseTimeoutRef.current) {
                clearTimeout(pauseTimeoutRef.current);
            }
            window.runtime.EventsOff("clear-results");
        };
    }, []);

    const handleStartCheck = async (params) => {
        // Merge in top controls
        const checkParams = {
            ...params,
            Endpoint: endpoint,
            ...(upstreamProxy && { UpstreamProxy: upstreamProxy, UpstreamType: upstreamType.toLowerCase() })
        };
        try {
            await StartCheck(checkParams);
            setIsChecking(true);
            setIsPaused(false);
            setIsPausing(false);
        } catch (error) {
            window.runtime.EventsEmit("log", `Error starting check: ${error.message}`);
        }
    };




    // Determine if buttons should be enabled
    const hasResults = results.length > 0;
    const hasWorkingProxies = stats.Live > 0;
    const canStartCheck = !isChecking && !isPausing && !isStopping;
    const canStopCheck = isChecking && !isPaused && !isPausing && !isStopping;
    const canPauseCheck = isChecking && !isPaused && !isPausing && !isStopping;
    const canResumeCheck = isPaused && !isPausing && !isStopping;
    const canClearResults = hasResults && !isChecking && !isPausing && !isStopping;
    const canExportProxies = hasWorkingProxies && !isPausing && !isStopping;


    // --- Responsive, professional stats bar ---
    function StatsNavbar() {
        return (
            <nav className="w-full flex flex-wrap items-center justify-between gap-2 rounded-lg bg-gray-800/80 px-6 py-3 mb-3 shadow ring-1 ring-gray-800 dark:ring-gray-700">
                <span className="text-xs md:text-sm font-semibold text-gray-400">Total: <span className="text-white">{stats.Total}</span></span>
                <span className="text-xs md:text-sm font-semibold text-yellow-400">Pending: <span className="text-white">{stats.Pending}</span></span>
                <span className="text-xs md:text-sm font-semibold text-green-400">Live: <span className="text-white">{stats.Live}</span></span>
                <span className="text-xs md:text-sm font-semibold text-red-400">Dead: <span className="text-white">{stats.Dead}</span></span>
                <span className="text-xs md:text-sm font-semibold text-pink-400">Errors: <span className="text-white">{stats.Errors}</span></span>
                <span className="text-xs md:text-sm font-semibold text-blue-300">HTTP: <span className="text-white">{stats.TypeCounts.http || 0}</span></span>
                <span className="text-xs md:text-sm font-semibold text-blue-300">HTTPS: <span className="text-white">{stats.TypeCounts.https || 0}</span></span>
                <span className="text-xs md:text-sm font-semibold text-blue-300">SOCKS4: <span className="text-white">{stats.TypeCounts.socks4 || 0}</span></span>
                <span className="text-xs md:text-sm font-semibold text-blue-300">SOCKS5: <span className="text-white">{stats.TypeCounts.socks5 || 0}</span></span>
            </nav>
        );
    }


    return (
        <div className="w-screen h-screen min-h-screen bg-gray-900 dark:bg-black text-gray-200 flex flex-col">
            {/* --- Top controls: Wide bar --- */}
            <div className="w-full px-6 py-3 flex flex-col gap-2 md:flex-row md:items-end md:gap-6 bg-gray-800/90 shadow-lg z-10">
                <div className="flex-1 flex flex-col md:flex-row gap-2">
                    {/* Upstream Proxy */}
                    <div className="flex-1">
                        <label className="block text-xs font-medium text-gray-400 dark:text-gray-400 mb-1">Upstream Proxy</label>
                        <input
                            type="text"
                            value={upstreamProxy}
                            onChange={e => setUpstreamProxy(e.target.value)}
                            placeholder="ip:port"
                            className="block w-full rounded-md border border-gray-700 bg-gray-900 text-gray-100 py-2 px-3 shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-900 transition font-mono"
                        />
                    </div>
                    {/* Upstream Type */}
                    <div className="w-26">
                        <label className="block text-xs font-medium text-gray-400 dark:text-gray-400 mb-1">Proxy Type</label>
                        <select
                            value={upstreamType}
                            onChange={e => setUpstreamType(e.target.value)}
                            className="block w-full rounded-md border border-gray-700 bg-gray-900 text-gray-100 py-2 px-3 shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-900 transition"
                        >
                            <option>HTTP</option>
                            <option>HTTPS</option>
                            <option>SOCKS4</option>
                            <option>SOCKS5</option>
                        </select>
                    </div>
                </div>
                {/* Endpoint */}
                <div className="flex-1">
                    <label className="block text-xs font-medium text-gray-400 dark:text-gray-400 mb-1">Endpoint</label>
                    <input
                        type="text"
                        value={endpoint}
                        onChange={e => setEndpoint(e.target.value)}
                        placeholder="https://api.ipify.org"
                        className="block w-full rounded-md border border-gray-700 bg-gray-900 text-gray-100 py-2 px-3 shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-900 transition font-mono"
                    />
                </div>
            </div>

            {/* --- Stats navbar under top bar --- */}
            <div className="flex flex-1 overflow-hidden">
                <div className="w-full max-w-xs flex-shrink-0 flex flex-col p-4">
                    <InputPanel
                        onStart={handleStartCheck}
                        onStop={handleStopCheck}
                        onPauseResume={handlePauseResumeCheck}
                        onExport={handleExport}
                        onClear={handleClearResults}
                        isChecking={isChecking}
                        isPaused={isPaused}
                        isPausing={isPausing}
                        isStopping={isStopping}
                        canStart={canStartCheck}
                        canStop={canStopCheck}
                        canPause={canPauseCheck}
                        canResume={canResumeCheck}
                        canClear={canClearResults}
                        canExport={canExportProxies}
                        proxyList={proxyList}      // Make sure to pass this
                        setProxyList={setProxyList} // And this
                    />
                </div>
                <div className="flex-1 flex flex-col overflow-hidden px-2 py-4 gap-2">
                    {/* STATS NAVBAR */}
                    <StatsNavbar />

                    {/* Results table fills space and scrolls */}
                    <div className="flex-1 min-h-0 min-w-0">
                        <ResultsTable results={results} />
                    </div>
                    {/* Bottom: log panel, fixed height, always at bottom, scrollable */}
                    <div className="h-48 min-h-[10rem]">
                        <LogPanel logs={logs} />
                    </div>
                </div>
            </div>
            {/* Export Dialog */}
            <ExportDialog
                isOpen={isExportDialogOpen}
                onClose={() => setIsExportDialogOpen(false)}
                workingProxies={workingProxies}
            />
        </div>
    );
}