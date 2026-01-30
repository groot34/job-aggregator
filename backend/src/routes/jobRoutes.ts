import { Router } from 'express';
import { createJobsBatch, getJobs } from '../controllers/jobController';

const router = Router();

router.post('/batch', createJobsBatch);
router.get('/', getJobs);

export default router;
