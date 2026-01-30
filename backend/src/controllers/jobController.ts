import { Request, Response } from 'express';
import Job from '../models/Job';

// POST /api/jobs/batch
// Receive a batch of jobs from the scraper
export const createJobsBatch = async (req: Request, res: Response) => {
  try {
    const jobs = req.body;
    if (!Array.isArray(jobs)) {
      return res.status(400).json({ message: 'Payload must be an array of jobs' });
    }

    const operations = jobs.map((job) => ({
      updateOne: {
        filter: { externalId: job.externalId }, // potential dup check by ID
        update: { $set: job },
        upsert: true,
      },
    }));

    if (operations.length > 0) {
      const result = await Job.bulkWrite(operations);
      console.log(`📦 Processed ${jobs.length} jobs: ${result.upsertedCount} inserted, ${result.modifiedCount} updated.`);
      return res.status(200).json({ 
        message: 'Batch processed', 
        inserted: result.upsertedCount, 
        updated: result.modifiedCount 
      });
    }

    res.status(200).json({ message: 'No jobs to process' });
  } catch (error) {
    console.error('Error in createJobsBatch:', error);
    res.status(500).json({ message: 'Internal Server Error' });
  }
};

// GET /api/jobs
export const getJobs = async (req: Request, res: Response) => {
  try {
    const { page = 1, limit = 20, tag, search } = req.query;
    
    const query: any = { active: true };
    if (tag) {
      query.tags = tag;
    }
    if (search) {
      query.$or = [
        { title: { $regex: search, $options: 'i' } },
        { company: { $regex: search, $options: 'i' } }
      ];
    }

    const jobs = await Job.find(query)
      .sort({ postedAt: -1 })
      .skip((Number(page) - 1) * Number(limit))
      .limit(Number(limit));

    const total = await Job.countDocuments(query);
   
    
    res.json({
      jobs,
      total,
      currentPage: Number(page),
      totalPages: Math.ceil(total / Number(limit)),
    });
  } catch (error) {
    res.status(500).json({ message: 'Error fetching jobs' });
  }
};
