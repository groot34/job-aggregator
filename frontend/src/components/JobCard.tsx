import React from 'react';
import { Job } from '../types/job';
import { MapPin, Calendar, ExternalLink, Building2 } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

interface JobCardProps {
  job: Job;
}

const JobCard: React.FC<JobCardProps> = ({ job }) => {
  return (
    <div className="bg-white border border-gray-200 rounded-xl p-6 hover:shadow-lg transition-shadow duration-300 flex flex-col justify-between h-full group relative overflow-hidden">
      {/* Decorative gradient blob */}
      <div className="absolute top-0 right-0 w-24 h-24 bg-blue-50 rounded-bl-full -mr-4 -mt-4 transition-transform group-hover:scale-110" />

      <div>
        <div className="flex justify-between items-start mb-4">
          <div>
            <span className="inline-block px-2 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded-full mb-2">
              {job.source}
            </span>
            <h3 className="text-lg font-bold text-gray-900 leading-tight group-hover:text-blue-600 transition-colors">
              {job.title}
            </h3>
            <div className="flex items-center text-gray-600 mt-1">
              <Building2 className="w-4 h-4 mr-1.5" />
              <span className="font-medium text-sm">{job.company}</span>
            </div>
          </div>
        </div>

        <div className="flex items-center text-sm text-gray-500 mb-4 space-x-4">
          <div className="flex items-center">
            <MapPin className="w-4 h-4 mr-1" />
            <span>{job.remote ? 'Remote' : job.location}</span>
          </div>
          <div className="flex items-center">
            <Calendar className="w-4 h-4 mr-1" />
            <span>
              {job.postedAt && !isNaN(new Date(job.postedAt).getTime()) 
                ? formatDistanceToNow(new Date(job.postedAt), { addSuffix: true })
                : 'Recently'}
            </span>
          </div>
        </div>

        <div className="flex flex-wrap gap-2 mb-6">
          {job.tags.slice(0, 4).map((tag, i) => (
            <span 
              key={i} 
              className="px-2.5 py-0.5 text-xs font-medium text-gray-600 bg-gray-100 rounded-md border border-gray-200"
            >
              {tag}
            </span>
          ))}
          {job.tags.length > 4 && (
             <span className="px-2.5 py-0.5 text-xs font-medium text-gray-400 bg-gray-50 rounded-md">
               +{job.tags.length - 4} more
             </span>
          )}
        </div>
      </div>

      <a 
        href={job.url} 
        target="_blank" 
        rel="noopener noreferrer"
        className="block w-full bg-gray-900 hover:bg-blue-600 text-white text-center font-semibold py-2.5 rounded-lg transition-colors flex items-center justify-center group-hover:shadow-md"
      >
        Apply Now
        <ExternalLink className="w-4 h-4 ml-2" />
      </a>
    </div>
  );
};

export default JobCard;
