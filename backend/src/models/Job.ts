import mongoose, { Schema, Document } from 'mongoose';

export interface IJob extends Document {
  externalId: string;
  title: string;
  company: string;
  location: string;
  description: string;
  url: string;
  source: string;
  postedAt: Date;
  scrapedAt: Date;
  remote: boolean;
  salary?: string;
  tags: string[];
  active: boolean;
}

const JobSchema: Schema = new Schema({
  externalId: { type: String, required: true, unique: true, index: true },
  title: { type: String, required: true },
  company: { type: String, required: true },
  location: { type: String, required: true },
  description: { type: String },
  url: { type: String, required: true },
  source: { type: String, required: true },
  postedAt: { type: Date, required: true },
  scrapedAt: { type: Date, default: Date.now },
  remote: { type: Boolean, default: false },
  salary: { type: String },
  tags: { type: [String], index: true },
  active: { type: Boolean, default: true },
});

// Composite index for deduplication safety
JobSchema.index({ title: 1, company: 1 }, { unique: false }); 

export default mongoose.model<IJob>('Job', JobSchema);
