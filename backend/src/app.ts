import express from 'express';
import cors from 'cors';
import helmet from 'helmet';

const app = express();

import jobRoutes from './routes/jobRoutes';

// Middleware
app.use(express.json({ limit: '50mb' })); // Increase limit for batches
app.use(express.urlencoded({ extended: true }));
app.use(cors());
app.use(helmet());

// Routes
app.use('/api/jobs', jobRoutes);

// Health Check
app.get('/health', (req, res) => {
  res.status(200).json({ status: 'UP', message: 'Backend is running' });
});

export default app;
