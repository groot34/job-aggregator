import { Request, Response } from 'express';
import Job from '../models/Job';

const SENIOR_TITLE_REGEX =
  /\b(senior|sr\.?|staff|principal|lead|manager|head|director|architect|vp|president|supervisor|expert|specialist|consultant|sde\s*(2|ii|3|iii)|software engineer\s*(ii|iii|2|3)|developer\s*(ii|iii|2|3)|engineer\s*(ii|iii|2|3))\b/i;

const HIGHER_EXPERIENCE_REGEX =
  /\b((?:2|3|4|5|6|7|8|9|10)\+?\s*(?:year|years)|minimum\s*(?:2|3|4|5|6|7|8|9|10)\s*(?:year|years)|at least\s*(?:2|3|4|5|6|7|8|9|10)\s*(?:year|years)|(?:2|3|4|5|6|7|8|9|10)\s*(?:year|years)\s+of experience|(?:2|3|4|5|6|7|8|9|10)\s*(?:year|years)\s+experience)\b/i;

// POST /api/jobs/batch
// Receive a batch of jobs from the scraper
export const createJobsBatch = async (req: Request, res: Response) => {
  try {
    const jobs = req.body;
    if (!Array.isArray(jobs)) {
      return res.status(400).json({ message: 'Payload must be an array of jobs' });
    }

    const operations = jobs.map((job) => {
      // Fields that should be written only on first insert and never overwritten
      // by subsequent re-scrapes (otherwise postedAt keeps jumping to "Now"
      // every run and we lose the real age of a listing).
      const setOnInsert: Record<string, unknown> = {
        externalId: job.externalId,
        postedAt: job.postedAt ?? new Date(),
      };

      // Fields that may legitimately change between scrapes (tags improve,
      // salary/description/location are refreshed, etc.)
      const $set: Record<string, unknown> = {
        title: job.title,
        company: job.company,
        location: job.location ?? '',
        description: job.description ?? '',
        url: job.url,
        source: job.source,
        scrapedAt: job.scrapedAt ?? new Date(),
        remote: job.remote ?? false,
        active: true,
      };
      if (job.salary !== undefined && job.salary !== null) $set.salary = job.salary;
      if (Array.isArray(job.tags)) $set.tags = job.tags;

      return {
        updateOne: {
          filter: { externalId: job.externalId },
          update: { $setOnInsert: setOnInsert, $set },
          upsert: true,
        },
      };
    });

    const now = new Date();
    const sourcesInBatch = [...new Set(jobs.map((j) => j.source).filter(Boolean) as string[])];
    const externalIdsInBatch = [...new Set(jobs.map((j) => j.externalId).filter(Boolean) as string[])];

    if (operations.length > 0) {
      const result = await Job.bulkWrite(operations);
      console.log(`📦 Processed ${jobs.length} jobs: ${result.upsertedCount} inserted, ${result.modifiedCount} updated.`);

      // ---------------------------------------------------------------------
      // 1. Mark-stale: any job from the scraped sources that was NOT re-seen
      //    in this batch AND has not been scraped in the last ~36h (≈3 cron
      //    cycles at 12h intervals) is very likely "no longer on the site".
      // ---------------------------------------------------------------------
      const recentlyScrapedCutoff = new Date(now.getTime() - 36 * 60 * 60 * 1000);
      let staled = 0;
      if (sourcesInBatch.length > 0) {
        const staleRes = await Job.updateMany(
          {
            source: { $in: sourcesInBatch },
            active: true,
            scrapedAt: { $lt: recentlyScrapedCutoff },
            ...(externalIdsInBatch.length > 0 ? { externalId: { $nin: externalIdsInBatch } } : {}),
          },
          { $set: { active: false } }
        );
        staled = staleRes.modifiedCount ?? 0;
        if (staled > 0) {
          console.log(`🧽 Marked ${staled} jobs INACTIVE (disappeared from source and not re-scraped since ${recentlyScrapedCutoff.toISOString()})`);
        }
      }

      // ---------------------------------------------------------------------
      // 2. Hard-delete ceiling: anything older than 45 days regardless, plus
      //    any already-inactive job that was last updated more than 14d ago.
      //    Prevents the DB from keeping listings that died months ago.
      // ---------------------------------------------------------------------
      const postedCutoff = new Date(now.getTime() - 45 * 24 * 60 * 60 * 1000);
      const inactiveUpdatedCutoff = new Date(now.getTime() - 14 * 24 * 60 * 60 * 1000);
      const deletedRes = await Job.deleteMany({
        $or: [
          { postedAt: { $lt: postedCutoff } },
          { active: false, updatedAt: { $lt: inactiveUpdatedCutoff } },
        ],
      });
      if (deletedRes.deletedCount) {
        console.log(`🗑️  Hard-deleted ${deletedRes.deletedCount} jobs (posted<45d OR inactive>14d)`);
      }

      return res.status(200).json({
        message: 'Batch processed',
        inserted: result.upsertedCount,
        updated: result.modifiedCount,
        markedInactive: staled,
        hardDeleted: deletedRes.deletedCount ?? 0,
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
    const { page = 1, limit = 20, tag, search, remote, location, savedIds } = req.query;
    
    const query: any = {
      active: true,
      title: { $not: SENIOR_TITLE_REGEX },
      description: { $not: HIGHER_EXPERIENCE_REGEX },
    };
    
    if (savedIds) {
      const ids = (savedIds as string).split(',').filter(Boolean);
      if (ids.length > 0) {
        query._id = { $in: ids };
      }
    }

    if (tag) {
      query.tags = tag; // Wait, actually should probably handle multiple tags, maybe later, but let's just stick to single tag or whatever UI sends
    }
    
    if (remote === 'true') {
      query.remote = true;
    }
    
    if (location) {
      query.location = { $regex: location, $options: 'i' };
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

// DELETE /api/jobs/cleanup
// Remove jobs older than 30 days
export const cleanupOldJobs = async (req: Request, res: Response) => {
  try {
    const thirtyDaysAgo = new Date();
    thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

    const result = await Job.deleteMany({ postedAt: { $lt: thirtyDaysAgo } });

    console.log(`🧹 Cleanup: Removed ${result.deletedCount} jobs older than 30 days`);
    res.status(200).json({
      message: 'Cleanup complete',
      deleted: result.deletedCount,
      cutoffDate: thirtyDaysAgo.toISOString()
    });
  } catch (error) {
    console.error('Error in cleanupOldJobs:', error);
    res.status(500).json({ message: 'Error during cleanup' });
  }
};
