/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  experimental: {
    webpackMemoryOptimizations: true,
  },
}

export default nextConfig
