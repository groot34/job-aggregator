import express from 'express';
import cors from 'cors';
import helmet from 'helmet';

const app = express();

import jobRoutes from './routes/jobRoutes';

// Middleware
app.use(express.json({ limit: '50mb' })); // Increase limit for batches
app.use(express.urlencoded({ extended: true }));
app.use(cors({
  origin: [
    'http://localhost:3000',
    'https://scraphire.vercel.app', // Your actual Vercel URL
    'https://*.vercel.app' // All Vercel preview deployments
  ],
  credentials: true
}));
app.use(helmet());

// Routes
app.use('/api/jobs', jobRoutes);

// Root Endpoint
app.get('/', (req, res) => {
  res.status(200).json({ 
    message: 'Job Aggregator API is running 🚀', 
    service: 'Job Aggregator Backend',
    endpoints: {
      health: '/health',
      jobs: '/api/jobs'
    }
  });
});

// Health Check
app.get('/health', (req, res) => {
  res.status(200).json({ status: 'UP', message: 'Backend is running' });
});

export default app;
