import { Router } from 'express';
import { createJobsBatch, getJobs, cleanupOldJobs } from '../controllers/jobController';

const router = Router();

router.post('/batch', createJobsBatch);
router.get('/', getJobs);
router.delete('/cleanup', cleanupOldJobs);

export default router;
