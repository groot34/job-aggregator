import app from './app';
import mongoose from 'mongoose';
import dotenv from 'dotenv';

dotenv.config();

const PORT = process.env.PORT || 5000;
const MONGO_URI = process.env.MONGO_URI || 'mongodb://localhost:27017/job-aggregator';

const HOST = process.env.HOST || process.env.NODE_ENV === 'production' ? '0.0.0.0' : 'localhost';

// MongoDB connection with caching for serverless + long-running hosts
let isConnected = false;

const connectDB = async () => {
  if (isConnected) {
    return;
  }

  try {
    await mongoose.connect(MONGO_URI);
    isConnected = true;
    console.log('✅ Connected to MongoDB');
  } catch (error) {
    console.error('❌ Failed to connect to MongoDB:', error);
    throw error;
  }
};

// Start a long-running HTTP server for Render / local dev / any classic Node host.
// (NODE_ENV=production on Render does NOT mean "serverless" — it still needs listen().)
const startServer = async () => {
  try {
    await connectDB();
    app.listen(Number(PORT), HOST, () => {
      console.log(`🚀 Server running on http://${HOST}:${PORT} (NODE_ENV=${process.env.NODE_ENV ?? 'development'})`);
    });
  } catch (error) {
    console.error('💥 Failed to start server:', error);
    process.exit(1);
  }
};

// Vercel serverless runtime provides `req`/`res` via ESM exports AND sets
// VERCEL=1. On any other host (Render, local, Docker, etc.) start listening.
if (!process.env.VERCEL) {
  startServer();
}

// Vercel serverless: Export handler
export default async (req: any, res: any) => {
  await connectDB();
  return app(req, res);
};
