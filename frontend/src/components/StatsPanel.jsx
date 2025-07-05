/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 * 
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

import React from 'react';

export default function StatsPanel({ stats }) {
    const {
        Total = 0,
        Pending = 0,
        Live = 0,
        Dead = 0,
        Errors = 0,
        TypeCounts = {}
    } = stats || {};
    const successRate = Total > 0 ? Math.round((Live / Total) * 100) : 0;

    return (
        <div className="rounded-2xl bg-white/90 dark:bg-gray-900/90 shadow-2xl ring-1 ring-gray-200 dark:ring-gray-700 p-3 flex flex-col items-center gap-1 text-sm">
            <div>
                <span className="text-gray-500 dark:text-gray-400">Total: </span>
                <span className="font-semibold">{Total}</span>
            </div>
            <div>
                <span className="text-gray-500 dark:text-gray-400">Pending: </span>
                <span className="font-semibold">{Pending}</span>
            </div>
            <div>
                <span className="text-gray-500 dark:text-gray-400">Live: </span>
                <span className="text-green-500 font-semibold">{Live}</span>
            </div>
            <div>
                <span className="text-gray-500 dark:text-gray-400">Dead: </span>
                <span className="text-red-500 font-semibold">{Dead}</span>
            </div>
            <div>
                <span className="text-gray-500 dark:text-gray-400">Errors: </span>
                <span className="text-red-500 font-semibold">{Errors}</span>
            </div>
            <div>
                <span className="text-gray-500 dark:text-gray-400">Success: </span>
                <span className="font-semibold">{successRate}%</span>
            </div>
            {/* Type counts */}
            <div className="mt-2 space-y-0.5 text-xs">
                <div>HTTP: <span className="font-mono">{TypeCounts.http || 0}</span></div>
                <div>HTTPS: <span className="font-mono">{TypeCounts.https || 0}</span></div>
                <div>SOCKS4: <span className="font-mono">{TypeCounts.socks4 || 0}</span></div>
                <div>SOCKS5: <span className="font-mono">{TypeCounts.socks5 || 0}</span></div>
            </div>
        </div>
    );
}