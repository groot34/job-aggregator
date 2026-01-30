'use client';

import { useEffect, useState } from 'react';
import api from '@/lib/api';
import { Job, JobResponse } from '@/types/job';
import JobCard from '@/components/JobCard';
import { Search, Loader2 } from 'lucide-react';

export default function Home() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');

  const fetchJobs = async () => {
    setLoading(true);
    try {
      const res = await api.get<JobResponse>('/jobs', {
        params: { page, limit: 12, search }
      });
      setJobs(res.data.jobs);
      setTotal(res.data.total);
      setTotalPages(res.data.totalPages);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchJobs();
  }, [page, search]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1); // Reset to page 1 on new search
    fetchJobs();
  };

  return (
    <main className="min-h-screen bg-gray-50 pb-20">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-10 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div className="flex items-center gap-2">
               <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-bold text-xl">
                 J
               </div>
               <h1 className="text-2xl font-bold text-gray-900 tracking-tight">JobAggregator</h1>
            </div>

            <form onSubmit={handleSearch} className="w-full md:w-96 relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Search jobs, companies, skills..."
                className="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all outline-none"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </form>
          </div>
        </div>
      </header>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex items-center justify-between mb-8">
          <h2 className="text-xl font-semibold text-gray-800">
            Fresh Opportunities
            <span className="ml-2 text-sm font-normal text-gray-500 bg-gray-200 px-2 py-0.5 rounded-full">{total} jobs found</span>
          </h2>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64">
            <Loader2 className="w-10 h-10 text-blue-600 animate-spin" />
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {jobs.map((job) => (
              <JobCard key={job._id} job={job} />
            ))}
          </div>
        )}

        {/* Pagination */}
        <div className="mt-12 flex justify-center items-center gap-2 flex-wrap">
          {/* Previous Button */}
          <button
            disabled={page === 1}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            className="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 disabled:opacity-50 hover:bg-gray-50 font-medium transition-colors"
          >
            Previous
          </button>

          {/* Page Numbers */}
          {totalPages > 0 && (() => {
            const pageNumbers = [];
            const maxVisiblePages = 7; // Show max 7 page buttons
            
            if (totalPages <= maxVisiblePages) {
              // Show all pages if total is small
              for (let i = 1; i <= totalPages; i++) {
                pageNumbers.push(i);
              }
            } else {
              // Smart pagination: 1 ... 4 5 6 ... 10
              if (page <= 3) {
                // Near start: 1 2 3 4 ... 10
                for (let i = 1; i <= 4; i++) pageNumbers.push(i);
                pageNumbers.push('...');
                pageNumbers.push(totalPages);
              } else if (page >= totalPages - 2) {
                // Near end: 1 ... 7 8 9 10
                pageNumbers.push(1);
                pageNumbers.push('...');
                for (let i = totalPages - 3; i <= totalPages; i++) pageNumbers.push(i);
              } else {
                // Middle: 1 ... 4 5 6 ... 10
                pageNumbers.push(1);
                pageNumbers.push('...');
                for (let i = page - 1; i <= page + 1; i++) pageNumbers.push(i);
                pageNumbers.push('...');
                pageNumbers.push(totalPages);
              }
            }

            return pageNumbers.map((num, idx) => {
              if (num === '...') {
                return (
                  <span key={`ellipsis-${idx}`} className="px-2 py-2 text-gray-400">
                    ...
                  </span>
                );
              }

              const isActive = num === page;
              return (
                <button
                  key={num}
                  onClick={() => setPage(num as number)}
                  className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                    isActive
                      ? 'bg-blue-600 text-white'
                      : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'
                  }`}
                >
                  {num}
                </button>
              );
            });
          })()}

          {/* Next Button */}
          <button
            disabled={page >= totalPages}
            onClick={() => setPage((p) => p + 1)}
            className="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 disabled:opacity-50 hover:bg-gray-50 font-medium transition-colors"
          >
            Next
          </button>
        </div>
      </div>
    </main>
  );
}
