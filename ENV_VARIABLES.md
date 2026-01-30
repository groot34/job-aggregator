# 📋 Environment Variables Guide

This document lists all environment variables needed for deployment.

## 🔙 Backend (Render)

Set these in **Render Dashboard** → Your Service → **Environment**:

```env
PORT=5000
NODE_ENV=production
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/job-aggregator?retryWrites=true&w=majority
```

**Where to get values:**
- `MONGO_URI`: From MongoDB Atlas → Connect → Drivers

---

## 🎨 Frontend (Vercel)

Set these in **Vercel Dashboard** → Your Project → **Settings** → **Environment Variables**:

```env
NEXT_PUBLIC_API_URL=https://job-aggregator-backend.onrender.com/api
```

**Where to get values:**
- `NEXT_PUBLIC_API_URL`: Your Render backend URL (from Step 6 of deployment guide)

---

## 🕷️ Scrapers (GitHub Actions - Optional)

Set these in **GitHub** → Repository Settings → **Secrets and Variables** → **Actions**:

```env
BACKEND_URL=https://job-aggregator-backend.onrender.com/api/jobs/batch
```

**Where to get values:**
- `BACKEND_URL`: Your Render backend URL + `/jobs/batch`

---

## 🚨 Important Notes

### ❌ DO NOT commit actual .env files
- `.env` (sensitive!)
- `.env.local` (sensitive!)
- `.env.production` (sensitive!)

### ✅ DO commit template files
- `.env.example` (safe - no secrets)
- `.env.production.example` (safe - placeholder values)

### 🔐 Security Best Practices
1. Never commit files with real API keys/passwords
2. Use `.env.example` files as templates
3. Store actual secrets in:
   - **Render**: Environment variables section
   - **Vercel**: Environment variables section
   - **GitHub**: Repository secrets
4. Each developer maintains their own `.env.local` file

---

## 📝 Local Development Setup

For new developers cloning the repo:

1. **Backend**:
   ```bash
   cd backend
   cp .env.example .env
   # Edit .env with your local MongoDB connection
   ```

2. **Frontend**:
   ```bash
   cd frontend
   # No .env needed for local dev (uses defaults)
   # Or create .env.local if you want to point to remote backend
   ```

3. **Scrapers**:
   ```bash
   cd scrapers
   # No .env needed (uses default localhost:5000)
   # Or create .env if testing against remote backend
   ```

---

## ✅ Verification Checklist

Before deploying, ensure:

- [ ] All `.env` files are in `.gitignore`
- [ ] `.env.example` files are committed
- [ ] No secrets in `.env.example` files
- [ ] Production values set in Render/Vercel dashboards
- [ ] GitHub secrets configured (if using Actions)

---

## 🆘 If You Accidentally Commit Secrets

1. **Immediately rotate all exposed credentials**:
   - Change MongoDB password in Atlas
   - Regenerate any API keys
   
2. **Remove from Git history**:
   ```bash
   git filter-branch --force --index-filter \
     "git rm --cached --ignore-unmatch backend/.env" \
     --prune-empty --tag-name-filter cat -- --all
   ```

3. **Force push** (⚠️ dangerous):
   ```bash
   git push origin --force --all
   ```

4. **Better**: Use [BFG Repo-Cleaner](https://rtyley.github.io/bfg-repo-cleaner/)
