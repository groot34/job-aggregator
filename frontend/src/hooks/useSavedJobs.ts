import { useState, useEffect } from 'react';

export function useSavedJobs() {
  const [savedJobIds, setSavedJobIds] = useState<string[]>([]);
  const [isLoaded, setIsLoaded] = useState(false);

  // Load from localStorage on mount
  useEffect(() => {
    try {
      const stored = localStorage.getItem('savedJobs');
      if (stored) {
        setSavedJobIds(JSON.parse(stored));
      }
    } catch (e) {
      console.error('Error parsing savedJobs from localStorage', e);
    }
    setIsLoaded(true);
  }, []);

  // Toggle save status
  const toggleSave = (id: string) => {
    setSavedJobIds((prev) => {
      let next;
      if (prev.includes(id)) {
        next = prev.filter((jobId) => jobId !== id);
      } else {
        next = [...prev, id];
      }
      try {
        localStorage.setItem('savedJobs', JSON.stringify(next));
      } catch (e) {
        console.error('Error saving to localStorage', e);
      }
      return next;
    });
  };

  const isSaved = (id: string) => savedJobIds.includes(id);

  return { savedJobIds, toggleSave, isSaved, isLoaded };
}
