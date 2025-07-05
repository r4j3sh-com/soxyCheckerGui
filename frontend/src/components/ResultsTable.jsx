/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 * 
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

import React, { useState, useMemo } from 'react';
import { ArrowDownIcon, ArrowUpIcon } from '@heroicons/react/20/solid';

function ResultsTable({ results = [] }) {
    const [sortConfig, setSortConfig] = useState({
        key: 'proxy',
        direction: 'ascending'
    });

    const requestSort = (key) => {
        let direction = 'ascending';
        if (sortConfig.key === key && sortConfig.direction === 'ascending') {
            direction = 'descending';
        }
        setSortConfig({ key, direction });
    };

    const sortedResults = useMemo(() => {
        const sortableResults = [...results];
        if (sortConfig.key) {
            sortableResults.sort((a, b) => {
                if (sortConfig.key === 'latency') {
                    const aValue = parseFloat(a[sortConfig.key]) || 0;
                    const bValue = parseFloat(b[sortConfig.key]) || 0;
                    return sortConfig.direction === 'ascending'
                        ? aValue - bValue
                        : bValue - aValue;
                }
                if (a[sortConfig.key] < b[sortConfig.key]) {
                    return sortConfig.direction === 'ascending' ? -1 : 1;
                }
                if (a[sortConfig.key] > b[sortConfig.key]) {
                    return sortConfig.direction === 'ascending' ? 1 : -1;
                }
                return 0;
            });
        }
        return sortableResults;
    }, [results, sortConfig]);

    const getSortIndicator = (name) => {
        if (sortConfig.key === name) {
            return sortConfig.direction === 'ascending'
                ? <ArrowUpIcon className="inline-block w-3 h-3 ml-1 -mt-1 text-gray-400" />
                : <ArrowDownIcon className="inline-block w-3 h-3 ml-1 -mt-1 text-gray-400" />;
        }
        return null;
    };

    const getStatusClass = (status) => {
        if (!status) return '';
        const statusLower = status.toLowerCase();
        if (statusLower === 'live' || statusLower === 'working') {
            return 'text-green-500 font-bold';
        } else if (statusLower === 'dead' || statusLower === 'failed') {
            return 'text-red-500 font-bold';
        } else if (statusLower.includes('error')) {
            return 'text-yellow-500 font-bold';
        }
        return '';
    };

    return (
        <div className="results-table-panel rounded-2xl bg-white/90 dark:bg-gray-900/90 shadow-2xl ring-1 ring-gray-200 dark:ring-gray-700 mt-1 mb-1 flex flex-col h-full min-h-[220px]">
            <div className="flex items-center justify-between px-4 py-2 border-b border-gray-200 dark:border-gray-700">
                <h3 className="text-base font-semibold text-gray-700 dark:text-gray-100">
                    Proxy Check Results
                </h3>
                <span className="text-xs text-gray-400">{results.length} Results</span>
            </div>
            <div className="flex-1 overflow-auto px-2 py-2">
                <table className="min-w-full text-xs md:text-sm table-auto">
                    <thead>
                        <tr className="text-gray-700 dark:text-gray-200 bg-gray-100 dark:bg-gray-800">
                            <th
                                className="px-3 py-2 cursor-pointer select-none"
                                onClick={() => requestSort('proxy')}
                            >
                                Proxy {getSortIndicator('proxy')}
                            </th>
                            <th
                                className="px-3 py-2 cursor-pointer select-none"
                                onClick={() => requestSort('type')}
                            >
                                Type {getSortIndicator('type')}
                            </th>
                            <th
                                className="px-3 py-2 cursor-pointer select-none"
                                onClick={() => requestSort('status')}
                            >
                                Status {getSortIndicator('status')}
                            </th>
                            <th
                                className="px-3 py-2 cursor-pointer select-none"
                                onClick={() => requestSort('latency')}
                            >
                                Latency {getSortIndicator('latency')}
                            </th>
                            <th
                                className="px-3 py-2 cursor-pointer select-none"
                                onClick={() => requestSort('outgoingIp')}
                            >
                                Outgoing IP {getSortIndicator('outgoingIp')}
                            </th>
                            <th
                                className="px-3 py-2 cursor-pointer select-none"
                                onClick={() => requestSort('geo')}
                            >
                                Geo {getSortIndicator('geo')}
                            </th>
                            <th
                                className="px-3 py-2 cursor-pointer select-none"
                                onClick={() => requestSort('error')}
                            >
                                Error {getSortIndicator('error')}
                            </th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200 dark:divide-gray-800">
                        {sortedResults.length === 0 ? (
                            <tr>
                                <td colSpan="7" className="py-10 text-center text-gray-400 italic">
                                    No results yet. Start a check to see proxy results here.
                                </td>
                            </tr>
                        ) : (
                            sortedResults.map((result, index) => (
                                <tr key={index} className="hover:bg-gray-50 dark:hover:bg-gray-800">
                                    <td className="px-3 py-1 break-all">{result.proxy}</td>
                                    <td className="px-3 py-1">{result.type}</td>
                                    <td className={`px-3 py-1 ${getStatusClass(result.status)}`}>
                                        {result.status}
                                    </td>
                                    <td className="px-3 py-1 text-right">
                                        {result.latency ? `${result.latency}ms` : '-'}
                                    </td>
                                    <td className="px-3 py-1">{result.outgoingIp || '-'}</td>
                                    <td className="px-3 py-1">{result.geo || '-'}</td>
                                    <td className="px-3 py-1 text-red-400 break-all">{result.error || '-'}</td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}

export default ResultsTable;