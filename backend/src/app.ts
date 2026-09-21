import cors from 'cors';
import express from 'express';
import helmet from 'helmet';

const app = express();

import jobRoutes from './routes/jobRoutes';

// Middleware
app.use(express.json({ limit: '50mb' })); // Increase limit for batches
app.use(express.urlencoded({ extended: true }));
app.use(cors({
  origin: (origin, callback) => {
    const allowList = [
      'http://localhost:3000',
      'http://localhost:5000',
      'https://scraphire.vercel.app',
      'https://job-aggregator-1vi4.onrender.com',
    ];
    if (!origin) {
      // curl / Render health ping / same-origin calls all have no Origin header — allow them
      return callback(null, true);
    }
    if (allowList.includes(origin)) {
      return callback(null, true);
    }
    // Allow any Vercel preview deployment (https://<name>-<hash>-<org>.vercel.app)
    if (/^https:\/\/[-a-z0-9]+\.vercel\.app$/i.test(origin)) {
      return callback(null, true);
    }
    // Deny unknown origins
    return callback(new Error(`CORS blocked origin: ${origin}`), false);
  },
  credentials: true,
  optionsSuccessStatus: 200,
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
