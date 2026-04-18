import React, { useEffect } from 'react';
import { Job } from '../types/job';
import { X, MapPin, Calendar, ExternalLink, Building2, Bookmark, BookmarkCheck } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

interface JobDrawerProps {
  job: Job | null;
  isOpen: boolean;
  onClose: () => void;
  isSaved: boolean;
  onToggleSave: () => void;
}

export default function JobDrawer({ job, isOpen, onClose, isSaved, onToggleSave }: JobDrawerProps) {
  // Prevent body scroll when drawer is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-hidden">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-black/40 transition-opacity"
        onClick={onClose}
      />
      
      {/* Drawer */}
      <div className="absolute inset-y-0 right-0 max-w-full flex">
        <div className="w-screen max-w-2xl transform transition ease-in-out duration-300 translate-x-0 h-full bg-white shadow-xl flex flex-col">
          {/* Header */}
          <div className="px-6 py-4 border-b border-gray-200 flex items-center justify-between bg-white sticky top-0 z-10">
            <h2 className="text-xl font-bold text-gray-900 truncate pr-4">Job Details</h2>
            <div className="flex items-center space-x-2">
              <button
                onClick={onToggleSave}
                className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-full transition-colors flex items-center"
                title={isSaved ? "Remove from saved" : "Save job"}
              >
                {isSaved ? <BookmarkCheck className="w-5 h-5 text-blue-600" /> : <Bookmark className="w-5 h-5" />}
              </button>
              <button
                onClick={onClose}
                className="p-2 text-gray-400 hover:text-gray-500 hover:bg-gray-100 rounded-full transition-colors"
              >
                <X className="w-6 h-6" />
              </button>
            </div>
          </div>

          {/* Content */}
          {job && (
            <div className="flex-1 overflow-y-auto p-6">
              <div className="mb-6 border-b border-gray-200 pb-6">
                 <span className="inline-block px-2.5 py-1 text-sm font-medium text-blue-600 bg-blue-50 rounded-full mb-3">
                  {job.source}
                </span>
                <h1 className="text-3xl font-bold text-gray-900 mb-3">{job.title}</h1>
                <div className="flex items-center text-gray-700 text-lg mb-4 cursor-pointer hover:underline">
                  <Building2 className="w-5 h-5 mr-2" />
                  <span className="font-semibold">{job.company}</span>
                </div>
                
                <div className="flex flex-col sm:flex-row sm:items-center text-gray-500 space-y-2 sm:space-y-0 sm:space-x-6">
                  <div className="flex items-center">
                    <MapPin className="w-5 h-5 mr-1.5" />
                    <span>{job.remote ? 'Remote' : job.location}</span>
                  </div>
                  <div className="flex items-center">
                    <Calendar className="w-5 h-5 mr-1.5" />
                    <span>
                      {job.postedAt && !isNaN(new Date(job.postedAt).getTime()) 
                        ? 'Posted ' + formatDistanceToNow(new Date(job.postedAt), { addSuffix: true })
                        : 'Recently posted'}
                    </span>
                  </div>
                </div>
              </div>

              <div className="mb-8">
                <h3 className="text-lg font-bold text-gray-900 mb-3">Skills & Tags</h3>
                <div className="flex flex-wrap gap-2">
                  {job.tags.map((tag, i) => (
                    <span 
                      key={i} 
                      className="px-3 py-1 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg border border-gray-200"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              </div>

              <div>
                <h3 className="text-lg font-bold text-gray-900 mb-4">Job Description</h3>
                <div 
                  className="prose prose-gray max-w-none prose-p:mb-4 prose-a:text-blue-600 whitespace-pre-wrap"
                  dangerouslySetInnerHTML={{ __html: job.description }} 
                />
              </div>
            </div>
          )}

          {/* Footer */}
          {job && (
             <div className="p-4 border-t border-gray-200 bg-gray-50 flex items-center justify-between sticky bottom-0">
               <div>
                 <p className="text-xs text-gray-500">Scraped on {new Date(job.scrapedAt).toLocaleDateString()}</p>
               </div>
               <a
                 href={job.url}
                 target="_blank"
                 rel="noopener noreferrer"
                 className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg shadow-sm transition-colors flex items-center"
               >
                 Apply Externally
                 <ExternalLink className="w-4 h-4 ml-2" />
               </a>
             </div>
          )}
        </div>
      </div>
    </div>
  );
}
