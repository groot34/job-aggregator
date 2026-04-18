'use client';

import { useEffect, useState } from 'react';
import api from '@/lib/api';
import { Job, JobResponse } from '@/types/job';
import JobCard from '@/components/JobCard';
import JobCardSkeleton from '@/components/JobCardSkeleton';
import JobDrawer from '@/components/JobDrawer';
import { useSavedJobs } from '@/hooks/useSavedJobs';
import { Search, MapPin, Briefcase } from 'lucide-react';

export default function Home() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [page, setPage] = useState(1);
  
  // Filters
  const [search, setSearch] = useState('');
  const [location, setLocation] = useState('');
  const [remoteOnly, setRemoteOnly] = useState(false);
  const [showSavedOnly, setShowSavedOnly] = useState(false);

  // Drawer
  const [selectedJob, setSelectedJob] = useState<Job | null>(null);

  const { savedJobIds, toggleSave, isSaved, isLoaded } = useSavedJobs();

  const fetchJobs = async () => {
    // If showSavedOnly is true and there are no saved jobs, don't fetch, just clear.
    if (showSavedOnly && savedJobIds.length === 0) {
      setJobs([]);
      setTotal(0);
      setTotalPages(0);
      setLoading(false);
      return;
    }

    setLoading(true);
    // Add artificial delay for skeleton demo
    await new Promise(resolve => setTimeout(resolve, 300));
    
    try {
      const params: any = { page, limit: 12, search };
      if (remoteOnly) params.remote = 'true';
      if (location) params.location = location;
      if (showSavedOnly && savedJobIds.length > 0) {
        params.savedIds = savedJobIds.join(',');
      }

      const res = await api.get<JobResponse>('/jobs', { params });
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
    if (!isLoaded) return;
    fetchJobs();
  }, [page, search, remoteOnly, location, showSavedOnly, isLoaded]); // Refetch when filters change

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1); // Reset to page 1 on new search
  };

  const handleToggleSave = (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    toggleSave(id);
  };

  return (
    <main className="min-h-screen bg-gray-50 flex flex-col">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-20 shadow-sm">
        <div className="max-w-[1400px] mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-2">
               <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-bold text-xl">
                 S
               </div>
               <h1 className="text-2xl font-bold text-gray-900 tracking-tight">ScrapHire</h1>
            </div>
            
            <form onSubmit={handleSearch} className="w-full max-w-lg relative hidden md:block">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Search jobs, companies, skills..."
                className="w-full pl-10 pr-4 py-2 bg-gray-100 border-transparent rounded-lg focus:bg-white focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all outline-none"
                value={search}
                onChange={(e) => {
                  setSearch(e.target.value);
                  setPage(1);
                }}
              />
            </form>
          </div>
        </div>
      </header>

      {/* Content Layout */}
      <div className="flex-1 max-w-[1400px] mx-auto w-full px-4 sm:px-6 lg:px-8 py-8 flex flex-col md:flex-row gap-8">
        
        {/* Sidebar */}
        <aside className="w-full md:w-64 flex-shrink-0">
          <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm sticky top-24">
            <h2 className="text-lg font-bold text-gray-900 mb-6 flex items-center gap-2">
              <Briefcase className="w-5 h-5 text-gray-500" />
              Filters
            </h2>
            
            {/* Mobile View Search */}
            <div className="mb-6 md:hidden">
              <label className="block text-sm font-medium text-gray-700 mb-2">Search</label>
              <input
                type="text"
                placeholder="Keywords..."
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                value={search}
                onChange={(e) => { setSearch(e.target.value); setPage(1); }}
              />
            </div>

            <div className="space-y-6">
               <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Location</label>
                  <div className="relative">
                    <MapPin className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                    <input
                      type="text"
                      placeholder="e.g. New York, Remote"
                      className="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      value={location}
                      onChange={(e) => { setLocation(e.target.value); setPage(1); }}
                    />
                  </div>
               </div>

               <div className="pt-4 border-t border-gray-100">
                 <label className="flex items-center cursor-pointer group">
                   <input
                     type="checkbox"
                     className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500 group-hover:border-blue-500"
                     checked={remoteOnly}
                     onChange={(e) => { setRemoteOnly(e.target.checked); setPage(1); }}
                   />
                   <span className="ml-2 text-sm font-medium text-gray-700 group-hover:text-gray-900">Remote Only</span>
                 </label>
               </div>

               <div className="pt-4 border-t border-gray-100">
                 <label className="flex items-center cursor-pointer group">
                   <input
                     type="checkbox"
                     className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500 group-hover:border-blue-500"
                     checked={showSavedOnly}
                     onChange={(e) => { setShowSavedOnly(e.target.checked); setPage(1); }}
                   />
                   <span className="ml-2 text-sm font-medium text-gray-700 group-hover:text-gray-900">Show Saved Jobs</span>
                 </label>
               </div>
            </div>
            
            <button 
              onClick={() => {
                setSearch('');
                setLocation('');
                setRemoteOnly(false);
                setShowSavedOnly(false);
                setPage(1);
              }}
              className="mt-8 w-full py-2 px-4 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
            >
              Clear Filters
            </button>
          </div>
        </aside>

        {/* Main Job List */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold text-gray-800">
              {showSavedOnly ? 'Your Saved Jobs' : 'Fresh Opportunities'}
              <span className="ml-3 text-sm font-medium text-gray-500 bg-gray-200/60 px-2.5 py-0.5 rounded-full">
                {total} {total === 1 ? 'job' : 'jobs'}
              </span>
            </h2>
          </div>

          {loading ? (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
              {Array.from({ length: 12 }).map((_, i) => (
                <JobCardSkeleton key={i} />
              ))}
            </div>
          ) : jobs.length === 0 ? (
            <div className="text-center py-20 bg-white border border-gray-200 rounded-xl">
              <h3 className="text-lg font-medium text-gray-900 mt-2">No jobs found</h3>
              <p className="mt-1 text-gray-500">Try adjusting your search or filters to find what you're looking for.</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
              {jobs.map((job) => (
                <JobCard 
                  key={job._id} 
                  job={job} 
                  onClick={setSelectedJob} 
                  isSaved={isSaved(job._id)} 
                  onToggleSave={handleToggleSave} 
                />
              ))}
            </div>
          )}

          {/* Pagination */}
          {!loading && totalPages > 1 && (
            <div className="mt-10 flex justify-center items-center gap-2 flex-wrap">
              <button
                disabled={page === 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                className="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 disabled:opacity-50 hover:bg-gray-50 font-medium transition-colors"
              >
                Previous
              </button>

              {(() => {
                const pageNumbers = [];
                const maxVisiblePages = 7;
                
                if (totalPages <= maxVisiblePages) {
                  for (let i = 1; i <= totalPages; i++) pageNumbers.push(i);
                } else {
                  if (page <= 3) {
                    for (let i = 1; i <= 4; i++) pageNumbers.push(i);
                    pageNumbers.push('...');
                    pageNumbers.push(totalPages);
                  } else if (page >= totalPages - 2) {
                    pageNumbers.push(1);
                    pageNumbers.push('...');
                    for (let i = totalPages - 3; i <= totalPages; i++) pageNumbers.push(i);
                  } else {
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
                      <span key={`ellipsis-${idx}`} className="px-2 py-2 text-gray-400">...</span>
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

              <button
                disabled={page >= totalPages}
                onClick={() => setPage((p) => p + 1)}
                className="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 disabled:opacity-50 hover:bg-gray-50 font-medium transition-colors"
              >
                Next
              </button>
            </div>
          )}
        </div>
      </div>

      <JobDrawer 
        job={selectedJob} 
        isOpen={selectedJob !== null} 
        onClose={() => setSelectedJob(null)} 
        isSaved={selectedJob ? isSaved(selectedJob._id) : false}
        onToggleSave={() => selectedJob && toggleSave(selectedJob._id)}
      />
    </main>
  );
}
