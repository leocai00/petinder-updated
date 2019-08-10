# Use `node` as the base image
FROM node

# Set the current working directory to `/app`
WORKDIR /messaging

# Copy the `package.json` and `package-lock.json` files to the working directory
COPY package.json package-lock.json ./

# Run `npm install` to install all dependent packages
RUN npm install

# Copy your JavaScript source files
COPY . .

# Declare that your server will listen on port 80
EXPOSE 80

# Set the `ENTRYPOINT` command
ENTRYPOINT ["node", "index.js"]