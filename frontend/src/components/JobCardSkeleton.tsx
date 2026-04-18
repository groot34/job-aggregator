import React from 'react';

export default function JobCardSkeleton() {
  return (
    <div className="bg-white border border-gray-200 rounded-xl p-6 h-full flex flex-col justify-between animate-pulse">
      <div>
        <div className="flex justify-between items-start mb-4">
          <div className="w-full">
            {/* Source Pill */}
            <div className="h-5 w-16 bg-gray-200 rounded-full mb-3"></div>
            {/* Title */}
            <div className="h-6 w-3/4 bg-gray-200 rounded-md mb-2"></div>
            <div className="h-6 w-1/2 bg-gray-200 rounded-md mb-3"></div>
            {/* Company */}
            <div className="flex items-center space-x-2">
              <div className="h-4 w-4 bg-gray-200 rounded"></div>
              <div className="h-4 w-32 bg-gray-200 rounded-md"></div>
            </div>
          </div>
        </div>

        {/* Location & Time */}
        <div className="flex items-center space-x-4 mb-5">
          <div className="flex items-center space-x-2">
            <div className="h-4 w-4 bg-gray-200 rounded"></div>
            <div className="h-4 w-20 bg-gray-200 rounded-md"></div>
          </div>
          <div className="flex items-center space-x-2">
            <div className="h-4 w-4 bg-gray-200 rounded"></div>
            <div className="h-4 w-20 bg-gray-200 rounded-md"></div>
          </div>
        </div>

        {/* Tags */}
        <div className="flex flex-wrap gap-2 mb-6">
          <div className="h-5 w-12 bg-gray-200 rounded-md"></div>
          <div className="h-5 w-16 bg-gray-200 rounded-md"></div>
          <div className="h-5 w-14 bg-gray-200 rounded-md"></div>
        </div>
      </div>

      {/* Button */}
      <div className="h-11 w-full bg-gray-200 rounded-lg"></div>
    </div>
  );
}
