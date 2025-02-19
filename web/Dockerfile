# Use Node 20 Alpine as base
FROM node:20-alpine

WORKDIR /app

# Copy package files and install ALL dependencies (including dev)
COPY package.json package-lock.json ./
RUN npm install

# Copy all source files
COPY . .

# Ensure node_modules/.bin is in the PATH (so CLI tools like react-router work)
ENV PATH="/app/node_modules/.bin:$PATH"
ENV NODE_ENV=development

# Expose React's development server port
EXPOSE 3000

# Start the development server
CMD ["npm", "run", "dev"]
